package permissions

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	System bool `json:"system"`

	Permissions []*Permission `json:"permissions"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewRole(
	name string,
	description string,
) *Role {

	now := time.Now().UTC()

	return &Role{
		ID:          uuid.NewString(),
		Name:        normalize(name),
		Description: strings.TrimSpace(description),
		Permissions: make([]*Permission, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewSystemRole(
	name string,
	description string,
) *Role {

	role := NewRole(
		name,
		description,
	)

	role.System = true

	return role
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

	r.UpdatedAt = time.Now().UTC()
}

func (r *Role) AddPermission(
	permission *Permission,
) error {

	if permission == nil {
		return ErrInvalidPermission
	}

	if err := permission.Validate(); err != nil {
		return err
	}

	for _, existing := range r.Permissions {

		if existing.Key() ==
			permission.Key() {

			return ErrDuplicatePermission
		}
	}

	r.Permissions = append(
		r.Permissions,
		permission,
	)

	r.Touch()

	return nil
}

func (r *Role) RemovePermission(
	resource string,
	action string,
) bool {

	key := normalize(resource) +
		":" +
		normalize(action)

	for index, permission := range r.Permissions {

		if permission.Key() == key {

			r.Permissions = append(
				r.Permissions[:index],
				r.Permissions[index+1:]...,
			)

			r.Touch()

			return true
		}
	}

	return false
}

func (r *Role) HasPermission(
	resource string,
	action string,
) bool {

	for _, permission := range r.Permissions {

		if permission.Matches(
			resource,
			action,
		) {
			return true
		}
	}

	return false
}

func (r *Role) PermissionCount() int {

	return len(
		r.Permissions,
	)
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

		clone.Permissions = append(
			clone.Permissions,
			permission.Clone(),
		)
	}

	return &clone
}

func (r *Role) IsSystem() bool {

	return r.System
}
