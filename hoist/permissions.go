package hoist

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Permissions struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

type Settings struct {
	Permissions Permissions `json:"permissions"`
}

func FindProjectSettings() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := filepath.Join(cwd, ".claude", "settings.local.json")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("no .claude/settings.local.json in current directory")
	}
	return p, nil
}

func UserSettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "settings.local.json"), nil
}

func ReadSettings(path string) (Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{}, fmt.Errorf("parsing %s: %w", path, err)
	}
	return s, nil
}

func WriteSettings(path string, s Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0600)
}

// Diff returns items in source that are not in target.
func Diff(source, target []string) []string {
	have := make(map[string]bool, len(target))
	for _, v := range target {
		have[v] = true
	}
	var result []string
	for _, v := range source {
		if !have[v] {
			result = append(result, v)
		}
	}
	return result
}

func Merge(user Settings, newAllow, newDeny []string) Settings {
	user.Permissions.Allow = dedup(append(user.Permissions.Allow, newAllow...))
	user.Permissions.Deny = dedup(append(user.Permissions.Deny, newDeny...))
	sort.Strings(user.Permissions.Allow)
	sort.Strings(user.Permissions.Deny)
	return user
}

func dedup(items []string) []string {
	seen := make(map[string]bool, len(items))
	var result []string
	for _, v := range items {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// LoadBoth reads project and user settings, returning both plus the computed diffs.
func LoadBoth() (project, user Settings, userPath string, newAllow, newDeny []string, err error) {
	projectPath, err := FindProjectSettings()
	if err != nil {
		return Settings{}, Settings{}, "", nil, nil, err
	}

	userPath, err = UserSettingsPath()
	if err != nil {
		return Settings{}, Settings{}, "", nil, nil, err
	}

	project, err = ReadSettings(projectPath)
	if err != nil {
		return Settings{}, Settings{}, "", nil, nil, fmt.Errorf("reading project settings: %w", err)
	}

	user, err = ReadSettings(userPath)
	if err != nil && !os.IsNotExist(err) {
		return Settings{}, Settings{}, "", nil, nil, fmt.Errorf("reading user settings: %w", err)
	}
	if os.IsNotExist(err) {
		user = Settings{}
	}

	newAllow = Diff(project.Permissions.Allow, user.Permissions.Allow)
	newDeny = Diff(project.Permissions.Deny, user.Permissions.Deny)

	return project, user, userPath, newAllow, newDeny, nil
}
