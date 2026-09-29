package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequestBuildsEndpointAndAuthentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prefix/api/v1/users" || r.URL.Query().Get("limit") != "25" {
			t.Errorf("unexpected URL %s", r.URL.String())
		}
		if got := r.Header.Get("Authorization"); got != "Bearer access-token" {
			t.Errorf("Authorization = %q", got)
		}
		if r.Header.Get("X-API-Key") != "" {
			t.Error("unexpected API key header")
		}
		if r.Header.Get("User-Agent") != userAgent {
			t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"message":"ok","data":{"users":[],"count":0}}`))
	}))
	defer server.Close()

	apiClient, err := New(server.URL+"/prefix/api/v1", WithBearerToken("access-token"))
	if err != nil {
		t.Fatal(err)
	}
	query := url.Values{"limit": []string{"25"}}
	var response struct {
		Users []any `json:"users"`
		Count int   `json:"count"`
	}
	if err := apiClient.Get(context.Background(), "/api/v1/users", query, &response); err != nil {
		t.Fatal(err)
	}
	if response.Users == nil || response.Count != 0 {
		t.Fatalf("response = %#v", response)
	}
}

func TestAPIKeyTakesPrecedenceOverBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "key-value" || r.Header.Get("Authorization") != "" {
			t.Errorf("unexpected credentials: api-key=%q authorization=%q", r.Header.Get("X-API-Key"), r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	apiClient, err := New(server.URL, WithAPIKey(" key-value "), WithBearerToken("access-token"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := apiClient.Request(context.Background(), http.MethodGet, "/health", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetListDecodesWrappedList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"count":1,"users":[{"id":"u1"}]}}`))
	}))
	defer server.Close()

	apiClient, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	var users []struct {
		ID string `json:"id"`
	}
	count, err := apiClient.GetList(context.Background(), "/api/v1/users", nil, "users", &users)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 || len(users) != 1 || users[0].ID != "u1" {
		t.Fatalf("count=%d users=%s", count, mustJSON(t, users))
	}
}

func TestRequestReturnsStructuredAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"success":false,"error":"authentication required"}`))
	}))
	defer server.Close()

	apiClient, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = apiClient.Request(context.Background(), http.MethodGet, "/api/v1/users", nil, nil)
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusUnauthorized || apiErr.Message != "authentication required" {
		t.Fatalf("error = %#v", err)
	}
}

func TestNewRejectsMalformedBaseURL(t *testing.T) {
	for _, baseURL := range []string{"", "ftp://example.com", "https://user:pass@example.com", "https://example.com?token=secret"} {
		if _, err := New(baseURL); !errors.Is(err, ErrInvalidBaseURL) {
			t.Errorf("New(%q) error = %v", baseURL, err)
		}
	}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
