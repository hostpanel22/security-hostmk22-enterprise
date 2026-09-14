package tenancy

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

/*
|--------------------------------------------------------------------------
| Middleware Configuration
|--------------------------------------------------------------------------
*/

type MiddlewareConfig struct {
	RequireTenant bool

	EnableValidation bool

	EnableIsolation bool

	Resolver *Resolver

	Validator *Validator

	Isolation *IsolationEngine
}

/*
|--------------------------------------------------------------------------
| Tenant Middleware
|--------------------------------------------------------------------------
|
| Resolves, validates and injects the tenant
| context into every protected request.
|
*/

type TenantMiddleware struct {
	config MiddlewareConfig
}

/*
|--------------------------------------------------------------------------
| Constructor
|--------------------------------------------------------------------------
*/

func NewMiddleware(
	config MiddlewareConfig,
) *TenantMiddleware {

	if config.RequireTenant == false {
		config.RequireTenant = true
	}

	if config.EnableValidation == false {
		config.EnableValidation = true
	}

	if config.EnableIsolation == false {
		config.EnableIsolation = true
	}

	if config.Validator == nil {
		config.Validator = NewValidator()
	}

	if config.Isolation == nil {
		config.Isolation = NewIsolationEngine()
	}

	return &TenantMiddleware{
		config: config,
	}
}

/*
|--------------------------------------------------------------------------
| HTTP Middleware
|--------------------------------------------------------------------------
*/

func (m *TenantMiddleware) Handle(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {

			startedAt := time.Now().UTC()

			if m.config.Resolver == nil {

				m.writeError(
					writer,
					http.StatusInternalServerError,
					"tenant resolver not configured",
				)

				return
			}

			tenantCtx, err := m.config.Resolver.ResolveRequest(
				request,
			)

			if err != nil {

				if m.config.RequireTenant {

					m.writeError(
						writer,
						http.StatusUnauthorized,
						err.Error(),
					)

					return
				}

				next.ServeHTTP(
					writer,
					request,
				)

				return
			}

			if m.config.EnableValidation {

				if err := m.config.Validator.ValidateContext(
					tenantCtx,
				); err != nil {

					m.writeError(
						writer,
						http.StatusForbidden,
						err.Error(),
					)

					return
				}
			}

			ctx := WithTenant(
				request.Context(),
				tenantCtx,
			)

			ctx = context.WithValue(
				ctx,
				contextKey("hostmk.request.started_at"),
				startedAt,
			)

			request = request.WithContext(
				ctx,
			)

			next.ServeHTTP(
				writer,
				request,
			)
		},
	)
}

/*
|--------------------------------------------------------------------------
| Tenant Extraction Helper
|--------------------------------------------------------------------------
*/

func Tenant(
	request *http.Request,
) (*TenantContext, error) {

	return MustTenantFromContext(
		request.Context(),
	)
}

/*
|--------------------------------------------------------------------------
| Tenant Required Middleware
|--------------------------------------------------------------------------
*/

func (m *TenantMiddleware) RequireTenant(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {

			tenantCtx, err := Tenant(
				request,
			)

			if err != nil {

				m.writeError(
					writer,
					http.StatusUnauthorized,
					err.Error(),
				)

				return
			}

			if !tenantCtx.IsValid() {

				m.writeError(
					writer,
					http.StatusForbidden,
					ErrInvalidTenantContext.Error(),
				)

				return
			}

			next.ServeHTTP(
				writer,
				request,
			)
		},
	)
}

/*
|--------------------------------------------------------------------------
| JSON Error Response
|--------------------------------------------------------------------------
*/

func (m *TenantMiddleware) writeError(
	writer http.ResponseWriter,
	status int,
	message string,
) {

	writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	writer.WriteHeader(
		status,
	)

	_ = json.NewEncoder(
		writer,
	).Encode(
		map[string]any{
			"success": false,
			"error":   message,
			"status":  status,
		},
	)
}
