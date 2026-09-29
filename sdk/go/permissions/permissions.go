package permissions

import (
	"context"
	"net/url"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Permission struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type CheckResult struct {
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
	Allowed    bool   `json:"allowed"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[Permission], error) {
	result := client.ListResult[Permission]{Items: []Permission{}}
	count, err := s.client.GetList(ctx, "/api/v1/permissions", nil, "permissions", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Create(ctx context.Context, permission Permission) (Permission, error) {
	var result Permission
	err := s.client.Post(ctx, "/api/v1/permissions", nil, permission, &result)
	return result, err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, "/api/v1/permissions", url.Values{"id": {id}}, nil)
}

func (s *Service) Check(ctx context.Context, userID, permission string) (CheckResult, error) {
	var result CheckResult
	body := map[string]string{"user_id": userID, "permission": permission}
	err := s.client.Post(ctx, "/api/v1/permissions/check", nil, body, &result)
	return result, err
}
