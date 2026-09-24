package models

import (
	"fmt"
	"strings"
	"time"
)

type FraudRule struct {
	ID            string         `json:"id" db:"id"`
	TenantID      string         `json:"tenant_id" db:"tenant_id"`
	Name          string         `json:"name" db:"name"`
	Enabled       bool           `json:"enabled" db:"enabled"`
	Field         string         `json:"field" db:"field"`
	Operator      string         `json:"operator" db:"operator"`
	Threshold     string         `json:"threshold" db:"threshold"`
	WindowSeconds *int           `json:"window_seconds,omitempty" db:"window_seconds"`
	Action        string         `json:"action" db:"action"`
	Score         int            `json:"score" db:"score"`
	Metadata      map[string]any `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

func (r *FraudRule) Normalize() {
	if r == nil {
		return
	}
	r.TenantID = strings.TrimSpace(r.TenantID)
	r.Name = strings.TrimSpace(r.Name)
	r.Field = strings.TrimSpace(r.Field)
	r.Operator = strings.TrimSpace(strings.ToLower(r.Operator))
	r.Threshold = strings.TrimSpace(r.Threshold)
	r.Action = strings.TrimSpace(strings.ToLower(r.Action))
	if r.Metadata == nil {
		r.Metadata = map[string]any{}
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = r.CreatedAt
	}
}

func (r *FraudRule) Validate() error {
	if r == nil {
		return fmt.Errorf("%w: fraud rule is required", ErrValidation)
	}
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("%w: fraud rule id is required", ErrValidation)
	}
	if strings.TrimSpace(r.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for fraud rule", ErrValidation)
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: fraud rule name is required", ErrValidation)
	}
	if strings.TrimSpace(r.Field) == "" {
		return fmt.Errorf("%w: fraud rule field is required", ErrValidation)
	}
	if strings.TrimSpace(r.Operator) == "" {
		return fmt.Errorf("%w: fraud rule operator is required", ErrValidation)
	}
	if strings.TrimSpace(r.Action) == "" {
		return fmt.Errorf("%w: fraud rule action is required", ErrValidation)
	}
	return nil
}

func (r *FraudRule) Clone() FraudRule {
	if r == nil {
		return FraudRule{}
	}
	clone := *r
	clone.Metadata = cloneMap(r.Metadata)
	if r.WindowSeconds != nil {
		w := *r.WindowSeconds
		clone.WindowSeconds = &w
	}
	return clone
}
