package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/n3xem/gh-otui/github"
)

func GetCachePath() string {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "gh", "extensions", "gh-otui")
	return filepath.Join(configDir, "cache.json")
}

func LoadCache() ([]github.Repository, error) {
	cacheData, err := os.ReadFile(GetCachePath())
	if err != nil {
		return nil, err
	}
	var repos []github.Repository
	if err := json.Unmarshal(cacheData, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func SaveCache(repos []github.Repository) error {
	cacheData, err := json.Marshal(repos)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	cacheDir := filepath.Dir(GetCachePath())
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(GetCachePath(), cacheData, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}
