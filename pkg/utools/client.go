package utools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/xCatch/xcatch/config"
)

// Client is the uTools API HTTP client with built-in auth, retry, and rate limiting.
type Client struct {
	baseURL    string
	apiKey     string
	authToken  string
	ct0        string
	httpClient *http.Client
	maxRetries int
	limiter    *rate.Limiter
}

// NewClient creates a new uTools API client from the given config.
func NewClient(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Client{
		baseURL:   strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:    cfg.APIKey,
		authToken: cfg.AuthToken,
		ct0:       cfg.CT0,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		maxRetries: cfg.MaxRetries,
		limiter:    rate.NewLimiter(rate.Limit(cfg.RateLimit), 1),
	}, nil
}

// Get performs a GET request to the given API path with query parameters.
// The response JSON is unmarshalled into result.
func (c *Client) Get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	return c.doWithRetry(ctx, http.MethodGet, path, params, result)
}

// Post performs a POST request to the given API path with form parameters.
// The response JSON is unmarshalled into result.
func (c *Client) Post(ctx context.Context, path string, params map[string]string, result interface{}) error {
	return c.doWithRetry(ctx, http.MethodPost, path, params, result)
}

// GetRaw performs a GET request and returns the raw response body bytes.
func (c *Client) GetRaw(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return c.doRawWithRetry(ctx, http.MethodGet, path, params)
}

func (c *Client) doWithRetry(ctx context.Context, method, path string, params map[string]string, result interface{}) error {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			log.Printf("[utools] retry %d/%d for %s %s (backoff %v)", attempt, c.maxRetries, method, path, backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Wait for rate limiter
		if err := c.limiter.Wait(ctx); err != nil {
			return fmt.Errorf("utools: rate limiter: %w", err)
		}

		lastErr = c.do(ctx, method, path, params, result)
		if lastErr == nil {
			return nil
		}

		if !isRetryableError(lastErr) {
			return lastErr
		}
	}
	return lastErr
}

func (c *Client) doRawWithRetry(ctx context.Context, method, path string, params map[string]string) ([]byte, error) {
	var (
		lastErr error
		body    []byte
	)

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
			log.Printf("[utools] retry %d/%d for %s %s (backoff %v)", attempt, c.maxRetries, method, path, backoff)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		if err := c.limiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("utools: rate limiter: %w", err)
		}

		body, lastErr = c.doRaw(ctx, method, path, params)
		if lastErr == nil {
			return body, nil
		}

		if !isRetryableError(lastErr) {
			return nil, lastErr
		}
	}

	return nil, lastErr
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return false
}

func (c *Client) doRaw(ctx context.Context, method, path string, params map[string]string) ([]byte, error) {
	reqURL := c.baseURL + path

	merged := make(map[string]string, len(params)+1)
	for k, v := range params {
		merged[k] = v
	}
	merged["apiKey"] = c.apiKey

	var req *http.Request
	var err error

	switch method {
	case http.MethodGet:
		u, parseErr := url.Parse(reqURL)
		if parseErr != nil {
			return nil, fmt.Errorf("utools: parse url: %w", parseErr)
		}
		q := u.Query()
		for k, v := range merged {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)

	case http.MethodPost:
		form := url.Values{}
		for k, v := range merged {
			form.Set(k, v)
		}
		req, err = http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(form.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

	default:
		return nil, fmt.Errorf("utools: unsupported method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("utools: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("utools: http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("utools: read body: %w", err)
	}

	if resetStr := resp.Header.Get("x-rate-limit-reset"); resetStr != "" {
		if resetVal, parseErr := strconv.Atoi(resetStr); parseErr == nil && resetVal < 9 {
			log.Printf("[utools] x-rate-limit-reset=%d, consider calling tokenSync", resetVal)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(body),
		}
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Msg     string `json:"msg"`
		}
		if json.Unmarshal(body, &errResp) == nil {
			apiErr.Code = errResp.Code
			if errResp.Message != "" {
				apiErr.Message = errResp.Message
			} else {
				apiErr.Message = errResp.Msg
			}
		}
		if apiErr.Message == "" {
			apiErr.Message = string(body)
		}
		return nil, apiErr
	}

	return body, nil
}

func (c *Client) do(ctx context.Context, method, path string, params map[string]string, result interface{}) error {
	// Build URL
	reqURL := c.baseURL + path

	// Copy params to avoid mutating the caller's map, and inject apiKey
	merged := make(map[string]string, len(params)+1)
	for k, v := range params {
		merged[k] = v
	}
	merged["apiKey"] = c.apiKey

	var req *http.Request
	var err error

	switch method {
	case http.MethodGet:
		u, parseErr := url.Parse(reqURL)
		if parseErr != nil {
			return fmt.Errorf("utools: parse url: %w", parseErr)
		}
		q := u.Query()
		for k, v := range merged {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)

	case http.MethodPost:
		form := url.Values{}
		for k, v := range merged {
			form.Set(k, v)
		}
		req, err = http.NewRequestWithContext(ctx, method, reqURL, strings.NewReader(form.Encode()))
		if err == nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

	default:
		return fmt.Errorf("utools: unsupported method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("utools: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("utools: http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("utools: read body: %w", err)
	}

	// Check x-rate-limit-reset header
	if resetStr := resp.Header.Get("x-rate-limit-reset"); resetStr != "" {
		if resetVal, parseErr := strconv.Atoi(resetStr); parseErr == nil && resetVal < 9 {
			log.Printf("[utools] x-rate-limit-reset=%d, consider calling tokenSync", resetVal)
		}
	}

	// Handle non-2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(body),
		}
		// Try to parse error details from body
		var errResp struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Msg     string `json:"msg"`
		}
		if json.Unmarshal(body, &errResp) == nil {
			apiErr.Code = errResp.Code
			if errResp.Message != "" {
				apiErr.Message = errResp.Message
			} else {
				apiErr.Message = errResp.Msg
			}
		}
		if apiErr.Message == "" {
			apiErr.Message = string(body)
		}
		return apiErr
	}

	// Unwrap the API envelope: {"code":1, "data":"<json_string>", "msg":"SUCCESS"}
	// The "data" field is a JSON-encoded string that needs double-unmarshal.
	if result != nil {
		var envelope struct {
			Code   int             `json:"code"`
			Data   json.RawMessage `json:"data"`
			Msg    string          `json:"msg"`
			Result json.RawMessage `json:"result"`
		}
		if err := json.Unmarshal(body, &envelope); err == nil && (len(envelope.Data) > 0 || envelope.Code != 0) {
			// Check for business-level errors (code != 1 means failure)
			if envelope.Code != 0 && envelope.Code != 1 {
				return &APIError{
					StatusCode: resp.StatusCode,
					Code:       envelope.Code,
					Message:    envelope.Msg,
					RawBody:    string(body),
				}
			}

			// Check if data is a JSON string (starts with `"`)
			if len(envelope.Data) > 0 && envelope.Data[0] == '"' {
				var dataStr string
				if err := json.Unmarshal(envelope.Data, &dataStr); err == nil {
					// dataStr is the inner JSON — unmarshal it into result
					if err := json.Unmarshal([]byte(dataStr), result); err != nil {
						return fmt.Errorf("utools: unmarshal inner data: %w (data: %s)", err, Truncate(dataStr, 500))
					}
					return nil
				}
			}
			// data is already a JSON object/array, use it directly
			if err := json.Unmarshal(envelope.Data, result); err != nil {
				return fmt.Errorf("utools: unmarshal data field: %w (data: %s)", err, Truncate(string(envelope.Data), 500))
			}
			return nil
		}

		// Fallback: no envelope, unmarshal the whole body
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("utools: unmarshal response: %w (body: %s)", err, Truncate(string(body), 500))
		}
	}

	return nil
}

// TokenSync calls the tokenSync endpoint to refresh the robot token.
// Should be called when x-rate-limit-reset < 9 or persistent errors occur.
func (c *Client) TokenSync(ctx context.Context) error {
	params := map[string]string{}
	var result json.RawMessage
	return c.Get(ctx, "/api/base/apitools/tokenSync", params, &result)
}

// Truncate shortens a string to maxLen characters, appending "..." if truncated.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
