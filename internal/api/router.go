package api

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/handlers"
	"github.com/azezaz/HostMK-Security-Framework/internal/api/middleware"
	"github.com/azezaz/HostMK-Security-Framework/internal/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewRouter(api *handlers.API) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", api.Health)
	mux.HandleFunc("/ready", api.Health)
	mux.HandleFunc("/api/v1/version", api.Version)

	mux.HandleFunc("/api/v1/auth/bootstrap", api.AuthBootstrap)
	mux.HandleFunc("/api/v1/auth/login", api.AuthLogin)
	mux.HandleFunc("/api/v1/auth/refresh", api.AuthRefresh)

	protect := func(resource string, fn http.HandlerFunc) http.Handler {
		return middleware.RequireCRUD(api.Store, resource)(fn)
	}

	mux.Handle("/api/v1/users", protect("users", api.Users))
	mux.Handle("/api/v1/users/action", protect("users", api.UserAction))
	mux.Handle("/api/v1/tenants", protect("tenants", api.Tenants))
	mux.Handle("/api/v1/tenants/action", protect("tenants", api.TenantAction))
	mux.Handle("/api/v1/rbac/roles", protect("roles", api.Roles))
	mux.Handle("/api/v1/permissions", protect("permissions", api.Permissions))
	mux.Handle("/api/v1/permissions/check", protect("permissions", api.PermissionCheck))
	mux.Handle("/api/v1/api-keys", protect("api_keys", api.APIKeys))
	mux.Handle("/api/v1/secrets", protect("secrets", api.Secrets))
	mux.Handle("/api/v1/audit/events", protect("audit", api.AuditEvents))
	mux.Handle("/api/v1/fraud/rules", protect("fraud", api.FraudRules))
	mux.Handle("/api/v1/fraud/analyze", protect("fraud", api.FraudAnalyze))
	mux.Handle("/api/v1/fraud/violations", protect("fraud", api.FraudViolations))
	mux.Handle("/api/v1/settings", protect("settings", api.Settings))

	dashboard, err := fs.Sub(web.Dashboard, "dashboard")
	if err == nil {
		fileServer := http.FileServer(http.FS(dashboard))
		mux.Handle("/dashboard/", http.StripPrefix("/dashboard/", fileServer))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/dashboard/", http.StatusFound)
				return
			}
			http.NotFound(w, r)
		})
	}

	var handler http.Handler = mux
	handler = middleware.Authenticate(api.Tokens, api.Store)(handler)
	handler = middleware.RateLimit(middleware.SettingsRateLimitProvider{Reader: api.Store})(handler)
	handler = middleware.StructuredLogging(api.Log)(handler)
	handler = middleware.Recovery(handler)
	handler = otelhttp.NewHandler(handler, "hostmk.http")
	handler = securityHeaders(handler)
	return handler
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}
