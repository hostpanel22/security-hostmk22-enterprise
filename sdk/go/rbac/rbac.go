package rbac

import (
	"context"
	"net/url"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Roles struct{ client *client.Client }

type Role struct {
	ID          string   `json:"id"`
	TenantID    string   `json:"tenant_id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func NewRoles(apiClient *client.Client) *Roles { return &Roles{client: apiClient} }

func (s *Roles) List(ctx context.Context) (client.ListResult[Role], error) {
	result := client.ListResult[Role]{Items: []Role{}}
	count, err := s.client.GetList(ctx, "/api/v1/rbac/roles", nil, "roles", &result.Items)
	result.Count = count
	return result, err
}

func (s *Roles) Create(ctx context.Context, role Role) (Role, error) {
	return s.write(ctx, "POST", role)
}

func (s *Roles) Update(ctx context.Context, role Role) (Role, error) {
	return s.write(ctx, "PUT", role)
}

func (s *Roles) write(ctx context.Context, method string, role Role) (Role, error) {
	var result Role
	var err error
	if method == "PUT" {
		err = s.client.Put(ctx, "/api/v1/rbac/roles", nil, role, &result)
	} else {
		err = s.client.Post(ctx, "/api/v1/rbac/roles", nil, role, &result)
	}
	return result, err
}

func (s *Roles) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, "/api/v1/rbac/roles", url.Values{"id": {id}}, nil)
}
