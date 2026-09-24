package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Setting struct {
	ID        string          `json:"id" db:"id"`
	TenantID  string          `json:"tenant_id,omitempty" db:"tenant_id"`
	Key       string          `json:"key" db:"key"`
	Value     json.RawMessage `json:"value" db:"value"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

func (s *Setting) Normalize() {
	if s == nil {
		return
	}
	s.TenantID = strings.TrimSpace(s.TenantID)
	s.Key = strings.TrimSpace(s.Key)
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = time.Now().UTC()
	}
}

func (s *Setting) Validate() error {
	if s == nil {
		return fmt.Errorf("%w: setting is required", ErrValidation)
	}
	if strings.TrimSpace(s.Key) == "" {
		return fmt.Errorf("%w: setting key is required", ErrValidation)
	}
	if len(s.Value) == 0 {
		return fmt.Errorf("%w: setting value is required", ErrValidation)
	}
	return nil
}
