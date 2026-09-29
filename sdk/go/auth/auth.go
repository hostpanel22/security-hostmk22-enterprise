package auth

import (
	"context"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type LoginRequest struct {
	TenantSlug string `json:"tenant_slug"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	MFACode    string `json:"mfa_code,omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type User struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type Tenant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

type BootstrapRequest struct {
	TenantName string `json:"tenant_name"`
	TenantSlug string `json:"tenant_slug"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type BootstrapResult struct {
	Tenant Tenant `json:"tenant"`
	User   User   `json:"user"`
}

func (s *Service) Login(ctx context.Context, input LoginRequest) (Tokens, error) {
	var result Tokens
	err := s.client.Post(ctx, "/api/v1/auth/login", nil, input, &result)
	if err != nil {
		return Tokens{}, err
	}
	s.client.SetAPIKey("")
	s.client.SetBearerToken(result.AccessToken)
	return result, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	var result Tokens
	err := s.client.Post(ctx, "/api/v1/auth/refresh", nil, map[string]string{"refresh_token": refreshToken}, &result)
	if err != nil {
		return Tokens{}, err
	}
	s.client.SetAPIKey("")
	s.client.SetBearerToken(result.AccessToken)
	return result, nil
}

func (s *Service) Bootstrap(ctx context.Context, input BootstrapRequest) (BootstrapResult, error) {
	var result BootstrapResult
	err := s.client.Post(ctx, "/api/v1/auth/bootstrap", nil, input, &result)
	return result, err
}

func (s *Service) ClearToken() { s.client.SetBearerToken("") }
