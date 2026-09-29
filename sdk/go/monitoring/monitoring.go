package monitoring

import (
	"context"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Health struct {
	Status      string `json:"status"`
	Service     string `json:"service"`
	Version     string `json:"version"`
	Environment string `json:"env"`
	Components  []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"components"`
}

type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
}

func (s *Service) Health(ctx context.Context) (Health, error) {
	var result Health
	err := s.client.Get(ctx, "/health", nil, &result)
	return result, err
}

func (s *Service) Ready(ctx context.Context) (Health, error) {
	var result Health
	err := s.client.Get(ctx, "/ready", nil, &result)
	return result, err
}

func (s *Service) Version(ctx context.Context) (Version, error) {
	var result Version
	err := s.client.Get(ctx, "/api/v1/version", nil, &result)
	return result, err
}
