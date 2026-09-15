package permissions_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/permissions"
)

func TestNewRole(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"Management role",
	)

	if role == nil {
		t.Fatal("role should not be nil")
	}

	if role.Name != "manager" {
		t.Fatalf(
			"unexpected role name: %s",
			role.Name,
		)
	}

	if role.Description != "Management role" {
		t.Fatalf(
			"unexpected description: %s",
			role.Description,
		)
	}

	if role.ID == "" {
		t.Fatal("role id should not be empty")
	}
}

func TestNewSystemRole(t *testing.T) {

	role := permissions.NewSystemRole(
		"SuperAdmin",
		"System role",
	)

	if !role.IsSystem() {
		t.Fatal("role should be system")
	}
}

func TestRoleValidateSuccess(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	if err := role.Validate(); err != nil {
		t.Fatalf(
			"unexpected validation error: %v",
			err,
		)
	}
}

func TestRoleValidateNil(t *testing.T) {

	var role *permissions.Role

	if err := role.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRoleAddPermission(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	if err := role.AddPermission(
		permission,
	); err != nil {

		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if role.PermissionCount() != 1 {
		t.Fatal("permission should be added")
	}
}

func TestRoleAddDuplicatePermission(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	if err := role.AddPermission(
		permission,
	); err != nil {

		t.Fatal(err)
	}

	if err := role.AddPermission(
		permission,
	); err == nil {

		t.Fatal("expected duplicate permission error")
	}
}

func TestRoleRemovePermission(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	_ = role.AddPermission(
		permission,
	)

	if !role.RemovePermission(
		"users",
		"read",
	) {
		t.Fatal("permission should be removed")
	}

	if role.PermissionCount() != 0 {
		t.Fatal("permission count should be zero")
	}
}

func TestRoleRemovePermissionNotFound(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	if role.RemovePermission(
		"users",
		"read",
	) {
		t.Fatal("permission should not exist")
	}
}

func TestRoleHasPermission(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	_ = role.AddPermission(
		permission,
	)

	if !role.HasPermission(
		"users",
		"read",
	) {
		t.Fatal("permission should exist")
	}
}

func TestRoleHasPermissionWildcard(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"All",
		"*",
		"*",
		"",
	)

	_ = role.AddPermission(
		permission,
	)

	if !role.HasPermission(
		"orders",
		"delete",
	) {
		t.Fatal("wildcard permission should match")
	}
}

func TestRolePermissionCount(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	for i := 0; i < 3; i++ {

		p := permissions.NewPermission(
			"Permission",
			"resource",
			string(rune('a'+i)),
			"",
		)

		_ = role.AddPermission(
			p,
		)
	}

	if role.PermissionCount() != 3 {
		t.Fatalf(
			"expected 3 permissions got %d",
			role.PermissionCount(),
		)
	}
}

func TestRoleClone(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	_ = role.AddPermission(
		permission,
	)

	clone := role.Clone()

	if clone == nil {
		t.Fatal("clone should not be nil")
	}

	if clone == role {
		t.Fatal("clone should be a new instance")
	}

	if clone.PermissionCount() !=
		role.PermissionCount() {

		t.Fatal("permission count mismatch")
	}

	if clone.Permissions[0] ==
		role.Permissions[0] {

		t.Fatal("permissions should be cloned")
	}
}

func TestRoleTouch(t *testing.T) {

	role := permissions.NewRole(
		"Manager",
		"",
	)

	updated := role.UpdatedAt

	time.Sleep(
		10 * time.Millisecond,
	)

	role.Touch()

	if !role.UpdatedAt.After(
		updated,
	) {
		t.Fatal("updated_at should change")
	}
}

func TestRoleIsSystem(t *testing.T) {

	role := permissions.NewSystemRole(
		"SuperAdmin",
		"",
	)

	if !role.IsSystem() {
		t.Fatal("expected system role")
	}
}
