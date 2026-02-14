package utools

import (
	"context"
	"encoding/json"
)

// ============================================================
// Social Relationship APIs
// ============================================================

// GetFollowers retrieves the followers list for a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetFollowers(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/followersListV2", params, &result)
	return result, err
}

// GetFollowings retrieves the followings list for a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetFollowings(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/followingsListV2", params, &result)
	return result, err
}

// GetFollowerIDs retrieves the follower IDs for a user.
// cursor can be empty for the first page.
func (c *Client) GetFollowerIDs(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/followersIds", params, &result)
	return result, err
}

// GetFollowingIDs retrieves the following IDs for a user.
// cursor can be empty for the first page.
func (c *Client) GetFollowingIDs(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/followingsIds", params, &result)
	return result, err
}

// GetRelationship retrieves the relationship information between two users.
func (c *Client) GetRelationship(ctx context.Context, sourceID, targetID string) (json.RawMessage, error) {
	params := map[string]string{
		"sourceId": sourceID,
		"targetId": targetID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/getFriendshipsShow", params, &result)
	return result, err
}

// GetFollowersYouKnow retrieves mutual followers (followers you know) for a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetFollowersYouKnow(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/followersYouKnowV2", params, &result)
	return result, err
}

// GetBlueVerifiedFollowers retrieves blue-verified followers for a user (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetBlueVerifiedFollowers(ctx context.Context, userID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"userId": userID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/blueVerifiedFollowersV2", params, &result)
	return result, err
}

// ============================================================
// List APIs
// ============================================================

// GetListByUser retrieves lists owned by a user.
func (c *Client) GetListByUser(ctx context.Context, userID, screenName string) (json.RawMessage, error) {
	params := map[string]string{}
	if userID != "" {
		params["userId"] = userID
	}
	if screenName != "" {
		params["screenName"] = screenName
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/getListByUserIdOrScreenName", params, &result)
	return result, err
}

// GetListMembers retrieves members of a Twitter list (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetListMembers(ctx context.Context, listID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"listId": listID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/listMembersByListIdV2", params, &result)
	return result, err
}

// GetListTimeline retrieves the latest tweets from a Twitter list (V2 endpoint).
// cursor can be empty for the first page.
func (c *Client) GetListTimeline(ctx context.Context, listID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"listId": listID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/listLatestTweetsTimeline", params, &result)
	return result, err
}

// ============================================================
// Communities APIs
// ============================================================

// GetCommunitiesByScreenName retrieves communities for a user by screen name.
func (c *Client) GetCommunitiesByScreenName(ctx context.Context, screenName string) (json.RawMessage, error) {
	params := map[string]string{
		"screenName": screenName,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/getCommunitiesByScreenName", params, &result)
	return result, err
}

// GetCommunityInfo retrieves detailed information about a community.
func (c *Client) GetCommunityInfo(ctx context.Context, communityID string) (json.RawMessage, error) {
	params := map[string]string{
		"communityId": communityID,
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/communitiesFetchOneQuery", params, &result)
	return result, err
}

// GetCommunityTweets retrieves tweets from a community timeline.
// cursor can be empty for the first page.
func (c *Client) GetCommunityTweets(ctx context.Context, communityID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"communityId": communityID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/communitiesTweetsTimelineV2", params, &result)
	return result, err
}

// GetCommunityMembers retrieves members of a community.
// cursor can be empty for the first page.
func (c *Client) GetCommunityMembers(ctx context.Context, communityID string, cursor string) (json.RawMessage, error) {
	params := map[string]string{
		"communityId": communityID,
	}
	if cursor != "" {
		params["cursor"] = cursor
	}
	var result json.RawMessage
	err := c.Get(ctx, "/api/base/apitools/communitiesMemberV2", params, &result)
	return result, err
}
