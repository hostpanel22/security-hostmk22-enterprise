package security

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

// Custom error types for file upload security
var (
	ErrFileTooLarge           = errors.New("file too large")
	ErrMIMETypeBlocked        = errors.New("MIME type blocked")
	ErrMIMETypeNotAllowed     = errors.New("MIME type not allowed")
	ErrNoFileExtension        = errors.New("no file extension")
	ErrDangerousFileExtension = errors.New("dangerous file extension")
	ErrPHPCodeDetected        = errors.New("PHP code detected")
	ErrScriptTagDetected      = errors.New("script tag detected")
	ErrEvalDetected           = errors.New("eval statement detected")
	ErrSuspiciousContent      = errors.New("suspicious content detected")
)

// RequestSizeLimit returns a middleware that limits the size of incoming requests
func RequestSizeLimit(config RequestSizeConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Limit body size
			if config.MaxBodySize > 0 {
				r.Body = http.MaxBytesReader(w, r.Body, config.MaxBodySize)
			}

			// Check content length header
			if r.ContentLength > 0 && config.MaxBodySize > 0 && r.ContentLength > config.MaxBodySize {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte(`{"error":"request body too large"}`))
				return
			}

			// Check for multipart forms
			if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				// This will be handled by the multipart parser
				// We'll add additional checks in the handler
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SimpleRequestSizeLimit returns a simple request size limit middleware
func SimpleRequestSizeLimit(maxBodySize int64) func(http.Handler) http.Handler {
	config := DefaultRequestSizeConfig()
	config.MaxBodySize = maxBodySize
	return RequestSizeLimit(config)
}

// FileUploadSecurity provides security checks for file uploads
type FileUploadSecurity struct {
	config RequestSizeConfig
}

// NewFileUploadSecurity creates a new file upload security checker
func NewFileUploadSecurity(config RequestSizeConfig) *FileUploadSecurity {
	return &FileUploadSecurity{config: config}
}

// CheckFileSize checks if a file size is within limits
func (f *FileUploadSecurity) CheckFileSize(size int64) error {
	if f.config.MaxFileSize > 0 && size > f.config.MaxFileSize {
		return ErrFileTooLarge
	}
	return nil
}

// CheckMIMEType checks if a MIME type is allowed
func (f *FileUploadSecurity) CheckMIMEType(mimeType string) error {
	mimeType = strings.ToLower(mimeType)

	// Check against blocked types
	for _, blocked := range f.config.BlockedMIMETypes {
		if strings.Contains(mimeType, blocked) {
			return ErrMIMETypeBlocked
		}
	}

	// Check against allowed types (if specified)
	if len(f.config.AllowedMIMETypes) > 0 {
		for _, allowed := range f.config.AllowedMIMETypes {
			if mimeType == allowed {
				return nil
			}
		}
		return ErrMIMETypeNotAllowed
	}

	return nil
}

// CheckFileExtension checks if a file extension is allowed
func (f *FileUploadSecurity) CheckFileExtension(filename string) error {
	// Extract extension
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ErrNoFileExtension
	}
	ext := strings.ToLower(parts[len(parts)-1])

	// Block dangerous extensions
	dangerousExtensions := map[string]bool{
		"php":     true,
		"php5":    true,
		"php7":    true,
		"phtml":   true,
		"asp":     true,
		"aspx":    true,
		"jsp":     true,
		"jsx":     true,
		"sh":      true,
		"bash":    true,
		"py":      true,
		"pl":      true,
		"cgi":     true,
		"exe":     true,
		"dll":     true,
		"so":      true,
		"bat":     true,
		"cmd":     true,
		"ps1":     true,
		"vbs":     true,
		"js":      true,
		"html":    true,
		"htm":     true,
		"svg":     true,
		"json":    true,
		"xml":     true,
		"yaml":    true,
		"yml":     true,
		"ini":     true,
		"conf":    true,
		"config":  true,
		"db":      true,
		"sqlite":  true,
		"sql":     true,
		"bak":     true,
		"backup":  true,
		"zip":     true,
		"tar":     true,
		"gz":      true,
		"rar":     true,
		"7z":      true,
	}

	if dangerousExtensions[ext] {
		return ErrDangerousFileExtension
	}

	return nil
}

// CheckFileContent checks the content of a file for malicious patterns
func (f *FileUploadSecurity) CheckFileContent(content []byte) error {
	contentStr := string(content)

	// Check for PHP tags
	if strings.Contains(contentStr, "<?php") || strings.Contains(contentStr, "<?=") {
		return ErrPHPCodeDetected
	}

	// Check for script tags
	if strings.Contains(contentStr, "<script") {
		return ErrScriptTagDetected
	}

	// Check for eval statements
	if strings.Contains(contentStr, "eval(") {
		return ErrEvalDetected
	}

	// Check for common attack patterns
	attackPatterns := []string{
		"../",
		"../",
		"<iframe",
		"<object",
		"<embed",
		"javascript:",
		"data:text/html",
		"onerror=",
		"onload=",
		"onclick=",
		"<svg onload",
		"<?",
		"<%",
		"${{",
		"{{{",
	}

	for _, pattern := range attackPatterns {
		if strings.Contains(strings.ToLower(contentStr), strings.ToLower(pattern)) {
			return ErrSuspiciousContent
		}
	}

	return nil
}

// SecureMultipartReader wraps a multipart reader with security checks
type SecureMultipartReader struct {
	reader    io.Reader
	security  *FileUploadSecurity
	maxMemory int64
}

// NewSecureMultipartReader creates a new secure multipart reader
func NewSecureMultipartReader(reader io.Reader, security *FileUploadSecurity, maxMemory int64) *SecureMultipartReader {
	return &SecureMultipartReader{
		reader:    reader,
		security:  security,
		maxMemory: maxMemory,
	}
}

// Read implements io.Reader
func (r *SecureMultipartReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// FileUploadError represents a file upload security error
type FileUploadError struct {
	Message string
}

func (e *FileUploadError) Error() string {
	return e.Message
}

// IsFileUploadError checks if an error is a file upload security error
func IsFileUploadError(err error) bool {
	_, ok := err.(*FileUploadError)
	return ok
}
