// +build ignore

// Example main.go demonstrating how to use the security package
// This file is for documentation purposes and can be used as a template
// To use it, copy to main.go and customize as needed

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/security"
)

func main() {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Secure World!"))
	})

	// Option 1: Use a predefined security level (recommended)
	// secureHandler := security.WrapWithSecurity(handler, security.SecurityLevelProduction)

	// Option 2: Use custom configuration (advanced)
	config := security.SecurityConfig{
		// CORS Configuration
		CORS: security.CORSConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"https://myapp.com", "https://*.myapp.com"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
			AllowCredentials: true,
			MaxAge:           24 * time.Hour,
		},

		// CSRF Configuration
		CSRF: security.CSRFConfig{
			Enabled:         true,
			TokenLength:     32,
			TokenExpiry:     24 * time.Hour,
			CookieName:      "csrf_token",
			CookieSecure:    true,
			CookieHttpOnly:  true,
			CookieSameSite:  http.SameSiteLaxMode,
			HeaderName:      "X-CSRF-Token",
			TrustProxy:      true,
		},

		// Rate Limit Configuration
		RateLimit: security.RateLimitConfig{
			Enabled:          true,
			RequestsPerMinute: 100,
			BurstSize:        20,
			WindowDuration:   time.Minute,
			TrustProxy:       true,
			Whitelist:        []string{"127.0.0.1", "::1"},
		},

		// IP Filter Configuration
		IPFilter: security.IPFilterConfig{
			Enabled:       false, // Typically handled at network level
			TrustProxy:    true,
			DefaultAction: "allow",
		},

		// Request Size Configuration
		RequestSize: security.RequestSizeConfig{
			MaxBodySize:   10 * 1024 * 1024, // 10 MB
			MaxFileSize:   50 * 1024 * 1024, // 50 MB
			AllowedMIMETypes: []string{"application/json", "multipart/form-data"},
			BlockedMIMETypes: []string{"text/html", "application/xhtml+xml"},
		},

		// Security Headers Configuration
		Headers: security.HeadersConfig{
			Enabled:               true,
			ContentSecurityPolicy: security.StrictCSP(),
			HSTS:                   "max-age=63072000; includeSubDomains; preload",
			XContentTypeOptions:    "nosniff",
			XFrameOptions:           "DENY",
			XSSProtection:           "1; mode=block",
			ReferrerPolicy:         "strict-origin-when-cross-origin",
			PermissionsPolicy:      "geolocation=(), microphone=(), camera=()",
		},

		// Brute Force Configuration
		BruteForce: security.BruteForceConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
			WindowDuration:  5 * time.Minute,
			ProtectedPaths:  []string{"/api/v1/auth/login", "/api/v1/auth/bootstrap"},
			TrustProxy:      true,
		},

		// Session Configuration
		Session: security.SessionConfig{
			SessionName:     "hostmk_session",
			SessionExpiry:   24 * time.Hour,
			SessionSecure:   true,
			SessionHttpOnly: true,
			SessionSameSite: http.SameSiteLaxMode,
		},
	}

	// Build the security middleware
	secureHandler := security.BuildSecureMiddleware(config)(handler)

	// Start the server
	log.Println("Server starting on :8080 with security middleware...")
	log.Fatal(http.ListenAndServe(":8080", secureHandler))
}
