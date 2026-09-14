package tenancy

import (
	"context"
	"time"
)

/*
|--------------------------------------------------------------------------
| Internal Context Keys
|--------------------------------------------------------------------------
|
| Typed keys prevent collisions with external packages.
|
*/

type contextKey string

const (
	TenantContextKey contextKey = "hostmk.tenant.context"
)

/*
|--------------------------------------------------------------------------
| Tenant Context
|--------------------------------------------------------------------------
|
| Represents the active tenant attached to a request.
| This structure is propagated through the application
| lifecycle and used by security, audit and authorization
| layers.
|
*/

type TenantContext struct {
	TenantID string `json:"tenant_id"`

	TenantSlug string `json:"tenant_slug"`

	Domain string `json:"domain"`

	Source string `json:"source"`

	RequestID string `json:"request_id"`

	ActorID string `json:"actor_id"`

	IPAddress string `json:"ip_address"`

	UserAgent string `json:"user_agent"`

	ResolvedAt time.Time `json:"resolved_at"`
}

/*
|--------------------------------------------------------------------------
| Context Management
|--------------------------------------------------------------------------
*/

/*
WithTenant stores a tenant context into a request context.
*/
func WithTenant(
	ctx context.Context,
	tenantCtx *TenantContext,
) context.Context {
	return context.WithValue(
		ctx,
		TenantContextKey,
		tenantCtx,
	)
}

/*
TenantFromContext retrieves the tenant context.
*/
func TenantFromContext(
	ctx context.Context,
) (*TenantContext, bool) {

	tenantCtx, ok := ctx.Value(
		TenantContextKey,
	).(*TenantContext)

	return tenantCtx, ok
}

/*
MustTenantFromContext retrieves the tenant context
or returns a security error.
*/
func MustTenantFromContext(
	ctx context.Context,
) (*TenantContext, error) {

	tenantCtx, ok := TenantFromContext(ctx)

	if !ok || tenantCtx == nil {
		return nil, ErrMissingTenantContext
	}

	return tenantCtx, nil
}

/*
|--------------------------------------------------------------------------
| Validation
|--------------------------------------------------------------------------
*/

/*
IsValid validates minimum security requirements.
*/
func (t *TenantContext) IsValid() bool {
	if t == nil {
		return false
	}

	if t.TenantID == "" {
		return false
	}

	if t.ResolvedAt.IsZero() {
		return false
	}

	return true
}

/*
BelongsTo verifies that a resource belongs
to the active tenant.
*/
func (t *TenantContext) BelongsTo(
	tenantID string,
) bool {
	return t.TenantID == tenantID
}

/*
|--------------------------------------------------------------------------
| Security Helpers
|--------------------------------------------------------------------------
*/

/*
HasActor indicates whether the request
is associated with an authenticated actor.
*/
func (t *TenantContext) HasActor() bool {
	return t.ActorID != ""
}

/*
HasRequestID indicates whether a request
tracking identifier exists.
*/
func (t *TenantContext) HasRequestID() bool {
	return t.RequestID != ""
}

/*
HasDomain indicates whether a tenant
domain was resolved.
*/
func (t *TenantContext) HasDomain() bool {
	return t.Domain != ""
}

/*
|--------------------------------------------------------------------------
| Serialization
|--------------------------------------------------------------------------
*/

/*
ToMap converts the context into a structure
suitable for audit logging.
*/
func (t *TenantContext) ToMap() map[string]any {
	return map[string]any{
		"tenant_id":   t.TenantID,
		"tenant_slug": t.TenantSlug,
		"domain":      t.Domain,
		"source":      t.Source,
		"request_id":  t.RequestID,
		"actor_id":    t.ActorID,
		"ip_address":  t.IPAddress,
		"user_agent":  t.UserAgent,
		"resolved_at": t.ResolvedAt,
	}
}
