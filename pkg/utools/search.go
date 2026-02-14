package utools

import (
	"context"
	"encoding/json"
)

// SearchOptions contains optional query parameters for advanced search.
// Field names follow the uTools generated API parameter names.
type SearchOptions struct {
	Type       string
	Cursor     string
	Any        string
	From       string
	Likes      string
	Mentioning string
	None       string
	Phrase     string
	Replies    string
	Retweets   string
	Since      string
	Tag        string
	To         string
	Until      string
}

// ============================================================
// Search APIs
// ============================================================

// Search performs an advanced search on Twitter.
// query is the search keyword/expression.
// searchType can be "Latest", "Top", "People", "Photos", "Videos" etc.
// cursor can be empty for the first page.
func (c *Client) Search(ctx context.Context, query, searchType, cursor string) (json.RawMessage, error) {
	return c.SearchWithOptions(ctx, query, SearchOptions{
		Type:   searchType,
		Cursor: cursor,
	})
}

// SearchWithOptions performs advanced search with optional filters.
func (c *Client) SearchWithOptions(ctx context.Context, query string, opts SearchOptions) (json.RawMessage, error) {
	params := map[string]string{
		"words": query,
	}
	if opts.Type != "" {
		params["type"] = opts.Type
	}
	if opts.Cursor != "" {
		params["cursor"] = opts.Cursor
	}
	if opts.Any != "" {
		params["any"] = opts.Any
	}
	if opts.From != "" {
		params["from"] = opts.From
	}
	if opts.Likes != "" {
		params["likes"] = opts.Likes
	}
	if opts.Mentioning != "" {
		params["mentioning"] = opts.Mentioning
	}
	if opts.None != "" {
		params["none"] = opts.None
	}
	if opts.Phrase != "" {
		params["phrase"] = opts.Phrase
	}
	if opts.Replies != "" {
		params["replies"] = opts.Replies
	}
	if opts.Retweets != "" {
		params["retweets"] = opts.Retweets
	}
	if opts.Since != "" {
		params["since"] = opts.Since
	}
	if opts.Tag != "" {
		params["tag"] = opts.Tag
	}
	if opts.To != "" {
		params["to"] = opts.To
	}
	if opts.Until != "" {
		params["until"] = opts.Until
	}
	var result json.RawMessage
	err := c.Get(ctx, "/search", params, &result)
	return result, err
}

// SearchBox performs a search box query (typeahead / autocomplete).
func (c *Client) SearchBox(ctx context.Context, query string) (json.RawMessage, error) {
	params := map[string]string{
		"words": query,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/searchBox", params, &result)
	return result, err
}

// GetTrends retrieves trending topics for a given location.
// woeid is the "Where On Earth ID" (e.g. "1" for worldwide, "23424977" for US).
func (c *Client) GetTrends(ctx context.Context, woeid string) (json.RawMessage, error) {
	params := map[string]string{}
	if woeid != "" {
		// Official generated API uses "id"; keep "woeid" as a compatibility alias.
		params["id"] = woeid
		params["woeid"] = woeid
	}
	var result json.RawMessage
	err := c.Get(ctx, "/trends", params, &result)
	return result, err
}

// GetTrending retrieves the current trending topics.
func (c *Client) GetTrending(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/trending", params, &result)
	return result, err
}

// GetNews retrieves news content from Twitter.
func (c *Client) GetNews(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/news", params, &result)
	return result, err
}

// GetExplorePage retrieves the explore page content.
func (c *Client) GetExplorePage(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/explore", params, &result)
	return result, err
}

// GetSports retrieves sports-related trending content.
func (c *Client) GetSports(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/sports", params, &result)
	return result, err
}

// GetEntertainment retrieves entertainment-related trending content.
func (c *Client) GetEntertainment(ctx context.Context) (json.RawMessage, error) {
	params := map[string]string{}
	var result json.RawMessage
	err := c.Get(ctx, "/entertainment", params, &result)
	return result, err
}
