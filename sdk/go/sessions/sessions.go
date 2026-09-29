package sessions

import (
	"context"
	"net/url"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Session struct {
	ID        string  `json:"id"`
	TenantID  string  `json:"tenant_id"`
	UserID    string  `json:"user_id"`
	UserAgent string  `json:"user_agent,omitempty"`
	IPAddress string  `json:"ip_address,omitempty"`
	ExpiresAt string  `json:"expires_at"`
	RevokedAt *string `json:"revoked_at,omitempty"`
	CreatedAt string  `json:"created_at"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[Session], error) {
	result := client.ListResult[Session]{Items: []Session{}}
	count, err := s.client.GetList(ctx, "/api/v1/sessions", nil, "sessions", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Revoke(ctx context.Context, id string) error {
	query := url.Values{"id": {id}, "action": {"revoke"}}
	return s.client.Post(ctx, "/api/v1/sessions", query, nil, nil)
}
