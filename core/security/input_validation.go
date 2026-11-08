package security

import (
	"net/http"
	"regexp"
	"strings"
)

// InputValidator provides methods for validating and sanitizing input
type InputValidator struct {
	config InputValidationConfig
}

// InputValidationConfig defines input validation settings
type InputValidationConfig struct {
	// MaxLength is the maximum length for string inputs
	MaxLength int
	// AllowHTML determines if HTML tags are allowed
	AllowHTML bool
	// AllowScriptTags determines if script tags are allowed
	AllowScriptTags bool
	// AllowSQLKeywords determines if SQL keywords are allowed
	AllowSQLKeywords bool
	// CustomPatterns is a list of custom regex patterns to validate against
	CustomPatterns []*regexp.Regexp
	// BlockedPatterns is a list of patterns to block
	BlockedPatterns []*regexp.Regexp
}

// DefaultInputValidationConfig returns a secure input validation configuration
func DefaultInputValidationConfig() InputValidationConfig {
	return InputValidationConfig{
		MaxLength:          10000,
		AllowHTML:          false,
		AllowScriptTags:    false,
		AllowSQLKeywords:   false,
		CustomPatterns:     []*regexp.Regexp{},
		BlockedPatterns:    []*regexp.Regexp{},
	}
}

// NewInputValidator creates a new input validator
func NewInputValidator(config InputValidationConfig) *InputValidator {
	return &InputValidator{config: config}
}

// ValidateString validates and sanitizes a string input
func (v *InputValidator) ValidateString(input string) (string, error) {
	if len(input) > v.config.MaxLength {
		return "", ErrInputTooLong
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Check for blocked patterns
	for _, pattern := range v.config.BlockedPatterns {
		if pattern.MatchString(input) {
			return "", ErrInputBlockedPattern
		}
	}

	// Check for HTML
	if !v.config.AllowHTML {
		if containsHTML(input) {
			input = sanitizeHTML(input)
		}
	}

	// Check for script tags
	if !v.config.AllowScriptTags {
		if containsScriptTags(input) {
			return "", ErrScriptTagsNotAllowed
		}
	}

	// Check for SQL keywords
	if !v.config.AllowSQLKeywords {
		if containsSQLKeywords(input) {
			return "", ErrSQLKeywordsNotAllowed
		}
	}

	// Validate against custom patterns
	for _, pattern := range v.config.CustomPatterns {
		if !pattern.MatchString(input) {
			return "", ErrInputValidationFailed
		}
	}

	return input, nil
}

// containsHTML checks if a string contains HTML tags
func containsHTML(input string) bool {
	return strings.Contains(input, "<") && strings.Contains(input, ">")
}

// sanitizeHTML removes HTML tags from a string
func sanitizeHTML(input string) string {
	// Simple HTML tag removal - for production, use a proper HTML sanitizer
	return regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
}

// containsScriptTags checks if a string contains script tags
func containsScriptTags(input string) bool {
	lower := strings.ToLower(input)
	return strings.Contains(lower, "<script") || strings.Contains(lower, "</script")
}

// containsSQLKeywords checks if a string contains SQL keywords
func containsSQLKeywords(input string) bool {
	sqlKeywords := []string{
		"select", "insert", "update", "delete", "drop", "truncate", "alter",
		"create", "exec", "execute", "union", "or 1=1", "or '1'='1'",
		"--", "/*", "*/", ";", "'", "\"", "xp_", "sp_",
	}

	lower := strings.ToLower(input)
	for _, keyword := range sqlKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// ValidateQueryParams validates query parameters
func (v *InputValidator) ValidateQueryParams(r *http.Request) error {
	for key, values := range r.URL.Query() {
		for _, value := range values {
			if _, err := v.ValidateString(value); err != nil {
				return NewErrInvalidQueryParam(key, err)
			}
		}
	}
	return nil
}

// ValidateFormValues validates form values
func (v *InputValidator) ValidateFormValues(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return ErrInvalidFormData
	}

	for key, values := range r.Form {
		for _, value := range values {
			if _, err := v.ValidateString(value); err != nil {
				return NewErrInvalidFormValue(key, err)
			}
		}
	}
	return nil
}

// ValidateJSON validates JSON body (simplified - for full validation, use a JSON schema validator)
func (v *InputValidator) ValidateJSON(body []byte) error {
	bodyStr := string(body)
	if len(bodyStr) > v.config.MaxLength {
		return ErrInputTooLong
	}

	// Check for blocked patterns in JSON
	for _, pattern := range v.config.BlockedPatterns {
		if pattern.MatchString(bodyStr) {
			return ErrInputBlockedPattern
		}
	}

	return nil
}

// InputValidationMiddleware returns a middleware that validates input
func InputValidationMiddleware(config InputValidationConfig) func(http.Handler) http.Handler {
	validator := NewInputValidator(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate query parameters
			if err := validator.ValidateQueryParams(r); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"` + err.Error() + `"}`))
				return
			}

			// Validate form values if present
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
					if err := validator.ValidateFormValues(r); err != nil {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(`{"error":"` + err.Error() + `"}`))
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SimpleInputValidation returns a simple input validation middleware
func SimpleInputValidation() func(http.Handler) http.Handler {
	return InputValidationMiddleware(DefaultInputValidationConfig())
}

// SanitizeInput sanitizes a string input by removing potentially dangerous content
func SanitizeInput(input string) string {
	// Remove HTML tags
	input = sanitizeHTML(input)
	// Remove script tags
	input = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`).ReplaceAllString(input, "")
	// Remove SQL comments
	input = regexp.MustCompile(`(?i)--.*?$|/\*.*?\*/`).ReplaceAllString(input, "")
	// Remove common attack patterns
	input = regexp.MustCompile(`(?i)(select|insert|update|delete|drop|truncate|alter|create|exec|execute|union|or 1=1|or '1'='1')`).ReplaceAllString(input, "")
	// Trim whitespace
	input = strings.TrimSpace(input)
	return input
}

// SanitizeAll sanitizes all input sources in a request
func SanitizeAll(r *http.Request) {
	// Sanitize query parameters
	for key, values := range r.URL.Query() {
		for i, value := range values {
			r.URL.Query()[key][i] = SanitizeInput(value)
		}
	}

	// Sanitize form values
	if r.Form != nil {
		for key, values := range r.Form {
			for i, value := range values {
				r.Form[key][i] = SanitizeInput(value)
			}
		}
	}
}

// SanitizeMiddleware returns a middleware that sanitizes all input
func SanitizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SanitizeAll(r)
		next.ServeHTTP(w, r)
	})
}

// Custom error types
var (
	ErrInputTooLong          = &InputValidationError{Message: "input too long"}
	ErrInputBlockedPattern   = &InputValidationError{Message: "input contains blocked pattern"}
	ErrScriptTagsNotAllowed  = &InputValidationError{Message: "script tags not allowed"}
	ErrSQLKeywordsNotAllowed = &InputValidationError{Message: "SQL keywords not allowed"}
	ErrInputValidationFailed  = &InputValidationError{Message: "input validation failed"}
	ErrInvalidFormData       = &InputValidationError{Message: "invalid form data"}
)

// InputValidationError represents an input validation error
type InputValidationError struct {
	Message string
}

func (e *InputValidationError) Error() string {
	return e.Message
}

// ErrInvalidQueryParam creates an error for invalid query parameter
type ErrInvalidQueryParam struct {
	param string
	err   error
}

func (e *ErrInvalidQueryParam) Error() string {
	return "invalid query parameter '" + e.param + "': " + e.err.Error()
}

// NewErrInvalidQueryParam creates a new ErrInvalidQueryParam
func NewErrInvalidQueryParam(param string, err error) *ErrInvalidQueryParam {
	return &ErrInvalidQueryParam{param: param, err: err}
}

// ErrInvalidFormValue creates an error for invalid form value
type ErrInvalidFormValue struct {
	field string
	err   error
}

func (e *ErrInvalidFormValue) Error() string {
	return "invalid form value for '" + e.field + "': " + e.err.Error()
}

// NewErrInvalidFormValue creates a new ErrInvalidFormValue
func NewErrInvalidFormValue(field string, err error) *ErrInvalidFormValue {
	return &ErrInvalidFormValue{field: field, err: err}
}
