package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

type RateLimitProvider interface {
	RateLimit(ctx context.Context, tenantID string) RateLimitConfig
}

type IntSettingReader interface {
	IntSetting(ctx context.Context, tenantID, key string, fallback int) int
}

type SettingsRateLimitProvider struct {
	Reader IntSettingReader
}

func (p SettingsRateLimitProvider) RateLimit(ctx context.Context, tenantID string) RateLimitConfig {
	limit := 60
	if p.Reader != nil {
		limit = p.Reader.IntSetting(ctx, tenantID, "engine.rate_limit.requests_per_minute", 60)
	}
	if limit <= 0 {
		limit = 60
	}
	return RateLimitConfig{Limit: limit, Window: time.Minute}
}

type StaticRateLimitProvider struct {
	Limit  int
	Window time.Duration
}

func (p StaticRateLimitProvider) RateLimit(context.Context, string) RateLimitConfig {
	limit := p.Limit
	window := p.Window
	if limit <= 0 {
		limit = 60
	}
	if window <= 0 {
		window = time.Minute
	}
	return RateLimitConfig{Limit: limit, Window: window}
}

type clientRateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
}

func newClientLimiter() *clientRateLimiter {
	return &clientRateLimiter{clients: make(map[string][]time.Time)}
}

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

func RateLimit(provider RateLimitProvider) func(http.Handler) http.Handler {
	if provider == nil {
		provider = StaticRateLimitProvider{Limit: 60, Window: time.Minute}
	}
	limiter := newClientLimiter()
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanExpired(time.Now(), time.Minute)
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
			cfg := provider.RateLimit(r.Context(), tenantID)
			clientID := clientIdentifier(r)
			if tenantID != "" {
				clientID = tenantID + ":" + clientID
			}
			if !limiter.allow(clientID, cfg.Limit, cfg.Window) {
				w.Header().Set("Retry-After", strconv.Itoa(int(cfg.Window.Seconds())))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
