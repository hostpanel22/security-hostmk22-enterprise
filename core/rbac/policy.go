package rbac

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type PolicyEffect string

const (
	PolicyAllow PolicyEffect = "allow"
	PolicyDeny  PolicyEffect = "deny"
)

type Policy struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	Name string `json:"name"`

	Description string `json:"description"`

	Effect PolicyEffect `json:"effect"`

	Resources []string `json:"resources"`

	Actions []string `json:"actions"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewPolicy(
	tenantID string,
	name string,
	description string,
	effect PolicyEffect,
) *Policy {

	now := time.Now().UTC()

	if effect != PolicyAllow &&
		effect != PolicyDeny {

		effect = PolicyAllow
	}

	return &Policy{
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

		Effect: effect,

		Resources: make(
			[]string,
			0,
		),

		Actions: make(
			[]string,
			0,
		),

		CreatedAt: now,

		UpdatedAt: now,
	}
}

func (p *Policy) Validate() error {

	if p == nil {
		return ErrInvalidPolicy
	}

	if strings.TrimSpace(
		p.ID,
	) == "" {

		return ErrInvalidPolicy
	}

	if strings.TrimSpace(
		p.Name,
	) == "" {

		return ErrEmptyPolicyName
	}

	if p.Effect != PolicyAllow &&
		p.Effect != PolicyDeny {

		return ErrInvalidPolicy
	}

	return nil
}

func (p *Policy) Clone() *Policy {

	if p == nil {
		return nil
	}

	clone := *p

	clone.Resources = append(
		[]string{},
		p.Resources...,
	)

	clone.Actions = append(
		[]string{},
		p.Actions...,
	)

	return &clone
}

func (p *Policy) Touch() {

	if p == nil {
		return
	}

	p.UpdatedAt = time.Now().UTC()
}

func (p *Policy) IsAllow() bool {

	if p == nil {
		return false
	}

	return p.Effect == PolicyAllow
}

func (p *Policy) IsDeny() bool {

	if p == nil {
		return false
	}

	return p.Effect == PolicyDeny
}

func (p *Policy) AddResource(
	resource string,
) {

	if p == nil {
		return
	}

	resource = normalize(
		resource,
	)

	if resource == "" {
		return
	}

	if containsString(
		p.Resources,
		resource,
	) {

		return
	}

	p.Resources = append(
		p.Resources,
		resource,
	)

	p.Touch()
}

func (p *Policy) RemoveResource(
	resource string,
) {

	if p == nil {
		return
	}

	p.Resources = removeString(
		p.Resources,
		resource,
	)

	p.Touch()
}

func (p *Policy) HasResource(
	resource string,
) bool {

	if p == nil {
		return false
	}

	return containsString(
		p.Resources,
		resource,
	)
}

func (p *Policy) AddAction(
	action string,
) {

	if p == nil {
		return
	}

	action = normalize(
		action,
	)

	if action == "" {
		return
	}

	if containsString(
		p.Actions,
		action,
	) {

		return
	}

	p.Actions = append(
		p.Actions,
		action,
	)

	p.Touch()
}

func (p *Policy) RemoveAction(
	action string,
) {

	if p == nil {
		return
	}

	p.Actions = removeString(
		p.Actions,
		action,
	)

	p.Touch()
}

func (p *Policy) HasAction(
	action string,
) bool {

	if p == nil {
		return false
	}

	return containsString(
		p.Actions,
		action,
	)
}

func (p *Policy) Matches(
	resource string,
	action string,
) bool {

	if p == nil {
		return false
	}

	return p.HasResource(
		resource,
	) &&
		p.HasAction(
			action,
		)
}

func (p *Policy) BelongsToTenant(
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

func (p *Policy) UpdateDescription(
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

func (p *Policy) ResourceCount() int {

	if p == nil {
		return 0
	}

	return len(
		p.Resources,
	)
}

func (p *Policy) ActionCount() int {

	if p == nil {
		return 0
	}

	return len(
		p.Actions,
	)
}

func (p *Policy) ClearResources() {

	if p == nil {
		return
	}

	p.Resources = make(
		[]string,
		0,
	)

	p.Touch()
}

func (p *Policy) ClearActions() {

	if p == nil {
		return
	}

	p.Actions = make(
		[]string,
		0,
	)

	p.Touch()
}
