package rbac

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	Name string `json:"name"`

	Resource string `json:"resource"`

	Action string `json:"action"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewPermission(
	tenantID string,
	name string,
	resource string,
	action string,
	description string,
) *Permission {

	now := time.Now().UTC()

	return &Permission{
		ID: uuid.NewString(),

		TenantID: strings.TrimSpace(
			tenantID,
		),

		Name: normalize(
			name,
		),

		Resource: normalize(
			resource,
		),

		Action: normalize(
			action,
		),

		Description: strings.TrimSpace(
			description,
		),

		CreatedAt: now,

		UpdatedAt: now,
	}
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

		return ErrEmptyPermissionResource
	}

	if strings.TrimSpace(
		p.Action,
	) == "" {

		return ErrEmptyPermissionAction
	}

	return nil
}

func (p *Permission) Touch() {

	if p == nil {
		return
	}

	p.UpdatedAt = time.Now().UTC()
}

func (p *Permission) Clone() *Permission {

	if p == nil {
		return nil
	}

	clone := *p

	return &clone
}

func (p *Permission) Key() string {

	if p == nil {
		return ""
	}

	return normalize(
		p.Resource,
	) + ":" + normalize(
		p.Action,
	)
}

func (p *Permission) Equals(
	other *Permission,
) bool {

	if p == nil ||
		other == nil {

		return false
	}

	return p.Key() ==
		other.Key()
}

func (p *Permission) BelongsToTenant(
	tenantID string,
) bool {

	if p == nil {
		return false
	}

	return strings.TrimSpace(
		p.TenantID,
	) == strings.TrimSpace(
		tenantID,
	)
}

func (p *Permission) UpdateDescription(
	description string,
) {

	if p == nil {
		return
	}

	p.Description = strings.TrimSpace(
		description,
	)

	p.Touch()
}

func (p *Permission) IsResource(
	resource string,
) bool {

	if p == nil {
		return false
	}

	return normalize(
		p.Resource,
	) == normalize(
		resource,
	)
}

func (p *Permission) IsAction(
	action string,
) bool {

	if p == nil {
		return false
	}

	return normalize(
		p.Action,
	) == normalize(
		action,
	)
}

func (p *Permission) Matches(
	resource string,
	action string,
) bool {

	if p == nil {
		return false
	}

	return p.IsResource(
		resource,
	) && p.IsAction(
		action,
	)
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
