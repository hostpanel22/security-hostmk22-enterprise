package tenancy

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

/*
|--------------------------------------------------------------------------
| Resolution Sources
|--------------------------------------------------------------------------
*/

const (
	SourceSubdomain = "subdomain"
	SourceDomain    = "domain"
	SourceHeader    = "header"
	SourceJWT       = "jwt"
	SourceAPIKey    = "api_key"
	SourceManual    = "manual"
)

/*
|--------------------------------------------------------------------------
| Tenant Repository
|--------------------------------------------------------------------------
|
| Interface implemented by storage layers.
| PostgreSQL, MySQL, Redis or custom adapters
| can implement this contract.
|
*/

type Repository interface {
	FindByID(
		ctx context.Context,
		tenantID string,
	) (*Tenant, error)

	FindBySlug(
		ctx context.Context,
		slug string,
	) (*Tenant, error)

	FindByDomain(
		ctx context.Context,
		domain string,
	) (*Tenant, error)
}

/*
|--------------------------------------------------------------------------
| Resolver Configuration
|--------------------------------------------------------------------------
*/

type ResolverConfig struct {
	CentralDomain string

	HeaderName string

	TrustedProxies []string
}

/*
|--------------------------------------------------------------------------
| Resolver
|--------------------------------------------------------------------------
*/

type Resolver struct {
	repository Repository

	config ResolverConfig
}

/*
|--------------------------------------------------------------------------
| Constructor
|--------------------------------------------------------------------------
*/

func NewResolver(
	repository Repository,
	config ResolverConfig,
) *Resolver {

	if config.HeaderName == "" {
		config.HeaderName = "X-Tenant-ID"
	}

	return &Resolver{
		repository: repository,
		config:     config,
	}
}

/*
|--------------------------------------------------------------------------
| Resolve Request
|--------------------------------------------------------------------------
*/

func (r *Resolver) ResolveRequest(
	request *http.Request,
) (*TenantContext, error) {

	if request == nil {
		return nil, ErrTenantResolutionFailed
	}

	if tenantCtx, err := r.resolveFromHeader(request); err == nil {
		return tenantCtx, nil
	}

	if tenantCtx, err := r.resolveFromDomain(request); err == nil {
		return tenantCtx, nil
	}

	return nil, ErrTenantNotFound
}

/*
|--------------------------------------------------------------------------
| Header Resolution
|--------------------------------------------------------------------------
*/

func (r *Resolver) resolveFromHeader(
	request *http.Request,
) (*TenantContext, error) {

	tenantID := strings.TrimSpace(
		request.Header.Get(
			r.config.HeaderName,
		),
	)

	if tenantID == "" {
		return nil, ErrMissingTenantHeader
	}

	tenant, err := r.repository.FindByID(
		request.Context(),
		tenantID,
	)

	if err != nil {
		return nil, ErrTenantNotFound
	}

	return r.buildContext(
		request,
		tenant,
		SourceHeader,
	), nil
}

/*
|--------------------------------------------------------------------------
| Domain Resolution
|--------------------------------------------------------------------------
*/

func (r *Resolver) resolveFromDomain(
	request *http.Request,
) (*TenantContext, error) {

	host := request.Host

	if host == "" {
		return nil, ErrInvalidTenantDomain
	}

	host = strings.Split(
		host,
		":",
	)[0]

	tenant, err := r.repository.FindByDomain(
		request.Context(),
		host,
	)

	if err == nil {

		return r.buildContext(
			request,
			tenant,
			SourceDomain,
		), nil
	}

	subdomain := r.extractSubdomain(host)

	if subdomain == "" {
		return nil, ErrTenantNotFound
	}

	tenant, err = r.repository.FindBySlug(
		request.Context(),
		subdomain,
	)

	if err != nil {
		return nil, ErrTenantNotFound
	}

	return r.buildContext(
		request,
		tenant,
		SourceSubdomain,
	), nil
}

/*
|--------------------------------------------------------------------------
| Subdomain Extraction
|--------------------------------------------------------------------------
*/

func (r *Resolver) extractSubdomain(
	host string,
) string {

	if r.config.CentralDomain == "" {
		return ""
	}

	if !strings.HasSuffix(
		host,
		r.config.CentralDomain,
	) {
		return ""
	}

	subdomain := strings.TrimSuffix(
		host,
		r.config.CentralDomain,
	)

	subdomain = strings.TrimSuffix(
		subdomain,
		".",
	)

	return strings.TrimSpace(
		subdomain,
	)
}

/*
|--------------------------------------------------------------------------
| Context Builder
|--------------------------------------------------------------------------
*/

func (r *Resolver) buildContext(
	request *http.Request,
	tenant *Tenant,
	source string,
) *TenantContext {

	return &TenantContext{
		TenantID:   tenant.ID,
		TenantSlug: tenant.Slug,
		Domain:     tenant.Domain,
		Source:     source,

		RequestID: request.Header.Get(
			"X-Request-ID",
		),

		IPAddress: r.clientIP(
			request,
		),

		UserAgent: request.UserAgent(),

		ResolvedAt: time.Now().UTC(),
	}
}

/*
|--------------------------------------------------------------------------
| Client IP Resolution
|--------------------------------------------------------------------------
*/

func (r *Resolver) clientIP(
	request *http.Request,
) string {

	if forwarded := request.Header.Get(
		"X-Forwarded-For",
	); forwarded != "" {

		parts := strings.Split(
			forwarded,
			",",
		)

		return strings.TrimSpace(
			parts[0],
		)
	}

	if realIP := request.Header.Get(
		"X-Real-IP",
	); realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(
		request.RemoteAddr,
	)

	if err != nil {
		return request.RemoteAddr
	}

	return ip
}
