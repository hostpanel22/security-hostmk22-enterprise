package tenancy

import (
	"time"
)

/*
TenantStatus defines the allowed lifecycle states
for a tenant within the platform.
*/
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusSuspended = "suspended"
	StatusExpired   = "expired"
	StatusDisabled  = "disabled"
	StatusDeleted   = "deleted"
)

/*
Tenant represents a logical customer isolated
inside a multi-tenant environment.

A Tenant must never access resources that belong
to another Tenant.
*/
type Tenant struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Domain string `json:"domain"`

	Status string `json:"status"`

	OwnerID string `json:"owner_id"`

	Timezone string `json:"timezone"`

	Country string `json:"country"`

	Plan string `json:"plan"`

	Metadata map[string]any `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	SuspendedAt *time.Time `json:"suspended_at,omitempty"`

	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

/*
IsActive determines whether the tenant
is allowed to operate normally.
*/
func (t Tenant) IsActive() bool {
	return t.Status == StatusActive
}

/*
IsPending indicates that the tenant
has not yet been fully activated.
*/
func (t Tenant) IsPending() bool {
	return t.Status == StatusPending
}

/*
IsSuspended indicates that access
has been administratively suspended.
*/
func (t Tenant) IsSuspended() bool {
	return t.Status == StatusSuspended
}

/*
IsExpired indicates that the tenant
subscription or contract has expired.
*/
func (t Tenant) IsExpired() bool {
	if t.Status == StatusExpired {
		return true
	}

	if t.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*t.ExpiresAt)
}

/*
IsDisabled indicates that the tenant
is permanently disabled.
*/
func (t Tenant) IsDisabled() bool {
	return t.Status == StatusDisabled
}

/*
IsDeleted indicates that the tenant
has been logically removed.
*/
func (t Tenant) IsDeleted() bool {
	return t.Status == StatusDeleted
}

/*
HasCustomDomain determines whether
the tenant is using a custom domain.
*/
func (t Tenant) HasCustomDomain() bool {
	return t.Domain != ""
}

/*
CanAccessPlatform determines whether
the tenant can access protected resources.
*/
func (t Tenant) CanAccessPlatform() bool {
	if t.IsDeleted() {
		return false
	}

	if t.IsDisabled() {
		return false
	}

	if t.IsSuspended() {
		return false
	}

	if t.IsExpired() {
		return false
	}

	return t.IsActive()
}

/*
RequiresReview determines whether
security validation should be reinforced.
*/
func (t Tenant) RequiresReview() bool {
	return t.IsPending() ||
		t.IsSuspended() ||
		t.IsExpired()
}

/*
ToMap converts the tenant into a generic map
for logging, auditing or serialization layers.
*/
func (t Tenant) ToMap() map[string]any {
	return map[string]any{
		"id":         t.ID,
		"name":       t.Name,
		"slug":       t.Slug,
		"domain":     t.Domain,
		"status":     t.Status,
		"owner_id":   t.OwnerID,
		"timezone":   t.Timezone,
		"country":    t.Country,
		"plan":       t.Plan,
		"created_at": t.CreatedAt,
		"updated_at": t.UpdatedAt,
	}
}
