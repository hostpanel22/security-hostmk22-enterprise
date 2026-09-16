package rbac_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestRequirePermissionSuccess(t *testing.T) {
	store := rbac.NewMemoryStore()

	manager, err := rbac.NewManager(store)
	if err != nil {
		t.Fatal(err)
	}

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
		"",
	)

	if err := manager.CreateRole(nil, role); err != nil {
		t.Fatal(err)
	}

	if err := manager.CreatePermission(nil, permission); err != nil {
		t.Fatal(err)
	}

	if err := manager.AssignPermissionToRole(
		nil,
		"admin",
		"users.read",
	); err != nil {
		t.Fatal(err)
	}

	protected := rbac.RequirePermission(
		manager,
		"users.read",
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{
			RequireRole: true,
		},
	)(
		protected,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		rbac.RoleHeader,
		"admin",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
