package secrets_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/secrets"
)

func TestSecretMiddlewareWithoutRequirements(
	t *testing.T,
) {

	handler := secrets.SecretMiddleware(
		secrets.MiddlewareConfig{},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(
		rr,
		req,
	)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			rr.Code,
		)
	}
}

func TestSecretMiddlewareRequireTenant(
	t *testing.T,
) {

	handler := secrets.SecretMiddleware(
		secrets.MiddlewareConfig{
			RequireTenant: true,
		},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(
		rr,
		req,
	)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			rr.Code,
		)
	}
}

func TestSecretMiddlewareRequireUser(
	t *testing.T,
) {

	handler := secrets.SecretMiddleware(
		secrets.MiddlewareConfig{
			RequireUser: true,
		},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {
				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(
		rr,
		req,
	)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			rr.Code,
		)
	}
}

func TestSecretMiddlewareInjectContext(
	t *testing.T,
) {

	handler := secrets.SecretMiddleware(
		secrets.MiddlewareConfig{
			RequireTenant: true,
			RequireUser:   true,
		},
	)(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				tenantID, ok := secrets.TenantIDFromContext(
					r.Context(),
				)

				if !ok {
					t.Fatal(
						"tenant missing from context",
					)
				}

				userID, ok := secrets.UserIDFromContext(
					r.Context(),
				)

				if !ok {
					t.Fatal(
						"user missing from context",
					)
				}

				if tenantID != "tenant-1" {
					t.Fatalf(
						"unexpected tenant %s",
						tenantID,
					)
				}

				if userID != "user-1" {
					t.Fatalf(
						"unexpected user %s",
						userID,
					)
				}

				w.WriteHeader(
					http.StatusOK,
				)
			},
		),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	req.Header.Set(
		secrets.TenantHeader,
		"tenant-1",
	)

	req.Header.Set(
		secrets.UserHeader,
		"user-1",
	)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(
		rr,
		req,
	)

	if rr.Code != http.StatusOK {
		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			rr.Code,
		)
	}
}

func TestTenantIDFromContext(
	t *testing.T,
) {

	ctx := context.WithValue(
		context.Background(),
		secrets.TenantContextKey,
		"tenant-123",
	)

	value, ok := secrets.TenantIDFromContext(
		ctx,
	)

	if !ok {
		t.Fatal(
			"tenant not found",
		)
	}

	if value != "tenant-123" {
		t.Fatalf(
			"unexpected value %s",
			value,
		)
	}
}

func TestUserIDFromContext(
	t *testing.T,
) {

	ctx := context.WithValue(
		context.Background(),
		secrets.UserContextKey,
		"user-123",
	)

	value, ok := secrets.UserIDFromContext(
		ctx,
	)

	if !ok {
		t.Fatal(
			"user not found",
		)
	}

	if value != "user-123" {
		t.Fatalf(
			"unexpected value %s",
			value,
		)
	}
}

func TestTenantIDFromContextNil(
	t *testing.T,
) {

	_, ok := secrets.TenantIDFromContext(
		nil,
	)

	if ok {
		t.Fatal(
			"expected false",
		)
	}
}

func TestUserIDFromContextNil(
	t *testing.T,
) {

	_, ok := secrets.UserIDFromContext(
		nil,
	)

	if ok {
		t.Fatal(
			"expected false",
		)
	}
}