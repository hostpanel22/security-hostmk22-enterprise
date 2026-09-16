package rbac_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestNewPermission(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant-1",
		"Users Create",
		"Users",
		"Create",
		"Create users",
	)

	if permission == nil {
		t.Fatal("permission should not be nil")
	}

	if permission.ID == "" {
		t.Fatal("permission id should not be empty")
	}

	if permission.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if permission.Name != "users create" {
		t.Fatal("name should be normalized")
	}

	if permission.Resource != "users" {
		t.Fatal("resource should be normalized")
	}

	if permission.Action != "create" {
		t.Fatal("action should be normalized")
	}

	if permission.CreatedAt.IsZero() {
		t.Fatal("createdAt should be set")
	}

	if permission.UpdatedAt.IsZero() {
		t.Fatal("updatedAt should be set")
	}
}

func TestPermissionValidate(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant-1",
		"users:create",
		"users",
		"create",
		"create users",
	)

	if err := permission.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPermissionValidateNil(t *testing.T) {
	var permission *rbac.Permission

	if err := permission.Validate(); err != rbac.ErrInvalidPermission {
		t.Fatal("expected ErrInvalidPermission")
	}
}

func TestPermissionValidateEmptyName(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"",
		"users",
		"create",
		"",
	)

	if err := permission.Validate(); err != rbac.ErrEmptyPermissionName {
		t.Fatal("expected ErrEmptyPermissionName")
	}
}

func TestPermissionValidateEmptyResource(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"",
		"create",
		"",
	)

	if err := permission.Validate(); err != rbac.ErrEmptyPermissionResource {
		t.Fatal("expected ErrEmptyPermissionResource")
	}
}

func TestPermissionValidateEmptyAction(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"",
		"",
	)

	if err := permission.Validate(); err != rbac.ErrEmptyPermissionAction {
		t.Fatal("expected ErrEmptyPermissionAction")
	}
}

func TestPermissionKey(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	if permission.Key() != "users:create" {
		t.Fatal("invalid permission key")
	}
}

func TestPermissionEquals(t *testing.T) {
	first := rbac.NewPermission(
		"tenant",
		"a",
		"users",
		"create",
		"",
	)

	second := rbac.NewPermission(
		"tenant",
		"b",
		"users",
		"create",
		"",
	)

	if !first.Equals(second) {
		t.Fatal("permissions should be equal")
	}
}

func TestPermissionClone(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	clone := permission.Clone()

	if clone == nil {
		t.Fatal("clone should not be nil")
	}

	if clone == permission {
		t.Fatal("clone should be different pointer")
	}

	if clone.ID != permission.ID {
		t.Fatal("clone id mismatch")
	}
}

func TestPermissionTouch(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	before := permission.UpdatedAt

	time.Sleep(time.Millisecond)

	permission.Touch()

	if !permission.UpdatedAt.After(before) {
		t.Fatal("updatedAt should change")
	}
}

func TestPermissionBelongsToTenant(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant-123",
		"users:create",
		"users",
		"create",
		"",
	)

	if !permission.BelongsToTenant("tenant-123") {
		t.Fatal("tenant should match")
	}

	if permission.BelongsToTenant("tenant-999") {
		t.Fatal("tenant should not match")
	}
}

func TestPermissionUpdateDescription(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"old",
	)

	permission.UpdateDescription("new")

	if permission.Description != "new" {
		t.Fatal("description not updated")
	}
}

func TestPermissionUpdateName(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"old",
		"users",
		"create",
		"",
	)

	permission.UpdateName("New Name")

	if permission.Name != "new name" {
		t.Fatal("name should be normalized")
	}
}

func TestPermissionUpdateResource(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	permission.UpdateResource("Billing")

	if permission.Resource != "billing" {
		t.Fatal("resource should be normalized")
	}
}

func TestPermissionUpdateAction(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	permission.UpdateAction("Delete")

	if permission.Action != "delete" {
		t.Fatal("action should be normalized")
	}
}

func TestPermissionMatches(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	if !permission.Matches("users", "create") {
		t.Fatal("permission should match")
	}

	if permission.Matches("billing", "create") {
		t.Fatal("permission should not match")
	}
}

func TestPermissionReadOnly(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:read",
		"users",
		"read",
		"",
	)

	if !permission.IsReadOnly() {
		t.Fatal("permission should be read-only")
	}
}

func TestPermissionWrite(t *testing.T) {
	actions := []string{
		"create",
		"update",
		"delete",
		"write",
	}

	for _, action := range actions {
		permission := rbac.NewPermission(
			"tenant",
			"test",
			"users",
			action,
			"",
		)

		if !permission.IsWrite() {
			t.Fatalf("%s should be write action", action)
		}
	}
}

func TestPermissionString(t *testing.T) {
	permission := rbac.NewPermission(
		"tenant",
		"users:create",
		"users",
		"create",
		"",
	)

	if permission.String() != "users:create" {
		t.Fatal("invalid string representation")
	}
}