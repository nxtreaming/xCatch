package utools

import (
	"context"
	"encoding/json"
)

// ============================================================
// Search APIs
// ============================================================

// Search performs an advanced search on Twitter.
// query is the search keyword/expression.
// searchType can be "Latest", "Top", "People", "Photos", "Videos" etc.
// cursor can be empty for the first page.
func (c *Client) Search(ctx context.Context, query, searchType, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"words": query,
	}
	if searchType != "" {
		params["type"] = searchType
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/search", params, &result)
	return result, err
}

// SearchBox performs a search box query (typeahead / autocomplete).
func (c *Client) SearchBox(ctx context.Context, query string) (json.RawMessage, error) {
	params := map[string]string{
		"words": query,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/searchBox", params, &result)
	return result, err
}

// GetTrends retrieves trending topics for a given location.
// woeid is the "Where On Earth ID" (e.g. "1" for worldwide, "23424977" for US).
func (c *Client) GetTrends(ctx context.Context, woeid string) (json.RawMessage, error) {
	params := map[string]string{}
	if woeid != "" {
		params["woeid"] = woeid
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/trends", params, &result)
	return result, err
}

// GetTrending retrieves the current trending topics.
func (c *Client) GetTrending(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/trending", params, &result)
	return result, err
}

// GetNews retrieves news content from Twitter.
func (c *Client) GetNews(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/news", params, &result)
	return result, err
}

// GetExplorePage retrieves the explore page content.
func (c *Client) GetExplorePage(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/explore", params, &result)
	return result, err
}

// GetSports retrieves sports-related trending content.
func (c *Client) GetSports(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/sports", params, &result)
	return result, err
}

// GetEntertainment retrieves entertainment-related trending content.
func (c *Client) GetEntertainment(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/entertainment", params, &result)
	return result, err
}
