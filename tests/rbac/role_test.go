package rbac_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestNewRole(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator role",
	)

	if role == nil {
		t.Fatal("expected role")
	}

	if role.ID == "" {
		t.Fatal("expected id")
	}

	if role.Name != "admin" {
		t.Fatalf(
			"expected admin got %s",
			role.Name,
		)
	}

	if role.PermissionCount() != 0 {
		t.Fatal("expected empty permissions")
	}
}

func TestRoleValidate(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	if err := role.Validate(); err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestRoleValidateNil(t *testing.T) {

	var role *rbac.Role

	if err := role.Validate(); err != rbac.ErrInvalidRole {
		t.Fatal("expected ErrInvalidRole")
	}
}

func TestRoleValidateEmptyID(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	role.ID = ""

	if err := role.Validate(); err != rbac.ErrEmptyRoleID {
		t.Fatal("expected ErrEmptyRoleID")
	}
}

func TestRoleValidateEmptyName(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"",
		"administrator",
	)

	if err := role.Validate(); err != rbac.ErrEmptyRoleName {
		t.Fatal("expected ErrEmptyRoleName")
	}
}

func TestRoleClone(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddPermission(
		"users:create",
	)

	clone := role.Clone()

	if clone == nil {
		t.Fatal("expected clone")
	}

	if clone == role {
		t.Fatal("clone should not be same pointer")
	}

	if clone.PermissionCount() != 1 {
		t.Fatal("expected cloned permissions")
	}
}

func TestRoleTouch(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	oldTime := role.UpdatedAt

	time.Sleep(
		10 * time.Millisecond,
	)

	role.Touch()

	if !role.UpdatedAt.After(
		oldTime,
	) {
		t.Fatal("updated time not changed")
	}
}

func TestRoleBelongsToTenant(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	if !role.BelongsToTenant(
		"tenant-1",
	) {
		t.Fatal("expected true")
	}

	if role.BelongsToTenant(
		"tenant-2",
	) {
		t.Fatal("expected false")
	}
}

func TestRoleUpdateDescription(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	role.UpdateDescription(
		"new description",
	)

	if role.Description != "new description" {
		t.Fatal("description not updated")
	}
}

func TestRoleAddPermission(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	err := role.AddPermission(
		"users:create",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error %v",
			err,
		)
	}

	if !role.HasPermission(
		"users:create",
	) {
		t.Fatal("permission missing")
	}
}

func TestRoleAddDuplicatePermission(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddPermission(
		"users:create",
	)

	err := role.AddPermission(
		"users:create",
	)

	if err != rbac.ErrPermissionAlreadyAssigned {
		t.Fatal("expected duplicate error")
	}
}

func TestRoleRemovePermission(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddPermission(
		"users:create",
	)

	err := role.RemovePermission(
		"users:create",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error %v",
			err,
		)
	}

	if role.HasPermission(
		"users:create",
	) {
		t.Fatal("permission still exists")
	}
}

func TestRoleParentRole(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	err := role.AddParentRole(
		"manager",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error %v",
			err,
		)
	}

	if !role.HasParentRole(
		"manager",
	) {
		t.Fatal("parent role missing")
	}
}

func TestRoleParentLoop(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	err := role.AddParentRole(
		"admin",
	)

	if err != rbac.ErrRoleInheritanceLoop {
		t.Fatal("expected inheritance loop")
	}
}

func TestRoleRemoveParentRole(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddParentRole(
		"manager",
	)

	err := role.RemoveParentRole(
		"manager",
	)

	if err != nil {
		t.Fatalf(
			"unexpected error %v",
			err,
		)
	}

	if role.HasParentRole(
		"manager",
	) {
		t.Fatal("parent role still exists")
	}
}

func TestRoleClearPermissions(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddPermission(
		"users:create",
	)

	role.ClearPermissions()

	if role.PermissionCount() != 0 {
		t.Fatal("permissions not cleared")
	}
}

func TestRoleClearParentRoles(t *testing.T) {

	role := rbac.NewRole(
		"tenant-1",
		"admin",
		"administrator",
	)

	_ = role.AddParentRole(
		"manager",
	)

	role.ClearParentRoles()

	if role.ParentRoleCount() != 0 {
		t.Fatal("parent roles not cleared")
	}
}