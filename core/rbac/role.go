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

	Permissions []string `json:"permissions"`

	ParentRoles []string `json:"parent_roles"`

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
			[]string,
			0,
		),

		ParentRoles: make(
			[]string,
			0,
		),

		CreatedAt: now,

		UpdatedAt: now,
	}
}

func (r *Role) Validate() error {

	if r == nil {
		return ErrInvalidRole
	}

	if strings.TrimSpace(
		r.ID,
	) == "" {

		return ErrEmptyRoleID
	}

	if strings.TrimSpace(
		r.Name,
	) == "" {

		return ErrEmptyRoleName
	}

	return nil
}

func (r *Role) Clone() *Role {

	if r == nil {
		return nil
	}

	clone := *r

	clone.Permissions = append(
		[]string{},
		r.Permissions...,
	)

	clone.ParentRoles = append(
		[]string{},
		r.ParentRoles...,
	)

	return &clone
}

func (r *Role) Touch() {

	if r == nil {
		return
	}

	r.UpdatedAt = time.Now().UTC()
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

func (r *Role) UpdateDescription(
	description string,
) {

	if r == nil {
		return
	}

	r.Description = strings.TrimSpace(
		description,
	)

	r.Touch()
}

func (r *Role) AddPermission(
	permission string,
) error {

	if r == nil {
		return ErrInvalidRole
	}

	permission = normalize(
		permission,
	)

	if permission == "" {
		return ErrInvalidPermission
	}

	if containsString(
		r.Permissions,
		permission,
	) {

		return ErrPermissionAlreadyAssigned
	}

	r.Permissions = append(
		r.Permissions,
		permission,
	)

	r.Touch()

	return nil
}

func (r *Role) RemovePermission(
	permission string,
) error {

	if r == nil {
		return ErrInvalidRole
	}

	if !containsString(
		r.Permissions,
		permission,
	) {

		return ErrPermissionNotAssigned
	}

	r.Permissions = removeString(
		r.Permissions,
		permission,
	)

	r.Touch()

	return nil
}

func (r *Role) HasPermission(
	permission string,
) bool {

	if r == nil {
		return false
	}

	return containsString(
		r.Permissions,
		permission,
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

func (r *Role) AddParentRole(
	roleName string,
) error {

	if r == nil {
		return ErrInvalidRole
	}

	roleName = normalize(
		roleName,
	)

	if roleName == "" {
		return ErrInvalidRole
	}

	if normalize(
		r.Name,
	) == roleName {

		return ErrRoleInheritanceLoop
	}

	if containsString(
		r.ParentRoles,
		roleName,
	) {

		return ErrRoleAlreadyAssigned
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

	if r == nil {
		return ErrInvalidRole
	}

	if !containsString(
		r.ParentRoles,
		roleName,
	) {

		return ErrRoleNotAssigned
	}

	r.ParentRoles = removeString(
		r.ParentRoles,
		roleName,
	)

	r.Touch()

	return nil
}

func (r *Role) HasParentRole(
	roleName string,
) bool {

	if r == nil {
		return false
	}

	return containsString(
		r.ParentRoles,
		roleName,
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

func (r *Role) ClearPermissions() {

	if r == nil {
		return
	}

	r.Permissions = make(
		[]string,
		0,
	)

	r.Touch()
}

func (r *Role) ClearParentRoles() {

	if r == nil {
		return
	}

	r.ParentRoles = make(
		[]string,
		0,
	)

	r.Touch()
}
