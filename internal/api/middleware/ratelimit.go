package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type clientRateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string][]time.Time
}

func newRateLimiter(limit int, window time.Duration) *clientRateLimiter {
	if limit <= 0 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &clientRateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string][]time.Time),
	}
}

var rateLimiter = newRateLimiter(60, time.Minute)

func init() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rateLimiter.cleanExpired(time.Now())
		}
	}()
}

func (rl *clientRateLimiter) cleanExpired(now time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := now.Add(-rl.window)
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

func (rl *clientRateLimiter) allow(clientID string) bool {
	if clientID == "" {
		clientID = "unknown"
	}

	now := time.Now()
	cutoff := now.Add(-rl.window)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	timestamps := rl.clients[clientID]
	filtered := timestamps[:0]
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}

	if len(filtered) >= rl.limit {
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

func clientIdentifier(r *http.Request) string {
	if r == nil {
		return "unknown"
	}

	if r.Header != nil {
		if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
			for _, candidate := range strings.Split(forwarded, ",") {
				candidate = strings.TrimSpace(candidate)
				if candidate != "" {
					return candidate
				}
			}
		}

		if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
			return realIP
		}
	}

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

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := clientIdentifier(r)
		if !rateLimiter.allow(clientID) {
			w.Header().Set("Retry-After", strconv.Itoa(int(rateLimiter.window.Seconds())))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
