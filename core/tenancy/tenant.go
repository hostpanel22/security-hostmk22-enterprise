package tenancy

import (
	"time"
)

/*
|--------------------------------------------------------------------------
| Tenant Status
|--------------------------------------------------------------------------
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
|--------------------------------------------------------------------------
| Tenant
|--------------------------------------------------------------------------
|
| Represents a logical isolated customer
| inside a multi-tenant platform.
|
*/

type Tenant struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Domain string `json:"domain"`

	Status string `json:"status"`

	OwnerID string `json:"owner_id"`

	Plan string `json:"plan"`

	Timezone string `json:"timezone"`

	Country string `json:"country"`

	Metadata map[string]any `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	SuspendedAt *time.Time `json:"suspended_at,omitempty"`

	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

/*
|--------------------------------------------------------------------------
| Status Helpers
|--------------------------------------------------------------------------
*/

func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive
}

func (t *Tenant) IsPending() bool {
	return t.Status == StatusPending
}

func (t *Tenant) IsSuspended() bool {
	return t.Status == StatusSuspended
}

func (t *Tenant) IsDisabled() bool {
	return t.Status == StatusDisabled
}

func (t *Tenant) IsDeleted() bool {
	return t.Status == StatusDeleted
}

func (t *Tenant) IsExpired() bool {
	if t.Status == StatusExpired {
		return true
	}

	if t.ExpiresAt == nil {
		return false
	}

	return time.Now().UTC().After(
		*t.ExpiresAt,
	)
}

/*
|--------------------------------------------------------------------------
| Platform Access
|--------------------------------------------------------------------------
*/

func (t *Tenant) CanAccessPlatform() bool {

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
|--------------------------------------------------------------------------
| Domain Helpers
|--------------------------------------------------------------------------
*/

func (t *Tenant) HasCustomDomain() bool {
	return t.Domain != ""
}

/*
|--------------------------------------------------------------------------
| Lifecycle Helpers
|--------------------------------------------------------------------------
*/

func (t *Tenant) Suspend() {
	now := time.Now().UTC()

	t.Status = StatusSuspended
	t.SuspendedAt = &now
}

func (t *Tenant) Activate() {
	t.Status = StatusActive
	t.SuspendedAt = nil
}

func (t *Tenant) Disable() {
	t.Status = StatusDisabled
}

func (t *Tenant) MarkDeleted() {
	t.Status = StatusDeleted
}

func (t *Tenant) MarkExpired() {
	t.Status = StatusExpired
}

/*
|--------------------------------------------------------------------------
| Validation
|--------------------------------------------------------------------------
*/

func (t *Tenant) Validate() error {

	if t.ID == "" {
		return ErrInvalidTenantIdentifier
	}

	if t.Name == "" {
		return ErrTenantValidationFailed
	}

	if t.Slug == "" {
		return ErrTenantValidationFailed
	}

	return nil
}

/*
|--------------------------------------------------------------------------
| Audit Serialization
|--------------------------------------------------------------------------
*/

func (t *Tenant) ToMap() map[string]any {

	return map[string]any{
		"id":           t.ID,
		"name":         t.Name,
		"slug":         t.Slug,
		"domain":       t.Domain,
		"status":       t.Status,
		"owner_id":     t.OwnerID,
		"plan":         t.Plan,
		"timezone":     t.Timezone,
		"country":      t.Country,
		"created_at":   t.CreatedAt,
		"updated_at":   t.UpdatedAt,
		"suspended_at": t.SuspendedAt,
		"expires_at":   t.ExpiresAt,
	}
}
