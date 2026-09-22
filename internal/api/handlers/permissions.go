package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/azezaz/HostMK-Security-Framework/internal/api/response"
)

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type PermissionCheckRequest struct {
	Role       string `json:"role"`
	Permission string `json:"permission"`
}

type PermissionCheckResponse struct {
	Role       string `json:"role"`
	Permission string `json:"permission"`
	Allowed    bool   `json:"allowed"`
}

func Permissions(
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

	permissions := systemPermissions()

	response.OK(
		w,
		"permissions retrieved",
		map[string]any{
			"count":       len(permissions),
			"permissions": permissions,
		},
	)
}

func PermissionCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {

		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"method not allowed",
		)

		return
	}

	var request PermissionCheckRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {

		response.BadRequest(
			w,
			"invalid json payload",
		)

		return
	}

	if err := validatePermissionCheck(request); err != nil {

		response.BadRequest(
			w,
			err.Error(),
		)

		return
	}

	result := PermissionCheckResponse{
		Role:       request.Role,
		Permission: request.Permission,
		Allowed: isPermissionAllowed(
			request.Role,
			request.Permission,
		),
	}

	response.OK(
		w,
		"permission check completed",
		result,
	)
}

func systemPermissions() []Permission {
	return []Permission{
		{
			ID:          "perm-001",
			Name:        "tenant.read",
			Resource:    "tenant",
			Action:      "read",
			Description: "Read tenant information",
		},
		{
			ID:          "perm-002",
			Name:        "tenant.write",
			Resource:    "tenant",
			Action:      "write",
			Description: "Modify tenant information",
		},
		{
			ID:          "perm-003",
			Name:        "user.read",
			Resource:    "user",
			Action:      "read",
			Description: "Read users",
		},
		{
			ID:          "perm-004",
			Name:        "user.write",
			Resource:    "user",
			Action:      "write",
			Description: "Modify users",
		},
		{
			ID:          "perm-005",
			Name:        "audit.read",
			Resource:    "audit",
			Action:      "read",
			Description: "Read audit logs",
		},
	}
}

func isPermissionAllowed(
	role string,
	permission string,
) bool {

	rolePermissions := map[string][]string{
		"super-admin": {
			"*",
		},
		"tenant-admin": {
			"tenant.read",
			"tenant.write",
			"user.read",
			"user.write",
			"audit.read",
		},
		"user": {
			"tenant.read",
			"user.read",
		},
	}

	permissions, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, value := range permissions {

		if value == "*" {
			return true
		}

		if value == permission {
			return true
		}
	}

	return false
}

func validatePermissionCheck(
	request PermissionCheckRequest,
) error {

	if strings.TrimSpace(request.Role) == "" {
		return newPermissionValidationError(
			"role is required",
		)
	}

	if strings.TrimSpace(request.Permission) == "" {
		return newPermissionValidationError(
			"permission is required",
		)
	}

	return nil
}

type permissionValidationError struct {
	message string
}

func (e permissionValidationError) Error() string {
	return e.message
}

func newPermissionValidationError(
	message string,
) error {
	return permissionValidationError{
		message: message,
	}
}
