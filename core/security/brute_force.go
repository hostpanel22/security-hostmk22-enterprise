package security

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

// BruteForceStore defines the interface for storing brute force attempt data
type BruteForceStore interface {
	// RecordAttempt records a failed authentication attempt
	RecordAttempt(identifier string) error
	// GetAttempts returns the number of attempts for an identifier
	GetAttempts(identifier string) (int, error)
	// IsLocked checks if an identifier is locked out
	IsLocked(identifier string) (bool, error)
	// Lock locks an identifier
	Lock(identifier string, duration time.Duration) error
	// Unlock unlocks an identifier
	Unlock(identifier string) error
	// Cleanup removes expired lockouts
	Cleanup()
}

// MemoryBruteForceStore is an in-memory implementation of BruteForceStore
type MemoryBruteForceStore struct {
	mu          sync.RWMutex
	attempts    map[string][]time.Time
	lockouts    map[string]time.Time
	maxAttempts int
	lockoutDur  time.Duration
	windowDur   time.Duration
}

// NewMemoryBruteForceStore creates a new in-memory brute force store
func NewMemoryBruteForceStore(maxAttempts int, lockoutDur, windowDur time.Duration) *MemoryBruteForceStore {
	store := &MemoryBruteForceStore{
		attempts:    make(map[string][]time.Time),
		lockouts:    make(map[string]time.Time),
		maxAttempts: maxAttempts,
		lockoutDur:  lockoutDur,
		windowDur:   windowDur,
	}
	go store.Cleanup()
	return store
}

// RecordAttempt records a failed authentication attempt
func (s *MemoryBruteForceStore) RecordAttempt(identifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.attempts[identifier] = append(s.attempts[identifier], now)

	// Clean old attempts
	cutoff := now.Add(-s.windowDur)
	var filtered []time.Time
	for _, ts := range s.attempts[identifier] {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	s.attempts[identifier] = filtered

	return nil
}

// GetAttempts returns the number of attempts for an identifier
func (s *MemoryBruteForceStore) GetAttempts(identifier string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-s.windowDur)
	count := 0

	for _, ts := range s.attempts[identifier] {
		if ts.After(cutoff) {
			count++
		}
	}

	return count, nil
}

// IsLocked checks if an identifier is locked out
func (s *MemoryBruteForceStore) IsLocked(identifier string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lockoutTime, exists := s.lockouts[identifier]
	if !exists {
		return false, nil
	}

	return time.Now().Before(lockoutTime), nil
}

// Lock locks an identifier
func (s *MemoryBruteForceStore) Lock(identifier string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockouts[identifier] = time.Now().Add(duration)
	return nil
}

// Unlock unlocks an identifier
func (s *MemoryBruteForceStore) Unlock(identifier string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.lockouts, identifier)
	return nil
}

// Cleanup removes expired lockouts periodically
func (s *MemoryBruteForceStore) Cleanup() {
	ticker := time.NewTicker(s.lockoutDur)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, lockoutTime := range s.lockouts {
			if now.After(lockoutTime) {
				delete(s.lockouts, id)
			}
		}
		s.mu.Unlock()
	}
}

// BruteForceProtection returns a middleware that protects against brute force attacks
func BruteForceProtection(config BruteForceConfig, store BruteForceStore) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this path should be protected
			path := r.URL.Path
			if !shouldProtectPath(path, config.ProtectedPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// Get client identifier
			identifier := getBruteForceIdentifier(r, config.TrustProxy)

			// Check if locked out
			locked, err := store.IsLocked(identifier)
			if err != nil || locked {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", formatDuration(config.LockoutDuration))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"too many attempts, please try again later"}`))
				return
			}

			// For POST requests to protected paths, we need to check after the request
			// But for now, we'll just pass through and let the handler record attempts
			
			// Add identifier to context for the handler to use
			ctx := WithBruteForceID(r.Context(), identifier)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// shouldProtectPath checks if a path should be protected
func shouldProtectPath(path string, protectedPaths []string) bool {
	for _, protected := range protectedPaths {
		if strings.HasPrefix(path, protected) || path == protected {
			return true
		}
	}
	return false
}

// getBruteForceIdentifier extracts a unique identifier for brute force tracking
func getBruteForceIdentifier(r *http.Request, trustProxy bool) string {
	if trustProxy {
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			if len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			return strings.TrimSpace(realIP)
		}
	}
	return r.RemoteAddr
}

// contextKey for brute force ID
const bfContextKey contextKey = "brute_force_id"

// WithBruteForceID adds the brute force identifier to the request context
func WithBruteForceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, bfContextKey, id)
}

// GetBruteForceID retrieves the brute force identifier from the request context
func GetBruteForceID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(bfContextKey).(string)
	return id, ok
}

// BruteForceHandler provides helper methods for handlers to use brute force protection
type BruteForceHandler struct {
	store   BruteForceStore
	config  BruteForceConfig
	getID   func(r *http.Request) string
}

// NewBruteForceHandler creates a new brute force handler
func NewBruteForceHandler(store BruteForceStore, config BruteForceConfig) *BruteForceHandler {
	return &BruteForceHandler{
		store:  store,
		config: config,
		getID: func(r *http.Request) string {
			return getBruteForceIdentifier(r, config.TrustProxy)
		},
	}
}

// CheckAndRecord checks if the request should be allowed and records attempts
func (h *BruteForceHandler) CheckAndRecord(w http.ResponseWriter, r *http.Request, isSuccess bool) bool {
	identifier := h.getID(r)

	// Check if locked out
	locked, _ := h.store.IsLocked(identifier)
	if locked {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", formatDuration(h.config.LockoutDuration))
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":"too many attempts, please try again later"}`))
		return false
	}

	// If authentication failed, record attempt
	if !isSuccess {
		h.store.RecordAttempt(identifier)

		// Check if we've exceeded max attempts
		attempts, _ := h.store.GetAttempts(identifier)
		if attempts >= h.config.MaxAttempts {
			h.store.Lock(identifier, h.config.LockoutDuration)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", formatDuration(h.config.LockoutDuration))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"too many attempts, please try again later"}`))
			return false
		}
	}

	return true
}

// RecordFailedAttempt records a failed authentication attempt
func (h *BruteForceHandler) RecordFailedAttempt(r *http.Request) {
	identifier := h.getID(r)
	h.store.RecordAttempt(identifier)

	// Check if we've exceeded max attempts
	attempts, _ := h.store.GetAttempts(identifier)
	if attempts >= h.config.MaxAttempts {
		h.store.Lock(identifier, h.config.LockoutDuration)
	}
}

// RecordSuccessfulAttempt clears failed attempts after a successful authentication
func (h *BruteForceHandler) RecordSuccessfulAttempt(r *http.Request) {
	identifier := h.getID(r)
	// Clear attempts on success
	h.store.Unlock(identifier)
}

// GetRemainingAttempts returns the number of remaining attempts before lockout
func (h *BruteForceHandler) GetRemainingAttempts(r *http.Request) int {
	identifier := h.getID(r)
	attempts, _ := h.store.GetAttempts(identifier)
	return h.config.MaxAttempts - attempts
}

// SimpleBruteForceProtection returns a simple brute force protection middleware
func SimpleBruteForceProtection(maxAttempts int, lockoutDuration time.Duration) func(http.Handler) http.Handler {
	config := DefaultBruteForceConfig()
	config.MaxAttempts = maxAttempts
	config.LockoutDuration = lockoutDuration
	store := NewMemoryBruteForceStore(maxAttempts, lockoutDuration, config.WindowDuration)
	return BruteForceProtection(config, store)
}

// BruteForceResponseWriter wraps http.ResponseWriter to track authentication results
type BruteForceResponseWriter struct {
	http.ResponseWriter
	handler  *BruteForceHandler
	request  *http.Request
	status   int
	isAuth   bool
}

// NewBruteForceResponseWriter creates a new brute force response writer
func NewBruteForceResponseWriter(w http.ResponseWriter, handler *BruteForceHandler, r *http.Request) *BruteForceResponseWriter {
	return &BruteForceResponseWriter{
		ResponseWriter: w,
		handler:        handler,
		request:        r,
		status:         http.StatusOK,
		isAuth:         false,
	}
}

// WriteHeader implements http.ResponseWriter
func (w *BruteForceResponseWriter) WriteHeader(code int) {
	w.status = code
	w.isAuth = code == http.StatusOK || code == http.StatusCreated
	w.ResponseWriter.WriteHeader(code)
}

// Write implements http.ResponseWriter
func (w *BruteForceResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Finish records the attempt result
func (w *BruteForceResponseWriter) Finish() {
	if w.isAuth {
		w.handler.RecordSuccessfulAttempt(w.request)
	} else if w.status == http.StatusUnauthorized {
		w.handler.RecordFailedAttempt(w.request)
	}
}
