package settings

import (
	"context"
	"encoding/json"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Setting struct {
	ID        string          `json:"id"`
	TenantID  string          `json:"tenant_id,omitempty"`
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	UpdatedAt string          `json:"updated_at"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[Setting], error) {
	result := client.ListResult[Setting]{Items: []Setting{}}
	count, err := s.client.GetList(ctx, "/api/v1/settings", nil, "settings", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Save(ctx context.Context, setting Setting) (Setting, error) {
	var result Setting
	err := s.client.Put(ctx, "/api/v1/settings", nil, setting, &result)
	return result, err
}
