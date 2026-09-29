package tenancy

import (
	"context"
	"net/url"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Tenant struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Domain      string         `json:"domain,omitempty"`
	Status      string         `json:"status"`
	OwnerID     string         `json:"owner_id,omitempty"`
	Plan        string         `json:"plan"`
	Timezone    string         `json:"timezone,omitempty"`
	Country     string         `json:"country,omitempty"`
	Quotas      map[string]any `json:"quotas,omitempty"`
	Branding    map[string]any `json:"branding,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	SuspendedAt *string        `json:"suspended_at,omitempty"`
	ExpiresAt   *string        `json:"expires_at,omitempty"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[Tenant], error) {
	result := client.ListResult[Tenant]{Items: []Tenant{}}
	count, err := s.client.GetList(ctx, "/api/v1/tenants", nil, "tenants", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Create(ctx context.Context, input Tenant) (Tenant, error) {
	var result Tenant
	err := s.client.Post(ctx, "/api/v1/tenants", nil, input, &result)
	return result, err
}

func (s *Service) Update(ctx context.Context, id string, input Tenant) (Tenant, error) {
	input.ID = id
	var result Tenant
	err := s.client.Post(ctx, "/api/v1/tenants/action", url.Values{"id": {id}}, input, &result)
	return result, err
}

func (s *Service) Activate(ctx context.Context, id string) (Tenant, error) {
	return s.setStatus(ctx, id, "activate")
}

func (s *Service) Suspend(ctx context.Context, id string) (Tenant, error) {
	return s.setStatus(ctx, id, "suspend")
}

func (s *Service) setStatus(ctx context.Context, id, action string) (Tenant, error) {
	var result Tenant
	query := url.Values{"id": {id}, "action": {action}}
	err := s.client.Post(ctx, "/api/v1/tenants/action", query, nil, &result)
	return result, err
}
