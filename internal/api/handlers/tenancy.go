package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Plan      string    `json:"plan"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTenantRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Plan string `json:"plan"`
}

func Tenants(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		listTenants(w)

	case http.MethodPost:
		createTenant(w, r)

	default:
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
	}
}

func listTenants(
	w http.ResponseWriter,
) {
	now := time.Now().UTC()

	tenants := []Tenant{
		{
			ID:        "tenant-001",
			Name:      "HostMK Demo",
			Slug:      "hostmk-demo",
			Plan:      "enterprise",
			Status:    "active",
			CreatedAt: now,
		},
		{
			ID:        "tenant-002",
			Name:      "Acme Corp",
			Slug:      "acme-corp",
			Plan:      "pro",
			Status:    "active",
			CreatedAt: now,
		},
	}

	response.OK(
		w,
		"tenants retrieved",
		map[string]any{
			"count":   len(tenants),
			"tenants": tenants,
		},
	)
}

func createTenant(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateTenantRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.BadRequest(
			w,
			"invalid json payload",
		)
		return
	}

	request.Name = strings.TrimSpace(
		request.Name,
	)

	request.Slug = strings.TrimSpace(
		request.Slug,
	)

	request.Plan = strings.TrimSpace(
		request.Plan,
	)

	if request.Name == "" {
		response.BadRequest(
			w,
			"name is required",
		)
		return
	}

	if request.Slug == "" {
		response.BadRequest(
			w,
			"slug is required",
		)
		return
	}

	if request.Plan == "" {
		response.BadRequest(
			w,
			"plan is required",
		)
		return
	}

	tenant := Tenant{
		ID:        uuid.NewString(),
		Name:      request.Name,
		Slug:      strings.ToLower(request.Slug),
		Plan:      request.Plan,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
	}

	response.Created(
		w,
		"tenant created",
		tenant,
	)
}
