package models

import (
	"fmt"
	"strings"
	"time"
)

type Permission struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Permission) Normalize() {
	if p == nil {
		return
	}
	p.Name = strings.TrimSpace(strings.ToLower(p.Name))
	p.Resource = strings.TrimSpace(strings.ToLower(p.Resource))
	p.Action = strings.TrimSpace(strings.ToLower(p.Action))
	p.Description = strings.TrimSpace(p.Description)
	p.TenantID = strings.TrimSpace(p.TenantID)
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
	}
}

func (p *Permission) Validate() error {
	if p == nil {
		return fmt.Errorf("%w: permission is required", ErrValidation)
	}
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("%w: permission id is required", ErrValidation)
	}
	if strings.TrimSpace(p.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for permission", ErrValidation)
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("%w: permission name is required", ErrValidation)
	}
	if strings.TrimSpace(p.Resource) == "" {
		return fmt.Errorf("%w: permission resource is required", ErrValidation)
	}
	if strings.TrimSpace(p.Action) == "" {
		return fmt.Errorf("%w: permission action is required", ErrValidation)
	}
	return nil
}

func (p *Permission) Key() string {
	if p == nil {
		return ""
	}
	return strings.ToLower(p.Resource) + ":" + strings.ToLower(p.Action)
}

func (p *Permission) Clone() Permission {
	if p == nil {
		return Permission{}
	}
	clone := *p
	return clone
}

func (p *Permission) Matches(resource, action string) bool {
	if p == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(p.Resource), strings.TrimSpace(resource)) &&
		strings.EqualFold(strings.TrimSpace(p.Action), strings.TrimSpace(action))
}

func (p *Permission) IsReadOnly() bool {
	if p == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(p.Action), "read")
}

func (p *Permission) IsWrite() bool {
	if p == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(p.Action)) {
	case "create", "update", "delete", "write":
		return true
	default:
		return false
	}
}

func (p *Permission) UpdateDescription(description string) {
	if p == nil {
		return
	}
	p.Description = strings.TrimSpace(description)
	p.UpdatedAt = time.Now().UTC()
}
