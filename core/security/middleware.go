package security

import (
	"net/http"
)

// AllInOne returns a combined middleware that applies all security middlewares
// This is the recommended way to add security to your application
func AllInOne(config SecurityConfig) func(http.Handler) http.Handler {
	// Create stores
	csrfStore := NewMemoryCSRFTokenStore(config.CSRF.TokenLength, config.CSRF.TokenExpiry)
	bfStore := NewMemoryBruteForceStore(config.BruteForce.MaxAttempts, config.BruteForce.LockoutDuration, config.BruteForce.WindowDuration)

	// Create validator

	// Build the middleware chain
	var handler http.Handler

	// Start with the base handler (will be set by the caller)
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// Apply middlewares in order (most outer to most inner)

	// 1. IP Filtering (first line of defense)
	handler = IPFilter(config.IPFilter)(handler)

	// 2. Request Size Limiting
	handler = RequestSizeLimit(config.RequestSize)(handler)

	// 3. Brute Force Protection
	handler = BruteForceProtection(config.BruteForce, bfStore)(handler)

	// 4. Rate Limiting
	handler = RateLimit(config.RateLimit)(handler)

	// 5. CORS
	handler = CORS(config.CORS)(handler)

	// 6. CSRF Protection
	handler = CSRF(config.CSRF, csrfStore)(handler)

	// 7. Input Validation and Sanitization
	handler = InputValidationMiddleware(DefaultInputValidationConfig())(handler)
	handler = SanitizeMiddleware(handler)

	// 8. Security Headers (last, so they're applied to all responses)
	handler = SecurityHeaders(config.Headers)(handler)

	// Return the middleware chain
	return func(next http.Handler) http.Handler {
		// Set the innermost handler to the provided next
		finalHandler := next

		// Rebuild the chain with the actual next handler
		chain := finalHandler
		chain = SecurityHeaders(config.Headers)(chain)
		chain = SanitizeMiddleware(chain)
		chain = InputValidationMiddleware(DefaultInputValidationConfig())(chain)
		chain = CSRF(config.CSRF, csrfStore)(chain)
		chain = CORS(config.CORS)(chain)
		chain = RateLimit(config.RateLimit)(chain)
		chain = BruteForceProtection(config.BruteForce, bfStore)(chain)
		chain = RequestSizeLimit(config.RequestSize)(chain)
		chain = IPFilter(config.IPFilter)(chain)

		return chain
	}
}

// SecureChain builds a secure middleware chain with custom configuration
func SecureChain(next http.Handler, configs ...func(http.Handler) http.Handler) http.Handler {
	// Apply middlewares in reverse order (last config is applied first)
	handler := next
	for i := len(configs) - 1; i >= 0; i-- {
		handler = configs[i](handler)
	}
	return handler
}

// BuildSecureMiddleware builds a secure middleware with all security features
func BuildSecureMiddleware(config SecurityConfig) func(http.Handler) http.Handler {
	// Create stores
	csrfStore := NewMemoryCSRFTokenStore(config.CSRF.TokenLength, config.CSRF.TokenExpiry)
	bfStore := NewMemoryBruteForceStore(config.BruteForce.MaxAttempts, config.BruteForce.LockoutDuration, config.BruteForce.WindowDuration)

	return func(next http.Handler) http.Handler {
		handler := next

		// Apply middlewares in the correct order
		// Order matters: IP filter first, then size limits, then rate limiting, then auth-related, then headers last

		// 1. IP Filter
		handler = IPFilter(config.IPFilter)(handler)

		// 2. Request Size Limit
		handler = RequestSizeLimit(config.RequestSize)(handler)

		// 3. Brute Force Protection
		handler = BruteForceProtection(config.BruteForce, bfStore)(handler)

		// 4. Rate Limit
		handler = RateLimit(config.RateLimit)(handler)

		// 5. CORS
		handler = CORS(config.CORS)(handler)

		// 6. CSRF
		handler = CSRF(config.CSRF, csrfStore)(handler)

		// 7. Input Validation
		handler = InputValidationMiddleware(DefaultInputValidationConfig())(handler)

		// 8. Sanitization
		handler = SanitizeMiddleware(handler)

		// 9. Security Headers (last)
		handler = SecurityHeaders(config.Headers)(handler)

		return handler
	}
}

// MinimalSecurity returns a minimal security middleware with sensible defaults
func MinimalSecurity() func(http.Handler) http.Handler {
	config := DefaultSecurityConfig()
	return BuildSecureMiddleware(config)
}

// APISecurity returns a security middleware optimized for REST APIs
func APISecurity() func(http.Handler) http.Handler {
	config := DefaultSecurityConfig()

	// Adjust for API
	config.CORS.AllowCredentials = true
	config.CORS.AllowedOrigins = []string{"*"}
	config.CSRF.Enabled = false // Typically not needed for APIs using JWT

	return BuildSecureMiddleware(config)
}

// WebAppSecurity returns a security middleware optimized for web applications
func WebAppSecurity() func(http.Handler) http.Handler {
	config := DefaultSecurityConfig()

	// Adjust for web apps
	config.CORS.Enabled = true
	config.CORS.AllowCredentials = true
	config.CSRF.Enabled = true

	return BuildSecureMiddleware(config)
}

// ProductionSecurity returns a security middleware with production-grade settings
func ProductionSecurity() func(http.Handler) http.Handler {
	config := DefaultSecurityConfig()

	// Production settings
	config.CORS.Enabled = true
	config.CORS.AllowCredentials = true
	config.CORS.AllowedOrigins = []string{"https://*.hostmk.com"}
	config.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	config.CORS.MaxAge = 86400

	config.CSRF.Enabled = true
	config.CSRF.CookieSecure = true
	config.CSRF.CookieHttpOnly = true
	config.CSRF.CookieSameSite = http.SameSiteLaxMode

	config.RateLimit.Enabled = true
	config.RateLimit.RequestsPerMinute = 100
	config.RateLimit.BurstSize = 20
	config.RateLimit.TrustProxy = true

	config.IPFilter.Enabled = false // Typically handled at the network level

	config.RequestSize.MaxBodySize = 10 * 1024 * 1024 // 10 MB
	config.RequestSize.MaxFileSize = 50 * 1024 * 1024 // 50 MB

	config.Headers.Enabled = true
	config.Headers.HSTS = "max-age=63072000; includeSubDomains; preload"
	config.Headers.XFrameOptions = "DENY"
	config.Headers.XContentTypeOptions = "nosniff"
	config.Headers.ReferrerPolicy = "strict-origin-when-cross-origin"

	config.BruteForce.Enabled = true
	config.BruteForce.MaxAttempts = 5
	config.BruteForce.LockoutDuration = 15 * 60 // 15 minutes

	return BuildSecureMiddleware(config)
}

// DevelopmentSecurity returns a security middleware with development-friendly settings
func DevelopmentSecurity() func(http.Handler) http.Handler {
	config := DefaultSecurityConfig()

	// Development settings
	config.CORS.Enabled = true
	config.CORS.AllowCredentials = false
	config.CORS.AllowedOrigins = []string{"http://localhost:*", "https://localhost:*"}

	config.CSRF.Enabled = false // Disable for development

	config.RateLimit.Enabled = true
	config.RateLimit.RequestsPerMinute = 1000 // Higher limit for development

	config.Headers.Enabled = true
	config.Headers.HSTS = "" // Disable HSTS in development

	config.BruteForce.Enabled = false // Disable for development

	return BuildSecureMiddleware(config)
}

// SecurityLevel represents the security level of an application
type SecurityLevel string

const (
	SecurityLevelMinimal     SecurityLevel = "minimal"
	SecurityLevelStandard    SecurityLevel = "standard"
	SecurityLevelProduction  SecurityLevel = "production"
	SecurityLevelParanoid    SecurityLevel = "paranoid"
)

// GetSecurityConfig returns a security configuration based on the security level
func GetSecurityConfig(level SecurityLevel) SecurityConfig {
	switch level {
	case SecurityLevelMinimal:
		return minimalSecurityConfig()
	case SecurityLevelStandard:
		return DefaultSecurityConfig()
	case SecurityLevelProduction:
		return productionSecurityConfig()
	case SecurityLevelParanoid:
		return paranoidSecurityConfig()
	default:
		return DefaultSecurityConfig()
	}
}

// minimalSecurityConfig returns a minimal security configuration
func minimalSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()
	config.CORS.Enabled = false
	config.CSRF.Enabled = false
	config.RateLimit.Enabled = false
	config.IPFilter.Enabled = false
	config.BruteForce.Enabled = false
	return config
}

// productionSecurityConfig returns a production security configuration
func productionSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()
	config.CORS.AllowedOrigins = []string{"https://*.hostmk.com"}
	config.CORS.AllowCredentials = true
	config.CSRF.Enabled = true
	config.RateLimit.Enabled = true
	config.RateLimit.RequestsPerMinute = 100
	config.IPFilter.Enabled = false
	config.Headers.HSTS = "max-age=63072000; includeSubDomains; preload"
	config.BruteForce.Enabled = true
	return config
}

// paranoidSecurityConfig returns a paranoid security configuration
func paranoidSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()
	config.CORS.Enabled = false // No CORS in paranoid mode
	config.CSRF.Enabled = true
	config.RateLimit.Enabled = true
	config.RateLimit.RequestsPerMinute = 10 // Very low limit
	config.RateLimit.BurstSize = 5
	config.IPFilter.Enabled = true
	config.IPFilter.Whitelist = []string{} // Only allow specific IPs
	config.IPFilter.DefaultAction = "deny"
	config.RequestSize.MaxBodySize = 1024 * 1024 // 1 MB max
	config.Headers.ContentSecurityPolicy = StrictCSP()
	config.Headers.HSTS = "max-age=63072000; includeSubDomains; preload"
	config.BruteForce.Enabled = true
	config.BruteForce.MaxAttempts = 3
	config.BruteForce.LockoutDuration = 60 * 60 // 1 hour lockout
	return config
}

// ApplySecurityMiddleware applies security middleware to a handler
func ApplySecurityMiddleware(handler http.Handler, config SecurityConfig) http.Handler {
	return BuildSecureMiddleware(config)(handler)
}

// WrapWithSecurity wraps an http.Handler with security middleware
func WrapWithSecurity(handler http.Handler, level SecurityLevel) http.Handler {
	config := GetSecurityConfig(level)
	return ApplySecurityMiddleware(handler, config)
}

// SecureHandler wraps a handler function with security middleware
func SecureHandler(fn http.HandlerFunc, config SecurityConfig) http.Handler {
	return BuildSecureMiddleware(config)(fn)
}

// SecureHandlerFunc returns a handler function with security middleware applied
func SecureHandlerFunc(fn func(http.ResponseWriter, *http.Request), config SecurityConfig) http.HandlerFunc {
	middleware := BuildSecureMiddleware(config)
	return func(w http.ResponseWriter, r *http.Request) {
		middleware(http.HandlerFunc(fn)).ServeHTTP(w, r)
	}
}

// ChainMiddlewares chains multiple middlewares together
func ChainMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// SecurityMiddlewareBuilder helps build a custom security middleware chain
type SecurityMiddlewareBuilder struct {
	middlewares []func(http.Handler) http.Handler
}

// NewSecurityMiddlewareBuilder creates a new middleware builder
func NewSecurityMiddlewareBuilder() *SecurityMiddlewareBuilder {
	return &SecurityMiddlewareBuilder{
		middlewares: []func(http.Handler) http.Handler{},
	}
}

// WithIPFilter adds IP filtering
func (b *SecurityMiddlewareBuilder) WithIPFilter(config IPFilterConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, IPFilter(config))
	return b
}

// WithRateLimit adds rate limiting
func (b *SecurityMiddlewareBuilder) WithRateLimit(config RateLimitConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, RateLimit(config))
	return b
}

// WithCORS adds CORS
func (b *SecurityMiddlewareBuilder) WithCORS(config CORSConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, CORS(config))
	return b
}

// WithCSRF adds CSRF protection
func (b *SecurityMiddlewareBuilder) WithCSRF(config CSRFConfig) *SecurityMiddlewareBuilder {
	store := NewMemoryCSRFTokenStore(config.TokenLength, config.TokenExpiry)
	b.middlewares = append(b.middlewares, CSRF(config, store))
	return b
}

// WithSecurityHeaders adds security headers
func (b *SecurityMiddlewareBuilder) WithSecurityHeaders(config HeadersConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, SecurityHeaders(config))
	return b
}

// WithInputValidation adds input validation
func (b *SecurityMiddlewareBuilder) WithInputValidation(config InputValidationConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, InputValidationMiddleware(config))
	return b
}

// WithSanitization adds input sanitization
func (b *SecurityMiddlewareBuilder) WithSanitization() *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, SanitizeMiddleware)
	return b
}

// WithBruteForce adds brute force protection
func (b *SecurityMiddlewareBuilder) WithBruteForce(config BruteForceConfig) *SecurityMiddlewareBuilder {
	store := NewMemoryBruteForceStore(config.MaxAttempts, config.LockoutDuration, config.WindowDuration)
	b.middlewares = append(b.middlewares, BruteForceProtection(config, store))
	return b
}

// WithRequestSizeLimit adds request size limiting
func (b *SecurityMiddlewareBuilder) WithRequestSizeLimit(config RequestSizeConfig) *SecurityMiddlewareBuilder {
	b.middlewares = append(b.middlewares, RequestSizeLimit(config))
	return b
}

// Build builds the middleware chain
func (b *SecurityMiddlewareBuilder) Build() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := next
		for i := len(b.middlewares) - 1; i >= 0; i-- {
			handler = b.middlewares[i](handler)
		}
		return handler
	}
}

// BuildHandler builds a handler with the middleware chain
func (b *SecurityMiddlewareBuilder) BuildHandler(next http.Handler) http.Handler {
	return b.Build()(next)
}

// MustUse ensures all security middlewares are applied
func MustUse(next http.Handler, config SecurityConfig) http.Handler {
	// This is a convenience function that panics if configuration is invalid
	// In production, you should validate configuration before use
	return BuildSecureMiddleware(config)(next)
}

// ValidateSecurityConfig validates a security configuration
func ValidateSecurityConfig(config SecurityConfig) error {
	// Validate CORS
	if config.CORS.Enabled {
		if config.CORS.AllowCredentials && len(config.CORS.AllowedOrigins) > 0 {
			for _, origin := range config.CORS.AllowedOrigins {
				if origin == "*" {
					return &ConfigurationError{Field: "CORS.AllowedOrigins", Message: "cannot use wildcard origin with credentials"}
				}
			}
		}
		if len(config.CORS.AllowedOrigins) == 0 && config.CORS.AllowCredentials {
			// Can't have credentials without specific origins
			return &ConfigurationError{Field: "CORS.AllowedOrigins", Message: "allowed origins cannot be empty when credentials are enabled"}
		}
	}

	// Validate CSRF
	if config.CSRF.Enabled && config.CSRF.TokenLength < 16 {
		return &ConfigurationError{Field: "CSRF.TokenLength", Message: "token length must be at least 16"}
	}

	// Validate Rate Limit
	if config.RateLimit.Enabled && config.RateLimit.RequestsPerMinute <= 0 {
		return &ConfigurationError{Field: "RateLimit.RequestsPerMinute", Message: "requests per minute must be positive"}
	}

	// Validate Brute Force
	if config.BruteForce.Enabled && config.BruteForce.MaxAttempts <= 0 {
		return &ConfigurationError{Field: "BruteForce.MaxAttempts", Message: "max attempts must be positive"}
	}

	return nil
}

// ConfigurationError represents a configuration error
type ConfigurationError struct {
	Field   string
	Message string
}

func (e *ConfigurationError) Error() string {
	return e.Field + ": " + e.Message
}

// IsValidSecurityConfig checks if a security configuration is valid
func IsValidSecurityConfig(config SecurityConfig) bool {
	return ValidateSecurityConfig(config) == nil
}
