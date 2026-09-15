package permissions

import (
	"context"
	"sync"
)

type Store interface {
	SaveRole(
		ctx context.Context,
		role *Role,
	) error

	GetRole(
		ctx context.Context,
		name string,
	) (*Role, error)

	DeleteRole(
		ctx context.Context,
		name string,
	) error

	ListRoles(
		ctx context.Context,
	) ([]*Role, error)

	SavePermission(
		ctx context.Context,
		permission *Permission,
	) error

	GetPermission(
		ctx context.Context,
		key string,
	) (*Permission, error)

	DeletePermission(
		ctx context.Context,
		key string,
	) error

	ListPermissions(
		ctx context.Context,
	) ([]*Permission, error)
}

type MemoryStore struct {
	mu sync.RWMutex

	roles map[string]*Role

	permissions map[string]*Permission
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		roles: make(
			map[string]*Role,
		),
		permissions: make(
			map[string]*Permission,
		),
	}
}

func (s *MemoryStore) SaveRole(
	ctx context.Context,
	role *Role,
) error {

	if role == nil {
		return ErrInvalidRole
	}

	if err := role.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.roles[normalize(role.Name)] = role.Clone()

	return nil
}

func (s *MemoryStore) GetRole(
	ctx context.Context,
	name string,
) (*Role, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	role, ok := s.roles[normalize(name)]

	if !ok {
		return nil, ErrRoleNotFound
	}

	return role.Clone(), nil
}

func (s *MemoryStore) DeleteRole(
	ctx context.Context,
	name string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.roles[normalize(name)]; !ok {
		return ErrRoleNotFound
	}

	delete(
		s.roles,
		normalize(name),
	)

	return nil
}

func (s *MemoryStore) ListRoles(
	ctx context.Context,
) ([]*Role, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Role,
		0,
		len(s.roles),
	)

	for _, role := range s.roles {
		result = append(
			result,
			role.Clone(),
		)
	}

	return result, nil
}

func (s *MemoryStore) SavePermission(
	ctx context.Context,
	permission *Permission,
) error {

	if permission == nil {
		return ErrInvalidPermission
	}

	if err := permission.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.permissions[permission.Key()] = permission.Clone()

	return nil
}

func (s *MemoryStore) GetPermission(
	ctx context.Context,
	key string,
) (*Permission, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	permission, ok := s.permissions[normalize(key)]

	if !ok {
		return nil, ErrPermissionNotFound
	}

	return permission.Clone(), nil
}

func (s *MemoryStore) DeletePermission(
	ctx context.Context,
	key string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.permissions[normalize(key)]; !ok {
		return ErrPermissionNotFound
	}

	delete(
		s.permissions,
		normalize(key),
	)

	return nil
}

func (s *MemoryStore) ListPermissions(
	ctx context.Context,
) ([]*Permission, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Permission,
		0,
		len(s.permissions),
	)

	for _, permission := range s.permissions {
		result = append(
			result,
			permission.Clone(),
		)
	}

	return result, nil
}
