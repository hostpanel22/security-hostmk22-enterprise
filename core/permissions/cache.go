package permissions

import (
	"sync"
	"time"
)

type CacheItem[T any] struct {
	Value T

	ExpiresAt time.Time
}

func (i *CacheItem[T]) Expired() bool {
	return time.Now().UTC().After(
		i.ExpiresAt,
	)
}

type Cache struct {
	mu sync.RWMutex

	ttl time.Duration

	roles map[string]*CacheItem[*Role]

	permissions map[string]*CacheItem[*Permission]
}

func NewCache(
	ttl time.Duration,
) *Cache {

	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &Cache{
		ttl: ttl,
		roles: make(
			map[string]*CacheItem[*Role],
		),
		permissions: make(
			map[string]*CacheItem[*Permission],
		),
	}
}

func (c *Cache) SetRole(
	role *Role,
) {

	if role == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.roles[normalize(role.Name)] = &CacheItem[*Role]{
		Value: role.Clone(),
		ExpiresAt: time.Now().
			UTC().
			Add(c.ttl),
	}
}

func (c *Cache) GetRole(
	name string,
) (*Role, bool) {

	c.mu.RLock()
	item, ok := c.roles[normalize(name)]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if item.Expired() {

		c.mu.Lock()
		delete(
			c.roles,
			normalize(name),
		)
		c.mu.Unlock()

		return nil, false
	}

	return item.Value.Clone(), true
}

func (c *Cache) DeleteRole(
	name string,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(
		c.roles,
		normalize(name),
	)
}

func (c *Cache) SetPermission(
	permission *Permission,
) {

	if permission == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.permissions[permission.Key()] = &CacheItem[*Permission]{
		Value: permission.Clone(),
		ExpiresAt: time.Now().
			UTC().
			Add(c.ttl),
	}
}

func (c *Cache) GetPermission(
	key string,
) (*Permission, bool) {

	c.mu.RLock()
	item, ok := c.permissions[normalize(key)]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if item.Expired() {

		c.mu.Lock()

		delete(
			c.permissions,
			normalize(key),
		)

		c.mu.Unlock()

		return nil, false
	}

	return item.Value.Clone(), true
}

func (c *Cache) DeletePermission(
	key string,
) {

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(
		c.permissions,
		normalize(key),
	)
}

func (c *Cache) Flush() {

	c.mu.Lock()
	defer c.mu.Unlock()

	c.roles = make(
		map[string]*CacheItem[*Role],
	)

	c.permissions = make(
		map[string]*CacheItem[*Permission],
	)
}

func (c *Cache) RoleCount() int {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(
		c.roles,
	)
}

func (c *Cache) PermissionCount() int {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(
		c.permissions,
	)
}
