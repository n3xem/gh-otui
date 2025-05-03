package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/n3xem/gh-otui/github"
	"github.com/n3xem/gh-otui/models"
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

func root() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "gh", "extensions", "gh-otui")
}

func hostPath(host string) string {
	return filepath.Join(root(), host)
}

func path(host, org string) string {
	return filepath.Join(hostPath(host), org+".json")
}

func FetchRepositories(ctx context.Context) ([]*models.RepositoryGroup, error) {
	dirs, err := os.ReadDir(root())
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	groups := make([]*models.RepositoryGroup, 0, len(dirs))
	for _, dir := range dirs {
		files, err := os.ReadDir(filepath.Join(root(), dir.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read cache directory: %w", err)
		}
		for _, file := range files {
			host := dir.Name()
			org := strings.TrimSuffix(file.Name(), ".json")
			repos, err := Load(ctx, host, org)
			if err != nil {
				return nil, fmt.Errorf("failed to load cache for %s/%s: %w", host, org, err)
			}
			groups = append(groups, repos)
		}
	}
	return groups, nil
}

type cacheDTO struct {
	Repositories []models.Repository `json:"repositories"`
}

func Load(ctx context.Context, host, org string) (*models.RepositoryGroup, error) {
	p := path(host, org)
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var dto cacheDTO
	if err := json.Unmarshal(b, &dto); err != nil {
		return nil, err
	}
	repos := make([]models.Repository, 0, len(dto.Repositories))
	for _, repo := range dto.Repositories {
		repos = append(repos, models.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			Host:    repo.Host,
			HtmlUrl: repo.HtmlUrl,
		})
	}
	g, err := models.NewRepositoryGroup(repos...)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository group: %w", err)
	}
	return g, nil
}

func Save(ctx context.Context, g *models.RepositoryGroup) error {
	repos := g.Repositories()
	dto := cacheDTO{
		Repositories: repos,
	}
	b, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	dir := hostPath(g.Host())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	p := path(g.Host(), g.Organization())
	if err := os.WriteFile(p, b, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
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
