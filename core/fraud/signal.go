package fraud

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Signal struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	UserID string `json:"user_id"`

	Type string `json:"type"`

	Value float64 `json:"value"`

	Source string `json:"source"`

	Metadata map[string]string `json:"metadata"`

	CreatedAt time.Time `json:"created_at"`
}

func NewSignal(
	tenantID string,
	userID string,
	signalType string,
	value float64,
	source string,
) *Signal {

	return &Signal{
		ID: uuid.NewString(),

		TenantID: strings.TrimSpace(tenantID),

		UserID: strings.TrimSpace(userID),

		Type: normalize(signalType),

		Value: value,

		Source: strings.TrimSpace(source),

		Metadata: make(
			map[string]string,
		),

		CreatedAt: time.Now().UTC(),
	}
}

func (s *Signal) Validate() error {

	if s == nil {
		return ErrInvalidSignal
	}

	if strings.TrimSpace(
		s.TenantID,
	) == "" {
		return ErrEmptyTenantID
	}

	if strings.TrimSpace(
		s.UserID,
	) == "" {
		return ErrEmptyUserID
	}

	if strings.TrimSpace(
		s.Type,
	) == "" {
		return ErrEmptySignalType
	}

	return nil
}

func (s *Signal) AddMetadata(
	key string,
	value string,
) {

	if s.Metadata == nil {

		s.Metadata = make(
			map[string]string,
		)
	}

	s.Metadata[key] = value
}

func (s *Signal) GetMetadata(
	key string,
) string {

	if s.Metadata == nil {
		return ""
	}

	return s.Metadata[key]
}

func (s *Signal) Clone() *Signal {

	if s == nil {
		return nil
	}

	clone := *s

	clone.Metadata = make(
		map[string]string,
	)

	for key, value := range s.Metadata {

		clone.Metadata[key] = value
	}

	return &clone
}
