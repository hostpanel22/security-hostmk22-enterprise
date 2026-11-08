// Package security provides comprehensive host-level security middleware and utilities
// for protecting Go applications against common web vulnerabilities and attacks.
//
// This package includes protection against:
//   - Cross-Site Request Forgery (CSRF)
//   - Cross-Origin Resource Sharing (CORS) misconfigurations
//   - Cross-Site Scripting (XSS)
//   - SQL Injection
//   - Brute Force attacks
//   - Denial of Service (DoS/DDoS)
//   - File upload vulnerabilities
//   - Information disclosure through security headers
//
// Basic usage:
//
//	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	})
//
//	// Use default production security
//	secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelProduction)
//
//	// Or use custom configuration
//	config := security.DefaultSecurityConfig()
//	config.CORS.AllowedOrigins = []string{"https://myapp.com"}
//	secureHandler := security.BuildSecureMiddleware(config)(handler)
//
//	http.ListenAndServe(":8080", secureHandler)
//
// Advanced usage with middleware builder:
//
//	builder := security.NewSecurityMiddlewareBuilder()
//	middleware := builder.
//	    WithCORS(security.DefaultCORSConfig()).
//	    WithCSRF(security.DefaultCSRFConfig()).
//	    WithSecurityHeaders(security.DefaultHeadersConfig()).
//	    WithRateLimit(security.DefaultRateLimitConfig()).
//	    Build()
//
//	handler := middleware(http.HandlerFunc(myHandler))
//
// Security levels:
//   - SecurityLevelMinimal: Basic protection for development
//   - SecurityLevelStandard: Recommended for most applications
//   - SecurityLevelProduction: Maximum protection for production
//   - SecurityLevelParanoid: Extreme security with strict defaults
package security
