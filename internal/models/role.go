package models

import (
	"fmt"
	"strings"
	"time"
)

type Role struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (r *Role) Normalize() {
	if r == nil {
		return
	}
	r.Name = strings.TrimSpace(strings.ToLower(r.Name))
	r.Description = strings.TrimSpace(r.Description)
	r.TenantID = strings.TrimSpace(r.TenantID)
	permissions := make([]string, 0, len(r.Permissions))
	seen := make(map[string]struct{}, len(r.Permissions))
	for _, permission := range r.Permissions {
		permission = strings.TrimSpace(strings.ToLower(permission))
		if permission == "" {
			continue
		}
		if _, exists := seen[permission]; exists {
			continue
		}
		seen[permission] = struct{}{}
		permissions = append(permissions, permission)
	}
	r.Permissions = permissions
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = r.CreatedAt
	}
}

func (r *Role) Validate() error {
	if r == nil {
		return fmt.Errorf("%w: role is required", ErrValidation)
	}
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("%w: role id is required", ErrValidation)
	}
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for role", ErrValidation)
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: role name is required", ErrValidation)
	}
	return nil
}

func (r *Role) Touch() {
	if r == nil {
		return
	}
	r.UpdatedAt = time.Now().UTC()
}

func (r *Role) Clone() Role {
	if r == nil {
		return Role{}
	}
	clone := *r
	clone.Permissions = append([]string{}, r.Permissions...)
	return clone
}

func (r *Role) AddPermission(permission string) error {
	if r == nil {
		return fmt.Errorf("%w: role is required", ErrValidation)
	}
	permission = strings.TrimSpace(strings.ToLower(permission))
	if permission == "" {
		return fmt.Errorf("%w: permission is required", ErrValidation)
	}
	for _, current := range r.Permissions {
		if current == permission {
			return nil
		}
	}
	r.Permissions = append(r.Permissions, permission)
	r.Touch()
	return nil
}

func (r *Role) RemovePermission(permission string) error {
	if r == nil {
		return fmt.Errorf("%w: role is required", ErrValidation)
	}
	permission = strings.TrimSpace(strings.ToLower(permission))
	for i, current := range r.Permissions {
		if current == permission {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.Touch()
			return nil
		}
	}
	return fmt.Errorf("%w: permission not assigned to role", ErrNotFound)
}

func (r *Role) HasPermission(permission string) bool {
	if r == nil {
		return false
	}
	permission = strings.TrimSpace(strings.ToLower(permission))
	for _, current := range r.Permissions {
		if current == permission {
			return true
		}
	}
	return false
}

func (r *Role) PermissionCount() int {
	if r == nil {
		return 0
	}
	return len(r.Permissions)
}
