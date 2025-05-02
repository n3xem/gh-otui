package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/qawatake/gh-otui/internal/cache"
	"github.com/qawatake/gh-otui/internal/cmd"
	"github.com/qawatake/gh-otui/internal/github"
	"github.com/qawatake/gh-otui/internal/models"
	"github.com/sourcegraph/conc/pool"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	debug := os.Getenv("GH_OTUI_DEBUG") == "1"
	logger := newLogger(debug)
	ctx = newContextWithLogger(ctx, logger)

	if err := run(ctx); err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := loggerFromCtx(ctx)

	core, err := NewCore(ctx)
	if err != nil {
		return fmt.Errorf("failed to create core: %w", err)
	}

	if err := cmd.CheckRequiredCommands(); err != nil {
		return fmt.Errorf("required commands not found: %w", err)
	}

	ghqRoot, err := cmd.GetGhqRoot(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ghq root: %w", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "clear" {
		cacher := cache.NewCacheStorage()
		if err := cacher.ClearAll(ctx); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
		logger.Debugln("キャッシュをクリアしました")
		return nil
	}

	cacher := cache.NewCacheStorage()
	// キャッシュの読み込み
	exists, err := cacher.Exists()
	if err != nil {
		return fmt.Errorf("failed to check cache existence: %w", err)
	}
	// キャッシュが見つからない場合、または、syncが必要な場合は、キャッシュを作成する
	sync := len(os.Args) > 1 && os.Args[1] == "sync"
	if !exists || sync {
		if !exists {
			logger.Println("初回実行のため、キャッシュを作成しています...")
		} else {
			logger.Println("キャッシュを更新しています...")
		}
		owners, err := fetchRepositoryOwners(ctx)
		if err != nil {
			return fmt.Errorf("failed to get repository owners: %w", err)
		}
		p := pool.New().WithErrors().WithContext(ctx)
		for _, owner := range owners {
			p.Go(func(ctx context.Context) error {
				repos, err := owner.FetchRepositories(ctx)
				if err != nil {
					return fmt.Errorf("failed to fetch repositories for %s: %w", owner.ID(), err)
				}
				if err := cacher.Save(ctx, owner.ID(), repos); err != nil {
					return fmt.Errorf("failed to save cache: %w", err)
				}
				return nil
			})
		}
		if err := p.Wait(); err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Debugln("キャッシュの作成がキャンセルされました")
				return nil
			}
			logger.Errorln("キャッシュの作成中にエラーが発生しました:", err)
		}
		me, err := core.Me(ctx)
		if err != nil {
			return fmt.Errorf("failed to get Me: %w", err)
		}
		if err := cacher.SaveMD(ctx, me.ID(), time.Now()); err != nil {
			return fmt.Errorf("failed to save cache: %w", err)
		}
		logger.Debugln("キャッシュの作成が完了しました")
	}

	md, err := cacher.LoadMD(ctx)
	if err != nil {
		return fmt.Errorf("failed to load last updated time: %w", err)
	}
	if cache.IsCacheStale(md.LastUpdated, cacheMaxAge) {
		ctx, cancel := context.WithCancel(ctx)
		p := pool.New().WithErrors().WithContext(ctx)
		logger.Debugln("バックグラウンドでキャッシュを更新しています...")
		p.Go(func(ctx context.Context) error {
			owners, err := fetchRepositoryOwners(ctx)
			if err != nil {
				return fmt.Errorf("failed to get repository owners: %w", err)
			}
			cachePool := pool.New().WithErrors().WithContext(ctx)
			for _, owner := range owners {
				cachePool.Go(func(ctx context.Context) error {
					repos, err := owner.FetchRepositories(ctx)
					if err != nil {
						return fmt.Errorf("failed to fetch repositories for %s: %w", owner.ID(), err)
					}
					if err := cacher.Save(ctx, owner.ID(), repos); err != nil {
						return fmt.Errorf("failed to save cache: %w", err)
					}
					return nil
				})
			}
			if err := cachePool.Wait(); err != nil {
				return err
			}
			return nil
		})
		defer func() {
			cancel()
			logger.Debugln("キャッシュの更新を待っています...")
			if err := p.Wait(); err != nil {
				if errors.Is(err, context.Canceled) {
					logger.Debugln("キャッシュの更新がキャンセルされました")
					return
				}
				logger.Errorln("キャッシュの更新中にエラーが発生しました:", err)
			}
			logger.Debugln("キャッシュの更新が完了しました")
		}()
	}

	rx, err := cacher.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load all caches: %w", err)
	}

	// Flatten the cache data
	var allRepos []models.Repository
	for _, repos := range rx {
		allRepos = append(allRepos, repos...)
	}

	// Add local repositories from ghq
	repos, err := (&GHQ{}).FetchRepositories(ctx)
	if err != nil {
		return err
	}
	allRepos = append(allRepos, repos...)

	// Remove duplicates
	allRepos = deduplicateRepositories(allRepos)

	allRepos = checkCloneStatus(allRepos, ghqRoot)

	slices.SortFunc(allRepos, func(a, b models.Repository) int {
		// 自分を後ろに配置
		if a.OrgName == md.Me {
			if b.OrgName != md.Me {
				return -1
			}
		} else if b.OrgName == md.Me {
			return 1
		}
		// Sort by FullPath
		return strings.Compare(a.FullPath(), b.FullPath())
	})

	selected, err := Select(ctx, allRepos)
	if err != nil {
		if errors.Is(err, errRepositoryNotSelected) {
			return nil
		}
		return fmt.Errorf("failed to select repository: %w", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "get" {
		if err := cmd.CloneRepository(ctx, selected.GetGitURL()); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		logger.Debugln("Cloned:", selected.GetGitURL())
		return nil
	}

	fmt.Println(selected.Name)
	return nil
}

type stdlogger struct {
	debug bool
}

type logger interface {
	Println(a ...any)
	Errorln(a ...any)
	Debugln(a ...any)
}

func newLogger(debug bool) *stdlogger {
	return &stdlogger{debug: debug}
}

type loggerKey struct{}

func newContextWithLogger(ctx context.Context, logger logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func loggerFromCtx(ctx context.Context) logger {
	return ctx.Value(loggerKey{}).(*stdlogger)
}

func (l *stdlogger) Debugln(a ...any) {
	if l.debug {
		fmt.Fprintln(os.Stderr, a...)
	}
}

func (l *stdlogger) Println(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func (l *stdlogger) Errorln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

// キャッシュの有効期限（24時間）
const cacheMaxAge = 1 * time.Second

func checkCloneStatus(repos []models.Repository, ghqRoot string) []models.Repository {
	for i, repo := range repos {
		path, _ := repo.GetClonePath(ghqRoot)
		if _, err := os.Stat(path); err == nil {
			repos[i].Cloned = true
		}
	}
	return repos
}

func deduplicateRepositories(repos []models.Repository) []models.Repository {
	seen := make(map[string]bool)
	var result []models.Repository

	for _, repo := range repos {
		// Create a unique key for each repository
		key := repo.FullPath()
		if !seen[key] {
			seen[key] = true
			result = append(result, repo)
		}
	}
	return result
}

func fetchRepositoryOwners(ctx context.Context) ([]RepositoryOwner, error) {
	owners := make([]RepositoryOwner, 0, 100)
	hosts := auth.KnownHosts()

	for _, host := range hosts {
		client, err := github.NewClient(api.ClientOptions{
			Host: host,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub client: %w", err)
		}

		me, err := NewMe(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create Me: %w", err)
		}
		owners = append(owners, me)

		orgs, err := client.FetchOrganizations(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch organizations: %w", err)
		}

		for _, org := range orgs {
			owners = append(owners, NewOrganization(client, org.Login))
		}
	}
	return owners, nil
}

type Core struct {
	mutex         sync.RWMutex
	hosts         []string
	me            *Me
	organizations []*Organization
	ghq           *GHQ
}

func NewCore(ctx context.Context) (*Core, error) {
	hosts := auth.KnownHosts()
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no authenticated hosts found")
	}
	return &Core{
		hosts: hosts,
		ghq:   &GHQ{},
	}, nil
}

func (c *Core) Me(ctx context.Context) (*Me, error) {
	if c.me == nil {
		client, err := github.NewClient(api.ClientOptions{
			Host: c.hosts[0],
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub client: %w", err)
		}
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.me, err = NewMe(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create Me: %w", err)
		}
	}
	return c.me, nil
}

func (c *Core) Organizations(ctx context.Context) ([]*Organization, error) {
	if c.organizations == nil {
		organizations := make([]*Organization, 0, 100)
		hosts := auth.KnownHosts()

		for _, host := range hosts {
			client, err := github.NewClient(api.ClientOptions{
				Host: host,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create GitHub client: %w", err)
			}

			orgs, err := client.FetchOrganizations(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch organizations: %w", err)
			}

			for _, org := range orgs {
				organizations = append(organizations, NewOrganization(client, org.Login))
			}
		}
		c.mutex.Lock()
		c.organizations = organizations
		c.mutex.Unlock()
	}
	return c.organizations, nil
}

func (c *Core) GHQ(ctx context.Context) (*GHQ, error) {
	return c.ghq, nil
}

func (c *Core) RepositoryOwners(ctx context.Context) ([]RepositoryOwner, error) {
	me, err := c.Me(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Me: %w", err)
	}
	organizations, err := c.Organizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}
	owners := make([]RepositoryOwner, 0, len(organizations)+1)
	owners = append(owners, me)
	for _, org := range organizations {
		owners = append(owners, org)
	}
	return owners, nil
}

type RepositoryOwner interface {
	ID() string
	RepositoryFetcher
}

type RepositoryFetcher interface {
	FetchRepositories(ctx context.Context) ([]models.Repository, error)
}

type GHQ struct {
}

func (g *GHQ) FetchRepositories(ctx context.Context) ([]models.Repository, error) {
	// Get all repositories in ghq root
	ghqPaths, err := cmd.ListGhqRepositories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ghq repositories: %w", err)
	}

	allRepos := make([]models.Repository, 0, len(ghqPaths))

	// Add local repositories from ghq
	for _, ghqRepo := range ghqPaths {
		repo, err := ghqRepo.ToRepository()
		if err != nil {
			continue
		}
		allRepos = append(allRepos, repo)
	}
	return allRepos, nil
}

type Me struct {
	name   string
	client *github.Client
}

func NewMe(ctx context.Context, client *github.Client) (*Me, error) {
	user, err := client.FetchUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &Me{
		name:   user,
		client: client,
	}, nil
}

func (m *Me) ID() string {
	return m.name
}

func (m *Me) FetchRepositories(ctx context.Context) ([]models.Repository, error) {
	allRepos := make([]models.Repository, 0, 10000)
	page := 1
	maxAttempts := 100

	for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
		repos, nextPage, err := m.client.FetchUserRepositories(ctx, page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user repositories: %w", err)
		}
		for _, repo := range repos {
			allRepos = append(allRepos, models.Repository{
				Name:    repo.Name,
				OrgName: repo.OrgName,
				Host:    repo.Host,
				HtmlUrl: repo.HtmlUrl,
			})
		}
		page = nextPage
		maxAttempts--
	}
	return allRepos, nil
}

type Organization struct {
	client *github.Client
	org    string
}

func NewOrganization(client *github.Client, org string) *Organization {
	return &Organization{
		client: client,
		org:    org,
	}
}

func (o *Organization) ID() string {
	return o.org
}

func (o *Organization) FetchRepositories(ctx context.Context) ([]models.Repository, error) {
	time.Sleep(1 * time.Second) // Rate limit handling
	allRepos := make([]models.Repository, 0, 100)
	page := 1
	maxAttempts := 100

	for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
		repos, nextPage, err := o.client.FetchRepositories(ctx, github.Organization{Login: o.org}, page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch organization repositories: %w", err)
		}
		for _, repo := range repos {
			allRepos = append(allRepos, models.Repository{
				Name:    repo.Name,
				OrgName: repo.OrgName,
				Host:    repo.Host,
				HtmlUrl: repo.HtmlUrl,
			})
		}
		page = nextPage
		maxAttempts--
	}
	return allRepos, nil
}

func Select(ctx context.Context, repos []models.Repository) (*models.Repository, error) {
	var lines []string
	for _, repo := range repos {
		lines = append(lines, repo.FormattedLine())
	}

	selected, err := cmd.RunSelector(ctx, lines)
	if err != nil {
		return nil, err
	}

	if selected == "" {
		return nil, errRepositoryNotSelected
	}
	for _, repo := range repos {
		if strings.Contains(repo.FormattedLine(), selected) {
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("repository not found")
}

var errRepositoryNotSelected = fmt.Errorf("repository not selected")
