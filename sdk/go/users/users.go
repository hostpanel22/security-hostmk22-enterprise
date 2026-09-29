package users

import (
	"context"
	"net/url"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type User struct {
	ID              string   `json:"id"`
	TenantID        string   `json:"tenant_id"`
	Email           string   `json:"email"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	EmailVerifiedAt *string  `json:"email_verified_at,omitempty"`
	MFAEnabled      bool     `json:"mfa_enabled"`
	RoleIDs         []string `json:"role_ids,omitempty"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	LastLogin       *string  `json:"last_login,omitempty"`
}

type CreateRequest struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Status   string   `json:"status"`
	RoleIDs  []string `json:"role_ids,omitempty"`
}

type UpdateRequest struct {
	Name    string   `json:"name,omitempty"`
	Status  string   `json:"status,omitempty"`
	RoleIDs []string `json:"role_ids,omitempty"`
}

type ActionResult struct {
	ID         string `json:"id,omitempty"`
	Secret     string `json:"secret,omitempty"`
	OTPAuthURL string `json:"otpauth_url,omitempty"`
}

func (s *Service) List(ctx context.Context) (client.ListResult[User], error) {
	result := client.ListResult[User]{Items: []User{}}
	count, err := s.client.GetList(ctx, "/api/v1/users", nil, "users", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Create(ctx context.Context, input CreateRequest) (User, error) {
	var result User
	err := s.client.Post(ctx, "/api/v1/users", nil, input, &result)
	return result, err
}

func (s *Service) Update(ctx context.Context, id string, input UpdateRequest) (User, error) {
	var result User
	err := s.client.Post(ctx, "/api/v1/users/action", url.Values{"id": {id}}, input, &result)
	return result, err
}

func (s *Service) Action(ctx context.Context, id, action string) (ActionResult, error) {
	var result ActionResult
	query := url.Values{"id": {id}, "action": {action}}
	err := s.client.Post(ctx, "/api/v1/users/action", query, nil, &result)
	return result, err
}
