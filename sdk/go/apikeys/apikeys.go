package apikeys

import (
	"context"
	"net/url"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type APIKey struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Key        string     `json:"key,omitempty"`
	Scopes     []string   `json:"scopes"`
	RateLimit  *int       `json:"rate_limit,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CreateRequest struct {
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	RateLimit *int       `json:"rate_limit,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[APIKey], error) {
	result := client.ListResult[APIKey]{Items: []APIKey{}}
	count, err := s.client.GetList(ctx, "/api/v1/api-keys", nil, "api_keys", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Create(ctx context.Context, input CreateRequest) (APIKey, error) {
	var result APIKey
	err := s.client.Post(ctx, "/api/v1/api-keys", nil, input, &result)
	return result, err
}

func (s *Service) Rotate(ctx context.Context, id string) (APIKey, error) {
	var result APIKey
	query := url.Values{"action": {"rotate"}, "id": {id}}
	err := s.client.Post(ctx, "/api/v1/api-keys", query, nil, &result)
	return result, err
}

func (s *Service) Revoke(ctx context.Context, id string) error {
	return s.client.Delete(ctx, "/api/v1/api-keys", url.Values{"id": {id}}, nil)
}
