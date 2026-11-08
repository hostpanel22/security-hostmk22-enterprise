package security

import (
	"context"
	"io"
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCORS(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	corsHandler := CORS(DefaultCORSConfig())(handler)

	// Test preflight request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	corsHandler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 for preflight, got %d", w.Code)
	}

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected Access-Control-Allow-Origin header")
	}

	// Test non-preflight request
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w = httptest.NewRecorder()
	corsHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	headersHandler := SecurityHeaders(DefaultHeadersConfig())(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	headersHandler.ServeHTTP(w, req)

	// Check security headers
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}

	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("Expected X-Frame-Options header")
	}

	if w.Header().Get("X-XSS-Protection") != "1; mode=block" {
		t.Error("Expected X-XSS-Protection header")
	}

	if w.Header().Get("Referrer-Policy") == "" {
		t.Error("Expected Referrer-Policy header")
	}

	if w.Header().Get("Content-Security-Policy") == "" {
		t.Error("Expected Content-Security-Policy header")
	}
}

func TestRateLimit(t *testing.T) {
	config := RateLimitConfig{
		Enabled:          true,
		RequestsPerMinute: 2,
		BurstSize:        2,
		WindowDuration:   time.Minute,
		TrustProxy:       false,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	rateLimitHandler := RateLimit(config)(handler)

	// First two requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		w := httptest.NewRecorder()
		rateLimitHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, w.Code)
		}
	}

	// Third request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	rateLimitHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}

	// Check Retry-After header
	if w.Header().Get("Retry-After") == "" {
		t.Error("Expected Retry-After header")
	}
}

func TestIPFilter(t *testing.T) {
	// Test whitelist
	config := IPFilterConfig{
		Enabled:       true,
		Whitelist:     []string{"192.168.1.1", "10.0.0.0/8"},
		DefaultAction: "deny",
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	ipFilterHandler := IPFilter(config)(handler)

	// Test allowed IP
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	ipFilterHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for whitelisted IP, got %d", w.Code)
	}

	// Test blocked IP
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "8.8.8.8:12345"
	w = httptest.NewRecorder()
	ipFilterHandler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 for non-whitelisted IP, got %d", w.Code)
	}

	// Test CIDR range
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.5:12345"
	w = httptest.NewRecorder()
	ipFilterHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for IP in CIDR range, got %d", w.Code)
	}
}

func TestCSRF(t *testing.T) {
	config := CSRFConfig{
		Enabled:         true,
		TokenLength:     16,
		TokenExpiry:     time.Hour,
		CookieName:      "csrf_token",
		CookieSecure:    false,
		CookieHttpOnly:  false,
		CookieSameSite:  http.SameSiteLaxMode,
		HeaderName:      "X-CSRF-Token",
		TrustProxy:      false,
	}

	store := NewMemoryCSRFTokenStore(config.TokenLength, config.TokenExpiry)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	csrfHandler := CSRF(config, store)(handler)

	// Test GET request (should set CSRF cookie)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()
	csrfHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for GET, got %d", w.Code)
	}

	// Get CSRF token from cookie
	cookies := w.Result().Cookies()
	var csrfToken string
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			csrfToken = cookie.Value
			break
		}
	}

	if csrfToken == "" {
		t.Error("Expected CSRF cookie to be set")
	}

	// Test POST without CSRF token (should fail)
	req = httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w = httptest.NewRecorder()
	csrfHandler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 for POST without CSRF, got %d", w.Code)
	}

	// Test POST with valid CSRF token (should succeed)
	req = httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	req.Header.Set(config.HeaderName, csrfToken)
	w = httptest.NewRecorder()
	csrfHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for POST with valid CSRF, got %d", w.Code)
	}
}

func TestInputValidation(t *testing.T) {
	validator := NewInputValidator(DefaultInputValidationConfig())

	// Test valid input
	validInput := "Hello, World!"
	safe, err := validator.ValidateString(validInput)
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
	if safe != validInput {
		t.Errorf("Expected input to be unchanged, got: %s", safe)
	}

	// Test input with HTML
	htmlInput := "<b>Bold</b>"
	safe, err = validator.ValidateString(htmlInput)
	if err != nil {
		t.Errorf("Expected no error for HTML input (sanitized), got: %v", err)
	}
	if safe == htmlInput {
		t.Error("Expected HTML to be sanitized")
	}

	// Test long input
	longInput := string(make([]byte, 10001))
	_, err = validator.ValidateString(longInput)
	if err == nil {
		t.Error("Expected error for long input")
	}
}

func TestBruteForceProtection(t *testing.T) {
	config := BruteForceConfig{
		Enabled:         true,
		MaxAttempts:     3,
		LockoutDuration: time.Minute,
		WindowDuration:  time.Minute,
		ProtectedPaths:  []string{"/login"},
		TrustProxy:      false,
	}

	store := NewMemoryBruteForceStore(config.MaxAttempts, config.LockoutDuration, config.WindowDuration)

	// Record failed attempts
	identifier := "192.168.1.1:12345"
	for i := 0; i < 3; i++ {
		store.RecordAttempt(identifier)
	}

	// Lock the identifier
	store.Lock(identifier, config.LockoutDuration)

	// Verify it is locked
	locked, _ := store.IsLocked(identifier)
	if !locked {
		t.Error("Expected IP to be locked after 3 attempts")
	}
}

func TestRequestSizeLimit(t *testing.T) {
	config := RequestSizeConfig{
		MaxBodySize: 100, // 100 bytes
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body to trigger size check
		io.ReadAll(r.Body)
		w.Write([]byte("OK"))
	})

	sizeLimitHandler := RequestSizeLimit(config)(handler)

	// Test request within limit
	req := httptest.NewRequest("POST", "/test", strings.NewReader("test"))
	w := httptest.NewRecorder()
	sizeLimitHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for small request, got %d", w.Code)
	}
}

func TestAllInOne(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	allInOneHandler := AllInOne(DefaultSecurityConfig())(handler)

	// Test that all middleware is applied
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	allInOneHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check that security headers are present
	if w.Header().Get("X-Content-Type-Options") == "" {
		t.Error("Expected X-Content-Type-Options header")
	}
}

func TestSecurityConfigValidation(t *testing.T) {
	// Test valid configuration
	config := DefaultSecurityConfig()
	if err := ValidateSecurityConfig(config); err != nil {
		t.Errorf("Expected valid config to pass validation, got: %v", err)
	}

	// Test invalid CSRF configuration
	config = DefaultSecurityConfig()
	config.CSRF.Enabled = true
	config.CSRF.TokenLength = 8
	if err := ValidateSecurityConfig(config); err == nil {
		t.Error("Expected validation error for short CSRF token length")
	}

	// Test invalid rate limit configuration
	config = DefaultSecurityConfig()
	config.RateLimit.Enabled = true
	config.RateLimit.RequestsPerMinute = 0
	if err := ValidateSecurityConfig(config); err == nil {
		t.Error("Expected validation error for zero rate limit")
	}
}

func TestIPUtilities(t *testing.T) {
	// Test IsPrivateIP
	if !IsPrivateIP("192.168.1.1") {
		t.Error("Expected 192.168.1.1 to be private")
	}

	if !IsPrivateIP("10.0.0.1") {
		t.Error("Expected 10.0.0.1 to be private")
	}

	if !IsPrivateIP("127.0.0.1") {
		t.Error("Expected 127.0.0.1 to be private")
	}

	if IsPrivateIP("8.8.8.8") {
		t.Error("Expected 8.8.8.8 to be public")
	}

	// Test IPInRange
	if !IPInRange("192.168.1.5", "192.168.1.0/24") {
		t.Error("Expected 192.168.1.5 to be in 192.168.1.0/24")
	}

	if IPInRange("10.0.0.1", "192.168.1.0/24") {
		t.Error("Expected 10.0.0.1 to not be in 192.168.1.0/24")
	}
}

func TestCSPBuilder(t *testing.T) {
	directives := map[string][]string{
		"default-src": {"'self'"},
		"script-src": {"'self'", "https://cdn.example.com"},
	}

	csp := BuildCSP(directives)

	if !strings.Contains(csp, "default-src 'self'") {
		t.Error("Expected default-src directive in CSP")
	}

	if !strings.Contains(csp, "script-src 'self' https://cdn.example.com") {
		t.Error("Expected script-src directive in CSP")
	}
}

func TestSanitizeInput(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"<b>Bold</b>", "Bold"},
	}
	for _, tc := range testCases {
		result := SanitizeInput(tc.input)
		if result != tc.expected {
			t.Errorf("SanitizeInput(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestFileUploadSecurity(t *testing.T) {
	config := DefaultRequestSizeConfig()
	config.MaxFileSize = 1024 * 1024 // 1 MB

	fileSecurity := NewFileUploadSecurity(config)

	// Test file size check
	if err := fileSecurity.CheckFileSize(500 * 1024); err != nil {
		t.Errorf("Expected no error for 500KB file, got: %v", err)
	}

	if err := fileSecurity.CheckFileSize(2 * 1024 * 1024); err == nil {
		t.Error("Expected error for 2MB file")
	}

	// Test MIME type check
	if err := fileSecurity.CheckMIMEType("application/json"); err != nil {
		t.Errorf("Expected no error for application/json, got: %v", err)
	}

	// Test file extension check
	if err := fileSecurity.CheckFileExtension("image.png"); err != nil {
		t.Errorf("Expected no error for .png, got: %v", err)
	}

	if err := fileSecurity.CheckFileExtension("script.php"); err == nil {
		t.Error("Expected error for .php")
	}

	// Test content check
	if err := fileSecurity.CheckFileContent([]byte("Hello, World!")); err != nil {
		t.Errorf("Expected no error for safe content, got: %v", err)
	}

	if err := fileSecurity.CheckFileContent([]byte("<?php echo 'hello'; ?>")); err == nil {
		t.Error("Expected error for PHP code")
	}
}

func TestSecurityLevels(t *testing.T) {
	levels := []SecurityLevel{
		SecurityLevelMinimal,
		SecurityLevelStandard,
		SecurityLevelProduction,
		SecurityLevelParanoid,
	}

	for _, level := range levels {
		config := GetSecurityConfig(level)
		if err := ValidateSecurityConfig(config); err != nil {
			t.Errorf("Security level %s: %v", level, err)
		}
	}
}

func TestMiddlewareBuilder(t *testing.T) {
	builder := NewSecurityMiddlewareBuilder()

	middleware := builder.
		WithCORS(DefaultCORSConfig()).
		WithSecurityHeaders(DefaultHeadersConfig()).
		WithRateLimit(DefaultRateLimitConfig()).
		Build()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	secureHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check that headers are applied
	if w.Header().Get("X-Content-Type-Options") == "" {
		t.Error("Expected security headers to be applied")
	}
}

func TestContextUtilities(t *testing.T) {
	ctx := context.Background()

	// Test WithIP and GetIP
	ctx = WithIP(ctx, "192.168.1.1")
	ip, ok := GetIP(ctx)
	if !ok || ip != "192.168.1.1" {
		t.Error("Expected to get IP from context")
	}

	// Test WithBruteForceID and GetBruteForceID
	ctx = WithBruteForceID(ctx, "bf_123")
	bfID, ok := GetBruteForceID(ctx)
	if !ok || bfID != "bf_123" {
		t.Error("Expected to get brute force ID from context")
	}
}
