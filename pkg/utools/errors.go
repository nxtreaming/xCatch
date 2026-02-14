package utools

import "fmt"

// APIError represents an error returned by the uTools API.
type APIError struct {
	StatusCode int
	Code       int    // Twitter error code (e.g. 88 = rate limit)
	Message    string
	RawBody    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("utools: HTTP %d, code=%d, message=%s", e.StatusCode, e.Code, e.Message)
}

// IsRateLimited returns true if the error is a rate limit error.
func (e *APIError) IsRateLimited() bool {
	return e.Code == 88 || e.StatusCode == 429
}

// IsForbidden returns true if the error is a 403 Forbidden.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsUnauthorized returns true if the error is a 401 Unauthorized.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsRetryable returns true if the request should be retried.
func (e *APIError) IsRetryable() bool {
	return e.IsRateLimited() || e.IsForbidden()
}
