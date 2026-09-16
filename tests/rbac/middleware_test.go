package rbac_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/rbac"
)

func TestRBACMiddlewareSuccess(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{
			RequireTenant: true,
			RequireUser:   true,
			RequireRole:   true,
		},
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

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		rbac.TenantHeader,
		"tenant-1",
	)

	req.Header.Set(
		rbac.UserHeader,
		"user-1",
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

func TestRBACMiddlewareMissingTenant(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{
			RequireTenant: true,
		},
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

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRBACMiddlewareMissingUser(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{
			RequireUser: true,
		},
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

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}
}

func TestRBACMiddlewareMissingRole(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{
			RequireRole: true,
		},
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

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected %d got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestTenantIDFromContext(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				tenantID, ok := rbac.TenantIDFromContext(
					r.Context(),
				)

				if !ok {
					t.Fatal("tenant not found")
				}

				if tenantID != "tenant-1" {
					t.Fatalf(
						"expected tenant-1 got %s",
						tenantID,
					)
				}

				w.WriteHeader(http.StatusOK)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		rbac.TenantHeader,
		"tenant-1",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)
}

func TestUserIDFromContext(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				userID, ok := rbac.UserIDFromContext(
					r.Context(),
				)

				if !ok {
					t.Fatal("user not found")
				}

				if userID != "user-1" {
					t.Fatalf(
						"expected user-1 got %s",
						userID,
					)
				}

				w.WriteHeader(http.StatusOK)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		rbac.UserHeader,
		"user-1",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)
}

func TestRoleFromContext(t *testing.T) {
	handler := rbac.RBACMiddleware(
		rbac.MiddlewareConfig{},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				role, ok := rbac.RoleFromContext(
					r.Context(),
				)

				if !ok {
					t.Fatal("role not found")
				}

				if role != "admin" {
					t.Fatalf(
						"expected admin got %s",
						role,
					)
				}

				w.WriteHeader(http.StatusOK)
			},
		),
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
}

func TestRequirePermissionSuccess(t *testing.T) {
	store := rbac.NewMemoryStore()

	manager, err := rbac.NewManager(
		store,
	)

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

func TestRequirePermissionDenied(t *testing.T) {
	store := rbac.NewMemoryStore()

	manager, _ := rbac.NewManager(
		store,
	)

	role := rbac.NewRole(
		"tenant",
		"admin",
		"administrator",
	)

	_ = manager.CreateRole(
		nil,
		role,
	)

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

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected %d got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestRequirePermissionNilManager(t *testing.T) {
	handler := rbac.RequirePermission(
		nil,
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

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(
		rec,
		req,
	)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected %d got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}
