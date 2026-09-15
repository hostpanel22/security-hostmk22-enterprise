package audit

import (
	"context"
	"sync"
	"time"
)

type Store interface {
	Save(ctx context.Context, event *Event) error

	FindByID(
		ctx context.Context,
		id string,
	) (*Event, error)

	List(
		ctx context.Context,
		filter Filter,
	) ([]*Event, error)

	Delete(
		ctx context.Context,
		id string,
	) error

	DeleteOlderThan(
		ctx context.Context,
		timestamp time.Time,
	) (int64, error)

	Count(
		ctx context.Context,
		filter Filter,
	) (int64, error)
}

type MemoryStore struct {
	mu sync.RWMutex

	events map[string]*Event
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		events: make(map[string]*Event),
	}
}

func (s *MemoryStore) Save(
	ctx context.Context,
	event *Event,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if event == nil {
		return ErrNilEvent
	}

	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event.Clone()

	return nil
}

func (s *MemoryStore) FindByID(
	ctx context.Context,
	id string,
) (*Event, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if id == "" {
		return nil, ErrInvalidID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]

	if !exists {
		return nil, ErrStoreReadFailed
	}

	return event.Clone(), nil
}

func (s *MemoryStore) List(
	ctx context.Context,
	filter Filter,
) ([]*Event, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*Event, 0)

	for _, event := range s.events {

		if filter.Match(event) {
			events = append(
				events,
				event.Clone(),
			)
		}
	}

	return events, nil
}

func (s *MemoryStore) Delete(
	ctx context.Context,
	id string,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if id == "" {
		return ErrInvalidID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(
		s.events,
		id,
	)

	return nil
}

func (s *MemoryStore) DeleteOlderThan(
	ctx context.Context,
	timestamp time.Time,
) (int64, error) {

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var deleted int64

	for id, event := range s.events {

		if event.Timestamp.Before(timestamp) {

			delete(
				s.events,
				id,
			)

			deleted++
		}
	}

	return deleted, nil
}

func (s *MemoryStore) Count(
	ctx context.Context,
	filter Filter,
) (int64, error) {

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var total int64

	for _, event := range s.events {

		if filter.Match(event) {
			total++
		}
	}

	return total, nil
}

func (s *MemoryStore) Exists(
	id string,
) bool {

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.events[id]

	return exists
}

func (s *MemoryStore) Clear() {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(
		map[string]*Event,
	)
}

func (s *MemoryStore) Size() int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(
		s.events,
	)
}
