package rbac

import (
	"context"
	"net/http"
	"strings"
)

const (
	TenantHeader = "X-Tenant-ID"
	UserHeader   = "X-User-ID"
	RoleHeader   = "X-Role"
)

type contextKey string

const (
	TenantContextKey contextKey = "rbac_tenant_id"
	UserContextKey   contextKey = "rbac_user_id"
	RoleContextKey   contextKey = "rbac_role"
)

type MiddlewareConfig struct {
	RequireTenant bool
	RequireUser   bool
	RequireRole   bool
}

func RBACMiddleware(
	config MiddlewareConfig,
) func(http.Handler) http.Handler {

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				tenantID := strings.TrimSpace(
					r.Header.Get(
						TenantHeader,
					),
				)

				userID := strings.TrimSpace(
					r.Header.Get(
						UserHeader,
					),
				)

				role := strings.TrimSpace(
					r.Header.Get(
						RoleHeader,
					),
				)

				if config.RequireTenant &&
					tenantID == "" {

					http.Error(
						w,
						ErrMissingTenantContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if config.RequireUser &&
					userID == "" {

					http.Error(
						w,
						ErrMissingUserContext.Error(),
						http.StatusUnauthorized,
					)

					return
				}

				if config.RequireRole &&
					role == "" {

					http.Error(
						w,
						ErrUnauthorized.Error(),
						http.StatusForbidden,
					)

					return
				}

				ctx := r.Context()

				if tenantID != "" {

					ctx = context.WithValue(
						ctx,
						TenantContextKey,
						tenantID,
					)
				}

				if userID != "" {

					ctx = context.WithValue(
						ctx,
						UserContextKey,
						userID,
					)
				}

				if role != "" {

					ctx = context.WithValue(
						ctx,
						RoleContextKey,
						normalize(role),
					)
				}

				next.ServeHTTP(
					w,
					r.WithContext(
						ctx,
					),
				)
			},
		)
	}
}

func TenantIDFromContext(
	ctx context.Context,
) (string, bool) {

	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(
		TenantContextKey,
	).(string)

	if !ok {
		return "", false
	}

	value = strings.TrimSpace(
		value,
	)

	if value == "" {
		return "", false
	}

	return value, true
}

func UserIDFromContext(
	ctx context.Context,
) (string, bool) {

	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(
		UserContextKey,
	).(string)

	if !ok {
		return "", false
	}

	value = strings.TrimSpace(
		value,
	)

	if value == "" {
		return "", false
	}

	return value, true
}

func RoleFromContext(
	ctx context.Context,
) (string, bool) {

	if ctx == nil {
		return "", false
	}

	value, ok := ctx.Value(
		RoleContextKey,
	).(string)

	if !ok {
		return "", false
	}

	value = normalize(
		value,
	)

	if value == "" {
		return "", false
	}

	return value, true
}

func RequirePermission(
	manager *Manager,
	permission string,
) func(http.Handler) http.Handler {

	permission = normalize(
		permission,
	)

	return func(
		next http.Handler,
	) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if manager == nil {

					http.Error(
						w,
						ErrRBACUnavailable.Error(),
						http.StatusInternalServerError,
					)

					return
				}

				roleName, ok := RoleFromContext(
					r.Context(),
				)

				if !ok {

					http.Error(
						w,
						ErrUnauthorized.Error(),
						http.StatusForbidden,
					)

					return
				}

				allowed, err := manager.RoleHasPermission(
					r.Context(),
					roleName,
					permission,
				)

				if err != nil {

					http.Error(
						w,
						err.Error(),
						http.StatusForbidden,
					)

					return
				}

				if !allowed {

					http.Error(
						w,
						ErrAccessDenied.Error(),
						http.StatusForbidden,
					)

					return
				}

				next.ServeHTTP(
					w,
					r,
				)
			},
		)
	}
}
