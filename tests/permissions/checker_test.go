package permissions_test

import (
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/permissions"
)

func createSubject() *permissions.Subject {

	role := permissions.NewRole(
		"manager",
		"manager role",
	)

	_ = role.AddPermission(
		permissions.NewPermission(
			"Users Read",
			"users",
			"read",
			"",
		),
	)

	_ = role.AddPermission(
		permissions.NewPermission(
			"Users Write",
			"users",
			"write",
			"",
		),
	)

	return &permissions.Subject{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Roles: []*permissions.Role{
			role,
		},
	}
}

func TestNewChecker(t *testing.T) {

	checker := permissions.NewChecker()

	if checker == nil {
		t.Fatal("checker should not be nil")
	}
}

func TestCanSuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.Can(
		createSubject(),
		"users",
		"read",
	) {
		t.Fatal("permission should be granted")
	}
}

func TestCanDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.Can(
		createSubject(),
		"orders",
		"delete",
	) {
		t.Fatal("permission should be denied")
	}
}

func TestCanNilSubject(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.Can(
		nil,
		"users",
		"read",
	) {
		t.Fatal("nil subject should fail")
	}
}

func TestCanAnySuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.CanAny(
		createSubject(),
		"users",
		"delete",
		"write",
	) {
		t.Fatal("one permission should match")
	}
}

func TestCanAnyDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.CanAny(
		createSubject(),
		"users",
		"delete",
		"export",
	) {
		t.Fatal("no permission should match")
	}
}

func TestCanAllSuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.CanAll(
		createSubject(),
		"users",
		"read",
		"write",
	) {
		t.Fatal("all permissions should match")
	}
}

func TestCanAllDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.CanAll(
		createSubject(),
		"users",
		"read",
		"delete",
	) {
		t.Fatal("missing permission should fail")
	}
}

func TestHasRoleSuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.HasRole(
		createSubject(),
		"manager",
	) {
		t.Fatal("role should exist")
	}
}

func TestHasRoleDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.HasRole(
		createSubject(),
		"admin",
	) {
		t.Fatal("role should not exist")
	}
}

func TestHasPermissionSuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.HasPermission(
		createSubject(),
		"users:read",
	) {
		t.Fatal("permission should exist")
	}
}

func TestHasPermissionDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.HasPermission(
		createSubject(),
		"orders:delete",
	) {
		t.Fatal("permission should not exist")
	}
}

func TestRoleCount(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.RoleCount(
		createSubject(),
	) != 1 {

		t.Fatal("expected one role")
	}
}

func TestRoleCountNilSubject(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.RoleCount(nil) != 0 {
		t.Fatal("expected zero roles")
	}
}

func TestPermissionCount(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.PermissionCount(
		createSubject(),
	) != 2 {

		t.Fatal("expected two permissions")
	}
}

func TestPermissionCountNilSubject(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.PermissionCount(nil) != 0 {
		t.Fatal("expected zero permissions")
	}
}

func TestValidateTenantSuccess(t *testing.T) {

	checker := permissions.NewChecker()

	if !checker.ValidateTenant(
		createSubject(),
		"tenant-1",
	) {
		t.Fatal("tenant should match")
	}
}

func TestValidateTenantDenied(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.ValidateTenant(
		createSubject(),
		"tenant-2",
	) {
		t.Fatal("tenant should not match")
	}
}

func TestValidateTenantNilSubject(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.ValidateTenant(
		nil,
		"tenant-1",
	) {
		t.Fatal("nil subject should fail")
	}
}

func TestIsSuperAdminFalse(t *testing.T) {

	checker := permissions.NewChecker()

	if checker.IsSuperAdmin(
		createSubject(),
	) {
		t.Fatal("should not be super admin")
	}
}

func TestIsSuperAdminTrue(t *testing.T) {

	role := permissions.NewRole(
		"super-admin",
		"",
	)

	subject := &permissions.Subject{
		UserID:   "1",
		TenantID: "tenant-1",
		Roles: []*permissions.Role{
			role,
		},
	}

	checker := permissions.NewChecker()

	if !checker.IsSuperAdmin(
		subject,
	) {
		t.Fatal("should be super admin")
	}
}
