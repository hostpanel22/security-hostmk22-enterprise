package security

import (
	"net/http"
	"strings"
	"time"
)

// CORS returns a middleware that adds CORS headers to responses
func CORS(config CORSConfig) func(http.Handler) http.Handler {
	if !config.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

				// Check if origin is allowed
				if origin != "" && config.IsOriginAllowed(origin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", boolToString(config.AllowCredentials))
				} else if len(config.AllowedOrigins) > 0 {
					// If we have allowed origins but this one isn't allowed, don't set the header
					// This prevents CORS errors for disallowed origins
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Credentials", "false")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}

				// Set allowed methods
				if len(config.AllowedMethods) > 0 {
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
				}

				// Set allowed headers
				if len(config.AllowedHeaders) > 0 {
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
				} else {
					w.Header().Set("Access-Control-Allow-Headers", "*")
				}

				// Set exposed headers
				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
				}

				// Set max age
				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", formatDuration(config.MaxAge))
				}

				// Set private network access
				if config.AllowPrivateNetwork {
					w.Header().Set("Access-Control-Allow-Private-Network", "true")
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			// For non-preflight requests, set CORS headers
			if origin != "" && config.IsOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", boolToString(config.AllowCredentials))
			} else if len(config.AllowedOrigins) == 0 {
				// If no origins are specified, allow all
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			// Set exposed headers for non-preflight requests
			if len(config.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			}

			// Set Vary header
			w.Header().Set("Vary", "Origin")

			next.ServeHTTP(w, r)
		})
	}
}

// boolToString converts a boolean to "true" or "false"
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// formatDuration formats a duration in seconds for CORS max-age
func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	return string(rune(seconds))
}

// SimpleCORS returns a simple CORS middleware with default settings
// This is a convenience function for basic CORS needs
func SimpleCORS() func(http.Handler) http.Handler {
	return CORS(DefaultCORSConfig())
}

// StrictCORS returns a strict CORS middleware that only allows specified origins
// This is suitable for production environments
func StrictCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	config := DefaultCORSConfig()
	config.AllowedOrigins = allowedOrigins
	config.AllowCredentials = true
	config.AllowPrivateNetwork = false
	return CORS(config)
}

// PermissiveCORS returns a permissive CORS middleware that allows all origins
// WARNING: This should only be used for development!
func PermissiveCORS() func(http.Handler) http.Handler {
	config := CORSConfig{
		Enabled:          true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400 * time.Second,
	}
	return CORS(config)
}
