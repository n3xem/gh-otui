package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

type Metadata struct {
	lastUpdated time.Time
}

func (m *Metadata) IsStale() bool {
	return time.Since(m.lastUpdated) > 1*time.Hour
}

func (m *Metadata) Initialized() bool {
	return !m.lastUpdated.IsZero()
}

type metadataDTO struct {
	LastUpdated time.Time `json:"last_updated"`
}

func LoadMetadata(ctx context.Context) (*Metadata, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	p := metadataPath()
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Metadata{}, nil
		}
		return nil, err
	}

	var dto metadataDTO
	if err := json.Unmarshal(b, &dto); err != nil {
		return nil, err
	}
	md := Metadata{
		lastUpdated: dto.LastUpdated,
	}
	return &md, nil
}

func Done(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	dto := metadataDTO{
		LastUpdated: time.Now(),
	}
	b, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	dir := filepath.Join(root())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	p := metadataPath()
	if err := os.WriteFile(p, b, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}

func metadataPath() string {
	return filepath.Join(root(), "_md.json")
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

func Clear(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := os.RemoveAll(root()); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

func FetchRepositories(ctx context.Context) ([]*models.RepositoryGroup, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	dirs, err := os.ReadDir(root())
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	groups := make([]*models.RepositoryGroup, 0, len(dirs))
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
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
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
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
	if ctx.Err() != nil {
		return ctx.Err()
	}
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
