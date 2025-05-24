package util

import "testing"

func TestParseEventLine_SkipsComment(t *testing.T) {
	_, _, _, _, err := ParseEventLine("# just a comment")
	if err == nil || err.Error()[:7] != "skipped" {
		t.Fatalf("expected skip error, got %v", err)
	}
}
