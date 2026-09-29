package hostmk_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	hostmk "github.com/azezaz/HostMK-Security-Framework/sdk/go"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/apikeys"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/audit"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/auth"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/fraud"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/permissions"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/rbac"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/secrets"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/settings"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/tenancy"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/users"
)

func TestServicesUseRegisteredRoutesAndDecodeResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/health" && r.URL.Path != "/ready" && r.URL.Path != "/api/v1/version" && r.URL.Path != "/api/v1/auth/bootstrap" && r.URL.Path != "/api/v1/auth/login" && r.URL.Path != "/api/v1/auth/refresh" {
			if r.Header.Get("Authorization") != "Bearer access-2" {
				t.Errorf("%s %s Authorization = %q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
			}
		}

		data := any(map[string]any{})
		status := http.StatusOK
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/health":
			data = map[string]any{"status": "healthy", "components": []any{}}
		case r.Method == http.MethodGet && r.URL.Path == "/ready":
			data = map[string]any{"status": "healthy", "components": []any{}}
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/version":
			data = map[string]any{"name": "HostMK", "version": "test", "env": "test"}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/bootstrap":
			status = http.StatusCreated
			data = map[string]any{"tenant": map[string]any{"id": "t1", "name": "Tenant", "slug": "tenant", "status": "active"}, "user": map[string]any{"id": "u1", "tenant_id": "t1", "email": "owner@example.com", "name": "Owner", "status": "active"}}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/login":
			data = map[string]any{"access_token": "access-1", "refresh_token": "refresh-1", "token_type": "Bearer", "expires_in": 900}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/auth/refresh":
			data = map[string]any{"access_token": "access-2", "refresh_token": "refresh-2", "token_type": "Bearer", "expires_in": 900}
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/users":
			data = map[string]any{"count": 1, "users": []any{map[string]any{"id": "u1", "tenant_id": "t1", "email": "owner@example.com", "name": "Owner", "status": "active"}}}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/users":
			status = http.StatusCreated
			data = map[string]any{"id": "u2", "tenant_id": "t1", "email": "new@example.com", "name": "New", "status": "active"}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/users/action":
			if r.URL.Query().Get("action") == "delete" {
				data = map[string]string{"id": r.URL.Query().Get("id")}
			} else {
				data = map[string]any{"id": "u1", "tenant_id": "t1", "email": "owner@example.com", "name": "Owner", "status": "active"}
			}
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/tenants":
			data = map[string]any{"count": 1, "tenants": []any{map[string]any{"id": "t1", "name": "Tenant", "slug": "tenant", "status": "active"}}}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/tenants":
			status = http.StatusCreated
			data = map[string]any{"id": "t2", "name": "Tenant 2", "slug": "tenant-2", "status": "pending"}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/tenants/action":
			data = map[string]any{"id": "t1", "name": "Tenant", "slug": "tenant", "status": "suspended"}
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions":
			data = map[string]any{"count": 1, "sessions": []any{map[string]any{"id": "s1", "tenant_id": "t1", "user_id": "u1"}}}
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions":
			if r.URL.Query().Get("action") != "revoke" {
				t.Errorf("session action = %q", r.URL.Query().Get("action"))
			}
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/rbac/roles" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "roles": []any{map[string]any{"id": "r1", "tenant_id": "t1", "name": "operator", "permissions": []string{"audit.read"}}}}
		case r.URL.Path == "/api/v1/rbac/roles" && (r.Method == http.MethodPost || r.Method == http.MethodPut):
			if r.Method == http.MethodPost {
				status = http.StatusCreated
			}
			data = map[string]any{"id": "r1", "tenant_id": "t1", "name": "operator", "permissions": []string{"audit.read"}}
		case r.URL.Path == "/api/v1/rbac/roles" && r.Method == http.MethodDelete:
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/permissions" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "permissions": []any{map[string]any{"id": "p1", "tenant_id": "t1", "name": "audit.read", "resource": "audit", "action": "read"}}}
		case r.URL.Path == "/api/v1/permissions" && r.Method == http.MethodPost:
			status = http.StatusCreated
			data = map[string]any{"id": "p1", "tenant_id": "t1", "name": "audit.read", "resource": "audit", "action": "read"}
		case r.URL.Path == "/api/v1/permissions" && r.Method == http.MethodDelete:
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/permissions/check" && r.Method == http.MethodPost:
			data = map[string]any{"user_id": "u1", "permission": "audit.read", "allowed": true}
		case r.URL.Path == "/api/v1/api-keys" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "api_keys": []any{map[string]any{"id": "k1", "tenant_id": "t1", "name": "key", "prefix": "hmk_ab", "scopes": []string{"audit.read"}}}}
		case r.URL.Path == "/api/v1/api-keys" && r.Method == http.MethodPost:
			if r.URL.Query().Get("action") == "rotate" {
				data = map[string]any{"id": "k1", "prefix": "hmk_cd", "key": "hmk_rotated"}
			} else {
				status = http.StatusCreated
				data = map[string]any{"id": "k1", "tenant_id": "t1", "name": "key", "prefix": "hmk_ab", "key": "hmk_created", "scopes": []string{"audit.read"}}
			}
		case r.URL.Path == "/api/v1/api-keys" && r.Method == http.MethodDelete:
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/secrets" && r.Method == http.MethodGet && r.URL.Query().Get("history") == "true":
			data = []any{map[string]any{"id": "v1", "secret_id": "secret-1", "masked_value": "pass****", "created_by": "u1"}}
		case r.URL.Path == "/api/v1/secrets" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "secrets": []any{map[string]any{"id": "secret-1", "tenant_id": "t1", "name": "DB", "type": "custom", "masked_value": "pass****"}}}
		case r.URL.Path == "/api/v1/secrets" && r.Method == http.MethodPost:
			if r.URL.Query().Get("action") == "rotate" {
				if r.URL.Query().Get("value") == "" {
					t.Error("secret rotation value query is empty")
				}
				data = map[string]any{"id": "secret-1", "masked_value": "next****"}
			} else {
				status = http.StatusCreated
				data = map[string]any{"id": "secret-1", "tenant_id": "t1", "name": "DB", "type": "custom", "masked_value": "pass****"}
			}
		case r.URL.Path == "/api/v1/secrets" && r.Method == http.MethodDelete:
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/audit/events" && r.Method == http.MethodGet:
			if r.URL.Query().Get("limit") != "100" {
				t.Errorf("audit limit = %q", r.URL.Query().Get("limit"))
			}
			data = map[string]any{"count": 1, "events": []any{map[string]any{"id": "e1", "action": "login", "resource": "auth"}}}
		case r.URL.Path == "/api/v1/audit/events" && r.Method == http.MethodPost:
			status = http.StatusCreated
			data = map[string]any{"id": "e2", "action": "policy.changed", "resource": "roles", "severity": "info", "status": "success"}
		case r.URL.Path == "/api/v1/fraud/rules" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "rules": []any{map[string]any{"id": "f1", "name": "velocity", "enabled": true}}}
		case r.URL.Path == "/api/v1/fraud/rules" && (r.Method == http.MethodPost || r.Method == http.MethodPut):
			if r.Method == http.MethodPost {
				status = http.StatusCreated
			}
			data = map[string]any{"id": "f1", "name": "velocity", "enabled": true, "field": "count", "operator": "gt", "threshold": "10", "action": "block", "score": 50}
		case r.URL.Path == "/api/v1/fraud/rules" && r.Method == http.MethodDelete:
			data = map[string]string{"id": r.URL.Query().Get("id")}
		case r.URL.Path == "/api/v1/fraud/analyze" && r.Method == http.MethodPost:
			data = map[string]any{"risk_score": 50, "risk_level": "high", "blocked": false, "matched": []string{"velocity"}, "analyzed_at": "2026-09-29T00:00:00Z"}
		case r.URL.Path == "/api/v1/fraud/violations" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "events": []any{map[string]any{"id": "v1", "tenant_id": "t1", "rule_id": "f1", "score": 50, "level": "high"}}}
		case r.URL.Path == "/api/v1/settings" && r.Method == http.MethodGet:
			data = map[string]any{"count": 1, "settings": []any{map[string]any{"id": "setting-1", "tenant_id": "t1", "key": "feature.enabled", "value": true}}}
		case r.URL.Path == "/api/v1/settings" && (r.Method == http.MethodPut || r.Method == http.MethodPost):
			data = map[string]any{"id": "setting-1", "tenant_id": "t1", "key": "feature.enabled", "value": true}
		default:
			t.Errorf("unregistered or unexpected SDK request: %s %s", r.Method, r.URL.String())
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
	}))
	defer server.Close()

	sdk, err := hostmk.NewClient(server.URL, client.WithBearerToken("access-2"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	if result, err := sdk.Monitoring.Health(ctx); err != nil || result.Status != "healthy" {
		t.Fatalf("health = %#v, %v", result, err)
	}
	if result, err := sdk.Monitoring.Ready(ctx); err != nil || result.Status != "healthy" {
		t.Fatalf("ready = %#v, %v", result, err)
	}
	if result, err := sdk.Monitoring.Version(ctx); err != nil || result.Version != "test" {
		t.Fatalf("version = %#v, %v", result, err)
	}
	if result, err := sdk.Auth.Bootstrap(ctx, auth.BootstrapRequest{TenantName: "Tenant", TenantSlug: "tenant", Name: "Owner", Email: "owner@example.com", Password: "secret"}); err != nil || result.Tenant.ID != "t1" {
		t.Fatalf("bootstrap = %#v, %v", result, err)
	}
	if tokens, err := sdk.Auth.Login(ctx, auth.LoginRequest{TenantSlug: "tenant", Email: "owner@example.com", Password: "secret"}); err != nil || tokens.AccessToken != "access-1" {
		t.Fatalf("login = %#v, %v", tokens, err)
	}
	if tokens, err := sdk.Auth.Refresh(ctx, "refresh-1"); err != nil || tokens.AccessToken != "access-2" {
		t.Fatalf("refresh = %#v, %v", tokens, err)
	}
	if result, err := sdk.Users.List(ctx); err != nil || result.Count != 1 || result.Items[0].ID != "u1" {
		t.Fatalf("users list = %#v, %v", result, err)
	}
	if result, err := sdk.Users.Create(ctx, users.CreateRequest{Email: "new@example.com", Name: "New", Password: "secret", Status: "active"}); err != nil || result.ID != "u2" {
		t.Fatalf("user create = %#v, %v", result, err)
	}
	if result, err := sdk.Users.Update(ctx, "u1", users.UpdateRequest{Name: "Owner"}); err != nil || result.ID != "u1" {
		t.Fatalf("user update = %#v, %v", result, err)
	}
	if _, err := sdk.Users.Action(ctx, "u1", "verify-email"); err != nil {
		t.Fatal(err)
	}
	if tenants, err := sdk.Tenants.List(ctx); err != nil || tenants.Count != 1 {
		t.Fatalf("tenants = %#v, %v", tenants, err)
	}
	if tenant, err := sdk.Tenants.Create(ctx, tenancy.Tenant{Name: "Tenant 2", Slug: "tenant-2"}); err != nil || tenant.ID != "t2" {
		t.Fatalf("tenant create = %#v, %v", tenant, err)
	}
	if tenant, err := sdk.Tenants.Update(ctx, "t1", tenancy.Tenant{Name: "Tenant"}); err != nil || tenant.ID != "t1" {
		t.Fatalf("tenant update = %#v, %v", tenant, err)
	}
	if _, err := sdk.Tenants.Suspend(ctx, "t1"); err != nil {
		t.Fatal(err)
	}
	if _, err := sdk.Tenants.Activate(ctx, "t1"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Sessions.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("sessions = %#v, %v", result, err)
	}
	if err := sdk.Sessions.Revoke(ctx, "s1"); err != nil {
		t.Fatal(err)
	}
	role := rbac.Role{Name: "operator", Permissions: []string{"audit.read"}}
	if result, err := sdk.Roles.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("roles = %#v, %v", result, err)
	}
	if result, err := sdk.Roles.Create(ctx, role); err != nil || result.ID != "r1" {
		t.Fatalf("role create = %#v, %v", result, err)
	}
	role.ID = "r1"
	if result, err := sdk.Roles.Update(ctx, role); err != nil || result.ID != "r1" {
		t.Fatalf("role update = %#v, %v", result, err)
	}
	if err := sdk.Roles.Delete(ctx, "r1"); err != nil {
		t.Fatal(err)
	}
	permission := permissions.Permission{Name: "audit.read", Resource: "audit", Action: "read"}
	if result, err := sdk.Permissions.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("permissions = %#v, %v", result, err)
	}
	if result, err := sdk.Permissions.Create(ctx, permission); err != nil || result.ID != "p1" {
		t.Fatalf("permission create = %#v, %v", result, err)
	}
	if err := sdk.Permissions.Delete(ctx, "p1"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Permissions.Check(ctx, "u1", "audit.read"); err != nil || !result.Allowed {
		t.Fatalf("permission check = %#v, %v", result, err)
	}
	if result, err := sdk.APIKeys.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("API keys = %#v, %v", result, err)
	}
	if result, err := sdk.APIKeys.Create(ctx, apikeys.CreateRequest{Name: "key", Scopes: []string{"audit.read"}}); err != nil || result.Key != "hmk_created" {
		t.Fatalf("API key create = %#v, %v", result, err)
	}
	if result, err := sdk.APIKeys.Rotate(ctx, "k1"); err != nil || result.Key != "hmk_rotated" {
		t.Fatalf("API key rotate = %#v, %v", result, err)
	}
	if err := sdk.APIKeys.Revoke(ctx, "k1"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Secrets.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("secrets = %#v, %v", result, err)
	}
	if result, err := sdk.Secrets.Create(ctx, secrets.CreateRequest{Name: "DB", Type: "custom", Value: "secret"}); err != nil || result.MaskedValue != "pass****" {
		t.Fatalf("secret create = %#v, %v", result, err)
	}
	if _, err := sdk.Secrets.Rotate(ctx, "secret-1", "new-secret"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Secrets.History(ctx, "secret-1"); err != nil || result.Count != 1 || result.Items[0].MaskedValue != "pass****" {
		t.Fatalf("secret history = %#v, %v", result, err)
	}
	if err := sdk.Secrets.Revoke(ctx, "secret-1"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Audit.List(ctx, 100); err != nil || result.Count != 1 {
		t.Fatalf("audit list = %#v, %v", result, err)
	}
	if result, err := sdk.Audit.Create(ctx, audit.CreateRequest{Action: "policy.changed", Resource: "roles", Severity: "info", Status: "success"}); err != nil || result.ID != "e2" {
		t.Fatalf("audit create = %#v, %v", result, err)
	}
	if result, err := sdk.Fraud.ListRules(ctx); err != nil || result.Count != 1 {
		t.Fatalf("fraud rules = %#v, %v", result, err)
	}
	rule := fraud.Rule{Name: "velocity", Enabled: true, Field: "count", Operator: "gt", Threshold: "10", Action: "block", Score: 50}
	if result, err := sdk.Fraud.CreateRule(ctx, rule); err != nil || result.ID != "f1" {
		t.Fatalf("fraud create = %#v, %v", result, err)
	}
	rule.ID = "f1"
	if result, err := sdk.Fraud.UpdateRule(ctx, rule); err != nil || result.ID != "f1" {
		t.Fatalf("fraud update = %#v, %v", result, err)
	}
	if err := sdk.Fraud.DeleteRule(ctx, "f1"); err != nil {
		t.Fatal(err)
	}
	if result, err := sdk.Fraud.Analyze(ctx, map[string]any{"count": 11}); err != nil || result.RiskScore != 50 {
		t.Fatalf("fraud analysis = %#v, %v", result, err)
	}
	if result, err := sdk.Fraud.ListViolations(ctx); err != nil || result.Count != 1 {
		t.Fatalf("fraud violations = %#v, %v", result, err)
	}
	if result, err := sdk.Settings.List(ctx); err != nil || result.Count != 1 {
		t.Fatalf("settings = %#v, %v", result, err)
	}
	if result, err := sdk.Settings.Save(ctx, settings.Setting{Key: "feature.enabled", Value: json.RawMessage("true")}); err != nil || result.Key != "feature.enabled" {
		t.Fatalf("setting save = %#v, %v", result, err)
	}
}
