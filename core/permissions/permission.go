package permissions

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Resource string `json:"resource"`

	Action string `json:"action"`

	Description string `json:"description"`

	System bool `json:"system"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewPermission(
	name string,
	resource string,
	action string,
	description string,
) *Permission {

	now := time.Now().UTC()

	return &Permission{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(name),
		Resource:    normalize(resource),
		Action:      normalize(action),
		Description: strings.TrimSpace(description),
		System:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewSystemPermission(
	name string,
	resource string,
	action string,
	description string,
) *Permission {

	permission := NewPermission(
		name,
		resource,
		action,
		description,
	)

	permission.System = true

	return permission
}

func (p *Permission) Validate() error {

	if p == nil {
		return ErrInvalidPermission
	}

	if strings.TrimSpace(
		p.Name,
	) == "" {
		return ErrEmptyPermissionName
	}

	if strings.TrimSpace(
		p.Resource,
	) == "" {
		return ErrInvalidPermission
	}

	if strings.TrimSpace(
		p.Action,
	) == "" {
		return ErrInvalidPermission
	}

	return nil
}

func (p *Permission) Key() string {

	return normalize(
		p.Resource,
	) + ":" + normalize(
		p.Action,
	)
}

func (p *Permission) Equals(
	other *Permission,
) bool {

	if p == nil || other == nil {
		return false
	}

	return p.Key() == other.Key()
}

func (p *Permission) Clone() *Permission {

	if p == nil {
		return nil
	}

	clone := *p

	return &clone
}

func (p *Permission) Touch() {

	p.UpdatedAt = time.Now().UTC()
}

func (p *Permission) IsSystem() bool {

	return p.System
}

func (p *Permission) IsWildcard() bool {

	return p.Resource == "*" ||
		p.Action == "*"
}

func (p *Permission) Matches(
	resource string,
	action string,
) bool {

	resource = normalize(
		resource,
	)

	action = normalize(
		action,
	)

	if p.Resource == "*" &&
		p.Action == "*" {
		return true
	}

	if p.Resource == resource &&
		p.Action == "*" {
		return true
	}

	if p.Resource == "*" &&
		p.Action == action {
		return true
	}

	return p.Resource == resource &&
		p.Action == action
}

func normalize(
	value string,
) string {

	return strings.ToLower(
		strings.TrimSpace(
			value,
		),
	)
}
