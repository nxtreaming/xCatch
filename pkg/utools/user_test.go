package utools

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserEndpoints_RequestMapping(t *testing.T) {
	type tc struct {
		name          string
		expectedPath  string
		expectedQuery map[string]string
		call          func(c *Client) (json.RawMessage, error)
	}

	cases := []tc{
		{
			name:         "GetUserByScreenName",
			expectedPath: "/api/base/apitools/getUserByIdOrNameShow",
			expectedQuery: map[string]string{
				"screenName": "elonmusk",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserByScreenName(context.Background(), "elonmusk")
			},
		},
		{
			name:         "GetUserByID",
			expectedPath: "/api/base/apitools/usersByIdRestIds",
			expectedQuery: map[string]string{
				"userIds": "123",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserByID(context.Background(), "123")
			},
		},
		{
			name:         "GetUsersByIDs",
			expectedPath: "/api/base/apitools/usersByIdRestIds",
			expectedQuery: map[string]string{
				"userIds": "1,2,3",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUsersByIDs(context.Background(), []string{"1", "2", "3"})
			},
		},
		{
			name:         "GetUsernameChanges",
			expectedPath: "/api/base/apitools/usernameChanges",
			expectedQuery: map[string]string{
				"userId": "456",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUsernameChanges(context.Background(), "456")
			},
		},
		{
			name:         "LookupUser with both params",
			expectedPath: "/api/base/apitools/getUserByIdOrNameLookup",
			expectedQuery: map[string]string{
				"screenName": "jack",
				"userId":     "12",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.LookupUser(context.Background(), "jack", "12")
			},
		},
		{
			name:         "GetUserByScreenNameV2",
			expectedPath: "/api/base/apitools/userByScreenNameV2",
			expectedQuery: map[string]string{
				"screenName": "xcatch",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserByScreenNameV2(context.Background(), "xcatch")
			},
		},
		{
			name:         "GetUserByIDV2",
			expectedPath: "/api/base/apitools/uerByIdRestIdV2",
			expectedQuery: map[string]string{
				"userId": "999",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUserByIDV2(context.Background(), "999")
			},
		},
		{
			name:         "GetUsersByIDsV2",
			expectedPath: "/api/base/apitools/usersByIdRestIds",
			expectedQuery: map[string]string{
				"userIds": "11,22",
			},
			call: func(c *Client) (json.RawMessage, error) {
				return c.GetUsersByIDsV2(context.Background(), []string{"11", "22"})
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

func TestLookupUser_OptionalParams(t *testing.T) {
	t.Run("screenName only", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("screenName") != "jack" {
				t.Fatalf("expected screenName=jack, got %q", q.Get("screenName"))
			}
			if q.Get("userId") != "" {
				t.Fatalf("expected empty userId, got %q", q.Get("userId"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
		}))
		defer ts.Close()

		client := newTestClient(t, ts.URL)
		if _, err := client.LookupUser(context.Background(), "jack", ""); err != nil {
			t.Fatalf("LookupUser error: %v", err)
		}
	})

	t.Run("userId only", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("userId") != "42" {
				t.Fatalf("expected userId=42, got %q", q.Get("userId"))
			}
			if q.Get("screenName") != "" {
				t.Fatalf("expected empty screenName, got %q", q.Get("screenName"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
		}))
		defer ts.Close()

		client := newTestClient(t, ts.URL)
		if _, err := client.LookupUser(context.Background(), "", "42"); err != nil {
			t.Fatalf("LookupUser error: %v", err)
		}
	})
}

func TestGetAccountAnalytics_AuthRequired(t *testing.T) {
	client := newTestClient(t, "http://127.0.0.1:0")
	_, err := client.GetAccountAnalytics(context.Background())
	if !errors.Is(err, ErrAuthTokenRequired) {
		t.Fatalf("expected ErrAuthTokenRequired, got %v", err)
	}
}

func TestGetAccountAnalytics_PassesAuthTokenAndCT0(t *testing.T) {
	t.Run("auth_token only", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if r.URL.Path != "/api/base/apitools/accountAnalytics" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			if q.Get("auth_token") != "auth-only" {
				t.Fatalf("expected auth_token=auth-only, got %q", q.Get("auth_token"))
			}
			if q.Get("ct0") != "" {
				t.Fatalf("expected empty ct0, got %q", q.Get("ct0"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
		}))
		defer ts.Close()

		client := newTestClient(t, ts.URL)
		client.authToken = "auth-only"
		if _, err := client.GetAccountAnalytics(context.Background()); err != nil {
			t.Fatalf("GetAccountAnalytics error: %v", err)
		}
	})

	t.Run("auth_token + ct0", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("auth_token") != "auth-token" {
				t.Fatalf("expected auth_token=auth-token, got %q", q.Get("auth_token"))
			}
			if q.Get("ct0") != "ct0-token" {
				t.Fatalf("expected ct0=ct0-token, got %q", q.Get("ct0"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"data":{"ok":true},"msg":"SUCCESS"}`))
		}))
		defer ts.Close()

		client := newTestClient(t, ts.URL)
		client.authToken = "auth-token"
		client.ct0 = "ct0-token"
		if _, err := client.GetAccountAnalytics(context.Background()); err != nil {
			t.Fatalf("GetAccountAnalytics error: %v", err)
		}
	})
}
