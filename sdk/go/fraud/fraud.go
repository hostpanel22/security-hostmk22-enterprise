package fraud

import (
	"context"
	"net/url"
	"time"

	"github.com/azezaz/HostMK-Security-Framework/sdk/go/client"
)

type Service struct{ client *client.Client }

func New(apiClient *client.Client) *Service { return &Service{client: apiClient} }

type Rule struct {
	ID            string         `json:"id"`
	TenantID      string         `json:"tenant_id"`
	Name          string         `json:"name"`
	Enabled       bool           `json:"enabled"`
	Field         string         `json:"field"`
	Operator      string         `json:"operator"`
	Threshold     string         `json:"threshold"`
	WindowSeconds *int           `json:"window_seconds,omitempty"`
	Action        string         `json:"action"`
	Score         int            `json:"score"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Violation struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenant_id"`
	UserID    string         `json:"user_id,omitempty"`
	SignalID  string         `json:"signal_id"`
	RuleID    string         `json:"rule_id"`
	Score     float64        `json:"score"`
	Level     string         `json:"level"`
	Reason    string         `json:"reason"`
	Status    string         `json:"status"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Analysis struct {
	RiskScore  int       `json:"risk_score"`
	RiskLevel  string    `json:"risk_level"`
	Blocked    bool      `json:"blocked"`
	Matched    []string  `json:"matched"`
	AnalyzedAt time.Time `json:"analyzed_at"`
}

func (s *Service) ListRules(ctx context.Context) (client.ListResult[Rule], error) {
	result := client.ListResult[Rule]{Items: []Rule{}}
	count, err := s.client.GetList(ctx, "/api/v1/fraud/rules", nil, "rules", &result.Items)
	result.Count = count
	return result, err
}

func (s *Service) CreateRule(ctx context.Context, rule Rule) (Rule, error) {
	var result Rule
	err := s.client.Post(ctx, "/api/v1/fraud/rules", nil, rule, &result)
	return result, err
}

func (s *Service) UpdateRule(ctx context.Context, rule Rule) (Rule, error) {
	var result Rule
	err := s.client.Put(ctx, "/api/v1/fraud/rules", nil, rule, &result)
	return result, err
}

func (s *Service) DeleteRule(ctx context.Context, id string) error {
	return s.client.Delete(ctx, "/api/v1/fraud/rules", url.Values{"id": {id}}, nil)
}

func (s *Service) Analyze(ctx context.Context, attributes map[string]any) (Analysis, error) {
	var result Analysis
	err := s.client.Post(ctx, "/api/v1/fraud/analyze", nil, attributes, &result)
	return result, err
}

func (s *Service) ListViolations(ctx context.Context) (client.ListResult[Violation], error) {
	result := client.ListResult[Violation]{Items: []Violation{}}
	count, err := s.client.GetList(ctx, "/api/v1/fraud/violations", nil, "events", &result.Items)
	result.Count = count
	return result, err
}
