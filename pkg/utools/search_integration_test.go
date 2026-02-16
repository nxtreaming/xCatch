//go:build integration

package utools

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestSearchIntegration_RealAPI(t *testing.T) {
	client := requireIntegrationClient(t)
	query := integrationTestValue(t, "XCATCH_TEST_SEARCH_QUERY")
	if query == "" {
		query = "bitcoin"
	}
	woeid := integrationTestValue(t, "XCATCH_TEST_WOEID")
	if woeid == "" {
		woeid = "1"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("Search", func(t *testing.T) {
		requireIntegrationJSON(t, "Search", func() (json.RawMessage, error) {
			return client.Search(ctx, query, "Latest", "")
		})
	})

	t.Run("SearchWithOptions", func(t *testing.T) {
		requireIntegrationJSON(t, "SearchWithOptions", func() (json.RawMessage, error) {
			return client.SearchWithOptions(ctx, query, SearchOptions{
				Type:  "Latest",
				From:  "elonmusk",
				Until: "2026-01-01",
			})
		})
	})

	t.Run("SearchBox", func(t *testing.T) {
		requireIntegrationJSON(t, "SearchBox", func() (json.RawMessage, error) {
			return client.SearchBox(ctx, query)
		})
	})

	t.Run("GetTrends", func(t *testing.T) {
		requireIntegrationJSON(t, "GetTrends", func() (json.RawMessage, error) {
			return client.GetTrends(ctx, woeid)
		})
	})

	t.Run("GetTrending", func(t *testing.T) {
		requireIntegrationJSON(t, "GetTrending", func() (json.RawMessage, error) {
			return client.GetTrending(ctx)
		})
	})

	t.Run("GetNews", func(t *testing.T) {
		requireIntegrationJSON(t, "GetNews", func() (json.RawMessage, error) {
			return client.GetNews(ctx)
		})
	})

	t.Run("GetExplorePage", func(t *testing.T) {
		requireIntegrationJSON(t, "GetExplorePage", func() (json.RawMessage, error) {
			return client.GetExplorePage(ctx)
		})
	})

	t.Run("GetSports", func(t *testing.T) {
		requireIntegrationJSON(t, "GetSports", func() (json.RawMessage, error) {
			return client.GetSports(ctx)
		})
	})

	t.Run("GetEntertainment", func(t *testing.T) {
		requireIntegrationJSON(t, "GetEntertainment", func() (json.RawMessage, error) {
			return client.GetEntertainment(ctx)
		})
	})
}
