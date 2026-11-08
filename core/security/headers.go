package security

import (
	"net/http"
	"strings"
)

// SecurityHeaders returns a middleware that adds security headers to responses
func SecurityHeaders(config HeadersConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Content Security Policy
			if config.ContentSecurityPolicy != "" {
				w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
			}

			// Content Security Policy Report Only
			if config.ContentSecurityPolicyReportOnly != "" {
				w.Header().Set("Content-Security-Policy-Report-Only", config.ContentSecurityPolicyReportOnly)
			}

			// Strict Transport Security (HSTS)
			if config.HSTS != "" {
				w.Header().Set("Strict-Transport-Security", config.HSTS)
			} else if config.HSTSMaxAge > 0 {
				hsts := "max-age=" + string(rune(config.HSTSMaxAge))
				if config.HSTSIncludeSubdomains {
					hsts += "; includeSubDomains"
				}
				if config.HSTSPreload {
					hsts += "; preload"
				}
				w.Header().Set("Strict-Transport-Security", hsts)
			}

			// X-Content-Type-Options
			if config.XContentTypeOptions != "" {
				w.Header().Set("X-Content-Type-Options", config.XContentTypeOptions)
			} else {
				w.Header().Set("X-Content-Type-Options", "nosniff")
			}

			// X-Frame-Options
			if config.XFrameOptions != "" {
				w.Header().Set("X-Frame-Options", config.XFrameOptions)
			} else {
				w.Header().Set("X-Frame-Options", "DENY")
			}

			// X-XSS-Protection
			if config.XSSProtection != "" {
				w.Header().Set("X-XSS-Protection", config.XSSProtection)
			} else {
				w.Header().Set("X-XSS-Protection", "1; mode=block")
			}

			// Referrer Policy
			if config.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
			} else {
				w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			}

			// Permissions Policy
			if config.PermissionsPolicy != "" {
				w.Header().Set("Permissions-Policy", config.PermissionsPolicy)
			}

			// Expect-CT (Certificate Transparency)
			if config.ExpectCT != "" {
				w.Header().Set("Expect-CT", config.ExpectCT)
			}

			// Server header - empty to hide server information
			if config.Server != "" {
				w.Header().Set("Server", config.Server)
			} else {
				// Remove Server header if it exists
				w.Header().Del("Server")
			}

			// Additional security headers
			w.Header().Set("X-Powered-By", "")

			// Cache control for API endpoints
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
			}

			// Add security.txt link
			w.Header().Set("Link", "<https://hostmk.com/.well-known/security.txt>; rel=\"security-text\"")

			next.ServeHTTP(w, r)
		})
	}
}

// DefaultSecurityHeaders returns a middleware with default security headers
func DefaultSecurityHeaders() func(http.Handler) http.Handler {
	return SecurityHeaders(DefaultHeadersConfig())
}

// BuildCSP constructs a Content Security Policy string from directives
func BuildCSP(directives map[string][]string) string {
	var parts []string
	for directive, sources := range directives {
		if len(sources) > 0 {
			parts = append(parts, directive+" "+strings.Join(sources, " "))
		}
	}
	return strings.Join(parts, "; ")
}

// StrictCSP returns a strict Content Security Policy for maximum security
func StrictCSP() string {
	return BuildCSP(map[string][]string{
		"default-src":      {"'none'"},
		"script-src":       {"'self'"},
		"style-src":        {"'self'", "'unsafe-inline'"},
		"img-src":          {"'self'", "data:", "https:"},
		"font-src":         {"'self'"},
		"connect-src":      {"'self'"},
		"frame-src":        {"'none'"},
		"object-src":       {"'none'"},
		"base-uri":         {"'self'"},
		"form-action":      {"'self'"},
		"frame-ancestors": {"'none'"},
		"report-uri":       {"https://hostmk.com/csp-report"},
	})
}

// APICSP returns a Content Security Policy suitable for REST APIs
func APICSP() string {
	return BuildCSP(map[string][]string{
		"default-src":      {"'self'"},
		"script-src":       {"'self'"},
		"style-src":        {"'self'", "'unsafe-inline'"},
		"img-src":          {"'self'", "data:", "https:"},
		"font-src":         {"'self'"},
		"connect-src":      {"'self'"},
		"frame-src":        {"'none'"},
		"object-src":       {"'none'"},
		"base-uri":         {"'self'"},
		"form-action":      {"'self'"},
		"frame-ancestors": {"'none'"},
	})
}

// WebAppCSP returns a Content Security Policy suitable for web applications
func WebAppCSP() string {
	return BuildCSP(map[string][]string{
		"default-src":      {"'self'"},
		"script-src":       {"'self'", "'unsafe-inline'", "https://cdn.jsdelivr.net", "https://cdn.tailwindcss.com"},
		"style-src":        {"'self'", "'unsafe-inline'", "https://fonts.googleapis.com", "https://cdn.jsdelivr.net"},
		"img-src":          {"'self'", "data:", "https://*.googleapis.com", "https://*.gstatic.com"},
		"font-src":         {"'self'", "https://fonts.gstatic.com"},
		"connect-src":      {"'self'", "https://api.hostmk.com", "wss://api.hostmk.com"},
		"frame-src":        {"'none'"},
		"object-src":       {"'none'"},
		"base-uri":         {"'self'"},
		"form-action":      {"'self'"},
		"frame-ancestors": {"'none'"},
		"upgrade-insecure-requests": {},
	})
}

// AddSecurityHeaders adds security headers to an http.ResponseWriter
func AddSecurityHeaders(w http.ResponseWriter, config HeadersConfig) {
	// Content Security Policy
	if config.ContentSecurityPolicy != "" {
		w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
	}

	// HSTS
	if config.HSTS != "" {
		w.Header().Set("Strict-Transport-Security", config.HSTS)
	}

	// X-Content-Type-Options
	if config.XContentTypeOptions != "" {
		w.Header().Set("X-Content-Type-Options", config.XContentTypeOptions)
	} else {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// X-Frame-Options
	if config.XFrameOptions != "" {
		w.Header().Set("X-Frame-Options", config.XFrameOptions)
	} else {
		w.Header().Set("X-Frame-Options", "DENY")
	}

	// X-XSS-Protection
	if config.XSSProtection != "" {
		w.Header().Set("X-XSS-Protection", config.XSSProtection)
	} else {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
	}

	// Referrer Policy
	if config.ReferrerPolicy != "" {
		w.Header().Set("Referrer-Policy", config.ReferrerPolicy)
	} else {
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	}

	// Permissions Policy
	if config.PermissionsPolicy != "" {
		w.Header().Set("Permissions-Policy", config.PermissionsPolicy)
	}

	// Hide server information
	w.Header().Del("Server")
	w.Header().Set("X-Powered-By", "")
}
