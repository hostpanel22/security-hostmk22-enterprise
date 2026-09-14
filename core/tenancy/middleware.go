package tenancy

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

/*
|--------------------------------------------------------------------------
| Internal Context Keys
|--------------------------------------------------------------------------
*/

const (
	RequestStartedAtKey contextKey = "hostmk.request.started_at"
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
| Main Middleware
|--------------------------------------------------------------------------
|
| Resolves tenant information from the request,
| validates the context and injects it into
| the request context.
|
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

			tenantContext, err := m.config.Resolver.ResolveRequest(
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
					tenantContext,
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
				tenantContext,
			)

			ctx = context.WithValue(
				ctx,
				RequestStartedAtKey,
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
| Current Tenant Helper
|--------------------------------------------------------------------------
|
| Retrieves the active tenant context
| from the request.
|
*/

func CurrentTenant(
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

func (m *TenantMiddleware) RequireValidTenant(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {

			tenantContext, err := CurrentTenant(
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

			if !tenantContext.IsValid() {

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
| Ownership Validation Middleware
|--------------------------------------------------------------------------
|
| Can be used by applications before
| executing sensitive operations.
|
*/

func (m *TenantMiddleware) ValidateOwnership(
	request *http.Request,
	resourceTenantID string,
) error {

	tenantContext, err := CurrentTenant(
		request,
	)

	if err != nil {
		return err
	}

	if m.config.Isolation == nil {
		return ErrInvalidTenantContext
	}

	return m.config.Isolation.ValidateOwnership(
		tenantContext,
		resourceTenantID,
		"resource",
		"access",
	)
}

/*
|--------------------------------------------------------------------------
| Request Helpers
|--------------------------------------------------------------------------
*/

func RequestStartedAt(
	request *http.Request,
) (time.Time, bool) {

	value := request.Context().Value(
		RequestStartedAtKey,
	)

	startedAt, ok := value.(time.Time)

	return startedAt, ok
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
			"success":   false,
			"error":     message,
			"status":    status,
			"timestamp": time.Now().UTC(),
		},
	)
}
