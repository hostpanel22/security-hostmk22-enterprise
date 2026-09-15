package rbac

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	Name string `json:"name"`

	Description string `json:"description"`

	Permissions []*Permission `json:"permissions"`

	ParentRoles []string `json:"parent_roles"`

	System bool `json:"system"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewRole(
	tenantID string,
	name string,
	description string,
) *Role {

	now := time.Now().UTC()

	return &Role{
		ID: uuid.NewString(),

		TenantID: strings.TrimSpace(
			tenantID,
		),

		Name: normalize(
			name,
		),

		Description: strings.TrimSpace(
			description,
		),

		Permissions: make(
			[]*Permission,
			0,
		),

		ParentRoles: make(
			[]string,
			0,
		),

		System: false,

		CreatedAt: now,

		UpdatedAt: now,
	}
}

func (r *Role) Validate() error {

	if r == nil {
		return ErrInvalidRole
	}

	if strings.TrimSpace(
		r.Name,
	) == "" {

		return ErrEmptyRoleName
	}

	return nil
}

func (r *Role) Touch() {

	if r == nil {
		return
	}

	r.UpdatedAt = time.Now().UTC()
}

func (r *Role) Clone() *Role {

	if r == nil {
		return nil
	}

	clone := *r

	clone.Permissions = make(
		[]*Permission,
		0,
		len(r.Permissions),
	)

	for _, permission := range r.Permissions {

		if permission == nil {
			continue
		}

		clone.Permissions = append(
			clone.Permissions,
			permission.Clone(),
		)
	}

	clone.ParentRoles = append(
		[]string{},
		r.ParentRoles...,
	)

	return &clone
}

func (r *Role) AddPermission(
	permission *Permission,
) error {

	if r == nil {
		return ErrInvalidRole
	}

	if permission == nil {
		return ErrInvalidPermission
	}

	for _, existing := range r.Permissions {

		if existing == nil {
			continue
		}

		if existing.Equals(
			permission,
		) {

			return ErrPermissionAlreadyAssigned
		}
	}

	r.Permissions = append(
		r.Permissions,
		permission.Clone(),
	)

	r.Touch()

	return nil
}

func (r *Role) RemovePermission(
	resource string,
	action string,
) error {

	if r == nil {
		return ErrInvalidRole
	}

	filtered := make(
		[]*Permission,
		0,
	)

	removed := false

	for _, permission := range r.Permissions {

		if permission == nil {
			continue
		}

		if permission.Matches(
			resource,
			action,
		) {

			removed = true
			continue
		}

		filtered = append(
			filtered,
			permission,
		)
	}

	if !removed {
		return ErrPermissionNotAssigned
	}

	r.Permissions = filtered

	r.Touch()

	return nil
}

func (r *Role) HasPermission(
	resource string,
	action string,
) bool {

	if r == nil {
		return false
	}

	for _, permission := range r.Permissions {

		if permission == nil {
			continue
		}

		if permission.Matches(
			resource,
			action,
		) {

			return true
		}
	}

	return false
}

func (r *Role) AddParentRole(
	roleName string,
) error {

	roleName = normalize(
		roleName,
	)

	if roleName == "" {
		return ErrRoleNotFound
	}

	for _, existing := range r.ParentRoles {

		if normalize(
			existing,
		) == roleName {

			return ErrRoleAlreadyAssigned
		}
	}

	r.ParentRoles = append(
		r.ParentRoles,
		roleName,
	)

	r.Touch()

	return nil
}

func (r *Role) RemoveParentRole(
	roleName string,
) error {

	roleName = normalize(
		roleName,
	)

	filtered := make(
		[]string,
		0,
	)

	removed := false

	for _, existing := range r.ParentRoles {

		if normalize(
			existing,
		) == roleName {

			removed = true
			continue
		}

		filtered = append(
			filtered,
			existing,
		)
	}

	if !removed {
		return ErrRoleNotAssigned
	}

	r.ParentRoles = filtered

	r.Touch()

	return nil
}

func (r *Role) BelongsToTenant(
	tenantID string,
) bool {

	if r == nil {
		return false
	}

	return strings.TrimSpace(
		r.TenantID,
	) == strings.TrimSpace(
		tenantID,
	)
}

func (r *Role) PermissionCount() int {

	if r == nil {
		return 0
	}

	return len(
		r.Permissions,
	)
}

func (r *Role) ParentRoleCount() int {

	if r == nil {
		return 0
	}

	return len(
		r.ParentRoles,
	)
}

func (r *Role) MarkSystem() {

	if r == nil {
		return
	}

	r.System = true

	r.Touch()
}

func (r *Role) IsSystem() bool {

	if r == nil {
		return false
	}

	return r.System
}
