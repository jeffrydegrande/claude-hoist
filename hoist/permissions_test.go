package hoist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name     string
		source   []string
		target   []string
		expected []string
	}{
		{
			name:     "all new",
			source:   []string{"a", "b", "c"},
			target:   []string{},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "none new",
			source:   []string{"a", "b"},
			target:   []string{"a", "b", "c"},
			expected: nil,
		},
		{
			name:     "some new",
			source:   []string{"a", "b", "c"},
			target:   []string{"b"},
			expected: []string{"a", "c"},
		},
		{
			name:     "empty source",
			source:   []string{},
			target:   []string{"a"},
			expected: nil,
		},
		{
			name:     "both empty",
			source:   nil,
			target:   nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Diff(tt.source, tt.target)
			if len(got) != len(tt.expected) {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Fatalf("got[%d] = %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestDedup(t *testing.T) {
	got := dedup([]string{"a", "b", "a", "c", "b"})
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestMerge(t *testing.T) {
	user := Settings{
		Permissions: Permissions{
			Allow: []string{"Bash(ls:*)", "WebSearch"},
		},
	}
	newAllow := []string{"Bash(cargo test:*)", "Bash(git add:*)"}
	newDeny := []string{"Bash(rm -rf:*)"}

	got := Merge(user, newAllow, newDeny)

	wantAllow := []string{"Bash(cargo test:*)", "Bash(git add:*)", "Bash(ls:*)", "WebSearch"}
	if len(got.Permissions.Allow) != len(wantAllow) {
		t.Fatalf("allow: got %v, want %v", got.Permissions.Allow, wantAllow)
	}
	for i := range wantAllow {
		if got.Permissions.Allow[i] != wantAllow[i] {
			t.Fatalf("allow[%d]: got %q, want %q", i, got.Permissions.Allow[i], wantAllow[i])
		}
	}

	if len(got.Permissions.Deny) != 1 || got.Permissions.Deny[0] != "Bash(rm -rf:*)" {
		t.Fatalf("deny: got %v", got.Permissions.Deny)
	}
}

func TestMergeEmptyUser(t *testing.T) {
	user := Settings{}
	newAllow := []string{"Bash(ls:*)", "WebSearch"}

	got := Merge(user, newAllow, nil)

	if len(got.Permissions.Allow) != 2 {
		t.Fatalf("got %d allow rules, want 2", len(got.Permissions.Allow))
	}
	if len(got.Permissions.Deny) != 0 {
		t.Fatalf("got %d deny rules, want 0", len(got.Permissions.Deny))
	}
}

func TestReadWriteSettings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.local.json")

	original := Settings{
		Permissions: Permissions{
			Allow: []string{"Bash(ls:*)", "WebSearch"},
			Deny:  []string{"Bash(rm:*)"},
		},
	}

	if err := WriteSettings(path, original); err != nil {
		t.Fatal(err)
	}

	got, err := ReadSettings(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Permissions.Allow) != 2 {
		t.Fatalf("allow length: got %d, want 2", len(got.Permissions.Allow))
	}
	if len(got.Permissions.Deny) != 1 {
		t.Fatalf("deny length: got %d, want 1", len(got.Permissions.Deny))
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("file mode: got %o, want 600", info.Mode().Perm())
	}
}

func TestReadSettingsNotFound(t *testing.T) {
	_, err := ReadSettings("/nonexistent/path/settings.local.json")
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got: %v", err)
	}
}

func TestReadSettingsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := ReadSettings(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestWriteSettingsProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	s := Settings{
		Permissions: Permissions{
			Allow: []string{"Bash(echo \"hello world\":*)"},
		},
	}

	if err := WriteSettings(path, s); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	var check Settings
	if err := json.Unmarshal(data, &check); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(check.Permissions.Allow) != 1 {
		t.Fatal("roundtrip failed")
	}
}
