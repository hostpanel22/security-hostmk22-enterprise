package fraud

import (
	"context"
	"sync"
)

type Store interface {
	SaveSignal(
		ctx context.Context,
		signal *Signal,
	) error

	SaveViolation(
		ctx context.Context,
		violation *Violation,
	) error

	ListSignals(
		ctx context.Context,
	) ([]*Signal, error)

	ListViolations(
		ctx context.Context,
	) ([]*Violation, error)

	Clear(
		ctx context.Context,
	) error
}

type MemoryStore struct {
	mu sync.RWMutex

	signals []*Signal

	violations []*Violation
}

func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		signals: make(
			[]*Signal,
			0,
		),
		violations: make(
			[]*Violation,
			0,
		),
	}
}

func (s *MemoryStore) SaveSignal(
	_ context.Context,
	signal *Signal,
) error {

	if signal == nil {
		return ErrInvalidSignal
	}

	if err := signal.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.signals = append(
		s.signals,
		signal.Clone(),
	)

	return nil
}

func (s *MemoryStore) SaveViolation(
	_ context.Context,
	violation *Violation,
) error {

	if violation == nil {
		return ErrFraudDetected
	}

	if err := violation.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	copyViolation := *violation

	s.violations = append(
		s.violations,
		&copyViolation,
	)

	return nil
}

func (s *MemoryStore) ListSignals(
	_ context.Context,
) ([]*Signal, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Signal,
		0,
		len(s.signals),
	)

	for _, signal := range s.signals {

		result = append(
			result,
			signal.Clone(),
		)
	}

	return result, nil
}

func (s *MemoryStore) ListViolations(
	_ context.Context,
) ([]*Violation, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Violation,
		0,
		len(s.violations),
	)

	for _, violation := range s.violations {

		copyViolation := *violation

		result = append(
			result,
			&copyViolation,
		)
	}

	return result, nil
}

func (s *MemoryStore) Clear(
	_ context.Context,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.signals = make(
		[]*Signal,
		0,
	)

	s.violations = make(
		[]*Violation,
		0,
	)

	return nil
}

func (s *MemoryStore) SignalCount() int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(
		s.signals,
	)
}

func (s *MemoryStore) ViolationCount() int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(
		s.violations,
	)
}
