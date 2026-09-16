package rbac_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestNewManager(t *testing.T) {
	store := rbac.NewMemoryStore()

	manager, err := rbac.NewManager(store)

	if err != nil {
		t.Fatal(err)
	}

	if manager == nil {
		t.Fatal("manager is nil")
	}
}

func TestNewManagerNilStore(t *testing.T) {
	manager, err := rbac.NewManager(nil)

	if err != rbac.ErrStoreUnavailable {
		t.Fatalf("expected ErrStoreUnavailable, got %v", err)
	}

	if manager != nil {
		t.Fatal("manager should be nil")
	}
}

func TestManagerCreateRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	if err := manager.CreateRole(ctx, role); err != nil {
		t.Fatal(err)
	}
}

func TestManagerCreatePermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	permission := rbac.NewPermission(
		"tenant-1",
		"users.read",
		"users",
		"read",
		"read users",
	)

	if err := manager.CreatePermission(ctx, permission); err != nil {
		t.Fatal(err)
	}
}

func TestManagerGetRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	_ = manager.CreateRole(
		ctx,
		role,
	)

	result, err := manager.GetRole(
		ctx,
		"admin",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("role not found")
	}
}

func TestManagerGetPermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	_ = manager.CreatePermission(
		ctx,
		permission,
	)

	result, err := manager.GetPermission(
		ctx,
		"users.read",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("permission not found")
	}
}

func TestManagerDeleteRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	_ = manager.CreateRole(
		ctx,
		role,
	)

	if err := manager.DeleteRole(
		ctx,
		"admin",
	); err != nil {
		t.Fatal(err)
	}
}

func TestManagerDeletePermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	_ = manager.CreatePermission(
		ctx,
		permission,
	)

	if err := manager.DeletePermission(
		ctx,
		"users.read",
	); err != nil {
		t.Fatal(err)
	}
}

func TestManagerAssignPermissionToRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	_ = manager.CreateRole(
		ctx,
		role,
	)

	_ = manager.CreatePermission(
		ctx,
		permission,
	)

	err := manager.AssignPermissionToRole(
		ctx,
		"admin",
		"users.read",
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestManagerRoleHasPermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	_ = manager.CreateRole(
		ctx,
		role,
	)

	_ = manager.CreatePermission(
		ctx,
		permission,
	)

	_ = manager.AssignPermissionToRole(
		ctx,
		"admin",
		"users.read",
	)

	ok, err := manager.RoleHasPermission(
		ctx,
		"admin",
		"users.read",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("permission should exist")
	}
}

func TestManagerRemovePermissionFromRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	_ = manager.CreateRole(
		ctx,
		role,
	)

	_ = manager.CreatePermission(
		ctx,
		permission,
	)

	_ = manager.AssignPermissionToRole(
		ctx,
		"admin",
		"users.read",
	)

	if err := manager.RemovePermissionFromRole(
		ctx,
		"admin",
		"users.read",
	); err != nil {
		t.Fatal(err)
	}
}

func TestManagerAssignParentRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	admin := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	user := rbac.NewRole(
		"tenant",
		"user",
		"user role",
	)

	_ = manager.CreateRole(ctx, admin)
	_ = manager.CreateRole(ctx, user)

	if err := manager.AssignParentRole(
		ctx,
		"admin",
		"user",
	); err != nil {
		t.Fatal(err)
	}
}

func TestManagerRemoveParentRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	admin := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	user := rbac.NewRole(
		"tenant",
		"user",
		"user role",
	)

	_ = manager.CreateRole(ctx, admin)
	_ = manager.CreateRole(ctx, user)

	_ = manager.AssignParentRole(
		ctx,
		"admin",
		"user",
	)

	if err := manager.RemoveParentRole(
		ctx,
		"admin",
		"user",
	); err != nil {
		t.Fatal(err)
	}
}

func TestManagerListRoles(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	_ = manager.CreateRole(
		ctx,
		rbac.NewRole(
			"tenant",
			"admin",
			"administrator",
		),
	)

	roles, err := manager.ListRoles(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 1 {
		t.Fatal("expected one role")
	}
}

func TestManagerListPermissions(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	_ = manager.CreatePermission(
		ctx,
		rbac.NewPermission(
			"tenant",
			"users.read",
			"users",
			"read",
			"read users",
		),
	)

	permissions, err := manager.ListPermissions(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(permissions) != 1 {
		t.Fatal("expected one permission")
	}
}

func TestManagerClear(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(store)

	_ = manager.CreateRole(
		ctx,
		rbac.NewRole(
			"tenant",
			"admin",
			"administrator",
		),
	)

	if err := manager.Clear(ctx); err != nil {
		t.Fatal(err)
	}

	roles, _ := manager.ListRoles(ctx)

	if len(roles) != 0 {
		t.Fatal("store should be empty")
	}
}
