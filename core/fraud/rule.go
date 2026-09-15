package fraud

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Rule struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	SignalType string `json:"signal_type"`

	Threshold float64 `json:"threshold"`

	Weight int `json:"weight"`

	Enabled bool `json:"enabled"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

func NewRule(
	name string,
	description string,
	signalType string,
	threshold float64,
	weight int,
) *Rule {

	now := time.Now().UTC()

	return &Rule{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		SignalType:  normalize(signalType),
		Threshold:   threshold,
		Weight:      weight,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (r *Rule) Validate() error {

	if r == nil {
		return ErrInvalidRule
	}

	if strings.TrimSpace(
		r.Name,
	) == "" {
		return ErrEmptyRuleName
	}

	if strings.TrimSpace(
		r.SignalType,
	) == "" {
		return ErrEmptySignalType
	}

	if r.Threshold <= 0 {
		return ErrInvalidThreshold
	}

	if r.Weight <= 0 {
		return ErrInvalidWeight
	}

	return nil
}

func (r *Rule) Enable() {
	r.Enabled = true
	r.Touch()
}

func (r *Rule) Disable() {
	r.Enabled = false
	r.Touch()
}

func (r *Rule) Touch() {
	r.UpdatedAt = time.Now().UTC()
}

func (r *Rule) Match(
	signal *Signal,
) bool {

	if signal == nil {
		return false
	}

	if !r.Enabled {
		return false
	}

	if normalize(signal.Type) !=
		normalize(r.SignalType) {
		return false
	}

	return signal.Value >= r.Threshold
}

func (r *Rule) Clone() *Rule {

	if r == nil {
		return nil
	}

	clone := *r

	return &clone
}
