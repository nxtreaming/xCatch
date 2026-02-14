package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tidwall/gjson"

	"github.com/xCatch/xcatch/config"
	"github.com/xCatch/xcatch/pkg/utools"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cfg := config.Load("")
	if err := cfg.Validate(); err != nil {
		log.Fatalf("config error: %v", err)
	}

	client, err := utools.NewClient(cfg)
	if err != nil {
		log.Fatalf("create client error: %v", err)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	ctx := context.Background()
	cmd := os.Args[1]

	switch cmd {
	case "user":
		cmdUser(ctx, client, os.Args[2:])
	case "tweets":
		cmdTweets(ctx, client, os.Args[2:])
	case "tweet":
		cmdTweetDetail(ctx, client, os.Args[2:])
	case "search":
		cmdSearch(ctx, client, os.Args[2:])
	case "followers":
		cmdFollowers(ctx, client, os.Args[2:])
	case "followings":
		cmdFollowings(ctx, client, os.Args[2:])
	case "likes":
		cmdLikes(ctx, client, os.Args[2:])
	case "trending":
		cmdTrending(ctx, client)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`xCatch - X.com Content Scraper powered by uTools API

Usage:
  xcatch <command> [arguments]

Commands:
  user       <screen_name>              Get user profile by screen name
  tweets     <user_id> [max_pages]      Get user tweets (default 1 page)
  tweet      <tweet_id>                 Get tweet detail with replies
  search     <query> [type]             Search tweets (type: Latest|Top|People|Photos|Videos)
  followers  <user_id>                  Get user followers (first page)
  followings <user_id>                  Get user followings (first page)
  likes      <user_id>                  Get user liked tweets (first page)
  trending                              Get current trending topics

Configuration:
  Copy config.ini.example to config.ini and fill in your API key.
  Environment variables can override config.ini values.

  Config file keys (in [xcatch] section):
    api_key, auth_token, base_url, timeout_sec, max_retries, rate_limit

  Environment Variables:
    XCATCH_API_KEY       (required) uTools API key
    XCATCH_AUTH_TOKEN    (optional) Twitter auth_token for authenticated endpoints
    XCATCH_BASE_URL      (optional) API base URL, default https://fapi.uk
    XCATCH_TIMEOUT_SEC   (optional) HTTP timeout in seconds, default 30
    XCATCH_MAX_RETRIES   (optional) Max retries, default 3
    XCATCH_RATE_LIMIT    (optional) QPS limit, default 5`)
}

// ============================================================
// Command handlers
// ============================================================

func cmdUser(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch user <screen_name>")
	}
	screenName := args[0]

	log.Printf("Fetching user profile for @%s ...", screenName)
	data, err := client.GetUserByScreenNameV2(ctx, screenName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)

	// Print summary
	parsed := gjson.ParseBytes(data)
	name := findField(parsed, "name")
	desc := findField(parsed, "description")
	followers := findField(parsed, "followers_count")
	following := findField(parsed, "friends_count")
	tweets := findField(parsed, "statuses_count")

	fmt.Println("\n--- Summary ---")
	fmt.Printf("Name:       %s\n", name)
	fmt.Printf("Handle:     @%s\n", screenName)
	fmt.Printf("Bio:        %s\n", desc)
	fmt.Printf("Followers:  %s\n", followers)
	fmt.Printf("Following:  %s\n", following)
	fmt.Printf("Tweets:     %s\n", tweets)
}

func cmdTweets(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch tweets <user_id> [max_pages]")
	}
	userID := args[0]
	maxPages := 1
	if len(args) > 1 {
		if _, err := fmt.Sscanf(args[1], "%d", &maxPages); err != nil || maxPages <= 0 {
			log.Fatalf("invalid max_pages: %q (must be a positive integer)", args[1])
		}
	}

	log.Printf("Fetching tweets for user %s (max %d pages) ...", userID, maxPages)

	iter := client.NewPageIterator("/userTweetsV2", map[string]string{
		"userId": userID,
	}, maxPages)

	for iter.HasMore() {
		page, err := iter.Next(ctx)
		if err != nil {
			log.Fatalf("error on page %d: %v", iter.PageCount(), err)
		}
		if page == nil {
			break
		}

		fmt.Printf("\n=== Page %d ===\n", iter.PageCount())
		printJSON(page.RawData)

		if page.NextCursor != "" {
			fmt.Printf("\n[Next cursor: %s]\n", utools.Truncate(page.NextCursor, 50))
		}
	}

	fmt.Printf("\nTotal pages fetched: %d\n", iter.PageCount())
}

func cmdTweetDetail(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch tweet <tweet_id>")
	}
	tweetID := args[0]

	log.Printf("Fetching tweet detail for %s ...", tweetID)
	data, err := client.GetTweetDetail(ctx, tweetID, "")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

func cmdSearch(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch search <query> [type]")
	}
	query := args[0]
	searchType := "Latest"
	if len(args) > 1 {
		searchType = args[1]
	}

	log.Printf("Searching for '%s' (type: %s) ...", query, searchType)
	data, err := client.Search(ctx, query, searchType, "")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

func cmdFollowers(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch followers <user_id>")
	}
	userID := args[0]

	log.Printf("Fetching followers for user %s ...", userID)
	data, err := client.GetFollowers(ctx, userID, "")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

func cmdFollowings(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch followings <user_id>")
	}
	userID := args[0]

	log.Printf("Fetching followings for user %s ...", userID)
	data, err := client.GetFollowings(ctx, userID, "")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

func cmdLikes(ctx context.Context, client *utools.Client, args []string) {
	if len(args) < 1 {
		log.Fatal("usage: xcatch likes <user_id>")
	}
	userID := args[0]

	log.Printf("Fetching likes for user %s ...", userID)
	data, err := client.GetUserLikes(ctx, userID, "")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

func cmdTrending(ctx context.Context, client *utools.Client) {
	log.Println("Fetching trending topics ...")
	data, err := client.GetTrending(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	printJSON(data)
}

// ============================================================
// Helpers
// ============================================================

func printJSON(data json.RawMessage) {
	var pretty json.RawMessage
	if err := json.Unmarshal(data, &pretty); err != nil {
		fmt.Println(string(data))
		return
	}
	out, err := json.MarshalIndent(pretty, "", "  ")
	if err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(string(out))
}

func findField(result gjson.Result, field string) string {
	// Search recursively for the field
	val := result.Get(field)
	if val.Exists() {
		return val.String()
	}
	// Try deeper paths common in Twitter V2 responses
	paths := []string{
		"data.user.result.legacy." + field,
		"data.user.legacy." + field,
		"user.result.legacy." + field,
		"result.legacy." + field,
		"legacy." + field,
		"data." + field,
	}
	for _, p := range paths {
		val = result.Get(p)
		if val.Exists() {
			return val.String()
		}
	}
	return ""
}
