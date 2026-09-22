package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type Secret struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	MaskedValue string    `json:"masked_value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateSecretRequest struct {
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func Secrets(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		listSecrets(w)

	case http.MethodPost:
		createSecret(w, r)

	case http.MethodDelete:
		deleteSecret(w, r)

	default:
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
	}
}

func listSecrets(
	w http.ResponseWriter,
) {
	now := time.Now().UTC()

	secrets := []Secret{
		{
			ID:          "secret-001",
			TenantID:    "tenant-001",
			Name:        "STRIPE_API_KEY",
			MaskedValue: "sk_live************",
			Description: "Stripe payment gateway",
			CreatedAt:   now,
		},
		{
			ID:          "secret-002",
			TenantID:    "tenant-001",
			Name:        "JWT_SECRET",
			MaskedValue: "hostmk************",
			Description: "JWT signing key",
			CreatedAt:   now,
		},
	}

	response.OK(
		w,
		"secrets retrieved",
		map[string]any{
			"count":   len(secrets),
			"secrets": secrets,
		},
	)
}

func createSecret(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateSecretRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.BadRequest(
			w,
			"invalid json payload",
		)
		return
	}

	request.TenantID = strings.TrimSpace(
		request.TenantID,
	)

	request.Name = strings.TrimSpace(
		request.Name,
	)

	request.Value = strings.TrimSpace(
		request.Value,
	)

	if request.TenantID == "" {
		response.BadRequest(
			w,
			"tenant_id is required",
		)
		return
	}

	if request.Name == "" {
		response.BadRequest(
			w,
			"name is required",
		)
		return
	}

	if request.Value == "" {
		response.BadRequest(
			w,
			"value is required",
		)
		return
	}

	secret := Secret{
		ID:          uuid.NewString(),
		TenantID:    request.TenantID,
		Name:        request.Name,
		MaskedValue: maskSecret(request.Value),
		Description: request.Description,
		CreatedAt:   time.Now().UTC(),
	}

	response.Created(
		w,
		"secret created",
		secret,
	)
}

func deleteSecret(
	w http.ResponseWriter,
	r *http.Request,
) {
	secretID := strings.TrimSpace(
		r.URL.Query().Get("id"),
	)

	if secretID == "" {
		response.BadRequest(
			w,
			"secret id is required",
		)
		return
	}

	response.OK(
		w,
		"secret deleted",
		map[string]any{
			"id":         secretID,
			"deleted_at": time.Now().UTC(),
		},
	)
}

func maskSecret(
	value string,
) string {
	if len(value) <= 4 {
		return "****"
	}

	return value[:4] + strings.Repeat(
		"*",
		len(value)-4,
	)
}
