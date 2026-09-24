package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// TenantStatus is a technical lifecycle of a tenant workspace.
type TenantStatus string

const (
	TenantStatusPending   TenantStatus = "pending"
	TenantStatusActive    TenantStatus = "active"
	TenantStatusSuspended TenantStatus = "suspended"
	TenantStatusExpired   TenantStatus = "expired"
	TenantStatusDisabled  TenantStatus = "disabled"
	TenantStatusDeleted   TenantStatus = "deleted"
)

type Tenant struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Slug        string         `json:"slug" db:"slug"`
	Domain      string         `json:"domain,omitempty" db:"domain"`
	Status      TenantStatus   `json:"status" db:"status"`
	OwnerID     string         `json:"owner_id,omitempty" db:"owner_id"`
	Plan        string         `json:"plan" db:"plan"`
	Timezone    string         `json:"timezone,omitempty" db:"timezone"`
	Country     string         `json:"country,omitempty" db:"country"`
	Quotas      map[string]any `json:"quotas,omitempty" db:"quotas"`
	Branding    map[string]any `json:"branding,omitempty" db:"branding"`
	Metadata    map[string]any `json:"metadata,omitempty" db:"metadata"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
	SuspendedAt *time.Time     `json:"suspended_at,omitempty" db:"suspended_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty" db:"expires_at"`
}

var tenantSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func (t *Tenant) Normalize() {
	if t == nil {
		return
	}

	t.Name = strings.TrimSpace(t.Name)
	t.Slug = strings.TrimSpace(strings.ToLower(t.Slug))
	t.Domain = strings.TrimSpace(strings.ToLower(t.Domain))
	t.Timezone = strings.TrimSpace(t.Timezone)
	t.Country = strings.TrimSpace(strings.ToUpper(t.Country))
	t.Plan = strings.TrimSpace(t.Plan)

	if t.Status == "" {
		t.Status = TenantStatusPending
	}
	if t.Quotas == nil {
		t.Quotas = map[string]any{}
	}
	if t.Branding == nil {
		t.Branding = map[string]any{}
	}
	if t.Metadata == nil {
		t.Metadata = map[string]any{}
	}
}

func (t *Tenant) IsActive() bool {
	return t != nil && t.Status == TenantStatusActive
}

func (t *Tenant) IsPending() bool {
	return t != nil && t.Status == TenantStatusPending
}

func (t *Tenant) IsSuspended() bool {
	return t != nil && t.Status == TenantStatusSuspended
}

func (t *Tenant) IsExpired() bool {
	if t == nil {
		return false
	}

	if t.Status == TenantStatusExpired {
		return true
	}
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*t.ExpiresAt)
}

func (t *Tenant) IsDisabled() bool {
	return t != nil && t.Status == TenantStatusDisabled
}

func (t *Tenant) IsDeleted() bool {
	return t != nil && t.Status == TenantStatusDeleted
}

func (t *Tenant) HasCustomDomain() bool {
	return t != nil && strings.TrimSpace(t.Domain) != ""
}

func (t *Tenant) CanAccessPlatform() bool {
	if t == nil {
		return false
	}
	if t.IsDeleted() || t.IsDisabled() || t.IsSuspended() || t.IsExpired() {
		return false
	}
	return t.IsActive() || t.IsPending()
}

func (t *Tenant) IsOwner(userID string) bool {
	return t != nil && strings.TrimSpace(t.OwnerID) != "" && t.OwnerID == strings.TrimSpace(userID)
}

func (t *Tenant) Activate() {
	if t == nil {
		return
	}
	t.Status = TenantStatusActive
	t.SuspendedAt = nil
	t.UpdatedAt = time.Now().UTC()
}

func (t *Tenant) Suspend() {
	if t == nil {
		return
	}
	now := time.Now().UTC()
	t.Status = TenantStatusSuspended
	t.SuspendedAt = &now
	t.UpdatedAt = now
}

func (t *Tenant) Disable() {
	if t == nil {
		return
	}
	t.Status = TenantStatusDisabled
	t.UpdatedAt = time.Now().UTC()
}

func (t *Tenant) MarkDeleted() {
	if t == nil {
		return
	}
	t.Status = TenantStatusDeleted
	t.UpdatedAt = time.Now().UTC()
}

func (t *Tenant) MarkExpired() {
	if t == nil {
		return
	}
	t.Status = TenantStatusExpired
	t.UpdatedAt = time.Now().UTC()
}

func (t *Tenant) Validate() error {
	if t == nil {
		return fmt.Errorf("%w: tenant is required", ErrValidation)
	}

	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("%w: tenant id is required", ErrValidation)
	}
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("%w: tenant name is required", ErrValidation)
	}
	if strings.TrimSpace(t.Slug) == "" {
		return fmt.Errorf("%w: tenant slug is required", ErrValidation)
	}
	if !tenantSlugPattern.MatchString(t.Slug) {
		return fmt.Errorf("%w: tenant slug must be lowercase and may only contain letters, numbers, and hyphens", ErrValidation)
	}

	if t.Status != "" {
		switch t.Status {
		case TenantStatusPending, TenantStatusActive, TenantStatusSuspended, TenantStatusExpired, TenantStatusDisabled, TenantStatusDeleted:
		default:
			return fmt.Errorf("%w: unsupported tenant status %q", ErrValidation, t.Status)
		}
	}

	if t.ExpiresAt != nil && !t.CreatedAt.IsZero() && t.ExpiresAt.Before(t.CreatedAt) {
		return fmt.Errorf("%w: tenant expiration date cannot be before creation date", ErrValidation)
	}

	return nil
}

func (t *Tenant) Clone() Tenant {
	if t == nil {
		return Tenant{}
	}

	clone := *t
	clone.Quotas = cloneMap(t.Quotas)
	clone.Branding = cloneMap(t.Branding)
	clone.Metadata = cloneMap(t.Metadata)
	if t.SuspendedAt != nil {
		suspendedAt := *t.SuspendedAt
		clone.SuspendedAt = &suspendedAt
	}
	if t.ExpiresAt != nil {
		expiresAt := *t.ExpiresAt
		clone.ExpiresAt = &expiresAt
	}
	return clone
}

func cloneMap(src map[string]any) map[string]any {
	if src == nil {
		return map[string]any{}
	}
	out := make(map[string]any, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
