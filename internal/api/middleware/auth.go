package middleware

import (
	"net/http"
	"strings"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
	"github.com/azezaz/HostMK-Security-Framework/internal/auth"
	"github.com/azezaz/HostMK-Security-Framework/internal/identity"
	"github.com/azezaz/HostMK-Security-Framework/internal/persistence"
)

func Authenticate(tokens auth.Tokens, store *persistence.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			if apiKey := strings.TrimSpace(r.Header.Get("X-API-Key")); apiKey != "" {
				key, err := store.GetAPIKeyByPlain(r.Context(), apiKey)
				if err != nil || !key.IsActive() {
					response.Unauthorized(w, "invalid api key")
					return
				}
				_ = store.TouchAPIKey(r.Context(), key.ID)
				ctx := identity.WithIdentity(r.Context(), identity.Identity{
					TenantID: key.TenantID,
					APIKeyID: key.ID,
				})
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
				response.Unauthorized(w, "authentication required")
				return
			}
			raw := strings.TrimSpace(header[7:])
			claims, err := tokens.ParseAccess(raw)
			if err != nil {
				response.Unauthorized(w, "invalid access token")
				return
			}

			user, err := store.GetUser(r.Context(), claims.UserID)
			if err != nil || !user.CanAuthenticate() || claims.TenantID == "" || user.TenantID != claims.TenantID {
				response.Unauthorized(w, "invalid access token")
				return
			}

			ctx := identity.WithIdentity(r.Context(), identity.Identity{
				UserID:   claims.UserID,
				TenantID: claims.TenantID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireResourceAction(store *persistence.Store, resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, ok := identity.FromContext(r.Context())
			if !ok || id.TenantID == "" {
				response.Unauthorized(w, "authentication required")
				return
			}
			if id.APIKeyID != "" {
				key, err := store.GetAPIKeyByID(r.Context(), id.APIKeyID)
				if err != nil || !key.IsActive() || key.TenantID != id.TenantID {
					response.Unauthorized(w, "invalid api key")
					return
				}
				// Scopes are exact resource.action grants; no implicit resource-wide or wildcard grants.
				scope := resource + "." + action
				if action == "" || !key.HasScope(scope) {
					response.Forbidden(w, "missing api key scope")
					return
				}
				next.ServeHTTP(w, r)
				return
			}
			tenant, err := store.GetTenant(r.Context(), id.TenantID)
			if err == nil && tenant.IsOwner(id.UserID) {
				next.ServeHTTP(w, r)
				return
			}
			allowed, err := store.UserHasResourceActionForTenant(r.Context(), id.UserID, id.TenantID, resource, action)
			if err != nil || !allowed {
				allowedByName, nameErr := store.UserHasPermissionForTenant(r.Context(), id.UserID, id.TenantID, resource+"."+action)
				if nameErr != nil || !allowedByName {
					response.Forbidden(w, "missing permission")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireCRUD(store *persistence.Store, resource string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			action := "read"
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				action = "write"
			case http.MethodDelete:
				action = "delete"
			}
			RequireResourceAction(store, resource, action)(next).ServeHTTP(w, r)
		})
	}
}

func isPublicPath(path string) bool {
	switch {
	case path == "/health", path == "/ready", path == "/api/v1/version":
		return true
	case strings.HasPrefix(path, "/dashboard"), path == "/":
		return true
	case strings.HasPrefix(path, "/api/v1/auth/"):
		return true
	default:
		return false
	}
}
