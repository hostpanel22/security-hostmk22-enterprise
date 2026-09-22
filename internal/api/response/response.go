package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func JSON(
	w http.ResponseWriter,
	status int,
	message string,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Success:   status < http.StatusBadRequest,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func Error(
	w http.ResponseWriter,
	status int,
	err string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Success:   false,
		Error:     err,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func OK(
	w http.ResponseWriter,
	message string,
	data interface{},
) {
	JSON(
		w,
		http.StatusOK,
		message,
		data,
	)
}

func Created(
	w http.ResponseWriter,
	message string,
	data interface{},
) {
	JSON(
		w,
		http.StatusCreated,
		message,
		data,
	)
}

func BadRequest(
	w http.ResponseWriter,
	err string,
) {
	Error(
		w,
		http.StatusBadRequest,
		err,
	)
}

func Unauthorized(
	w http.ResponseWriter,
	err string,
) {
	Error(
		w,
		http.StatusUnauthorized,
		err,
	)
}

func Forbidden(
	w http.ResponseWriter,
	err string,
) {
	Error(
		w,
		http.StatusForbidden,
		err,
	)
}

func NotFound(
	w http.ResponseWriter,
	err string,
) {
	Error(
		w,
		http.StatusNotFound,
		err,
	)
}

func InternalServerError(
	w http.ResponseWriter,
	err string,
) {
	Error(
		w,
		http.StatusInternalServerError,
		err,
	)
}
