package secrets

import (
	"time"
)

const (
	DefaultRotationDays = 90
)

type RotationPolicy struct {
	Interval time.Duration
}

func NewRotationPolicy(
	days int,
) RotationPolicy {

	if days <= 0 {
		days = DefaultRotationDays
	}

	return RotationPolicy{
		Interval: time.Hour * 24 * time.Duration(days),
	}
}

type RotationManager struct {
	Policy RotationPolicy
}

func NewRotationManager(
	policy RotationPolicy,
) *RotationManager {

	if policy.Interval <= 0 {
		policy = NewRotationPolicy(
			DefaultRotationDays,
		)
	}

	return &RotationManager{
		Policy: policy,
	}
}

func (r *RotationManager) NeedsRotation(
	secret *Secret,
) bool {

	if secret == nil {
		return false
	}

	return time.Now().UTC().After(
		r.NextRotationDate(secret),
	)
}

func (r *RotationManager) NextRotationDate(
	secret *Secret,
) time.Time {

	if secret == nil {
		return time.Time{}
	}

	if secret.RotatedAt != nil {
		return secret.RotatedAt.Add(
			r.Policy.Interval,
		)
	}

	return secret.CreatedAt.Add(
		r.Policy.Interval,
	)
}

func (r *RotationManager) Rotate(
	secret *Secret,
	newValue string,
) error {

	if secret == nil {
		return ErrInvalidSecret
	}

	if !r.NeedsRotation(
		secret,
	) {
		return ErrRotationNotRequired
	}

	return secret.Rotate(
		newValue,
	)
}

func (r *RotationManager) DaysUntilRotation(
	secret *Secret,
) int {

	if secret == nil {
		return 0
	}

	next := r.NextRotationDate(
		secret,
	)

	duration := time.Until(
		next,
	)

	return int(
		duration.Hours() / 24,
	)
}

func (r *RotationManager) IsExpired(
	secret *Secret,
) bool {

	if secret == nil {
		return false
	}

	return r.NeedsRotation(
		secret,
	)
}

func (r *RotationManager) ValidatePolicy() error {

	if r == nil {
		return ErrRotationScheduleInvalid
	}

	if r.Policy.Interval <= 0 {
		return ErrRotationScheduleInvalid
	}

	return nil
}

func (r *RotationManager) PolicyDays() int {

	if r == nil {
		return 0
	}

	return int(
		r.Policy.Interval.Hours() / 24,
	)
}
