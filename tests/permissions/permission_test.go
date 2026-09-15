package permissions_test

import (
	"testing"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/core/permissions"
)

func TestNewPermission(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"Read users",
	)

	if permission == nil {
		t.Fatal("permission should not be nil")
	}

	if permission.Name != "Users Read" {
		t.Fatalf(
			"unexpected name: %s",
			permission.Name,
		)
	}

	if permission.Resource != "users" {
		t.Fatalf(
			"unexpected resource: %s",
			permission.Resource,
		)
	}

	if permission.Action != "read" {
		t.Fatalf(
			"unexpected action: %s",
			permission.Action,
		)
	}
}

func TestNewSystemPermission(t *testing.T) {

	permission := permissions.NewSystemPermission(
		"Super Admin",
		"*",
		"*",
		"System permission",
	)

	if !permission.IsSystem() {
		t.Fatal("permission should be system")
	}
}

func TestPermissionValidateSuccess(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	if err := permission.Validate(); err != nil {
		t.Fatalf(
			"unexpected validation error: %v",
			err,
		)
	}
}

func TestPermissionKey(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	if permission.Key() != "users:read" {
		t.Fatalf(
			"unexpected key: %s",
			permission.Key(),
		)
	}
}

func TestPermissionMatchesExact(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	if !permission.Matches(
		"users",
		"read",
	) {
		t.Fatal("expected permission match")
	}
}

func TestPermissionMatchesWildcardAll(t *testing.T) {

	permission := permissions.NewPermission(
		"All",
		"*",
		"*",
		"",
	)

	if !permission.Matches(
		"orders",
		"delete",
	) {
		t.Fatal("wildcard should match")
	}
}

func TestPermissionMatchesWildcardResource(t *testing.T) {

	permission := permissions.NewPermission(
		"Read Any",
		"*",
		"read",
		"",
	)

	if !permission.Matches(
		"users",
		"read",
	) {
		t.Fatal("resource wildcard should match")
	}
}

func TestPermissionMatchesWildcardAction(t *testing.T) {

	permission := permissions.NewPermission(
		"Users All",
		"users",
		"*",
		"",
	)

	if !permission.Matches(
		"users",
		"delete",
	) {
		t.Fatal("action wildcard should match")
	}
}

func TestPermissionClone(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	clone := permission.Clone()

	if clone == nil {
		t.Fatal("clone should not be nil")
	}

	if clone == permission {
		t.Fatal("clone should be different pointer")
	}

	if clone.Key() != permission.Key() {
		t.Fatal("clone key mismatch")
	}
}

func TestPermissionTouch(t *testing.T) {

	permission := permissions.NewPermission(
		"Users Read",
		"users",
		"read",
		"",
	)

	original := permission.UpdatedAt

	time.Sleep(
		10 * time.Millisecond,
	)

	permission.Touch()

	if !permission.UpdatedAt.After(
		original,
	) {
		t.Fatal("updated_at should change")
	}
}

func TestPermissionIsSystem(t *testing.T) {

	permission := permissions.NewSystemPermission(
		"System",
		"*",
		"*",
		"",
	)

	if !permission.IsSystem() {
		t.Fatal("expected system permission")
	}
}
