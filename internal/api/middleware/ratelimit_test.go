package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitRemoteAddr(t *testing.T) {
	handler := RateLimit(StaticRateLimitProvider{Limit: 3, Window: time.Minute})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("request %d returned status %d, want %d", i+1, res.Code, http.StatusOK)
		}
	}
}

func TestRateLimitXForwardedFor(t *testing.T) {
	handler := RateLimit(StaticRateLimitProvider{Limit: 2, Window: time.Minute})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		req.RemoteAddr = "10.0.0.1:1111"
		req.Header.Set("X-Forwarded-For", "203.0.113.15, 198.51.100.4")
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("request %d returned status %d, want %d", i+1, res.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	req.RemoteAddr = "10.0.0.1:1111"
	req.Header.Set("X-Forwarded-For", "203.0.113.15, 198.51.100.4")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusTooManyRequests {
		t.Fatalf("third request returned status %d, want %d", res.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimitXRealIP(t *testing.T) {
	handler := RateLimit(StaticRateLimitProvider{Limit: 1, Window: time.Minute})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	firstReq := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	firstReq.RemoteAddr = "10.0.0.9:5656"
	firstReq.Header.Set("X-Real-IP", "198.51.100.25")
	res1 := httptest.NewRecorder()
	handler.ServeHTTP(res1, firstReq)
	if res1.Code != http.StatusOK {
		t.Fatalf("first request returned status %d, want %d", res1.Code, http.StatusOK)
	}

	secondReq := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	secondReq.RemoteAddr = "10.0.0.9:5656"
	secondReq.Header.Set("X-Real-IP", "198.51.100.25")
	res2 := httptest.NewRecorder()
	handler.ServeHTTP(res2, secondReq)
	if res2.Code != http.StatusTooManyRequests {
		t.Fatalf("second request returned status %d, want %d", res2.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimitRejectsRequestsOverLimit(t *testing.T) {
	handler := RateLimit(StaticRateLimitProvider{Limit: 2, Window: time.Minute})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		req.RemoteAddr = "127.0.0.1:5678"
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("request %d returned status %d, want %d", i+1, res.Code, http.StatusOK)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	req.RemoteAddr = "127.0.0.1:5678"
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusTooManyRequests {
		t.Fatalf("third request returned status %d, want %d", res.Code, http.StatusTooManyRequests)
	}
	if res.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header on rate-limited response")
	}
}

func TestRateLimitCleansExpiredEntries(t *testing.T) {
	limiter := newClientLimiter()
	limiter.clients["expired-client"] = []time.Time{time.Now().Add(-time.Second)}
	limiter.cleanExpired(time.Now(), 50*time.Millisecond)
	if _, exists := limiter.clients["expired-client"]; exists {
		t.Fatal("expired client entry should be removed from the map")
	}
}
