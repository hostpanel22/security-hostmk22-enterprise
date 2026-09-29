package audit

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Event struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenant_id"`
	UserID    string         `json:"user_id"`
	Action    string         `json:"action"`
	Resource  string         `json:"resource"`
	Severity  string         `json:"severity"`
	Status    string         `json:"status"`
	IPAddress string         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
	Message   string         `json:"message,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type CreateRequest struct {
	Action   string         `json:"action"`
	Resource string         `json:"resource"`
	Severity string         `json:"severity"`
	Status   string         `json:"status"`
	Message  string         `json:"message,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func (s *Service) List(ctx context.Context, limit int) (client.ListResult[Event], error) {
	result := client.ListResult[Event]{Items: []Event{}}
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", formatLimit(limit))
	}
	count, err := s.client.GetList(ctx, "/api/v1/audit/events", query, "events", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) Create(ctx context.Context, input CreateRequest) (Event, error) {
	var result Event
	err := s.client.Post(ctx, "/api/v1/audit/events", nil, input, &result)
	return result, err
}

func formatLimit(limit int) string {
	return strconv.Itoa(limit)
}
