package secrets

import (
	"context"
	"net/http"
	"strings"
)

const (
	TenantHeader = "X-Tenant-ID"
	UserHeader   = "X-User-ID"
)

type contextKey string

const (
	TenantContextKey contextKey = "secret_tenant_id"
	UserContextKey   contextKey = "secret_user_id"
)

type MiddlewareConfig struct {
	RequireTenant bool
	RequireUser   bool
}

func SecretMiddleware(
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
