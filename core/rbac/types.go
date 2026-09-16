package rbac

import "strings"

// Effect represents the final authorization effect of a policy.
type Effect string

// Decision represents the result of an authorization evaluation.
type Decision struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

func (d Decision) IsAllowed() bool {
	return d.Allowed
}

func (d Decision) IsDenied() bool {
	return !d.Allowed
}

// Subject represents the identity requesting access to a protected resource.
type Subject struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
}

func (s *Subject) Validate() error {
	if s == nil {
		return ErrUnauthorized
	}

	if strings.TrimSpace(s.UserID) == "" {
		return ErrInvalidUserID
	}

	if strings.TrimSpace(s.TenantID) == "" {
		return ErrInvalidTenantID
	}

	return nil
}

func (s *Subject) Clone() *Subject {
	if s == nil {
		return nil
	}

	clone := *s
	clone.Roles = append([]string{}, s.Roles...)

	return &clone
}

func (s *Subject) HasRole(role string) bool {
	if s == nil {
		return false
	}

	return containsString(s.Roles, role)
}

func (s *Subject) RoleCount() int {
	if s == nil {
		return 0
	}

	return len(s.Roles)
}

func (s *Subject) AddRole(role string) error {
	if s == nil {
		return ErrUnauthorized
	}

	role = normalize(role)

	if role == "" {
		return ErrInvalidRole
	}

	if containsString(s.Roles, role) {
		return ErrRoleAlreadyAssigned
	}

	s.Roles = append(s.Roles, role)

	return nil
}

func (s *Subject) RemoveRole(role string) error {
	if s == nil {
		return ErrUnauthorized
	}

	if !containsString(s.Roles, role) {
		return ErrRoleNotAssigned
	}

	s.Roles = removeString(s.Roles, role)

	return nil
}

// ResourceAction represents an authorization target.
type ResourceAction struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func (ra ResourceAction) Validate() error {
	if strings.TrimSpace(ra.Resource) == "" {
		return ErrEmptyPermissionResource
	}

	if strings.TrimSpace(ra.Action) == "" {
		return ErrEmptyPermissionAction
	}

	return nil
}

func (ra ResourceAction) Key() string {
	return normalize(ra.Resource) + ":" + normalize(ra.Action)
}
