package permissions_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	permissions "github.com/azezaz/HostMK-Security-Framework/core/permissions"
)

func newSubjectWithPermission(
	resource string,
	action string,
) *permissions.Subject {

	role := permissions.NewRole(
		"admin",
		"administrator",
	)

	_ = role.AddPermission(
		permissions.NewPermission(
			"test",
			resource,
			action,
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

func newRequestWithSubject(
	subject *permissions.Subject,
) *http.Request {

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	ctx := permissions.WithSubject(
		req.Context(),
		subject,
	)

	return req.WithContext(ctx)
}

func okHandler() http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			w.WriteHeader(http.StatusOK)
		},
	)
}

func TestWithSubject(t *testing.T) {

	subject := &permissions.Subject{
		UserID: "user-1",
	}

	ctx := permissions.WithSubject(
		context.Background(),
		subject,
	)

	value, ok := permissions.SubjectFromContext(
		ctx,
	)

	if !ok {
		t.Fatal("expected subject")
	}

	if value.UserID != "user-1" {
		t.Fatal("unexpected user id")
	}
}

func TestSubjectFromContextMissing(t *testing.T) {

	_, ok := permissions.SubjectFromContext(
		context.Background(),
	)

	if ok {
		t.Fatal("subject should not exist")
	}
}

func TestRequirePermissionSuccess(t *testing.T) {

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	req := newRequestWithSubject(
		newSubjectWithPermission(
			"users",
			"read",
		),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequirePermission(
			"users",
			"read",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRequirePermissionForbidden(t *testing.T) {

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	req := newRequestWithSubject(
		newSubjectWithPermission(
			"users",
			"read",
		),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequirePermission(
			"users",
			"delete",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403 got %d",
			rec.Code,
		)
	}
}

func TestRequirePermissionUnauthorized(t *testing.T) {

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequirePermission(
			"users",
			"read",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401 got %d",
			rec.Code,
		)
	}
}

func TestRequireAnyPermission(t *testing.T) {

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	req := newRequestWithSubject(
		newSubjectWithPermission(
			"users",
			"read",
		),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireAnyPermission(
			"users",
			"write",
			"read",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRequireAllPermissionsForbidden(t *testing.T) {

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	req := newRequestWithSubject(
		newSubjectWithPermission(
			"users",
			"read",
		),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireAllPermissions(
			"users",
			"read",
			"write",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403 got %d",
			rec.Code,
		)
	}
}

func TestRequireRoleSuccess(t *testing.T) {

	role := permissions.NewRole(
		"admin",
		"",
	)

	subject := &permissions.Subject{
		UserID:   "user",
		TenantID: "tenant-1",
		Roles: []*permissions.Role{
			role,
		},
	}

	req := newRequestWithSubject(subject)

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireRole(
			"admin",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRequireTenantSuccess(t *testing.T) {

	subject := &permissions.Subject{
		UserID:   "user",
		TenantID: "tenant-1",
	}

	req := newRequestWithSubject(subject)

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireTenant(
			"tenant-1",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRequireTenantForbidden(t *testing.T) {

	subject := &permissions.Subject{
		UserID:   "user",
		TenantID: "tenant-a",
	}

	req := newRequestWithSubject(subject)

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireTenant(
			"tenant-b",
		)(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403 got %d",
			rec.Code,
		)
	}
}

func TestRequireSuperAdminSuccess(t *testing.T) {

	role := permissions.NewRole(
		"super-admin",
		"",
	)

	subject := &permissions.Subject{
		UserID: "admin",
		Roles: []*permissions.Role{
			role,
		},
	}

	req := newRequestWithSubject(subject)

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireSuperAdmin()(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200 got %d",
			rec.Code,
		)
	}
}

func TestRequireSuperAdminForbidden(t *testing.T) {

	role := permissions.NewRole(
		"user",
		"",
	)

	subject := &permissions.Subject{
		UserID: "user",
		Roles: []*permissions.Role{
			role,
		},
	}

	req := newRequestWithSubject(subject)

	middleware := permissions.NewMiddleware(
		permissions.NewChecker(),
	)

	rec := httptest.NewRecorder()

	handler := middleware.
		RequireSuperAdmin()(okHandler())

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected 403 got %d",
			rec.Code,
		)
	}
}
