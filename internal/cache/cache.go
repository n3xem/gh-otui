package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qawatake/gh-otui/internal/github"
	"github.com/qawatake/gh-otui/internal/models"
	"github.com/sourcegraph/conc/iter"
)

type CacheStorage struct {
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{}
}

func (c *CacheStorage) RootDir() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "gh", "extensions", "gh-otui")
}

func (c *CacheStorage) ClearAll(ctx context.Context) error {
	path := c.RootDir()
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}
	return nil
}

func (c *CacheStorage) RootFilePath() string {
	return filepath.Join(c.RootDir(), "cache.json")
}

func (c *CacheStorage) LoadLastUpdated(ctx context.Context) (*time.Time, error) {
	path := c.RootFilePath()
	cacheData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache file does not exist: %w", err)
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}
	var cache cacheDTO
	if err := json.Unmarshal(cacheData, &cache); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}
	return &cache.LastUpdated, nil
}

type MD struct {
	Me          string
	LastUpdated time.Time
}

type mdDTO struct {
	Me          string    `json:"me"`
	LastUpdated time.Time `json:"last_updated"`
}

func (c *CacheStorage) SaveMD(ctx context.Context, me string, lastUpdated time.Time) error {
	dto := mdDTO{
		Me:          me,
		LastUpdated: lastUpdated,
	}

	cacheData, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	if err := os.WriteFile(c.RootFilePath(), cacheData, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}

func (c *CacheStorage) LoadMD(ctx context.Context) (MD, error) {
	path := c.RootFilePath()
	cacheData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return MD{}, fmt.Errorf("cache file does not exist: %w", err)
		}
		return MD{}, fmt.Errorf("failed to read cache file: %w", err)
	}
	var dto mdDTO
	if err := json.Unmarshal(cacheData, &dto); err != nil {
		return MD{}, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}
	return MD(dto), nil
}

func (c *CacheStorage) SaveLastUpdated(ctx context.Context, lastUpdatedAt time.Time) error {
	cache := cacheDTO{
		LastUpdated: lastUpdatedAt,
	}

	cacheData, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	if err := os.WriteFile(c.RootFilePath(), cacheData, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}

func (c *CacheStorage) Exists() (bool, error) {
	path := c.RootFilePath()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	return true, nil
}

func (*CacheStorage) Path(id string) string {
	return filepath.Join(os.Getenv("HOME"), ".config", "gh", "extensions", "gh-otui", id+".json")
}

func (c *CacheStorage) Existsx(ctx context.Context, id string) (bool, error) {
	path := c.Path(id)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	return true, nil
}

type Cache struct {
	ID           string
	Repositories []models.Repository
	LastUpdated  time.Time
}

// キャッシュデータ構造を拡張して最終更新時間を含める
type cacheDTO struct {
	Repositories []github.Repository `json:"repositories"`
	LastUpdated  time.Time           `json:"last_updated"`
}

func (c *CacheStorage) LoadAll(ctx context.Context) ([][]models.Repository, error) {
	files, err := os.ReadDir(c.RootDir())
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	return iter.MapErr(files, func(fp *os.DirEntry) ([]models.Repository, error) {
		file := *fp
		return c.LoadCache(ctx, strings.TrimSuffix(file.Name(), ".json"))
	})
}

func (c *CacheStorage) LoadCache(ctx context.Context, id string) ([]models.Repository, error) {
	path := c.Path(id)
	cacheData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache cacheDTO
	if err := json.Unmarshal(cacheData, &cache); err != nil {
		return nil, err
	}
	repos := make([]models.Repository, 0, len(cache.Repositories))
	for _, repo := range cache.Repositories {
		repos = append(repos, models.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			Host:    repo.Host,
			HtmlUrl: repo.HtmlUrl,
		})
	}
	return repos, nil
}

func (c *CacheStorage) Save(ctx context.Context, id string, repos []models.Repository) error {
	dtoRepos := make([]github.Repository, 0, len(repos))
	for _, repo := range repos {
		dtoRepos = append(dtoRepos, github.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			Host:    repo.Host,
			HtmlUrl: repo.HtmlUrl,
		})
	}
	cache := cacheDTO{
		Repositories: dtoRepos,
		LastUpdated:  time.Now(),
	}

	cacheData, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	cacheDir := filepath.Dir(c.Path(id))
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(c.Path(id), cacheData, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}

/* func SaveCache(repos []models.Repository) error {
	dtoRepos := make([]github.Repository, 0, len(repos))
	for _, repo := range repos {
		dtoRepos = append(dtoRepos, github.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			Host:    repo.Host,
			HtmlUrl: repo.HtmlUrl,
		})
	}
	cache := cacheDTO{
		Repositories: dtoRepos,
		LastUpdated:  time.Now(),
	}

	cacheData, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	cacheDir := filepath.Dir(c.Path())
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(c.Path(), cacheData, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}
	return nil
}
*/
// キャッシュが古いかどうかをチェックする関数
// maxAge: キャッシュの最大有効期間（時間）
func IsCacheStale(lastUpdated time.Time, maxAge time.Duration) bool {
	// キャッシュが空の場合はtrueを返す
	if lastUpdated.IsZero() {
		return true
	}

	return time.Since(lastUpdated) > maxAge
}
