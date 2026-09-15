package permissions_test

import (
	"testing"
	"time"

	permissions "github.com/azezaz/HostMK-Security-Framework/core/permissions"
)

func TestNewCache(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	if cache == nil {
		t.Fatal("expected cache instance")
	}
}

func TestNewCacheDefaultTTL(t *testing.T) {
	cache := permissions.NewCache(0)

	if cache == nil {
		t.Fatal("expected cache instance")
	}
}

func TestSetAndGetRole(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	role := permissions.NewRole(
		"admin",
		"administrator",
	)

	cache.SetRole(role)

	result, found := cache.GetRole("admin")

	if !found {
		t.Fatal("expected role in cache")
	}

	if result.Name != role.Name {
		t.Fatalf(
			"expected %s got %s",
			role.Name,
			result.Name,
		)
	}
}

func TestGetRoleNotFound(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	_, found := cache.GetRole("unknown")

	if found {
		t.Fatal("role should not exist")
	}
}

func TestDeleteRole(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	role := permissions.NewRole(
		"admin",
		"",
	)

	cache.SetRole(role)

	cache.DeleteRole("admin")

	_, found := cache.GetRole("admin")

	if found {
		t.Fatal("role should be deleted")
	}
}

func TestRoleExpiration(t *testing.T) {
	cache := permissions.NewCache(
		10 * time.Millisecond,
	)

	role := permissions.NewRole(
		"admin",
		"",
	)

	cache.SetRole(role)

	time.Sleep(
		20 * time.Millisecond,
	)

	_, found := cache.GetRole("admin")

	if found {
		t.Fatal("role should be expired")
	}
}

func TestCacheRoleCount(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	cache.SetRole(
		permissions.NewRole(
			"admin",
			"",
		),
	)

	cache.SetRole(
		permissions.NewRole(
			"manager",
			"",
		),
	)

	if cache.RoleCount() != 2 {
		t.Fatalf(
			"expected 2 roles got %d",
			cache.RoleCount(),
		)
	}
}

func TestRoleCloneIsolation(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	role := permissions.NewRole(
		"admin",
		"",
	)

	cache.SetRole(role)

	cached, found := cache.GetRole("admin")

	if !found {
		t.Fatal("expected cached role")
	}

	cached.Name = "modified"

	again, found := cache.GetRole("admin")

	if !found {
		t.Fatal("expected cached role")
	}

	if again.Name != "admin" {
		t.Fatal("cache returned mutable reference")
	}
}

func TestSetAndGetPermission(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	cache.SetPermission(permission)

	result, found := cache.GetPermission(
		permission.Key(),
	)

	if !found {
		t.Fatal("expected permission")
	}

	if result.Key() != permission.Key() {
		t.Fatalf(
			"expected %s got %s",
			permission.Key(),
			result.Key(),
		)
	}
}

func TestGetPermissionNotFound(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	_, found := cache.GetPermission(
		"users:read",
	)

	if found {
		t.Fatal("permission should not exist")
	}
}

func TestDeletePermission(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	cache.SetPermission(permission)

	cache.DeletePermission(
		permission.Key(),
	)

	_, found := cache.GetPermission(
		permission.Key(),
	)

	if found {
		t.Fatal("permission should be deleted")
	}
}

func TestPermissionExpiration(t *testing.T) {
	cache := permissions.NewCache(
		10 * time.Millisecond,
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	cache.SetPermission(permission)

	time.Sleep(
		20 * time.Millisecond,
	)

	_, found := cache.GetPermission(
		permission.Key(),
	)

	if found {
		t.Fatal("permission should be expired")
	}
}

func TestCachePermissionCount(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	cache.SetPermission(
		permissions.NewPermission(
			"Users Read",
			"users",
			"read",
			"",
		),
	)

	cache.SetPermission(
		permissions.NewPermission(
			"Users Write",
			"users",
			"write",
			"",
		),
	)

	if cache.PermissionCount() != 2 {
		t.Fatalf(
			"expected 2 permissions got %d",
			cache.PermissionCount(),
		)
	}
}

func TestPermissionCloneIsolation(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	cache.SetPermission(permission)

	cached, found := cache.GetPermission(
		permission.Key(),
	)

	if !found {
		t.Fatal("expected permission")
	}

	cached.Name = "modified"

	again, found := cache.GetPermission(
		permission.Key(),
	)

	if !found {
		t.Fatal("expected permission")
	}

	if again.Name != "Users Read" {
		t.Fatal("cache returned mutable permission reference")
	}
}

func TestFlush(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	cache.SetRole(
		permissions.NewRole(
			"admin",
			"",
		),
	)

	cache.SetPermission(
		permissions.NewPermission(
			"Users Read",
			"users",
			"read",
			"",
		),
	)

	cache.Flush()

	if cache.RoleCount() != 0 {
		t.Fatal("roles should be flushed")
	}

	if cache.PermissionCount() != 0 {
		t.Fatal("permissions should be flushed")
	}
}

func TestNilRoleIgnored(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	cache.SetRole(nil)

	if cache.RoleCount() != 0 {
		t.Fatal("nil role should be ignored")
	}
}

func TestNilPermissionIgnored(t *testing.T) {
	cache := permissions.NewCache(time.Minute)

	cache.SetPermission(nil)

	if cache.PermissionCount() != 0 {
		t.Fatal("nil permission should be ignored")
	}
}
