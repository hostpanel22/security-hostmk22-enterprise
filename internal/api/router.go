package api

import (
	"net/http"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/handlers"
	"github.com/azezaz/HostMK-Security-Framework/internal/api/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// System
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/api/v1/version", handlers.Version)

	// Audit
	mux.HandleFunc("/api/v1/audit/events", handlers.AuditEvents)

	// Fraud
	mux.HandleFunc("/api/v1/fraud/analyze", handlers.FraudAnalyze)
	mux.HandleFunc("/api/v1/fraud/violations", handlers.FraudViolations)

	// RBAC
	mux.HandleFunc("/api/v1/rbac/roles", handlers.RBACRoles)

	// Permissions
	mux.HandleFunc("/api/v1/permissions", handlers.Permissions)
	mux.HandleFunc("/api/v1/permissions/check", handlers.PermissionCheck)

	// Secrets
	mux.HandleFunc("/api/v1/secrets", handlers.Secrets)

	// Tenancy
	mux.HandleFunc("/api/v1/tenants", handlers.Tenants)

	var handler http.Handler = mux

	handler = middleware.Recovery(handler)
	handler = middleware.Logging(handler)
	handler = middleware.RateLimit(handler)

	return handler
}
