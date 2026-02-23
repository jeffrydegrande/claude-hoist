package hoist

import (
	"strings"
	"testing"
)

func TestUnifiedDiffIdentical(t *testing.T) {
	got := UnifiedDiff("a", "b", "hello\nworld\n", "hello\nworld\n")
	if got != "" {
		t.Fatalf("expected empty diff for identical input, got:\n%s", got)
	}
}

func TestUnifiedDiffAddedLines(t *testing.T) {
	a := "line1\nline2\n"
	b := "line1\nline2\nline3\n"
	got := UnifiedDiff("a", "b", a, b)

	if !strings.Contains(got, "+line3") {
		t.Fatalf("expected +line3 in diff, got:\n%s", got)
	}
	if strings.Contains(got, "-line") {
		t.Fatalf("unexpected removal in diff:\n%s", got)
	}
}

func TestUnifiedDiffRemovedLines(t *testing.T) {
	a := "line1\nline2\nline3\n"
	b := "line1\nline3\n"
	got := UnifiedDiff("a", "b", a, b)

	if !strings.Contains(got, "-line2") {
		t.Fatalf("expected -line2 in diff, got:\n%s", got)
	}
}

func TestUnifiedDiffChanged(t *testing.T) {
	a := "{\n  \"allow\": [\"a\"]\n}\n"
	b := "{\n  \"allow\": [\"a\", \"b\"]\n}\n"
	got := UnifiedDiff("before", "after", a, b)

	if !strings.Contains(got, "--- before") {
		t.Fatalf("expected --- before header, got:\n%s", got)
	}
	if !strings.Contains(got, "+++ after") {
		t.Fatalf("expected +++ after header, got:\n%s", got)
	}
	if !strings.Contains(got, "-  \"allow\": [\"a\"]") {
		t.Fatalf("expected removal of old line, got:\n%s", got)
	}
	if !strings.Contains(got, "+  \"allow\": [\"a\", \"b\"]") {
		t.Fatalf("expected addition of new line, got:\n%s", got)
	}
}

func TestUnifiedDiffEmpty(t *testing.T) {
	got := UnifiedDiff("a", "b", "", "line1\n")
	if !strings.Contains(got, "+line1") {
		t.Fatalf("expected +line1, got:\n%s", got)
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"a\n", 1},
		{"a\nb\n", 2},
		{"a\nb", 2},
	}
	for _, tt := range tests {
		got := splitLines(tt.input)
		if len(got) != tt.want {
			t.Errorf("splitLines(%q) = %d lines, want %d", tt.input, len(got), tt.want)
		}
	}
}
