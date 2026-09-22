package models

import (
	"fmt"
	"strings"
	"time"
)

type FraudRiskLevel string

const (
	FraudRiskLow      FraudRiskLevel = "low"
	FraudRiskMedium   FraudRiskLevel = "medium"
	FraudRiskHigh     FraudRiskLevel = "high"
	FraudRiskCritical FraudRiskLevel = "critical"
)

type FraudSignalType string

const (
	FraudSignalTypeLogin    FraudSignalType = "login"
	FraudSignalTypePayment  FraudSignalType = "payment"
	FraudSignalTypeAPI      FraudSignalType = "api"
	FraudSignalTypeBehavior FraudSignalType = "behavior"
	FraudSignalTypeAccess   FraudSignalType = "access"
	FraudSignalTypeDevice   FraudSignalType = "device"
	FraudSignalTypeGeo      FraudSignalType = "geo"
	FraudSignalTypeCustom   FraudSignalType = "custom"
)

type FraudSignal struct {
	ID        string          `json:"id" db:"id"`
	TenantID  string          `json:"tenant_id" db:"tenant_id"`
	UserID    string          `json:"user_id,omitempty" db:"user_id"`
	Type      FraudSignalType `json:"type" db:"type"`
	Source    string          `json:"source,omitempty" db:"source"`
	Value     float64         `json:"value" db:"value"`
	Metadata  map[string]any  `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type FraudViolation struct {
	ID        string         `json:"id" db:"id"`
	TenantID  string         `json:"tenant_id" db:"tenant_id"`
	UserID    string         `json:"user_id,omitempty" db:"user_id"`
	SignalID  string         `json:"signal_id" db:"signal_id"`
	RuleID    string         `json:"rule_id" db:"rule_id"`
	Score     float64        `json:"score" db:"score"`
	Level     FraudRiskLevel `json:"level" db:"level"`
	Reason    string         `json:"reason" db:"reason"`
	Status    string         `json:"status" db:"status"`
	Metadata  map[string]any `json:"metadata,omitempty" db:"metadata"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

func (s *FraudSignal) Normalize() {
	if s == nil {
		return
	}
	if s.Type == "" {
		s.Type = FraudSignalTypeCustom
	}
	s.TenantID = strings.TrimSpace(s.TenantID)
	s.UserID = strings.TrimSpace(s.UserID)
	s.Source = strings.TrimSpace(s.Source)
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = s.CreatedAt
	}
}

func (s *FraudSignal) Validate() error {
	if s == nil {
		return fmt.Errorf("%w: fraud signal is required", ErrValidation)
	}
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("%w: fraud signal id is required", ErrValidation)
	}
	if strings.TrimSpace(s.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for fraud signal", ErrValidation)
	}
	if s.Type == "" {
		return fmt.Errorf("%w: fraud signal type is required", ErrValidation)
	}
	if s.Value < 0 {
		return fmt.Errorf("%w: fraud signal value cannot be negative", ErrValidation)
	}
	return nil
}

func (s *FraudSignal) Clone() FraudSignal {
	if s == nil {
		return FraudSignal{}
	}
	clone := *s
	if s.Metadata != nil {
		clone.Metadata = make(map[string]any, len(s.Metadata))
		for k, v := range s.Metadata {
			clone.Metadata[k] = v
		}
	}
	return clone
}

func (v *FraudViolation) Normalize() {
	if v == nil {
		return
	}
	v.TenantID = strings.TrimSpace(v.TenantID)
	v.UserID = strings.TrimSpace(v.UserID)
	v.SignalID = strings.TrimSpace(v.SignalID)
	v.RuleID = strings.TrimSpace(v.RuleID)
	v.Reason = strings.TrimSpace(v.Reason)
	v.Status = strings.TrimSpace(v.Status)
	if v.Level == "" {
		v.Level = FraudRiskLow
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now().UTC()
	}
	if v.UpdatedAt.IsZero() {
		v.UpdatedAt = v.CreatedAt
	}
}

func (v *FraudViolation) Validate() error {
	if v == nil {
		return fmt.Errorf("%w: fraud violation is required", ErrValidation)
	}
	if strings.TrimSpace(v.ID) == "" {
		return fmt.Errorf("%w: fraud violation id is required", ErrValidation)
	}
	if strings.TrimSpace(v.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required for fraud violation", ErrValidation)
	}
	if strings.TrimSpace(v.SignalID) == "" {
		return fmt.Errorf("%w: fraud signal id is required", ErrValidation)
	}
	if strings.TrimSpace(v.RuleID) == "" {
		return fmt.Errorf("%w: fraud rule id is required", ErrValidation)
	}
	if v.Score < 0 {
		return fmt.Errorf("%w: fraud violation score cannot be negative", ErrValidation)
	}
	switch v.Level {
	case FraudRiskLow, FraudRiskMedium, FraudRiskHigh, FraudRiskCritical:
	default:
		return fmt.Errorf("%w: unsupported fraud risk level %q", ErrValidation, v.Level)
	}
	return nil
}

func (v *FraudViolation) Clone() FraudViolation {
	if v == nil {
		return FraudViolation{}
	}
	clone := *v
	if v.Metadata != nil {
		clone.Metadata = make(map[string]any, len(v.Metadata))
		for k, value := range v.Metadata {
			clone.Metadata[k] = value
		}
	}
	return clone
}

func (v *FraudViolation) IsCritical() bool {
	return v != nil && v.Level == FraudRiskCritical
}
