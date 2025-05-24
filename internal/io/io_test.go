package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteReadLines(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "out.txt")
	lines := []string{"a", "b", "c"}
	if err := WriteLines(p, lines); err != nil {
		t.Fatalf("write error: %v", err)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(got) != "a\nb\nc\n" {
		t.Fatalf("content mismatch: %q", got)
	}
}
