package utools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xCatch/xcatch/config"
)

func TestAuthEndpointsPassCT0WhenConfigured(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("apiKey") != "test-key" {
			t.Fatalf("expected apiKey, got %q", q.Get("apiKey"))
		}
		if q.Get("auth_token") != "test-auth" {
			t.Fatalf("expected auth_token, got %q", q.Get("auth_token"))
		}
		if q.Get("ct0") != "test-ct0" {
			t.Fatalf("expected ct0, got %q", q.Get("ct0"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	cfg := &config.Config{
		BaseURL:    ts.URL,
		APIKey:     "test-key",
		AuthToken:  "test-auth",
		CT0:        "test-ct0",
		Timeout:    5 * time.Second,
		MaxRetries: 0,
		RateLimit:  100,
	}
	c, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	if _, err := c.GetHomeTimeline(context.Background(), ""); err != nil {
		t.Fatalf("GetHomeTimeline error: %v", err)
	}
	if _, err := c.GetMentionsTimeline(context.Background(), ""); err != nil {
		t.Fatalf("GetMentionsTimeline error: %v", err)
	}
}

func TestGetTrendsSendsIDCompatibilityParam(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("id") != "1" {
			t.Fatalf("expected id=1, got %q", q.Get("id"))
		}
		if q.Get("woeid") != "1" {
			t.Fatalf("expected woeid=1, got %q", q.Get("woeid"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	if _, err := c.GetTrends(context.Background(), "1"); err != nil {
		t.Fatalf("GetTrends error: %v", err)
	}
}

func TestSearchWithOptionsMapsGeneratedParams(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("words") != "bitcoin" {
			t.Fatalf("expected words, got %q", q.Get("words"))
		}
		if q.Get("type") != "Latest" {
			t.Fatalf("expected type, got %q", q.Get("type"))
		}
		if q.Get("from") != "elonmusk" || q.Get("until") != "2026-01-01" {
			t.Fatalf("expected advanced filters, got from=%q until=%q", q.Get("from"), q.Get("until"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	var result json.RawMessage
	result, err := c.SearchWithOptions(context.Background(), "bitcoin", SearchOptions{
		Type:  "Latest",
		From:  "elonmusk",
		Until: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("SearchWithOptions error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}
}
