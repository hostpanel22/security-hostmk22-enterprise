package secrets

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	Name string `json:"name"`

	Value string `json:"value"`

	Type string `json:"type"`

	Version int `json:"version"`

	Encrypted bool `json:"encrypted"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	RotatedAt *time.Time `json:"rotated_at,omitempty"`
}

func NewSecret(
	tenantID string,
	name string,
	value string,
	secretType string,
) *Secret {

	now := time.Now().UTC()

	return &Secret{
		ID: uuid.NewString(),

		TenantID: strings.TrimSpace(
			tenantID,
		),

		Name: strings.TrimSpace(
			name,
		),

		Value: value,

		Type: normalize(
			secretType,
		),

		Version: 1,

		Encrypted: false,

		CreatedAt: now,

		UpdatedAt: now,
	}
}

func (s *Secret) Validate() error {

	if s == nil {
		return ErrInvalidSecret
	}

	if strings.TrimSpace(
		s.TenantID,
	) == "" {
		return ErrEmptyTenantID
	}

	if strings.TrimSpace(
		s.Name,
	) == "" {
		return ErrEmptySecretName
	}

	if strings.TrimSpace(
		s.Value,
	) == "" {
		return ErrEmptySecretValue
	}

	if s.Version <= 0 {
		return ErrInvalidSecretVersion
	}

	return nil
}

func (s *Secret) Touch() {

	s.UpdatedAt = time.Now().UTC()
}

func (s *Secret) Rotate(
	newValue string,
) error {

	if strings.TrimSpace(
		newValue,
	) == "" {
		return ErrEmptySecretValue
	}

	s.Value = newValue

	// Important :
	// après rotation la nouvelle valeur
	// est en clair et devra être
	// rechiffrée par Vault.
	s.Encrypted = false

	s.Version++

	now := time.Now().UTC()

	s.RotatedAt = &now

	s.UpdatedAt = now

	return nil
}

func (s *Secret) Rename(
	name string,
) error {

	name = strings.TrimSpace(
		name,
	)

	if name == "" {
		return ErrEmptySecretName
	}

	s.Name = name

	s.Touch()

	return nil
}

func (s *Secret) ChangeType(
	secretType string,
) {

	s.Type = normalize(
		secretType,
	)

	s.Touch()
}

func (s *Secret) MarkEncrypted() {

	s.Encrypted = true

	s.Touch()
}

func (s *Secret) MarkDecrypted() {

	s.Encrypted = false

	s.Touch()
}

func (s *Secret) IsEncrypted() bool {

	return s.Encrypted
}

func (s *Secret) IsRotated() bool {

	return s.RotatedAt != nil
}

func (s *Secret) Age() time.Duration {

	return time.Since(
		s.CreatedAt,
	)
}

func (s *Secret) Clone() *Secret {

	if s == nil {
		return nil
	}

	clone := *s

	if s.RotatedAt != nil {

		t := *s.RotatedAt

		clone.RotatedAt = &t
	}

	return &clone
}

func (s *Secret) BelongsToTenant(
	tenantID string,
) bool {

	return strings.TrimSpace(
		s.TenantID,
	) == strings.TrimSpace(
		tenantID,
	)
}

func (s *Secret) IsType(
	secretType string,
) bool {

	return normalize(
		s.Type,
	) == normalize(
		secretType,
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
