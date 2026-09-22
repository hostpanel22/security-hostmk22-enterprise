package models

import (
	"fmt"
	"strings"
	"time"
)

type SecretType string

const (
	SecretTypeAPIKey     SecretType = "api_key"
	SecretTypeJWT        SecretType = "jwt"
	SecretTypeDatabase   SecretType = "database"
	SecretTypeEncryption SecretType = "encryption"
	SecretTypeWebhook    SecretType = "webhook"
	SecretTypeCustom     SecretType = "custom"
)

type Secret struct {
	ID          string     `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Name        string     `json:"name" db:"name"`
	Type        SecretType `json:"type" db:"type"`
	Value       string     `json:"value,omitempty" db:"value"`
	MaskedValue string     `json:"masked_value,omitempty" db:"masked_value"`
	Description string     `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

func (s *Secret) Normalize() {
	if s == nil {
		return
	}
	s.TenantID = strings.TrimSpace(s.TenantID)
	s.Name = strings.TrimSpace(strings.ToUpper(s.Name))
	s.Description = strings.TrimSpace(s.Description)
	if s.Type == "" {
		s.Type = SecretTypeCustom
	}
	if s.Value != "" {
		s.MaskedValue = MaskSecret(s.Value)
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = s.CreatedAt
	}
}

func (s *Secret) Validate() error {
	if s == nil {
		return fmt.Errorf("%w: secret is required", ErrValidation)
	}
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("%w: secret id is required", ErrValidation)
	}
	if strings.TrimSpace(s.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for secret", ErrValidation)
	}
	if strings.TrimSpace(s.Name) == "" {
		return fmt.Errorf("%w: secret name is required", ErrValidation)
	}
	if strings.TrimSpace(s.Value) == "" {
		return fmt.Errorf("%w: secret value is required", ErrValidation)
	}
	switch s.Type {
	case SecretTypeAPIKey, SecretTypeJWT, SecretTypeDatabase, SecretTypeEncryption, SecretTypeWebhook, SecretTypeCustom:
	default:
		return fmt.Errorf("%w: unsupported secret type %q", ErrValidation, s.Type)
	}
	return nil
}

func (s *Secret) Clone() Secret {
	if s == nil {
		return Secret{}
	}
	clone := *s
	return clone
}

func (s *Secret) Mask() string {
	if s == nil {
		return ""
	}
	if s.Value == "" {
		return ""
	}
	return MaskSecret(s.Value)
}

func MaskSecret(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:4] + strings.Repeat("*", len(value)-4)
}
