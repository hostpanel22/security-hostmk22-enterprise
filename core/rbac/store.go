package rbac

import (
	"context"
	"strings"
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

	ListRoles(
		ctx context.Context,
	) ([]*Role, error)

	DeleteRole(
		ctx context.Context,
		name string,
	) error

	SavePermission(
		ctx context.Context,
		permission *Permission,
	) error

	GetPermission(
		ctx context.Context,
		name string,
	) (*Permission, error)

	ListPermissions(
		ctx context.Context,
	) ([]*Permission, error)

	DeletePermission(
		ctx context.Context,
		name string,
	) error

	Clear(
		ctx context.Context,
	) error
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
	_ context.Context,
	role *Role,
) error {

	if role == nil {
		return ErrInvalidRole
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.roles[
		normalize(role.Name),
	] = role.Clone()

	return nil
}

func (s *MemoryStore) GetRole(
	_ context.Context,
	name string,
) (*Role, error) {

	name = normalize(name)

	if name == "" {
		return nil, ErrRoleNotFound
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	role, exists := s.roles[name]

	if !exists {
		return nil, ErrRoleNotFound
	}

	return role.Clone(), nil
}

func (s *MemoryStore) ListRoles(
	_ context.Context,
) ([]*Role, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Role,
		0,
		len(s.roles),
	)

	for _, role := range s.roles {

		if role == nil {
			continue
		}

		result = append(
			result,
			role.Clone(),
		)
	}

	return result, nil
}

func (s *MemoryStore) DeleteRole(
	_ context.Context,
	name string,
) error {

	name = normalize(name)

	if name == "" {
		return ErrRoleNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.roles[name]; !exists {
		return ErrRoleNotFound
	}

	delete(
		s.roles,
		name,
	)

	return nil
}

func (s *MemoryStore) SavePermission(
	_ context.Context,
	permission *Permission,
) error {

	if permission == nil {
		return ErrInvalidPermission
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.permissions[
		normalize(permission.Name),
	] = permission.Clone()

	return nil
}

func (s *MemoryStore) GetPermission(
	_ context.Context,
	name string,
) (*Permission, error) {

	name = normalize(name)

	if name == "" {
		return nil, ErrPermissionNotFound
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	permission, exists := s.permissions[name]

	if !exists {
		return nil, ErrPermissionNotFound
	}

	return permission.Clone(), nil
}

func (s *MemoryStore) ListPermissions(
	_ context.Context,
) ([]*Permission, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(
		[]*Permission,
		0,
		len(s.permissions),
	)

	for _, permission := range s.permissions {

		if permission == nil {
			continue
		}

		result = append(
			result,
			permission.Clone(),
		)
	}

	return result, nil
}

func (s *MemoryStore) DeletePermission(
	_ context.Context,
	name string,
) error {

	name = normalize(name)

	if name == "" {
		return ErrPermissionNotFound
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.permissions[name]; !exists {
		return ErrPermissionNotFound
	}

	delete(
		s.permissions,
		name,
	)

	return nil
}

func (s *MemoryStore) Clear(
	_ context.Context,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.roles = make(
		map[string]*Role,
	)

	s.permissions = make(
		map[string]*Permission,
	)

	return nil
}

func normalize(
	value string,
) string {

	return strings.ToLower(
		strings.TrimSpace(
			value,
		),
	)
}