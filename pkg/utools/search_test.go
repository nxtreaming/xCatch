package utools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchEndpoints_RequestMapping(t *testing.T) {
	type tc struct {
		name          string
		expectedPath  string
		expectedQuery map[string]string
		call          func(c *Client) (json.RawMessage, error)
	}

	cases := []tc{
		{
			name:         "Search",
			expectedPath: "/api/base/apitools/search",
			expectedQuery: map[string]string{
				"words":  "bitcoin",
				"type":   "Latest",
				"cursor": "cur-1",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.Search(context.Background(), "bitcoin", "Latest", "cur-1")
			},
		},
		{
			name:         "SearchWithOptions",
			expectedPath: "/api/base/apitools/search",
			expectedQuery: map[string]string{
				"words":      "bitcoin",
				"type":       "Top",
				"cursor":     "cur-2",
				"any":        "btc,eth",
				"from":       "elonmusk",
				"likes":      "100",
				"mentioning": "jack",
				"none":       "spam",
				"phrase":     "hello world",
				"replies":    "10",
				"retweets":   "20",
				"since":      "2025-01-01",
				"tag":        "crypto",
				"to":         "x",
				"until":      "2025-12-31",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.SearchWithOptions(context.Background(), "bitcoin", SearchOptions{
					Type:       "Top",
					Cursor:     "cur-2",
					Any:        "btc,eth",
					From:       "elonmusk",
					Likes:      "100",
					Mentioning: "jack",
					None:       "spam",
					Phrase:     "hello world",
					Replies:    "10",
					Retweets:   "20",
					Since:      "2025-01-01",
					Tag:        "crypto",
					To:         "x",
					Until:      "2025-12-31",
				})
			},
		},
		{
			name:         "SearchBox",
			expectedPath: "/api/base/apitools/searchBox",
			expectedQuery: map[string]string{
				"words": "elon",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.SearchBox(context.Background(), "elon")
			},
		},
		{
			name:         "GetTrends",
			expectedPath: "/api/base/apitools/trends",
			expectedQuery: map[string]string{
				"id":    "1",
				"woeid": "1",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetTrends(context.Background(), "1")
			},
		},
		{
			name:         "GetTrending",
			expectedPath: "/api/base/apitools/trending",
			expectedQuery: map[string]string{},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetTrending(context.Background())
			},
		},
		{
			name:         "GetNews",
			expectedPath: "/api/base/apitools/news",
			expectedQuery: map[string]string{},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetNews(context.Background())
			},
		},
		{
			name:         "GetExplorePage",
			expectedPath: "/api/base/apitools/explore",
			expectedQuery: map[string]string{},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetExplorePage(context.Background())
			},
		},
		{
			name:         "GetSports",
			expectedPath: "/api/base/apitools/sports",
			expectedQuery: map[string]string{},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetSports(context.Background())
			},
		},
		{
			name:         "GetEntertainment",
			expectedPath: "/api/base/apitools/entertainment",
			expectedQuery: map[string]string{},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetEntertainment(context.Background())
			},
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != cse.expectedPath {
					t.Fatalf("path mismatch: got %s want %s", r.URL.Path, cse.expectedPath)
				}
				q := r.URL.Query()
				if got := q.Get("apiKey"); got != "test-key" {
					t.Fatalf("missing apiKey, got %q", got)
				}
				for k, want := range cse.expectedQuery {
					if got := q.Get(k); got != want {
						t.Fatalf("query[%s] mismatch: got %q want %q", k, got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
			}))
			defer ts.Close()

			client := newTestClient(t, ts.URL)
			raw, err := cse.call(client)
			if err != nil {
				t.Fatalf("call returned error: %v", err)
			}
			if !json.Valid(raw) {
				t.Fatalf("expected valid JSON, got %s", string(raw))
			}
		})
	}
}
