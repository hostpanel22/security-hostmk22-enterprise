// security_config_example.go provides example configurations for different environments
// This file is for documentation and can be copied/imported as needed

package security

import (
	"time"
)

// DevelopmentSecurityConfig returns a security configuration suitable for development
// This configuration is less strict to allow for easier testing and development
func DevelopmentSecurityConfig() SecurityConfig {
	return SecurityConfig{
		// CORS: Allow local development origins
		CORS: CORSConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"http://localhost:*", "https://localhost:*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           86400,
		},

		// CSRF: Disabled for development (enable for API testing)
		CSRF: CSRFConfig{
			Enabled:         false,
			TokenLength:     32,
			TokenExpiry:     24 * time.Hour,
			CookieSecure:    false,
			CookieHttpOnly:  true,
			CookieSameSite:  http.SameSiteLaxMode,
			HeaderName:      "X-CSRF-Token",
			TrustProxy:      false,
		},

		// Rate Limit: Higher limits for development
		RateLimit: RateLimitConfig{
			Enabled:          true,
			RequestsPerMinute: 1000,
			BurstSize:        100,
			WindowDuration:   time.Minute,
			TrustProxy:       false,
		},

		// IP Filter: Disabled for development
		IPFilter: IPFilterConfig{
			Enabled:       false,
			TrustProxy:    false,
			DefaultAction: "allow",
		},

		// Request Size: Generous limits for development
		RequestSize: RequestSizeConfig{
			MaxBodySize:         50 * 1024 * 1024,  // 50 MB
			MaxHeaderSize:       16 * 1024,         // 16 KB
			MaxMultipartMemory:  128 * 1024 * 1024, // 128 MB
			MaxFileSize:         100 * 1024 * 1024, // 100 MB
		},

		// Security Headers: Basic protection
		Headers: HeadersConfig{
			Enabled:               true,
			ContentSecurityPolicy: APICSP(),
			HSTS:                   "", // Disable HSTS in development
			XContentTypeOptions:    "nosniff",
			XFrameOptions:           "SAMEORIGIN",
			XSSProtection:           "1; mode=block",
			ReferrerPolicy:         "no-referrer",
		},

		// Brute Force: Disabled for development
		BruteForce: BruteForceConfig{
			Enabled:         false,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
			WindowDuration:  5 * time.Minute,
		},

		// Session: Basic configuration
		Session: SessionConfig{
			SessionName:     "dev_session",
			SessionExpiry:   24 * time.Hour,
			SessionSecure:   false,
			SessionHttpOnly: true,
			SessionSameSite: http.SameSiteLaxMode,
		},
	}
}

// StagingSecurityConfig returns a security configuration suitable for staging environments
// This configuration is more strict than development but less than production
func StagingSecurityConfig() SecurityConfig {
	return SecurityConfig{
		// CORS: Specific origins only
		CORS: CORSConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"https://staging.myapp.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
			AllowCredentials: true,
			MaxAge:           86400,
		},

		// CSRF: Enabled
		CSRF: CSRFConfig{
			Enabled:         true,
			TokenLength:     32,
			TokenExpiry:     24 * time.Hour,
			CookieSecure:    true,
			CookieHttpOnly:  true,
			CookieSameSite:  http.SameSiteLaxMode,
			HeaderName:      "X-CSRF-Token",
			TrustProxy:      true,
		},

		// Rate Limit: Moderate limits
		RateLimit: RateLimitConfig{
			Enabled:          true,
			RequestsPerMinute: 200,
			BurstSize:        50,
			WindowDuration:   time.Minute,
			TrustProxy:       true,
		},

		// IP Filter: Disabled (handled at load balancer)
		IPFilter: IPFilterConfig{
			Enabled:       false,
			TrustProxy:    true,
			DefaultAction: "allow",
		},

		// Request Size: Reasonable limits
		RequestSize: RequestSizeConfig{
			MaxBodySize:         25 * 1024 * 1024,  // 25 MB
			MaxHeaderSize:       8 * 1024,          // 8 KB
			MaxMultipartMemory:  64 * 1024 * 1024,  // 64 MB
			MaxFileSize:         50 * 1024 * 1024,  // 50 MB
		},

		// Security Headers: Strong protection
		Headers: HeadersConfig{
			Enabled:               true,
			ContentSecurityPolicy: WebAppCSP(),
			HSTS:                   "max-age=31536000; includeSubDomains",
			XContentTypeOptions:    "nosniff",
			XFrameOptions:           "DENY",
			XSSProtection:           "1; mode=block",
			ReferrerPolicy:         "strict-origin-when-cross-origin",
			PermissionsPolicy:      "geolocation=(), microphone=(), camera=()",
		},

		// Brute Force: Enabled
		BruteForce: BruteForceConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
			WindowDuration:  5 * time.Minute,
			ProtectedPaths:  []string{"/api/v1/auth/login", "/api/v1/auth/bootstrap"},
			TrustProxy:      true,
		},

		// Session: Secure configuration
		Session: SessionConfig{
			SessionName:     "staging_session",
			SessionExpiry:   24 * time.Hour,
			SessionSecure:   true,
			SessionHttpOnly: true,
			SessionSameSite: http.SameSiteLaxMode,
		},
	}
}

// ProductionSecurityConfig returns a security configuration suitable for production
// This configuration provides maximum protection
func ProductionSecurityConfig() SecurityConfig {
	return SecurityConfig{
		// CORS: Specific production origins only
		CORS: CORSConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"https://myapp.com", "https://api.myapp.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With", "X-Tenant-ID"},
			AllowCredentials: true,
			MaxAge:           86400,
		},

		// CSRF: Enabled with secure cookies
		CSRF: CSRFConfig{
			Enabled:         true,
			TokenLength:     32,
			TokenExpiry:     24 * time.Hour,
			CookieSecure:    true,
			CookieHttpOnly:  true,
			CookieSameSite:  http.SameSiteLaxMode,
			HeaderName:      "X-CSRF-Token",
			TrustProxy:      true,
		},

		// Rate Limit: Strict limits
		RateLimit: RateLimitConfig{
			Enabled:          true,
			RequestsPerMinute: 100,
			BurstSize:        20,
			WindowDuration:   time.Minute,
			TrustProxy:       true,
			Whitelist:        []string{"127.0.0.1", "::1"},
		},

		// IP Filter: Disabled (handled at network level)
		IPFilter: IPFilterConfig{
			Enabled:       false,
			TrustProxy:    true,
			DefaultAction: "allow",
		},

		// Request Size: Conservative limits
		RequestSize: RequestSizeConfig{
			MaxBodySize:         10 * 1024 * 1024,  // 10 MB
			MaxHeaderSize:       8 * 1024,          // 8 KB
			MaxMultipartMemory:  32 * 1024 * 1024,  // 32 MB
			MaxFileSize:         50 * 1024 * 1024,  // 50 MB
			AllowedMIMETypes:   []string{"application/json", "multipart/form-data"},
			BlockedMIMETypes:   []string{"text/html", "application/xhtml+xml", "application/xml"},
		},

		// Security Headers: Maximum protection
		Headers: HeadersConfig{
			Enabled:               true,
			ContentSecurityPolicy: StrictCSP(),
			HSTS:                   "max-age=63072000; includeSubDomains; preload",
			XContentTypeOptions:    "nosniff",
			XFrameOptions:           "DENY",
			XSSProtection:           "1; mode=block",
			ReferrerPolicy:         "strict-origin-when-cross-origin",
			PermissionsPolicy:      "geolocation=(), microphone=(), camera=(), payment=(), usb=()",
			ExpectCT:               "max-age=86400, enforce",
		},

		// Brute Force: Strict protection
		BruteForce: BruteForceConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
			WindowDuration:  5 * time.Minute,
			ProtectedPaths:  []string{"/api/v1/auth/login", "/api/v1/auth/bootstrap", "/api/v1/auth/refresh"},
			TrustProxy:      true,
		},

		// Session: Maximum security
		Session: SessionConfig{
			SessionName:     "prod_session",
			SessionExpiry:   24 * time.Hour,
			SessionSecure:   true,
			SessionHttpOnly: true,
			SessionSameSite: http.SameSiteLaxMode,
		},
	}
}

// ParanoidSecurityConfig returns an extremely strict security configuration
// Use this for applications requiring maximum security (banking, healthcare, etc.)
func ParanoidSecurityConfig() SecurityConfig {
	return SecurityConfig{
		// CORS: Disabled for maximum security
		CORS: CORSConfig{
			Enabled: false,
		},

		// CSRF: Enabled with maximum security
		CSRF: CSRFConfig{
			Enabled:         true,
			TokenLength:     64,
			TokenExpiry:     time.Hour,
			CookieSecure:    true,
			CookieHttpOnly:  true,
			CookieSameSite:  http.SameSiteStrictMode,
			HeaderName:      "X-CSRF-Token",
			TrustProxy:      true,
		},

		// Rate Limit: Very strict
		RateLimit: RateLimitConfig{
			Enabled:          true,
			RequestsPerMinute: 10,
			BurstSize:        5,
			WindowDuration:   time.Minute,
			TrustProxy:       true,
			Whitelist:        []string{},
			Blacklist:        []string{},
		},

		// IP Filter: Whitelist only
		IPFilter: IPFilterConfig{
			Enabled:       true,
			Whitelist:     []string{}, // Add your trusted IPs here
			Blacklist:     []string{},
			TrustProxy:    true,
			DefaultAction: "deny",
		},

		// Request Size: Very conservative
		RequestSize: RequestSizeConfig{
			MaxBodySize:         1 * 1024 * 1024,   // 1 MB
			MaxHeaderSize:       4 * 1024,           // 4 KB
			MaxMultipartMemory:  8 * 1024 * 1024,   // 8 MB
			MaxFileSize:         5 * 1024 * 1024,   // 5 MB
			AllowedMIMETypes:   []string{"application/json"},
			BlockedMIMETypes:   []string{"*"},
		},

		// Security Headers: Maximum protection
		Headers: HeadersConfig{
			Enabled:               true,
			ContentSecurityPolicy: StrictCSP(),
			HSTS:                   "max-age=63072000; includeSubDomains; preload",
			XContentTypeOptions:    "nosniff",
			XFrameOptions:           "DENY",
			XSSProtection:           "1; mode=block",
			ReferrerPolicy:         "no-referrer",
			PermissionsPolicy:      "", // Block all permissions
			ExpectCT:               "max-age=86400, enforce",
			Server:                 "", // Hide server information
		},

		// Brute Force: Very strict
		BruteForce: BruteForceConfig{
			Enabled:         true,
			MaxAttempts:     3,
			LockoutDuration: time.Hour,
			WindowDuration:  time.Minute,
			ProtectedPaths:  []string{"*"}, // Protect all paths
			TrustProxy:      true,
		},

		// Session: Maximum security
		Session: SessionConfig{
			SessionName:     "paranoid_session",
			SessionExpiry:   time.Hour,
			SessionSecure:   true,
			SessionHttpOnly: true,
			SessionSameSite: http.SameSiteStrictMode,
		},
	}
}

// GetSecurityConfig returns a security configuration based on the environment
func GetSecurityConfig(level SecurityLevel) SecurityConfig {
	switch level {
	case SecurityLevelMinimal:
		return DevelopmentSecurityConfig()
	case SecurityLevelStandard:
		return DefaultSecurityConfig()
	case SecurityLevelProduction:
		return ProductionSecurityConfig()
	case SecurityLevelParanoid:
		return ParanoidSecurityConfig()
	default:
		return DefaultSecurityConfig()
	}
}
