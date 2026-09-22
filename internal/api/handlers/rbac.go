package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func RBACRoles(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		listRoles(w)

	case http.MethodPost:
		createRole(w, r)

	default:
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)
	}
}

func listRoles(
	w http.ResponseWriter,
) {
	roles := []Role{
		{
			ID:          "role-super-admin",
			Name:        "super-admin",
			Description: "HostMK platform administrator",
			CreatedAt:   time.Now().UTC(),
		},
		{
			ID:          "role-tenant-admin",
			Name:        "tenant-admin",
			Description: "Tenant administrator",
			CreatedAt:   time.Now().UTC(),
		},
		{
			ID:          "role-user",
			Name:        "user",
			Description: "Standard user",
			CreatedAt:   time.Now().UTC(),
		},
	}

	response.OK(
		w,
		"roles retrieved",
		map[string]any{
			"count": len(roles),
			"roles": roles,
		},
	)
}

func createRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		response.BadRequest(
			w,
			"invalid json payload",
		)

		return
	}

	if err := validateRoleRequest(request); err != nil {

		response.BadRequest(
			w,
			err.Error(),
		)

		return
	}

	role := Role{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(request.Name),
		Description: strings.TrimSpace(request.Description),
		CreatedAt:   time.Now().UTC(),
	}

	response.Created(
		w,
		"role created",
		role,
	)
}

func validateRoleRequest(
	request CreateRoleRequest,
) error {

	if strings.TrimSpace(request.Name) == "" {
		return newRBACValidationError(
			"name is required",
		)
	}

	if strings.TrimSpace(request.Description) == "" {
		return newRBACValidationError(
			"description is required",
		)
	}

	if len(request.Name) > 100 {
		return newRBACValidationError(
			"name exceeds 100 characters",
		)
	}

	if len(request.Description) > 500 {
		return newRBACValidationError(
			"description exceeds 500 characters",
		)
	}

	return nil
}

type rbacValidationError struct {
	message string
}

func (e rbacValidationError) Error() string {
	return e.message
}

func newRBACValidationError(
	message string,
) error {
	return rbacValidationError{
		message: message,
	}
}
