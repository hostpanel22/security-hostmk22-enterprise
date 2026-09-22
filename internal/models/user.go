package models

import (
	"fmt"
	"strings"
	"time"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
	UserStatusBlocked  UserStatus = "blocked"
)

type User struct {
	ID        string     `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"password,omitempty" db:"password"`
	Name      string     `json:"name" db:"name"`
	Role      UserRole   `json:"role" db:"role"`
	Status    UserStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty" db:"last_login"`
}

func (u *User) Normalize() {
	if u == nil {
		return
	}
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	u.TenantID = strings.TrimSpace(u.TenantID)
	if u.Role == "" {
		u.Role = UserRoleUser
	}
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
	switch u.Role {
	case UserRoleAdmin, UserRoleUser, UserRoleGuest:
	default:
		return fmt.Errorf("%w: unsupported user role %q", ErrValidation, u.Role)
	}
	switch u.Status {
	case UserStatusActive, UserStatusInactive, UserStatusPending, UserStatusBlocked:
	default:
		return fmt.Errorf("%w: unsupported user status %q", ErrValidation, u.Status)
	}
	return nil
}

func (u *User) IsActive() bool {
	return u != nil && u.Status == UserStatusActive
}

func (u *User) IsBlocked() bool {
	return u != nil && u.Status == UserStatusBlocked
}

func (u *User) Clone() User {
	if u == nil {
		return User{}
	}
	clone := *u
	if u.LastLogin != nil {
		copyLastLogin := *u.LastLogin
		clone.LastLogin = &copyLastLogin
	}
	return clone
}

func (u *User) Touch() {
	if u == nil {
		return
	}
	u.UpdatedAt = time.Now().UTC()
}
