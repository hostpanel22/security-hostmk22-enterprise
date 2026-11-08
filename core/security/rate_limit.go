package security

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RateLimit returns a middleware that limits the number of requests per client
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	limiter := newClientRateLimiter()

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(config.WindowDuration)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanExpired(time.Now(), config.WindowDuration)
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get tenant ID if available
			tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
			
			// Get client identifier
			clientID := getRateLimitClientID(r, config.TrustProxy)
			
			// Add tenant to client ID if present
			if tenantID != "" {
				clientID = tenantID + ":" + clientID
			}
			
			// Check if client is in whitelist
			if isInList(clientID, config.Whitelist) {
				next.ServeHTTP(w, r)
				return
			}
			
			// Check if client is in blacklist
			if isInList(clientID, config.Blacklist) {
				w.Header().Set("Retry-After", strconv.Itoa(int(config.WindowDuration.Seconds())))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			
			// Check rate limit
			if !limiter.allow(clientID, config.RequestsPerMinute, config.WindowDuration) {
				w.Header().Set("Retry-After", strconv.Itoa(int(config.WindowDuration.Seconds())))
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(config.WindowDuration).Unix())))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			
			// Add rate limit headers to successful responses
			remaining, reset := limiter.getStats(clientID, config.RequestsPerMinute, config.WindowDuration)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(reset.Unix())))
			
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitProvider defines the interface for providing rate limit configuration
type RateLimitProvider interface {
	RateLimit(ctx context.Context, tenantID string) RateLimitConfig
}

// clientRateLimiter manages rate limiting for multiple clients
type clientRateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
}

// newClientRateLimiter creates a new client rate limiter
func newClientRateLimiter() *clientRateLimiter {
	return &clientRateLimiter{
		clients: make(map[string][]time.Time),
	}
}

// allow checks if a request from the given client should be allowed
func (rl *clientRateLimiter) allow(clientID string, limit int, window time.Duration) bool {
	if clientID == "" {
		clientID = "unknown"
	}
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	now := time.Now()
	cutoff := now.Add(-window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	timestamps := rl.clients[clientID]
	filtered := timestamps[:0]
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= limit {
		if len(filtered) == 0 {
			delete(rl.clients, clientID)
		} else {
			rl.clients[clientID] = filtered
		}
		return false
	}

	filtered = append(filtered, now)
	rl.clients[clientID] = filtered
	return true
}

// getStats returns the current rate limit statistics for a client
func (rl *clientRateLimiter) getStats(clientID string, limit int, window time.Duration) (int, time.Time) {
	now := time.Now()
	cutoff := now.Add(-window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	timestamps := rl.clients[clientID]
	filtered := timestamps[:0]
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	remaining := limit - len(filtered)
	if remaining < 0 {
		remaining = 0
	}

	reset := now.Add(window)
	return remaining, reset
}

// cleanExpired removes expired timestamps from the limiter
func (rl *clientRateLimiter) cleanExpired(now time.Time, window time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := now.Add(-window)
	for clientID, timestamps := range rl.clients {
		filtered := timestamps[:0]
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				filtered = append(filtered, ts)
			}
		}
		if len(filtered) == 0 {
			delete(rl.clients, clientID)
			continue
		}
		rl.clients[clientID] = filtered
	}
}

// getRateLimitClientID extracts the client identifier for rate limiting
func getRateLimitClientID(r *http.Request, trustProxy bool) string {
	if trustProxy {
		// Check X-Forwarded-For
		if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
			for _, candidate := range strings.Split(forwarded, ",") {
				candidate = strings.TrimSpace(candidate)
				if candidate != "" {
					return candidate
				}
			}
		}
		// Check X-Real-IP
		if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
			return realIP
		}
	}
	
	// Fall back to RemoteAddr
	if r.RemoteAddr == "" {
		return "unknown"
	}
	
	if ip := net.ParseIP(r.RemoteAddr); ip != nil {
		return ip.String()
	}
	
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	
	return strings.TrimSpace(r.RemoteAddr)
}

// isInList checks if an item is in a list
func isInList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// SimpleRateLimit returns a simple rate limit middleware with default settings
func SimpleRateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	config := DefaultRateLimitConfig()
	config.RequestsPerMinute = requestsPerMinute
	return RateLimit(config)
}
