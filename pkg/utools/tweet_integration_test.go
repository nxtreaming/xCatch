//go:build integration

package utools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func isLikelyTweetID(v string) bool {
	v = strings.TrimSpace(v)
	if len(v) < 5 {
		return false
	}
	for _, r := range v {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func extractStringByKeys(node any, keys []string) (string, bool) {
	switch x := node.(type) {
	case map[string]any:
		for _, k := range keys {
			if raw, ok := x[k]; ok {
				if s, ok := raw.(string); ok && isLikelyTweetID(s) {
					return s, true
				}
			}
		}
		for _, v := range x {
			if got, ok := extractStringByKeys(v, keys); ok {
				return got, true
			}
		}
	case []any:
		for _, v := range x {
			if got, ok := extractStringByKeys(v, keys); ok {
				return got, true
			}
		}
	}
	return "", false
}

func findTweetIDFromRaw(raw json.RawMessage) string {
	if len(raw) == 0 || !json.Valid(raw) {
		return ""
	}

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}

	keys := []string{"tweetId", "tweet_id", "rest_id", "id_str", "id"}
	if id, ok := extractStringByKeys(payload, keys); ok {
		return id
	}
	return ""
}

func discoverTweetID(ctx context.Context, client *Client, userID string) string {
	type source struct {
		name string
		call func() (json.RawMessage, error)
	}

	sources := []source{
		{
			name: "GetUserTweets",
			call: func() (json.RawMessage, error) { return client.GetUserTweets(ctx, userID, "") },
		},
		{
			name: "GetUserTimeline",
			call: func() (json.RawMessage, error) { return client.GetUserTimeline(ctx, userID, "") },
		},
		{
			name: "GetUserReplies",
			call: func() (json.RawMessage, error) { return client.GetUserReplies(ctx, userID, "") },
		},
	}

	for _, src := range sources {
		raw, err := src.call()
		if err != nil {
			continue
		}
		if id := findTweetIDFromRaw(raw); id != "" {
			return id
		}
	}

	return ""
}

func TestTweetIntegration_RealAPI(t *testing.T) {
	client := requireIntegrationClient(t)
	userID := integrationTestValue(t, "XCATCH_TEST_USER_ID")
	if userID == "" {
		t.Skip("missing XCATCH_TEST_USER_ID (env or config.ini [xcatch])")
	}
	tweetID := integrationTestValue(t, "XCATCH_TEST_TWEET_ID")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if tweetID == "" {
		tweetID = discoverTweetID(ctx, client, userID)
		if tweetID != "" {
			t.Logf("auto discovered XCATCH_TEST_TWEET_ID=%s from real API payload", tweetID)
		}
	}

	t.Run("GetUserTweets", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserTweets", func() (json.RawMessage, error) {
			return client.GetUserTweets(ctx, userID, "")
		})
	})

	t.Run("GetUserTimeline", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserTimeline", func() (json.RawMessage, error) {
			return client.GetUserTimeline(ctx, userID, "")
		})
	})

	t.Run("GetUserReplies", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserReplies", func() (json.RawMessage, error) {
			return client.GetUserReplies(ctx, userID, "")
		})
	})

	t.Run("GetUserLikes", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserLikes", func() (json.RawMessage, error) {
			return client.GetUserLikes(ctx, userID, "")
		})
	})

	t.Run("GetUserLikesV2", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserLikesV2", func() (json.RawMessage, error) {
			return client.GetUserLikesV2(ctx, userID, "")
		})
	})

	t.Run("GetUserHighlights", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserHighlights", func() (json.RawMessage, error) {
			return client.GetUserHighlights(ctx, userID, "")
		})
	})

	t.Run("GetUserArticlesTweets", func(t *testing.T) {
		requireIntegrationJSON(t, "GetUserArticlesTweets", func() (json.RawMessage, error) {
			return client.GetUserArticlesTweets(ctx, userID, "")
		})
	})

	t.Run("GetHomeTimeline", func(t *testing.T) {
		if client.authToken == "" {
			t.Skip("missing auth token; set XCATCH_AUTH_TOKEN or auth_token in config.ini")
		}
		requireIntegrationJSON(t, "GetHomeTimeline", func() (json.RawMessage, error) {
			return client.GetHomeTimeline(ctx, "")
		})
	})

	t.Run("GetMentionsTimeline", func(t *testing.T) {
		if client.authToken == "" {
			t.Skip("missing auth token; set XCATCH_AUTH_TOKEN or auth_token in config.ini")
		}
		requireIntegrationJSON(t, "GetMentionsTimeline", func() (json.RawMessage, error) {
			return client.GetMentionsTimeline(ctx, "")
		})
	})

	t.Run("TweetID required group", func(t *testing.T) {
		if tweetID == "" {
			t.Skip("missing XCATCH_TEST_TWEET_ID and auto-discovery failed from user tweet endpoints")
		}

		t.Run("GetTweetDetail", func(t *testing.T) {
			requireIntegrationJSON(t, "GetTweetDetail", func() (json.RawMessage, error) {
				return client.GetTweetDetail(ctx, tweetID, "")
			})
		})

		t.Run("GetTweetSimple", func(t *testing.T) {
			requireIntegrationJSON(t, "GetTweetSimple", func() (json.RawMessage, error) {
				return client.GetTweetSimple(ctx, tweetID)
			})
		})

		t.Run("GetTweetsByIDs", func(t *testing.T) {
			requireIntegrationJSON(t, "GetTweetsByIDs", func() (json.RawMessage, error) {
				return client.GetTweetsByIDs(ctx, []string{tweetID})
			})
		})

		t.Run("GetRetweeters", func(t *testing.T) {
			requireIntegrationJSON(t, "GetRetweeters", func() (json.RawMessage, error) {
				return client.GetRetweeters(ctx, tweetID, "")
			})
		})

		t.Run("GetRetweetersIDs", func(t *testing.T) {
			requireIntegrationJSON(t, "GetRetweetersIDs", func() (json.RawMessage, error) {
				return client.GetRetweetersIDs(ctx, tweetID, "")
			})
		})

		t.Run("GetFavoriters", func(t *testing.T) {
			requireIntegrationJSON(t, "GetFavoriters", func() (json.RawMessage, error) {
				return client.GetFavoriters(ctx, tweetID, "")
			})
		})

		t.Run("GetQuotes", func(t *testing.T) {
			requireIntegrationJSON(t, "GetQuotes", func() (json.RawMessage, error) {
				return client.GetQuotes(ctx, tweetID, "")
			})
		})
	})
}
