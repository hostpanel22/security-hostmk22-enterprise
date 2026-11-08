package security

import (
	"net/http"
	"strings"
	"time"
)

// SecurityConfig holds all security-related configuration for the application
type SecurityConfig struct {
	// CORS configuration
	CORS CORSConfig

	// CSRF configuration
	CSRF CSRFConfig

	// Rate limiting configuration
	RateLimit RateLimitConfig

	// IP filtering configuration
	IPFilter IPFilterConfig

	// Request size limits
	RequestSize RequestSizeConfig

	// Security headers configuration
	Headers HeadersConfig

	// Brute force protection configuration
	BruteForce BruteForceConfig

	// Session configuration
	Session SessionConfig
}

// CORSConfig defines Cross-Origin Resource Sharing settings
type CORSConfig struct {
	// Enabled determines if CORS is active
	Enabled bool
	// AllowedOrigins is a list of allowed origins (supports wildcards)
	AllowedOrigins []string
	// AllowedMethods is a list of allowed HTTP methods
	AllowedMethods []string
	// AllowedHeaders is a list of allowed request headers
	AllowedHeaders []string
	// ExposedHeaders is a list of headers exposed to the client
	ExposedHeaders []string
	// AllowCredentials determines if credentials are allowed
	AllowCredentials bool
	// MaxAge is the maximum age of the CORS preflight cache
	MaxAge time.Duration
	// AllowPrivateNetwork determines if private network access is allowed
	AllowPrivateNetwork bool
}

// DefaultCORSConfig returns a secure CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		Enabled:           true,
		AllowedOrigins:    []string{"https://*.hostmk.com", "https://localhost:*"},
		AllowedMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:    []string{"Authorization", "Content-Type", "X-Requested-With", "X-Tenant-ID", "X-API-Key", "X-CSRF-Token"},
		ExposedHeaders:    []string{"X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials:  true,
		MaxAge:            86400 * time.Second, // 24 hours
		AllowPrivateNetwork: false,
	}
}

// CSRFConfig defines CSRF protection settings
type CSRFConfig struct {
	// Enabled determines if CSRF protection is active
	Enabled bool
	// TokenLength is the length of the CSRF token
	TokenLength int
	// TokenExpiry is the duration after which the token expires
	TokenExpiry time.Duration
	// CookieName is the name of the CSRF cookie
	CookieName string
	// CookieSecure determines if the cookie is secure
	CookieSecure bool
	// CookieHttpOnly determines if the cookie is HTTP-only
	CookieHttpOnly bool
	// CookieSameSite is the SameSite attribute for the cookie
	CookieSameSite http.SameSite
	// CookieDomain is the domain for the CSRF cookie
	CookieDomain string
	// HeaderName is the name of the header containing the CSRF token
	HeaderName string
	// SafeMethods are HTTP methods that don't require CSRF protection
	SafeMethods []string
	// TrustProxy determines if X-Forwarded-For should be trusted
	TrustProxy bool
}

// DefaultCSRFConfig returns a secure CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		Enabled:         true,
		TokenLength:     32,
		TokenExpiry:     24 * time.Hour,
		CookieName:      "csrf_token",
		CookieSecure:    true,
		CookieHttpOnly:  true,
		CookieSameSite:  http.SameSiteLaxMode,
		CookieDomain:    "",
		TrustProxy:      true,
		HeaderName:      "X-CSRF-Token",
		SafeMethods:     []string{"GET", "HEAD", "OPTIONS", "TRACE"},
	}
}

// RateLimitConfig defines rate limiting settings
type RateLimitConfig struct {
	// Enabled determines if rate limiting is active
	Enabled bool
	// RequestsPerMinute is the default rate limit
	RequestsPerMinute int
	// BurstSize is the maximum number of requests allowed in a burst
	BurstSize int
	// WindowDuration is the time window for rate limiting
	WindowDuration time.Duration
	// TrustProxy determines if X-Forwarded-For should be trusted
	TrustProxy bool
	// Whitelist is a list of IPs that are exempt from rate limiting
	Whitelist []string
	// Blacklist is a list of IPs that are always rate limited
	Blacklist []string
}

// DefaultRateLimitConfig returns a secure rate limit configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Enabled:          true,
		RequestsPerMinute: 60,
		BurstSize:        10,
		WindowDuration:   time.Minute,
		TrustProxy:       true,
		Whitelist:        []string{"127.0.0.1", "::1"},
		Blacklist:        []string{},
	}
}

// IPFilterConfig defines IP filtering settings
type IPFilterConfig struct {
	// Enabled determines if IP filtering is active
	Enabled bool
	// Whitelist is a list of allowed IPs (if non-empty, only these IPs are allowed)
	Whitelist []string
	// Blacklist is a list of blocked IPs
	Blacklist []string
	// TrustProxy determines if X-Forwarded-For should be trusted
	TrustProxy bool
	// DefaultAction is the default action when an IP is not in the lists ("allow" or "deny")
	DefaultAction string
}

// DefaultIPFilterConfig returns a secure IP filter configuration
func DefaultIPFilterConfig() IPFilterConfig {
	return IPFilterConfig{
		Enabled:       false,
		Whitelist:     []string{},
		Blacklist:     []string{},
		TrustProxy:    true,
		DefaultAction: "allow",
	}
}

// RequestSizeConfig defines request size limits
type RequestSizeConfig struct {
	// MaxBodySize is the maximum size of the request body in bytes
	MaxBodySize int64
	// MaxHeaderSize is the maximum size of headers in bytes
	MaxHeaderSize int64
	// MaxMultipartMemory is the maximum memory for multipart forms
	MaxMultipartMemory int64
	// MaxFileSize is the maximum size of a single uploaded file
	MaxFileSize int64
	// AllowedMIMETypes is a list of allowed MIME types for uploads
	AllowedMIMETypes []string
	// BlockedMIMETypes is a list of blocked MIME types
	BlockedMIMETypes []string
}

// DefaultRequestSizeConfig returns a secure request size configuration
func DefaultRequestSizeConfig() RequestSizeConfig {
	return RequestSizeConfig{
		MaxBodySize:         10 * 1024 * 1024,    // 10 MB
		MaxHeaderSize:       8 * 1024,           // 8 KB
		MaxMultipartMemory:  32 * 1024 * 1024,   // 32 MB
		MaxFileSize:         50 * 1024 * 1024,   // 50 MB
		AllowedMIMETypes:   []string{"application/json", "multipart/form-data", "application/x-www-form-urlencoded"},
		BlockedMIMETypes:   []string{"text/html", "application/xhtml+xml"},
	}
}

// HeadersConfig defines security headers settings
type HeadersConfig struct {
	// Enabled determines if security headers are added
	Enabled bool
	// ContentSecurityPolicy defines the CSP header
	ContentSecurityPolicy string
	// ContentSecurityPolicyReportOnly defines the CSP-Report-Only header
	ContentSecurityPolicyReportOnly string
	// HSTS defines the Strict-Transport-Security header
	HSTS string
	// HSTSMaxAge is the max-age for HSTS
	HSTSMaxAge int
	// HSTSIncludeSubdomains determines if subdomains are included in HSTS
	HSTSIncludeSubdomains bool
	// HSTSPreload determines if the preload directive is included
	HSTSPreload bool
	// XContentTypeOptions defines the X-Content-Type-Options header
	XContentTypeOptions string
	// XFrameOptions defines the X-Frame-Options header
	XFrameOptions string
	// XSSProtection defines the X-XSS-Protection header
	XSSProtection string
	// ReferrerPolicy defines the Referrer-Policy header
	ReferrerPolicy string
	// PermissionsPolicy defines the Permissions-Policy header
	PermissionsPolicy string
	// ExpectCT defines the Expect-CT header
	ExpectCT string
	// Server defines the Server header (empty to hide)
	Server string
}

// DefaultHeadersConfig returns a secure headers configuration
func DefaultHeadersConfig() HeadersConfig {
	return HeadersConfig{
		Enabled: true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; img-src 'self' data: https://*.googleapis.com; font-src 'self' https://fonts.gstatic.com; connect-src 'self' https://api.hostmk.com; frame-src 'none'; object-src 'none'; base-uri 'self'; form-action 'self';",
		HSTS:                     "max-age=63072000; includeSubDomains; preload",
		HSTSMaxAge:               63072000,
		HSTSIncludeSubdomains:    true,
		HSTSPreload:              true,
		XContentTypeOptions:      "nosniff",
		XFrameOptions:            "DENY",
		XSSProtection:            "1; mode=block",
		ReferrerPolicy:           "strict-origin-when-cross-origin",
		PermissionsPolicy:        "geolocation=(), microphone=(), camera=(), payment=(), usb=()",
		ExpectCT:                 "max-age=86400, enforce",
		Server:                   "",
	}
}

// BruteForceConfig defines brute force protection settings
type BruteForceConfig struct {
	// Enabled determines if brute force protection is active
	Enabled bool
	// MaxAttempts is the maximum number of failed attempts before lockout
	MaxAttempts int
	// LockoutDuration is the duration of the lockout
	LockoutDuration time.Duration
	// WindowDuration is the time window for counting attempts
	WindowDuration time.Duration
	// LockoutKeyPrefix is the prefix for lockout keys
	LockoutKeyPrefix string
	// ProtectedPaths is a list of paths to protect
	ProtectedPaths []string
	// TrustProxy determines if X-Forwarded-For should be trusted
	TrustProxy bool
}

// DefaultBruteForceConfig returns a secure brute force configuration
func DefaultBruteForceConfig() BruteForceConfig {
	return BruteForceConfig{
		Enabled:          true,
		MaxAttempts:      5,
		LockoutDuration:  15 * time.Minute,
		WindowDuration:   5 * time.Minute,
		LockoutKeyPrefix: "bf_lockout_",
		ProtectedPaths:   []string{"/api/v1/auth/login", "/api/v1/auth/bootstrap"},
		TrustProxy:       true,
	}
}

// SessionConfig defines session security settings
type SessionConfig struct {
	// SessionName is the name of the session cookie
	SessionName string
	// SessionExpiry is the duration after which the session expires
	SessionExpiry time.Duration
	// SessionSecure determines if the session cookie is secure
	SessionSecure bool
	// SessionHttpOnly determines if the session cookie is HTTP-only
	SessionHttpOnly bool
	// SessionSameSite is the SameSite attribute for the session cookie
	SessionSameSite string
	// SessionDomain is the domain for the session cookie
	SessionDomain string
	// SessionPath is the path for the session cookie
	SessionPath string
}

// DefaultSessionConfig returns a secure session configuration
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		SessionName:     "hostmk_session",
		SessionExpiry:   24 * time.Hour,
		SessionSecure:   true,
		SessionHttpOnly: true,
		SessionSameSite: "Lax",
		SessionDomain:   "",
		SessionPath:     "/",
	}
}

// DefaultSecurityConfig returns a complete default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		CORS:        DefaultCORSConfig(),
		CSRF:        DefaultCSRFConfig(),
		RateLimit:   DefaultRateLimitConfig(),
		IPFilter:    DefaultIPFilterConfig(),
		RequestSize: DefaultRequestSizeConfig(),
		Headers:     DefaultHeadersConfig(),
		BruteForce:  DefaultBruteForceConfig(),
		Session:     DefaultSessionConfig(),
	}
}

// IsMethodSafe checks if an HTTP method is safe (doesn't modify data)
func IsMethodSafe(method string) bool {
	safeMethods := map[string]bool{
		"GET":     true,
		"HEAD":    true,
		"OPTIONS": true,
		"TRACE":   true,
	}
	return safeMethods[strings.ToUpper(method)]
}

// IsOriginAllowed checks if an origin is allowed based on CORS configuration
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if !c.Enabled {
		return false
	}
	for _, allowed := range c.AllowedOrigins {
		if matchOrigin(allowed, origin) {
			return true
		}
	}
	return false
}

// matchOrigin matches an origin against a pattern (supports wildcards)
func matchOrigin(pattern, origin string) bool {
	pattern = strings.ToLower(pattern)
	origin = strings.ToLower(origin)

	// Exact match
	if pattern == origin {
		return true
	}

	// Wildcard match for subdomains
	if strings.HasPrefix(pattern, "https://*.") {
		domain := strings.TrimPrefix(pattern, "https://*.")
		if strings.HasSuffix(origin, domain) && strings.HasPrefix(origin, "https://") {
			// Check if it's a direct subdomain or the domain itself
			parts := strings.Split(strings.TrimPrefix(origin, "https://"), ".")
			patternParts := strings.Split(domain, ".")
			if len(parts) >= len(patternParts) {
				for i := 0; i < len(patternParts); i++ {
					if parts[len(parts)-len(patternParts)+i] != patternParts[i] {
						return false
					}
				}
				return true
			}
		}
	}

	// Wildcard for localhost with port
	if strings.Contains(pattern, "localhost:*") {
		if strings.HasPrefix(origin, "https://localhost:") || strings.HasPrefix(origin, "http://localhost:") {
			return true
		}
	}

	return false
}
