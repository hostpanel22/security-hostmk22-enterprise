package models

import (
	"fmt"
	"strings"
	"time"
)

type APIKey struct {
	ID         string     `json:"id" db:"id"`
	TenantID   string     `json:"tenant_id" db:"tenant_id"`
	Name       string     `json:"name" db:"name"`
	Prefix     string     `json:"prefix" db:"prefix"`
	KeyHash    string     `json:"-" db:"key_hash"`
	Plaintext  string     `json:"key,omitempty"`
	Scopes     []string   `json:"scopes"`
	RateLimit  *int       `json:"rate_limit,omitempty" db:"rate_limit"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

func (k *APIKey) Normalize() {
	if k == nil {
		return
	}
	k.TenantID = strings.TrimSpace(k.TenantID)
	k.Name = strings.TrimSpace(k.Name)
	k.Prefix = strings.TrimSpace(k.Prefix)
	if k.CreatedAt.IsZero() {
		k.CreatedAt = time.Now().UTC()
	}
	if k.UpdatedAt.IsZero() {
		k.UpdatedAt = k.CreatedAt
	}
}

func (k *APIKey) Validate() error {
	if k == nil {
		return fmt.Errorf("%w: api key is required", ErrValidation)
	}
	if strings.TrimSpace(k.ID) == "" {
		return fmt.Errorf("%w: api key id is required", ErrValidation)
	}
	if strings.TrimSpace(k.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for api key", ErrValidation)
	}
	if strings.TrimSpace(k.Name) == "" {
		return fmt.Errorf("%w: api key name is required", ErrValidation)
	}
	return nil
}

func (k *APIKey) IsActive() bool {
	if k == nil || k.RevokedAt != nil {
		return false
	}
	if k.ExpiresAt != nil && time.Now().UTC().After(*k.ExpiresAt) {
		return false
	}
	return true
}

func (k *APIKey) PublicView() APIKey {
	if k == nil {
		return APIKey{}
	}
	clone := *k
	clone.KeyHash = ""
	clone.Plaintext = ""
	clone.Scopes = append([]string{}, k.Scopes...)
	return clone
}
