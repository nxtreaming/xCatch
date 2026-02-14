package utools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/xCatch/xcatch/config"
)

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	cfg := &config.Config{
		BaseURL:    baseURL,
		APIKey:     "test-key",
		Timeout:    5 * time.Second,
		MaxRetries: 2,
		RateLimit:  100,
	}
	c, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c
}

func TestGetRawReturnsHTTPBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("apiKey"); got != "test-key" {
			t.Fatalf("missing apiKey, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"data":"{\"hello\":\"world\"}","msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)

	raw, err := c.GetRaw(context.Background(), "/raw", nil)
	if err != nil {
		t.Fatalf("GetRaw error: %v", err)
	}
	if string(raw) == "" || string(raw)[0] != '{' || !json.Valid(raw) {
		t.Fatalf("GetRaw should return raw JSON body, got: %q", string(raw))
	}
	if !strings.Contains(string(raw), `"code":1`) {
		t.Fatalf("GetRaw should include envelope fields, got: %s", string(raw))
	}

	var parsed map[string]string
	if err := c.Get(context.Background(), "/raw", nil, &parsed); err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if parsed["hello"] != "world" {
		t.Fatalf("expected unwrapped data, got %+v", parsed)
	}
}

func TestDoWithRetryDoesNotRetryUnmarshalError(t *testing.T) {
	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"data":"{bad}","msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	var result map[string]any
	err := c.Get(context.Background(), "/bad", nil, &result)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected no retry for deterministic unmarshal error, hits=%d", got)
	}
}

func TestDoWithRetryRetriesOnRateLimitThenSucceeds(t *testing.T) {
	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		if current == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":88,"msg":"rate limit"}`))
			return
		}
		_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	var result map[string]bool
	if err := c.Get(context.Background(), "/retry", nil, &result); err != nil {
		t.Fatalf("expected retry success, got error: %v", err)
	}
	if !result["ok"] {
		t.Fatalf("expected parsed result after retry, got %+v", result)
	}
	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected exactly one retry, hits=%d", got)
	}
}
