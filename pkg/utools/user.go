package utools

import (
	"context"
	"encoding/json"
	"strings"
)

// ============================================================
// User Information APIs
// ============================================================

// GetUserByScreenName retrieves user information by Twitter screen name (handle).
// e.g. GetUserByScreenName(ctx, "elonmusk")
func (c *Client) GetUserByScreenName(ctx context.Context, screenName string) (json.RawMessage, error) {
	params := map[string]string{
		"screenName": screenName,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/screenname", params, &result)
	return result, err
}

// GetUserByID retrieves user information by Twitter user ID (rest_id).
func (c *Client) GetUserByID(ctx context.Context, userID string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/id", params, &result)
	return result, err
}

// GetUsersByIDs retrieves user information for multiple user IDs in batch.
// userIDs should be a slice of Twitter user ID strings.
func (c *Client) GetUsersByIDs(ctx context.Context, userIDs []string) (json.RawMessage, error) {
	params := map[string]string{
		"userIds": strings.Join(userIDs, ","),
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/ids", params, &result)
	return result, err
}

// GetUsernameChanges retrieves the username change history for a user.
func (c *Client) GetUsernameChanges(ctx context.Context, userID string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/usernameChanges", params, &result)
	return result, err
}

// LookupUser retrieves user information by username or user ID.
// Pass either screenName or userID (the other can be empty).
func (c *Client) LookupUser(ctx context.Context, screenName, userID string) (json.RawMessage, error) {
	params := map[string]string{}
	if screenName != "" {
		params["screenName"] = screenName
	}
	if userID != "" {
		params["userId"] = userID
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/lookup", params, &result)
	return result, err
}

// GetUserByScreenNameV2 retrieves user info by screen name using the V2 endpoint.
func (c *Client) GetUserByScreenNameV2(ctx context.Context, screenName string) (json.RawMessage, error) {
	params := map[string]string{
		"screenName": screenName,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/userByScreenNameV2", params, &result)
	return result, err
}

// GetUserByIDV2 retrieves user info by user ID using the V2 endpoint.
func (c *Client) GetUserByIDV2(ctx context.Context, userID string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/uerByIdRestIdV2", params, &result)
	return result, err
}

// GetUsersByIDsV2 retrieves multiple users by IDs using the V2 endpoint.
func (c *Client) GetUsersByIDsV2(ctx context.Context, userIDs []string) (json.RawMessage, error) {
	params := map[string]string{
		"userIds": strings.Join(userIDs, ","),
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/usersByIdRestIds", params, &result)
	return result, err
}

// GetAccountAnalytics retrieves account analytics data.
// Requires auth_token to be set in the client config.
func (c *Client) GetAccountAnalytics(ctx context.Context) (json.RawMessage, error) {
	if c.authToken == "" {
		return nil, ErrAuthTokenRequired
	}

	params := map[string]string{}
	params["auth_token"] = c.authToken
	if c.ct0 != "" {
		params["ct0"] = c.ct0
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/accountAnalytics", params, &result)
	return result, err
}
