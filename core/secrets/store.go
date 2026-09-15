package secrets

import (
	"context"
	"strings"
	"sync"
)

type Store interface {
	Save(
		ctx context.Context,
		secret *Secret,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (*Secret, error)

	GetByName(
		ctx context.Context,
		name string,
	) (*Secret, error)

	List(
		ctx context.Context,
	) ([]*Secret, error)

	Delete(
		ctx context.Context,
		id string,
	) error

	Clear(
		ctx context.Context,
	) error

	Count() int
}

type MemoryStore struct {
	mu sync.RWMutex

	secrets map[string]*Secret
}

func NewMemoryStore() *MemoryStore {

	return &MemoryStore{
		secrets: make(
			map[string]*Secret,
		),
	}
}

func (s *MemoryStore) Save(
	_ context.Context,
	secret *Secret,
) error {

	if secret == nil {
		return ErrInvalidSecret
	}

	if err := secret.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.secrets[secret.ID] =
		secret.Clone()

	return nil
}

func (s *MemoryStore) GetByID(
	_ context.Context,
	id string,
) (*Secret, error) {

	id = strings.TrimSpace(id)

	if id == "" {
		return nil, ErrSecretNotFound
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	secret, exists := s.secrets[id]

	if !exists {
		return nil, ErrSecretNotFound
	}

	return secret.Clone(), nil
}

func (s *MemoryStore) GetByName(
	_ context.Context,
	name string,
) (*Secret, error) {

	name = normalize(name)

	if name == "" {
		return nil, ErrSecretNotFound
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, secret := range s.secrets {

		if secret == nil {
			continue
		}

		if normalize(
			secret.Name,
		) == name {

			return secret.Clone(),
				nil
		}
	}

	return nil,
		ErrSecretNotFound
}

func (s *MemoryStore) List(
	_ context.Context,
) ([]*Secret, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Secret,
		0,
		len(s.secrets),
	)

	for _, secret := range s.secrets {

		if secret == nil {
			continue
		}

		result = append(
			result,
			secret.Clone(),
		)
	}

	return result, nil
}

func (s *MemoryStore) Delete(
	_ context.Context,
	id string,
) error {

	id = strings.TrimSpace(id)

	if id == "" {
		return ErrSecretNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.secrets[id]; !exists {
		return ErrSecretNotFound
	}

	delete(
		s.secrets,
		id,
	)

	return nil
}

func (s *MemoryStore) Clear(
	_ context.Context,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.secrets = make(
		map[string]*Secret,
	)

	return nil
}

func (s *MemoryStore) Count() int {

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(
		s.secrets,
	)
}
