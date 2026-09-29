package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeout  = 30 * time.Second
	maxResponseSize = 8 << 20
	userAgent       = "hostmk-go-sdk/1.0"
)

var ErrInvalidBaseURL = errors.New("invalid API base URL")

type Option func(*Client)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	mu         sync.RWMutex
	apiKey     string
	token      string
}

type APIResponse struct {
	Success   bool            `json:"success"`
	Message   string          `json:"message,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Error     string          `json:"error,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

type ListResult[T any] struct {
	Count int
	Items []T
}

type Error struct {
	StatusCode int
	Status     string
	Message    string
	Body       []byte
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("hostmk API: %s: %s", e.Status, e.Message)
	}
	return fmt.Sprintf("hostmk API: %s", e.Status)
}

func WithAPIKey(key string) Option {
	return func(c *Client) { c.apiKey = strings.TrimSpace(key) }
}

func WithBearerToken(token string) Option {
	return func(c *Client) { c.token = strings.TrimSpace(token) }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

func New(baseURL string, options ...Option) (*Client, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsed == nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, ErrInvalidBaseURL
	}

	c := &Client{
		baseURL:    parsed,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	return c, nil
}

func (c *Client) SetAPIKey(key string) {
	c.mu.Lock()
	c.apiKey = strings.TrimSpace(key)
	c.mu.Unlock()
}

func (c *Client) SetBearerToken(token string) {
	c.mu.Lock()
	c.token = strings.TrimSpace(token)
	c.mu.Unlock()
}

func (c *Client) BaseURL() string {
	return strings.TrimRight(c.baseURL.String(), "/")
}

func (c *Client) Request(ctx context.Context, method, path string, query url.Values, body any) ([]byte, error) {
	if c == nil || c.baseURL == nil {
		return nil, ErrInvalidBaseURL
	}
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return nil, fmt.Errorf("API path must be absolute: %q", path)
	}

	target := *c.baseURL
	basePath := strings.TrimRight(target.Path, "/")
	requestPath := path
	if strings.HasSuffix(basePath, "/api/v1") && (requestPath == "/api/v1" || strings.HasPrefix(requestPath, "/api/v1/")) {
		target.Path = basePath + strings.TrimPrefix(requestPath, "/api/v1")
	} else if basePath != "" && (requestPath == basePath || strings.HasPrefix(requestPath, basePath+"/")) {
		target.Path = requestPath
	} else {
		target.Path = basePath + requestPath
	}
	target.RawQuery = ""
	if query != nil {
		target.RawQuery = query.Encode()
	}

	var requestBody io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
		requestBody = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(ctx, method, target.String(), requestBody)
	if err != nil {
		return nil, fmt.Errorf("create API request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", userAgent)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	c.mu.RLock()
	apiKey, token := c.apiKey, c.token
	c.mu.RUnlock()
	if apiKey != "" {
		request.Header.Set("X-API-Key", apiKey)
	} else if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send API request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, maxResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("read API response: %w", err)
	}
	if len(responseBody) > maxResponseSize {
		return nil, fmt.Errorf("API response exceeds %d bytes", maxResponseSize)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, newAPIError(response, responseBody)
	}
	return responseBody, nil
}

func (c *Client) Get(ctx context.Context, path string, query url.Values, out any) error {
	return c.requestJSON(ctx, http.MethodGet, path, query, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.requestJSON(ctx, http.MethodPost, path, query, body, out)
}

func (c *Client) Put(ctx context.Context, path string, query url.Values, body, out any) error {
	return c.requestJSON(ctx, http.MethodPut, path, query, body, out)
}

func (c *Client) Delete(ctx context.Context, path string, query url.Values, out any) error {
	return c.requestJSON(ctx, http.MethodDelete, path, query, nil, out)
}

func (c *Client) requestJSON(ctx context.Context, method, path string, query url.Values, body, out any) error {
	responseBody, err := c.Request(ctx, method, path, query, body)
	if err != nil || out == nil || len(responseBody) == 0 {
		return err
	}
	return DecodeData(responseBody, out)
}

func (c *Client) GetList(ctx context.Context, path string, query url.Values, key string, out any) (int, error) {
	responseBody, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return 0, err
	}
	return DecodeList(responseBody, key, out)
}

func DecodeData(responseBody []byte, out any) error {
	if out == nil || len(responseBody) == 0 {
		return nil
	}
	var envelope APIResponse
	if err := json.Unmarshal(responseBody, &envelope); err == nil && len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, out); err != nil {
			return fmt.Errorf("decode API response data: %w", err)
		}
		return nil
	}
	if err := json.Unmarshal(responseBody, out); err != nil {
		return fmt.Errorf("decode API response: %w", err)
	}
	return nil
}

func DecodeList(responseBody []byte, key string, out any) (int, error) {
	if out == nil {
		return 0, errors.New("list output must not be nil")
	}
	var envelope APIResponse
	if err := json.Unmarshal(responseBody, &envelope); err != nil {
		return 0, fmt.Errorf("decode API list envelope: %w", err)
	}
	data := envelope.Data
	if len(data) == 0 {
		return 0, errors.New("API response has no data field")
	}
	var payload struct {
		Count int             `json:"count"`
		Items json.RawMessage `json:"items"`
	}
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal(data, &payload); err != nil {
			return 0, fmt.Errorf("decode API list metadata: %w", err)
		}
		if key != "" {
			var keyed map[string]json.RawMessage
			if err := json.Unmarshal(data, &keyed); err != nil {
				return 0, fmt.Errorf("decode API list data: %w", err)
			}
			if value, ok := keyed[key]; ok {
				payload.Items = value
			}
		}
	} else {
		payload.Items = data
	}
	if len(payload.Items) == 0 || string(payload.Items) == "null" {
		payload.Items = []byte("[]")
	}
	if err := json.Unmarshal(payload.Items, out); err != nil {
		return 0, fmt.Errorf("decode API list items: %w", err)
	}
	if payload.Count == 0 {
		var list []json.RawMessage
		if err := json.Unmarshal(payload.Items, &list); err == nil {
			payload.Count = len(list)
		}
	}
	return payload.Count, nil
}

func newAPIError(response *http.Response, body []byte) *Error {
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
		Detail  string `json:"detail"`
	}
	_ = json.Unmarshal(body, &payload)
	message := payload.Message
	if message == "" {
		message = payload.Error
	}
	if message == "" {
		message = payload.Detail
	}
	if message == "" {
		message = strings.TrimSpace(string(body))
	}
	return &Error{StatusCode: response.StatusCode, Status: response.Status, Message: message, Body: append([]byte(nil), body...)}
}
