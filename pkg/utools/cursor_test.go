package utools

import "testing"

func TestExtractCursorsFromDirectFields(t *testing.T) {
	jsonStr := `{"next_cursor":"next-123","previous_cursor":"prev-456"}`
	next, prev := extractCursors(jsonStr)
	if next != "next-123" {
		t.Fatalf("expected next cursor, got %q", next)
	}
	if prev != "prev-456" {
		t.Fatalf("expected previous cursor, got %q", prev)
	}
}

func TestExtractCursorsFromTimelineEntries(t *testing.T) {
	jsonStr := `{
		"data": {
			"entries": [
				{"content": {"cursorType": "Top", "value": "top-cursor"}},
				{"content": {"cursorType": "Bottom", "value": "bottom-cursor"}}
			]
		}
	}`
	next, prev := extractCursors(jsonStr)
	if next != "bottom-cursor" {
		t.Fatalf("expected bottom cursor as next, got %q", next)
	}
	if prev != "top-cursor" {
		t.Fatalf("expected top cursor as previous, got %q", prev)
	}
}
