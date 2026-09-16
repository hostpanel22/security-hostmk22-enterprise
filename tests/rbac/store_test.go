package rbac_test

import (
	"context"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestNewMemoryStore(t *testing.T) {
	store := rbac.NewMemoryStore()

	if store == nil {
		t.Fatal("store is nil")
	}
}

func TestStoreSaveAndGetRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	if err := store.SaveRole(
		ctx,
		role,
	); err != nil {
		t.Fatal(err)
	}

	result, err := store.GetRole(
		ctx,
		"admin",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("role not found")
	}

	if result.Name != role.Name {
		t.Fatal("role mismatch")
	}
}

func TestStoreGetRoleNotFound(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	_, err := store.GetRole(
		ctx,
		"unknown",
	)

	if err != rbac.ErrRoleNotFound {
		t.Fatalf(
			"expected %v got %v",
			rbac.ErrRoleNotFound,
			err,
		)
	}
}

func TestStoreListRoles(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	_ = store.SaveRole(
		ctx,
		rbac.NewRole(
			"tenant",
			"admin",
			"administrator",
		),
	)

	_ = store.SaveRole(
		ctx,
		rbac.NewRole(
			"tenant",
			"user",
			"user role",
		),
	)

	roles, err := store.ListRoles(
		ctx,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 2 {
		t.Fatalf(
			"expected 2 roles got %d",
			len(roles),
		)
	}
}

func TestStoreDeleteRole(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	_ = store.SaveRole(
		ctx,
		role,
	)

	if err := store.DeleteRole(
		ctx,
		"admin",
	); err != nil {
		t.Fatal(err)
	}

	_, err := store.GetRole(
		ctx,
		"admin",
	)

	if err != rbac.ErrRoleNotFound {
		t.Fatal("role should be deleted")
	}
}

func TestStoreSaveAndGetPermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"read users",
	)

	if err := store.SavePermission(
		ctx,
		permission,
	); err != nil {
		t.Fatal(err)
	}

	result, err := store.GetPermission(
		ctx,
		"users.read",
	)

	if err != nil {
		t.Fatal(err)
	}

	if result == nil {
		t.Fatal("permission not found")
	}

	if result.Name != permission.Name {
		t.Fatal("permission mismatch")
	}
}

func TestStoreGetPermissionNotFound(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	_, err := store.GetPermission(
		ctx,
		"unknown",
	)

	if err != rbac.ErrPermissionNotFound {
		t.Fatalf(
			"expected %v got %v",
			rbac.ErrPermissionNotFound,
			err,
		)
	}
}

func TestStoreListPermissions(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	_ = store.SavePermission(
		ctx,
		rbac.NewPermission(
			"tenant",
			"users.read",
			"users",
			"read",
			"",
		),
	)

	_ = store.SavePermission(
		ctx,
		rbac.NewPermission(
			"tenant",
			"users.write",
			"users",
			"write",
			"",
		),
	)

	permissions, err := store.ListPermissions(
		ctx,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(permissions) != 2 {
		t.Fatalf(
			"expected 2 permissions got %d",
			len(permissions),
		)
	}
}

func TestStoreDeletePermission(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	permission := rbac.NewPermission(
		"tenant",
		"users.read",
		"users",
		"read",
		"",
	)

	_ = store.SavePermission(
		ctx,
		permission,
	)

	if err := store.DeletePermission(
		ctx,
		"users.read",
	); err != nil {
		t.Fatal(err)
	}

	_, err := store.GetPermission(
		ctx,
		"users.read",
	)

	if err != rbac.ErrPermissionNotFound {
		t.Fatal("permission should be deleted")
	}
}

func TestStoreClear(t *testing.T) {
	ctx := context.Background()

	store := rbac.NewMemoryStore()

	_ = store.SaveRole(
		ctx,
		rbac.NewRole(
			"tenant",
			"admin",
			"administrator",
		),
	)

	_ = store.SavePermission(
		ctx,
		rbac.NewPermission(
			"tenant",
			"users.read",
			"users",
			"read",
			"",
		),
	)

	if err := store.Clear(
		ctx,
	); err != nil {
		t.Fatal(err)
	}

	roles, _ := store.ListRoles(ctx)
	permissions, _ := store.ListPermissions(ctx)

	if len(roles) != 0 {
		t.Fatal("roles not cleared")
	}

	if len(permissions) != 0 {
		t.Fatal("permissions not cleared")
	}
}
