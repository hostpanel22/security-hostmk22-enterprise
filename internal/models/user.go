package models

import (
	"fmt"
	"strings"
	"time"
)

// UserStatus is a technical lifecycle state of an account, not a business role.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBlocked   UserStatus = "blocked"
)

type User struct {
	ID              string     `json:"id" db:"id"`
	TenantID        string     `json:"tenant_id" db:"tenant_id"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	Name            string     `json:"name" db:"name"`
	Status          UserStatus `json:"status" db:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	MFAEnabled      bool       `json:"mfa_enabled" db:"mfa_enabled"`
	MFASecretEnc    []byte     `json:"-" db:"mfa_secret_enc"`
	RoleIDs         []string   `json:"role_ids,omitempty"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin       *time.Time `json:"last_login,omitempty" db:"last_login"`
}

func (u *User) Normalize() {
	if u == nil {
		return
	}
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	u.TenantID = strings.TrimSpace(u.TenantID)
	if u.Status == "" {
		u.Status = UserStatusPending
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
	}
}

func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("%w: user is required", ErrValidation)
	}
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("%w: user id is required", ErrValidation)
	}
	if strings.TrimSpace(u.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for user", ErrValidation)
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("%w: user email is required", ErrValidation)
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("%w: invalid user email", ErrValidation)
	}
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("%w: user name is required", ErrValidation)
	}
	switch u.Status {
	case UserStatusActive, UserStatusInactive, UserStatusPending, UserStatusSuspended, UserStatusBlocked:
	default:
		return fmt.Errorf("%w: unsupported user status %q", ErrValidation, u.Status)
	}
	return nil
}

func (u *User) IsActive() bool {
	return u != nil && u.Status == UserStatusActive
}

func (u *User) CanAuthenticate() bool {
	return u != nil && u.Status == UserStatusActive
}

func (u *User) IsBlocked() bool {
	return u != nil && (u.Status == UserStatusBlocked || u.Status == UserStatusSuspended)
}

func (u *User) EmailVerified() bool {
	return u != nil && u.EmailVerifiedAt != nil
}

func (u *User) Clone() User {
	if u == nil {
		return User{}
	}
	clone := *u
	clone.RoleIDs = append([]string{}, u.RoleIDs...)
	if u.LastLogin != nil {
		copyLastLogin := *u.LastLogin
		clone.LastLogin = &copyLastLogin
	}
	if u.EmailVerifiedAt != nil {
		copyVerified := *u.EmailVerifiedAt
		clone.EmailVerifiedAt = &copyVerified
	}
	if u.MFASecretEnc != nil {
		clone.MFASecretEnc = append([]byte{}, u.MFASecretEnc...)
	}
	return clone
}

func (u *User) Touch() {
	if u == nil {
		return
	}
	u.UpdatedAt = time.Now().UTC()
}
