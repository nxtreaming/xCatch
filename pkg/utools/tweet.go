package utools

import (
	"context"
	"encoding/json"
	"strings"
)

// ============================================================
// Tweet Content APIs
// ============================================================

// GetUserTweets retrieves tweets posted by a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetUserTweets(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/userTweetsV2", params, &result)
	return result, err
}

// GetUserTimeline retrieves the user timeline (same data as UserTweetsV2).
// cursor can be empty for the first page.
func (c *Client) GetUserTimeline(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/userTimeline", params, &result)
	return result, err
}

// GetTweetDetail retrieves a tweet's full details including its reply thread.
// cursor can be empty for the first page of replies.
func (c *Client) GetTweetDetail(ctx context.Context, tweetID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetId": tweetID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/tweetTimeline", params, &result)
	return result, err
}

// GetTweetSimple retrieves brief information about a tweet.
func (c *Client) GetTweetSimple(ctx context.Context, tweetID string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetId": tweetID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/tweetSimple", params, &result)
	return result, err
}

// GetTweetsByIDs retrieves multiple tweets by their IDs in batch.
func (c *Client) GetTweetsByIDs(ctx context.Context, tweetIDs []string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetIds": strings.Join(tweetIDs, ","),
	}
	var result json.RawMessage
	err := c.Get(ctx, "/tweetResultsByRestIds", params, &result)
	return result, err
}

// GetUserReplies retrieves reply tweets posted by a user.
// cursor can be empty for the first page.
func (c *Client) GetUserReplies(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/userTweetReply", params, &result)
	return result, err
}

// GetUserLikes retrieves tweets liked by a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetUserLikes(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/userLikeV2", params, &result)
	return result, err
}

// GetUserHighlights retrieves a user's highlighted/pinned tweets (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetUserHighlights(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/highlightsV2", params, &result)
	return result, err
}

// GetUserArticlesTweets retrieves a user's article-type tweets.
// cursor can be empty for the first page.
func (c *Client) GetUserArticlesTweets(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/userArticlesTweets", params, &result)
	return result, err
}

// GetHomeTimeline retrieves the authenticated user's home timeline.
// Requires auth_token to be set in the client config.
// cursor can be empty for the first page.
func (c *Client) GetHomeTimeline(ctx context.Context, cursor string) (json.RawMessage, error) {
	if c.authToken == "" {
		return nil, ErrAuthTokenRequired
	}

	params := map[string]string{}
	params["auth_token"] = c.authToken
	if c.ct0 != "" {
		params["ct0"] = c.ct0
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/homeTimeline", params, &result)
	return result, err
}

// GetMentionsTimeline retrieves the authenticated user's mentions timeline.
// Requires auth_token to be set in the client config.
// cursor can be empty for the first page.
func (c *Client) GetMentionsTimeline(ctx context.Context, cursor string) (json.RawMessage, error) {
	if c.authToken == "" {
		return nil, ErrAuthTokenRequired
	}

	params := map[string]string{}
	params["auth_token"] = c.authToken
	if c.ct0 != "" {
		params["ct0"] = c.ct0
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/mentionsTimeline", params, &result)
	return result, err
}

// ============================================================
// Tweet Interaction Data APIs
// ============================================================

// GetRetweeters retrieves the list of users who retweeted a tweet (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetRetweeters(ctx context.Context, tweetID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetId": tweetID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/retweetersV2", params, &result)
	return result, err
}

// GetFavoriters retrieves the list of users who liked a tweet (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetFavoriters(ctx context.Context, tweetID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetId": tweetID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/favoritersV2", params, &result)
	return result, err
}

// GetQuotes retrieves quote tweets for a given tweet (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetQuotes(ctx context.Context, tweetID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"tweetId": tweetID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/quotesV2", params, &result)
	return result, err
}
