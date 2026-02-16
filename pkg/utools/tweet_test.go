package utools

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTweetEndpoints_RequestMapping(t *testing.T) {
	type tc struct {
		name          string
		expectedPath  string
		expectedQuery map[string]string
		prepare       func(c *Client)
		call          func(c *Client) (json.RawMessage, error)
	}

	cases := []tc{
		{
			name:         "GetUserTweets",
			expectedPath: "/api/base/apitools/userTweetsV2",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-1",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserTweets(context.Background(), "123", "cur-1")
			},
		},
		{
			name:         "GetUserTimeline",
			expectedPath: "/api/base/apitools/userTimeline",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-2",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserTimeline(context.Background(), "123", "cur-2")
			},
		},
		{
			name:         "GetTweetDetail",
			expectedPath: "/api/base/apitools/tweetTimeline",
			expectedQuery: map[string]string{
				"tweetId": "456",
				"cursor":  "cur-3",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetTweetDetail(context.Background(), "456", "cur-3")
			},
		},
		{
			name:         "GetTweetSimple",
			expectedPath: "/api/base/apitools/tweetSimple",
			expectedQuery: map[string]string{
				"tweetId": "456",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetTweetSimple(context.Background(), "456")
			},
		},
		{
			name:         "GetTweetsByIDs",
			expectedPath: "/api/base/apitools/tweetResultsByRestIds",
			expectedQuery: map[string]string{
				"tweetIds": "11,22",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetTweetsByIDs(context.Background(), []string{"11", "22"})
			},
		},
		{
			name:         "GetUserReplies",
			expectedPath: "/api/base/apitools/userTweetReply",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-4",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserReplies(context.Background(), "123", "cur-4")
			},
		},
		{
			name:         "GetUserLikes",
			expectedPath: "/api/base/apitools/favoritesList",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-5",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserLikes(context.Background(), "123", "cur-5")
			},
		},
		{
			name:         "GetUserLikesV2",
			expectedPath: "/api/base/apitools/userLikeV2",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-6",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserLikesV2(context.Background(), "123", "cur-6")
			},
		},
		{
			name:         "GetUserHighlights",
			expectedPath: "/api/base/apitools/highlightsV2",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-7",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserHighlights(context.Background(), "123", "cur-7")
			},
		},
		{
			name:         "GetUserArticlesTweets",
			expectedPath: "/api/base/apitools/userArticlesTweets",
			expectedQuery: map[string]string{
				"userId": "123",
				"cursor": "cur-8",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserArticlesTweets(context.Background(), "123", "cur-8")
			},
		},
		{
			name:         "GetRetweeters",
			expectedPath: "/api/base/apitools/retweetersV2",
			expectedQuery: map[string]string{
				"tweetId": "456",
				"cursor":  "cur-9",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetRetweeters(context.Background(), "456", "cur-9")
			},
		},
		{
			name:         "GetRetweetersIDs",
			expectedPath: "/api/base/apitools/retweetersIds",
			expectedQuery: map[string]string{
				"tweetId": "456",
				"cursor":  "cur-10",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetRetweetersIDs(context.Background(), "456", "cur-10")
			},
		},
		{
			name:         "GetFavoriters",
			expectedPath: "/api/base/apitools/favoritersV2",
			expectedQuery: map[string]string{
				"tweetId":    "456",
				"cursor":     "cur-11",
				"auth_token": "auth-token",
			},
			prepare: func(c *Client) {
				c.authToken = "auth-token"
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetFavoriters(context.Background(), "456", "cur-11")
			},
		},
		{
			name:         "GetQuotes",
			expectedPath: "/api/base/apitools/quotesV2",
			expectedQuery: map[string]string{
				"tweetId": "456",
				"cursor":  "cur-12",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetQuotes(context.Background(), "456", "cur-12")
			},
		},
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != cse.expectedPath {
					t.Fatalf("path mismatch: got %s want %s", r.URL.Path, cse.expectedPath)
				}
				if got := r.URL.Query().Get("apiKey"); got != "test-key" {
					t.Fatalf("missing apiKey, got %q", got)
				}
				for k, want := range cse.expectedQuery {
					if got := r.URL.Query().Get(k); got != want {
						t.Fatalf("query[%s] mismatch: got %q want %q", k, got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
			}))
			defer ts.Close()

			client := newTestClient(t, ts.URL)
			if cse.prepare != nil {
				cse.prepare(client)
			}
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

func TestTweetTimelines_AuthRequired(t *testing.T) {
	client := newTestClient(t, "http://127.0.0.1:0")

	if _, err := client.GetHomeTimeline(context.Background(), ""); !errors.Is(err, ErrAuthTokenRequired) {
		t.Fatalf("GetHomeTimeline expected ErrAuthTokenRequired, got %v", err)
	}
	if _, err := client.GetMentionsTimeline(context.Background(), ""); !errors.Is(err, ErrAuthTokenRequired) {
		t.Fatalf("GetMentionsTimeline expected ErrAuthTokenRequired, got %v", err)
	}
}

func TestTweetTimelines_PassesAuthTokenAndCT0(t *testing.T) {
	type tc struct {
		name         string
		expectedPath string
		call         func(c *Client) (json.RawMessage, error)
	}

	cases := []tc{
		{
			name:         "GetHomeTimeline",
			expectedPath: "/api/base/apitools/homeTimeline",
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetHomeTimeline(context.Background(), "cur-home")
			},
		},
		{
			name:         "GetMentionsTimeline",
			expectedPath: "/api/base/apitools/mentionsTimeline",
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetMentionsTimeline(context.Background(), "cur-mentions")
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
				if got := q.Get("auth_token"); got != "auth-token" {
					t.Fatalf("expected auth_token=auth-token, got %q", got)
				}
				if got := q.Get("ct0"); got != "ct0-token" {
					t.Fatalf("expected ct0=ct0-token, got %q", got)
				}
				if got := q.Get("apiKey"); got != "test-key" {
					t.Fatalf("missing apiKey, got %q", got)
				}
				if got := q.Get("cursor"); got == "" {
					t.Fatalf("expected cursor query")
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
			}))
			defer ts.Close()

			client := newTestClient(t, ts.URL)
			client.authToken = "auth-token"
			client.ct0 = "ct0-token"

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
