package utools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// PageResult holds a single page of results with cursor information.
type PageResult struct {
	// RawData is the raw JSON response data.
	RawData json.RawMessage

	// NextCursor is the cursor value for fetching the next page.
	// Empty string means no more pages.
	NextCursor string

	// PreviousCursor is the cursor value for the previous page (if available).
	PreviousCursor string
}

// PageIterator provides an iterator interface for paginated API results.
type PageIterator struct {
	client     *Client
	path       string
	baseParams map[string]string
	nextCursor string
	hasMore    bool
	pageCount  int
	maxPages   int // 0 = unlimited
}

// NewPageIterator creates a new PageIterator for the given API path.
// maxPages controls the maximum number of pages to fetch (0 = unlimited).
func (c *Client) NewPageIterator(path string, params map[string]string, maxPages int) *PageIterator {
	// Copy params to avoid mutation
	copied := make(map[string]string, len(params))
	for k, v := range params {
		copied[k] = v
	}

	return &PageIterator{
		client:     c,
		path:       path,
		baseParams: copied,
		hasMore:    true,
		maxPages:   maxPages,
	}
}

// HasMore returns true if there are more pages to fetch.
func (it *PageIterator) HasMore() bool {
	return it.hasMore
}

// PageCount returns the number of pages fetched so far.
func (it *PageIterator) PageCount() int {
	return it.pageCount
}

// Next fetches the next page of results.
// Returns the PageResult and an error. When no more pages are available,
// PageResult will be nil and error will be nil.
func (it *PageIterator) Next(ctx context.Context) (*PageResult, error) {
	if !it.hasMore {
		return nil, nil
	}

	if it.maxPages > 0 && it.pageCount >= it.maxPages {
		it.hasMore = false
		return nil, nil
	}

	// Build params for this request
	params := make(map[string]string, len(it.baseParams)+1)
	for k, v := range it.baseParams {
		params[k] = v
	}
	if it.nextCursor != "" {
		params["cursor"] = it.nextCursor
	}

	// Execute request
	var raw json.RawMessage
	if err := it.client.Get(ctx, it.path, params, &raw); err != nil {
		return nil, fmt.Errorf("page iterator: %w", err)
	}

	it.pageCount++

	// Extract cursors from response
	result := &PageResult{
		RawData: raw,
	}

	nextCursor, prevCursor := extractCursors(string(raw))
	result.NextCursor = nextCursor
	result.PreviousCursor = prevCursor

	if nextCursor == "" || nextCursor == it.nextCursor {
		it.hasMore = false
	} else {
		it.nextCursor = nextCursor
	}

	return result, nil
}

// extractCursors extracts the bottom (next) and top (previous) cursor values
// from the API response JSON. The cursor can be in different locations depending
// on the endpoint.
func extractCursors(jsonStr string) (next string, prev string) {
	// Strategy 1: Look for cursor entries in the standard timeline format
	// Cursors are typically in entries with entryType "TimelineTimelineCursor"
	entries := gjson.Get(jsonStr, "..entries")
	if entries.Exists() {
		entries.ForEach(func(_, entry gjson.Result) bool {
			// Check each entry for cursor type
			entry.ForEach(func(_, item gjson.Result) bool {
				cursorType := item.Get("cursorType").String()
				if cursorType == "" {
					cursorType = item.Get("content.cursorType").String()
				}
				value := item.Get("value").String()
				if value == "" {
					value = item.Get("content.value").String()
				}

				switch cursorType {
				case "Bottom":
					next = value
				case "Top":
					prev = value
				}
				return true
			})
			return true
		})
	}

	// Strategy 2: Direct cursor fields (some endpoints)
	if next == "" {
		next = gjson.Get(jsonStr, "cursor_bottom").String()
		if next == "" {
			next = gjson.Get(jsonStr, "next_cursor").String()
		}
		if next == "" {
			next = gjson.Get(jsonStr, "next_cursor_str").String()
		}
	}
	if prev == "" {
		prev = gjson.Get(jsonStr, "cursor_top").String()
		if prev == "" {
			prev = gjson.Get(jsonStr, "previous_cursor").String()
		}
		if prev == "" {
			prev = gjson.Get(jsonStr, "previous_cursor_str").String()
		}
	}

	// Strategy 3: Deep search for cursor objects
	if next == "" {
		gjson.Parse(jsonStr).ForEach(func(key, value gjson.Result) bool {
			return findCursorDeep(value, &next, &prev)
		})
	}

	return next, prev
}

// findCursorDeep recursively searches for cursor values in nested JSON.
func findCursorDeep(value gjson.Result, next, prev *string) bool {
	if !value.IsObject() && !value.IsArray() {
		return true
	}

	if value.IsObject() {
		cursorType := value.Get("cursorType").String()
		val := value.Get("value").String()
		if cursorType == "Bottom" && val != "" && *next == "" {
			*next = val
		}
		if cursorType == "Top" && val != "" && *prev == "" {
			*prev = val
		}
	}

	value.ForEach(func(_, child gjson.Result) bool {
		return findCursorDeep(child, next, prev)
	})

	return *next == "" || *prev == ""
}

// CollectAll is a convenience method that fetches all pages and collects raw results.
func (it *PageIterator) CollectAll(ctx context.Context) ([]json.RawMessage, error) {
	var pages []json.RawMessage
	for it.HasMore() {
		page, err := it.Next(ctx)
		if err != nil {
			return pages, err
		}
		if page == nil {
			break
		}
		pages = append(pages, page.RawData)
	}
	return pages, nil
}
