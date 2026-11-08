package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

// CSRFTokenStore defines the interface for storing CSRF tokens
type CSRFTokenStore interface {
	// Generate creates a new CSRF token for the given identifier
	Generate(identifier string) (string, error)
	// Validate checks if the provided token is valid for the identifier
	Validate(identifier, token string) bool
	// Delete removes the token for the identifier
	Delete(identifier string)
}

// MemoryCSRFTokenStore is an in-memory implementation of CSRFTokenStore
type MemoryCSRFTokenStore struct {
	mu         sync.RWMutex
	tokens     map[string]csrfToken
	tokenLen   int
	tokenExpiry time.Duration
}

type csrfToken struct {
	token     string
	expiresAt time.Time
}

// NewMemoryCSRFTokenStore creates a new in-memory CSRF token store
func NewMemoryCSRFTokenStore(tokenLen int, expiry time.Duration) *MemoryCSRFTokenStore {
	store := &MemoryCSRFTokenStore{
		tokens:     make(map[string]csrfToken),
		tokenLen:   tokenLen,
		tokenExpiry: expiry,
	}
	// Start cleanup goroutine
	go store.cleanup()
	return store
}

// Generate creates a new CSRF token
func (s *MemoryCSRFTokenStore) Generate(identifier string) (string, error) {
	tokenBytes := make([]byte, s.tokenLen)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[identifier] = csrfToken{
		token:     token,
		expiresAt: time.Now().Add(s.tokenExpiry),
	}

	return token, nil
}

// Validate checks if the provided token is valid
func (s *MemoryCSRFTokenStore) Validate(identifier, token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stored, exists := s.tokens[identifier]
	if !exists {
		return false
	}

	// Check if token expired
	if time.Now().After(stored.expiresAt) {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(stored.token), []byte(token)) == 1
}

// Delete removes the token
func (s *MemoryCSRFTokenStore) Delete(identifier string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, identifier)
}

// cleanup removes expired tokens periodically
func (s *MemoryCSRFTokenStore) cleanup() {
	ticker := time.NewTicker(s.tokenExpiry)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, token := range s.tokens {
			if now.After(token.expiresAt) {
				delete(s.tokens, id)
			}
		}
		s.mu.Unlock()
	}
}

// CSRF returns a middleware that provides CSRF protection
func CSRF(config CSRFConfig, store CSRFTokenStore) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF check for safe methods
			if IsMethodSafe(r.Method) {
				// For GET requests, we might want to set the CSRF token in a cookie
				if r.Method == http.MethodGet {
					identifier := getClientIdentifier(r, config.TrustProxy)
					if token, err := store.Generate(identifier); err == nil {
						cookie := &http.Cookie{
							Name:     config.CookieName,
							Value:    token,
							Path:     "/",
							Secure:   config.CookieSecure,
							HttpOnly: config.CookieHttpOnly,
							SameSite: http.SameSite(config.CookieSameSite),
							MaxAge:   int(config.TokenExpiry.Seconds()),
						}
						if config.CookieDomain != "" {
							cookie.Domain = config.CookieDomain
						}
						http.SetCookie(w, cookie)
					}
				}
				next.ServeHTTP(w, r)
				return
			}

			// For unsafe methods, validate CSRF token
			identifier := getClientIdentifier(r, config.TrustProxy)

			// Try to get token from header first
			headerToken := r.Header.Get(config.HeaderName)

			// If not in header, try cookie
			var cookieToken string
			if headerToken == "" {
				if cookie, err := r.Cookie(config.CookieName); err == nil {
					cookieToken = cookie.Value
				}
			}

			// Validate token
			token := headerToken
			if token == "" {
				token = cookieToken
			}

			if token == "" || !store.Validate(identifier, token) {
				// CSRF token missing or invalid
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error":"CSRF token missing or invalid"}`))
				return
			}

			// Token is valid, continue
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIdentifier extracts a unique identifier for the client
func getClientIdentifier(r *http.Request, trustProxy bool) string {
	if trustProxy {
		// Check X-Forwarded-For
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		// Check X-Real-IP
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			return strings.TrimSpace(realIP)
		}
	}
	return r.RemoteAddr
}

// GenerateCSRFToken generates a CSRF token and returns it as a string
func GenerateCSRFToken(store CSRFTokenStore, identifier string) (string, error) {
	return store.Generate(identifier)
}

// ValidateCSRFToken validates a CSRF token
func ValidateCSRFToken(store CSRFTokenStore, identifier, token string) bool {
	return store.Validate(identifier, token)
}

// CSRFTokenFromRequest extracts the CSRF token from a request
func CSRFTokenFromRequest(r *http.Request, config CSRFConfig) string {
	// Try header first
	if token := r.Header.Get(config.HeaderName); token != "" {
		return token
	}
	// Try cookie
	if cookie, err := r.Cookie(config.CookieName); err == nil {
		return cookie.Value
	}
	// Try form value
	if r.FormValue(config.CookieName) != "" {
		return r.FormValue(config.CookieName)
	}
	return ""
}

// ErrCSRFTokenInvalid is returned when CSRF token validation fails
var ErrCSRFTokenInvalid = errors.New("csrf token invalid")

// ErrCSRFTokenMissing is returned when CSRF token is missing
var ErrCSRFTokenMissing = errors.New("csrf token missing")

// SimpleCSRF returns a simple CSRF middleware with default settings
func SimpleCSRF() func(http.Handler) http.Handler {
	config := DefaultCSRFConfig()
	store := NewMemoryCSRFTokenStore(config.TokenLength, config.TokenExpiry)
	return CSRF(config, store)
}

// CSRFHandler is a helper to generate and validate CSRF tokens
// It can be used in HTTP handlers to provide CSRF token endpoints
type CSRFHandler struct {
	store   CSRFTokenStore
	config  CSRFConfig
	getID   func(r *http.Request) string
}

// NewCSRFHandler creates a new CSRF handler
func NewCSRFHandler(store CSRFTokenStore, config CSRFConfig) *CSRFHandler {
	return &CSRFHandler{
		store:  store,
		config: config,
		getID: func(r *http.Request) string {
			return getClientIdentifier(r, config.TrustProxy)
		},
	}
}

// GetToken generates and returns a CSRF token
func (h *CSRFHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	identifier := h.getID(r)
	token, err := h.store.Generate(identifier)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to generate CSRF token"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"csrf_token":"` + token + `"}`))
}

// ValidateToken validates a CSRF token
func (h *CSRFHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	identifier := h.getID(r)
	token := CSRFTokenFromRequest(r, h.config)

	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"CSRF token missing"}`))
		return
	}

	if !h.store.Validate(identifier, token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"CSRF token invalid"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid":true}`))
}
