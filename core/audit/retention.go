package audit

import (
	"context"
	"time"
)

type RetentionPolicy struct {
	Name string

	RetentionPeriod time.Duration

	Enabled bool
}

type RetentionResult struct {
	PolicyName string

	DeletedRecords int64

	ExecutedAt time.Time
}

type RetentionManager struct {
	store Store

	policies []RetentionPolicy
}

func NewRetentionManager(
	store Store,
) *RetentionManager {

	return &RetentionManager{
		store:    store,
		policies: make([]RetentionPolicy, 0),
	}
}

func (m *RetentionManager) AddPolicy(
	policy RetentionPolicy,
) error {

	if policy.RetentionPeriod <= 0 {
		return ErrInvalidRetentionPeriod
	}

	m.policies = append(
		m.policies,
		policy,
	)

	return nil
}

func (m *RetentionManager) Policies() []RetentionPolicy {

	result := make(
		[]RetentionPolicy,
		len(m.policies),
	)

	copy(
		result,
		m.policies,
	)

	return result
}

func (m *RetentionManager) Execute(
	ctx context.Context,
) ([]RetentionResult, error) {

	if m.store == nil {
		return nil, ErrStoreUnavailable
	}

	results := make(
		[]RetentionResult,
		0,
	)

	for _, policy := range m.policies {

		if !policy.Enabled {
			continue
		}

		limitDate := time.Now().
			UTC().
			Add(-policy.RetentionPeriod)

		deleted, err := m.store.DeleteOlderThan(
			ctx,
			limitDate,
		)

		if err != nil {
			return nil, err
		}

		results = append(
			results,
			RetentionResult{
				PolicyName:     policy.Name,
				DeletedRecords: deleted,
				ExecutedAt:     time.Now().UTC(),
			},
		)
	}

	return results, nil
}

func DefaultPolicies() []RetentionPolicy {

	return []RetentionPolicy{
		{
			Name:            "30-days",
			RetentionPeriod: 30 * 24 * time.Hour,
			Enabled:         true,
		},
		{
			Name:            "90-days",
			RetentionPeriod: 90 * 24 * time.Hour,
			Enabled:         true,
		},
		{
			Name:            "365-days",
			RetentionPeriod: 365 * 24 * time.Hour,
			Enabled:         true,
		},
	}
}

func (m *RetentionManager) LoadDefaults() error {

	for _, policy := range DefaultPolicies() {

		if err := m.AddPolicy(
			policy,
		); err != nil {
			return err
		}
	}

	return nil
}
