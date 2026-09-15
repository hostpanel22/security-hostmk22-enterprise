package tenancy_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azezaz/HostMK-Security-Framework/core/tenancy"
)

type MockRepository struct {
	tenant *tenancy.Tenant
	err    error
}

func (m *MockRepository) FindByID(
	ctx context.Context,
	tenantID string,
) (*tenancy.Tenant, error) {
	return m.tenant, m.err
}

func (m *MockRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*tenancy.Tenant, error) {
	return m.tenant, m.err
}

func (m *MockRepository) FindByDomain(
	ctx context.Context,
	domain string,
) (*tenancy.Tenant, error) {
	return m.tenant, m.err
}

func buildTenant() *tenancy.Tenant {
	return &tenancy.Tenant{
		ID:     "tenant-1",
		Name:   "Client One",
		Slug:   "client-one",
		Domain: "client-one.hostmk.test",
		Status: tenancy.StatusActive,
	}
}

func TestResolverFromHeader(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{
			HeaderName: "X-Tenant-ID",
		},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if ctx.Source != tenancy.SourceHeader {
		t.Fatal("invalid resolution source")
	}
}

func TestResolverFromDomain(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Host = "client-one.hostmk.test"

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx.TenantID != "tenant-1" {
		t.Fatal("invalid tenant id")
	}

	if ctx.Source != tenancy.SourceDomain {
		t.Fatal("expected domain source")
	}
}

func TestResolverRejectsNilRequest(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{},
	)

	_, err := resolver.ResolveRequest(
		nil,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolverMissingHeader(t *testing.T) {

	repository := &MockRepository{
		err: errors.New("tenant not found"),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{
			HeaderName: "X-Tenant-ID",
		},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	_, err := resolver.ResolveRequest(
		request,
	)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolverExtractsSubdomain(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{
			CentralDomain: "hostmk.test",
		},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Host = "client-one.hostmk.test"

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx == nil {
		t.Fatal("expected tenant context")
	}
}

func TestResolverStoresIPAddress(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Header.Set(
		"X-Forwarded-For",
		"192.168.1.10",
	)

	request.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx.IPAddress != "192.168.1.10" {
		t.Fatalf(
			"expected ip 192.168.1.10 got %s",
			ctx.IPAddress,
		)
	}
}

func TestResolverStoresUserAgent(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	request.Header.Set(
		"User-Agent",
		"HostMK-Test-Agent",
	)

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx.UserAgent != "HostMK-Test-Agent" {
		t.Fatal("invalid user agent")
	}
}

func TestResolverStoresRequestID(t *testing.T) {

	repository := &MockRepository{
		tenant: buildTenant(),
	}

	resolver := tenancy.NewResolver(
		repository,
		tenancy.ResolverConfig{},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	request.Header.Set(
		"X-Tenant-ID",
		"tenant-1",
	)

	request.Header.Set(
		"X-Request-ID",
		"req-123",
	)

	ctx, err := resolver.ResolveRequest(
		request,
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if ctx.RequestID != "req-123" {
		t.Fatal("invalid request id")
	}
}
