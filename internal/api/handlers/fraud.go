package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type FraudRequest struct {
	TenantID  string  `json:"tenant_id"`
	UserID    string  `json:"user_id"`
	IPAddress string  `json:"ip_address"`
	Country   string  `json:"country"`
	Amount    float64 `json:"amount"`
}

type FraudResult struct {
	RiskScore  int       `json:"risk_score"`
	RiskLevel  string    `json:"risk_level"`
	Violations []string  `json:"violations"`
	AnalyzedAt time.Time `json:"analyzed_at"`
}

type FraudViolation struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Rule        string    `json:"rule"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func FraudAnalyze(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:
		getFraudAnalysis(w)

	case http.MethodPost:
		postFraudAnalysis(w, r)

	default:
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
	}
}

func FraudViolations(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodGet {

		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)

		return
	}

	violations := []FraudViolation{
		{
			ID:          uuid.NewString(),
			TenantID:    "tenant-001",
			Rule:        "HIGH_AMOUNT",
			Severity:    "medium",
			Description: "transaction amount exceeded threshold",
			CreatedAt:   time.Now().UTC(),
		},
	}

	response.OK(
		w,
		"fraud violations retrieved",
		map[string]any{
			"count":      len(violations),
			"violations": violations,
		},
	)
}

func getFraudAnalysis(
	w http.ResponseWriter,
) {

	result := FraudResult{
		RiskScore:  0,
		RiskLevel:  "low",
		Violations: []string{},
		AnalyzedAt: time.Now().UTC(),
	}

	response.OK(
		w,
		"fraud analysis completed",
		result,
	)
}

func postFraudAnalysis(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request FraudRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		response.BadRequest(
			w,
			"invalid json payload",
		)

		return
	}

	result := calculateRisk(
		request,
	)

	response.OK(
		w,
		"fraud analysis completed",
		result,
	)
}

func calculateRisk(
	request FraudRequest,
) FraudResult {

	score := 0

	violations := make(
		[]string,
		0,
	)

	if request.Amount >= 1000 {

		score += 40

		violations = append(
			violations,
			"HIGH_AMOUNT",
		)
	}

	if request.Country != "" &&
		!strings.EqualFold(
			request.Country,
			"FR",
		) {

		score += 20

		violations = append(
			violations,
			"FOREIGN_COUNTRY",
		)
	}

	if !isValidIP(
		request.IPAddress,
	) {

		score += 30

		violations = append(
			violations,
			"INVALID_IP",
		)
	}

	level := riskLevel(
		score,
	)

	return FraudResult{
		RiskScore:  score,
		RiskLevel:  level,
		Violations: violations,
		AnalyzedAt: time.Now().UTC(),
	}
}

func riskLevel(
	score int,
) string {

	switch {

	case score >= 70:
		return "high"

	case score >= 40:
		return "medium"

	default:
		return "low"
	}
}

func isValidIP(
	ip string,
) bool {

	if ip == "" {
		return false
	}

	return net.ParseIP(ip) != nil
}
