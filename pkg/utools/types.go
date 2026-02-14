package utools

import "encoding/json"

// ============================================================
// Common / Wrapper types
// ============================================================

// RawResponse is the raw JSON response from the API, useful when the
// caller wants to do custom parsing.
type RawResponse = json.RawMessage

// ============================================================
// User types
//
// Note: API methods currently return json.RawMessage for maximum flexibility.
// These typed structs are provided for callers who want structured access:
//
//	var user UserResult
//	if err := json.Unmarshal(rawData, &user); err != nil { ... }
// ============================================================

// UserResult represents a Twitter user profile.
type UserResult struct {
	ID                  string   `json:"id_str"`
	RestID              string   `json:"rest_id"`
	Name                string   `json:"name"`
	ScreenName          string   `json:"screen_name"`
	Description         string   `json:"description"`
	Location            string   `json:"location"`
	URL                 string   `json:"url"`
	Protected           bool     `json:"protected"`
	Verified            bool     `json:"verified"`
	IsBlueVerified      bool     `json:"is_blue_verified"`
	FollowersCount      int      `json:"followers_count"`
	FriendsCount        int      `json:"friends_count"`
	ListedCount         int      `json:"listed_count"`
	FavouritesCount     int      `json:"favourites_count"`
	StatusesCount       int      `json:"statuses_count"`
	MediaCount          int      `json:"media_count"`
	CreatedAt           string   `json:"created_at"`
	ProfileImageURL     string   `json:"profile_image_url_https"`
	ProfileBannerURL    string   `json:"profile_banner_url"`
	PinnedTweetIdsStr   []string `json:"pinned_tweet_ids_str"`
	HasCustomTimelines  bool     `json:"has_custom_timelines"`
	CanDM               bool     `json:"can_dm"`
	DefaultProfile      bool     `json:"default_profile"`
	DefaultProfileImage bool     `json:"default_profile_image"`
}

// UserListResult represents a paginated list of users.
type UserListResult struct {
	Users      []UserResult `json:"users"`
	NextCursor string       `json:"next_cursor"`
}

// UsernameChange represents a single username change record.
type UsernameChange struct {
	OldName   string `json:"old_name"`
	NewName   string `json:"new_name"`
	ChangedAt string `json:"changed_at"`
}

// RelationshipResult represents the relationship between two users.
type RelationshipResult struct {
	Source RelationshipUser `json:"source"`
	Target RelationshipUser `json:"target"`
}

// RelationshipUser holds relationship flags for one side.
type RelationshipUser struct {
	ID                   string `json:"id_str"`
	ScreenName           string `json:"screen_name"`
	Following            bool   `json:"following"`
	FollowedBy           bool   `json:"followed_by"`
	Blocking             bool   `json:"blocking"`
	Muting               bool   `json:"muting"`
	CanDM                bool   `json:"can_dm"`
	WantRetweets         bool   `json:"want_retweets"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
}

// ============================================================
// Tweet types
// ============================================================

// TweetResult represents a single tweet.
type TweetResult struct {
	ID                  string            `json:"id_str"`
	RestID              string            `json:"rest_id"`
	FullText            string            `json:"full_text"`
	Text                string            `json:"text"`
	CreatedAt           string            `json:"created_at"`
	ConversationIDStr   string            `json:"conversation_id_str"`
	InReplyToStatusID   string            `json:"in_reply_to_status_id_str"`
	InReplyToUserID     string            `json:"in_reply_to_user_id_str"`
	InReplyToScreenName string            `json:"in_reply_to_screen_name"`
	Lang                string            `json:"lang"`
	Source              string            `json:"source"`
	RetweetCount        int               `json:"retweet_count"`
	FavoriteCount       int               `json:"favorite_count"`
	ReplyCount          int               `json:"reply_count"`
	QuoteCount          int               `json:"quote_count"`
	BookmarkCount       int               `json:"bookmark_count"`
	ViewCount           string            `json:"view_count"`
	IsQuoteStatus       bool              `json:"is_quote_status"`
	Retweeted           bool              `json:"retweeted"`
	Favorited           bool              `json:"favorited"`
	Bookmarked          bool              `json:"bookmarked"`
	User                *UserResult       `json:"user"`
	Entities            *TweetEntities    `json:"entities"`
	ExtendedEntities    *ExtendedEntities `json:"extended_entities"`
	QuotedStatus        *TweetResult      `json:"quoted_status"`
	RetweetedStatus     *TweetResult      `json:"retweeted_status"`
	Card                json.RawMessage   `json:"card"`
}

// GetText returns the best available text content of the tweet.
func (t *TweetResult) GetText() string {
	if t.FullText != "" {
		return t.FullText
	}
	return t.Text
}

// TweetEntities holds entity information extracted from tweet text.
type TweetEntities struct {
	URLs         []URLEntity     `json:"urls"`
	Hashtags     []HashtagEntity `json:"hashtags"`
	UserMentions []MentionEntity `json:"user_mentions"`
	Symbols      []SymbolEntity  `json:"symbols"`
	Media        []MediaEntity   `json:"media"`
}

// ExtendedEntities holds extended media information.
type ExtendedEntities struct {
	Media []MediaEntity `json:"media"`
}

// URLEntity represents a URL found in tweet text.
type URLEntity struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
}

// HashtagEntity represents a hashtag in tweet text.
type HashtagEntity struct {
	Text string `json:"text"`
}

// MentionEntity represents a user mention in tweet text.
type MentionEntity struct {
	ID         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

// SymbolEntity represents a cashtag/symbol in tweet text.
type SymbolEntity struct {
	Text string `json:"text"`
}

// MediaEntity represents a media attachment.
type MediaEntity struct {
	ID          string          `json:"id_str"`
	MediaURL    string          `json:"media_url_https"`
	URL         string          `json:"url"`
	ExpandedURL string          `json:"expanded_url"`
	Type        string          `json:"type"` // photo, video, animated_gif
	VideoInfo   *VideoInfo      `json:"video_info"`
	Sizes       json.RawMessage `json:"sizes"`
}

// VideoInfo holds video-specific media information.
type VideoInfo struct {
	DurationMillis int            `json:"duration_millis"`
	AspectRatio    []int          `json:"aspect_ratio"`
	Variants       []VideoVariant `json:"variants"`
}

// VideoVariant represents a single video encoding variant.
type VideoVariant struct {
	Bitrate     int    `json:"bitrate"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

// TweetListResult represents a paginated list of tweets.
type TweetListResult struct {
	Tweets     []TweetResult `json:"tweets"`
	NextCursor string        `json:"next_cursor"`
}

// TweetDetailResult represents a tweet with its conversation thread.
type TweetDetailResult struct {
	Tweet      TweetResult   `json:"tweet"`
	Replies    []TweetResult `json:"replies"`
	NextCursor string        `json:"next_cursor"`
}

// ============================================================
// Search types
// ============================================================

// SearchResult represents search results.
type SearchResult struct {
	Tweets     []TweetResult `json:"tweets"`
	Users      []UserResult  `json:"users"`
	NextCursor string        `json:"next_cursor"`
}

// TrendResult represents a single trending topic.
type TrendResult struct {
	Name       string `json:"name"`
	Query      string `json:"query"`
	URL        string `json:"url"`
	TweetCount int    `json:"tweet_volume"`
}

// TrendsResult represents a list of trends.
type TrendsResult struct {
	Trends []TrendResult `json:"trends"`
}
