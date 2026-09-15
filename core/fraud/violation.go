package fraud

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Violation struct {
	ID string `json:"id"`

	TenantID string `json:"tenant_id"`

	UserID string `json:"user_id"`

	SignalID string `json:"signal_id"`

	RuleID string `json:"rule_id"`

	Reason string `json:"reason"`

	Score int `json:"score"`

	Level RiskLevel `json:"level"`

	CreatedAt time.Time `json:"created_at"`
}

func NewViolation(
	tenantID string,
	userID string,
	signalID string,
	ruleID string,
	reason string,
	score int,
	level RiskLevel,
) *Violation {

	return &Violation{
		ID: uuid.NewString(),

		TenantID: strings.TrimSpace(
			tenantID,
		),

		UserID: strings.TrimSpace(
			userID,
		),

		SignalID: strings.TrimSpace(
			signalID,
		),

		RuleID: strings.TrimSpace(
			ruleID,
		),

		Reason: strings.TrimSpace(
			reason,
		),

		Score: score,

		Level: level,

		CreatedAt: time.Now().UTC(),
	}
}

func (v *Violation) Validate() error {

	if v == nil {
		return ErrFraudDetected
	}

	if strings.TrimSpace(
		v.TenantID,
	) == "" {
		return ErrEmptyTenantID
	}

	if strings.TrimSpace(
		v.UserID,
	) == "" {
		return ErrEmptyUserID
	}

	if strings.TrimSpace(
		v.Reason,
	) == "" {
		return ErrInvalidSignal
	}

	return nil
}

func (v *Violation) IsCritical() bool {

	return v.Level ==
		RiskCritical
}

func (v *Violation) IsHigh() bool {

	return v.Level ==
		RiskHigh
}

func (v *Violation) RequiresBlock() bool {

	return v.Level == RiskHigh ||
		v.Level == RiskCritical
}
