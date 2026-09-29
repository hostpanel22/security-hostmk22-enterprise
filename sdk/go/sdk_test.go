package hostmk_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/apikeys"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/auth"
	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

func TestLoginSwitchesAPIKeyClientToBearerAndDecodesServices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/auth/login":
			var body auth.LoginRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode login request: %v", err)
			}
			if body.TenantSlug != "acme" || body.Email != "operator@example.com" || body.MFACode != "123456" {
				t.Errorf("login request = %#v", body)
			}
			_, _ = w.Write([]byte(`{"success":true,"message":"authenticated","data":{"access_token":"access-1","refresh_token":"refresh-1","token_type":"Bearer","expires_in":900}}`))
		case "/api/v1/users":
			if got := r.Header.Get("Authorization"); got != "Bearer access-1" {
				t.Errorf("Authorization = %q", got)
			}
			if got := r.Header.Get("X-API-Key"); got != "" {
				t.Errorf("unexpected X-API-Key = %q", got)
			}
			_, _ = w.Write([]byte(`{"success":true,"message":"users retrieved","data":{"count":1,"users":[{"id":"u1","tenant_id":"t1","email":"operator@example.com","name":"Operator","status":"active"}]}}`))
		case "/api/v1/api-keys":
			if r.Method != http.MethodPost {
				t.Errorf("API key method = %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"success":true,"data":{"id":"k1","tenant_id":"t1","name":"automation","prefix":"hmk_ab","key":"hmk_one_time","scopes":["audit.read"]}}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	sdk, err := hostmk.NewClient(server.URL+"/api/v1", client.WithAPIKey("bootstrap-key"), client.WithBearerToken("old-token"))
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := sdk.Auth.Login(context.Background(), auth.LoginRequest{TenantSlug: "acme", Email: "operator@example.com", Password: "not-logged", MFACode: "123456"})
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken != "access-1" || tokens.RefreshToken != "refresh-1" {
		t.Fatalf("tokens = %#v", tokens)
	}
	users, err := sdk.Users.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if users.Count != 1 || len(users.Items) != 1 || users.Items[0].ID != "u1" {
		t.Fatalf("users = %#v", users)
	}
	key, err := sdk.APIKeys.Create(context.Background(), apikeys.CreateRequest{Name: "automation", Scopes: []string{"audit.read"}})
	if err != nil {
		t.Fatal(err)
	}
	if key.Key != "hmk_one_time" || key.ID != "k1" {
		t.Fatalf("created key = %#v", key)
	}
}
