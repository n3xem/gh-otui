package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/n3xem/gh-otui/cache"
	"github.com/n3xem/gh-otui/cmd"
	"github.com/n3xem/gh-otui/github"
	"github.com/n3xem/gh-otui/models"
	"github.com/sourcegraph/conc/pool"

	"github.com/briandowns/spinner"
)

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
		key := fmt.Sprintf("%s/%s/%s", repo.Host, repo.OrgName, repo.Name)
		if !seen[key] {
			seen[key] = true
			result = append(result, repo)
		}
	}
	return result
}

func updateCache(ctx context.Context) (updated bool, err error) {
	hosts := auth.KnownHosts()
	gihubClients := make([]*github.Client, 0, len(hosts))
	for _, host := range hosts {
		client, err := github.NewClient(api.ClientOptions{
			Host: host,
		})
		if err != nil {
			return false, err
		}
		gihubClients = append(gihubClients, client)
	}
	p := pool.NewWithResults[[]*models.RepositoryGroup]().WithErrors().WithContext(ctx)
	// 自分のリポジトリを取得
	p.Go(func(ctx context.Context) ([]*models.RepositoryGroup, error) {
		gp := pool.NewWithResults[*models.RepositoryGroup]().WithErrors().WithContext(ctx)
		for _, client := range gihubClients {
			gp.Go(func(ctx context.Context) (*models.RepositoryGroup, error) {
				g, err := github.FetchUserRepositories(ctx, client)
				if err != nil {
					return nil, err
				}
				if err := cache.Save(ctx, g); err != nil {
					return nil, err
				}
				return g, nil
			})
		}
		return gp.Wait()
	})
	// organizationsのリポジトリを取得
	for _, client := range gihubClients {
		p.Go(func(ctx context.Context) ([]*models.RepositoryGroup, error) {
			orgs, err := github.NewOrganizations(ctx, client)
			if err != nil {
				return nil, err
			}
			gp := pool.NewWithResults[*models.RepositoryGroup]().WithErrors().WithContext(ctx)
			for _, org := range orgs {
				gp.Go(func(ctx context.Context) (*models.RepositoryGroup, error) {
					g, err := org.FetchRepositories(ctx)
					if err != nil {
						return nil, err
					}
					if err := cache.Save(ctx, g); err != nil {
						return nil, err
					}
					return g, nil
				})
			}
			return gp.Wait()
		})
	}
	// 自分がcollaboratorであるリポジトリを取得
	for _, client := range gihubClients {
		p.Go(func(ctx context.Context) ([]*models.RepositoryGroup, error) {
			gs, err := github.FetchCollaboratingRepositories(ctx, client)
			if err != nil {
				return nil, err
			}
			gp := pool.NewWithResults[*models.RepositoryGroup]().WithErrors().WithContext(ctx)
			for g := range gs {
				gp.Go(func(ctx context.Context) (*models.RepositoryGroup, error) {
					if err := cache.Save(ctx, g); err != nil {
						return nil, err
					}
					return g, nil
				})
			}
			return gp.Wait()
		})
	}

	gg, err := p.Wait()
	gs := flatten(gg)
	someCached := slices.ContainsFunc(gs, func(g *models.RepositoryGroup) bool {
		return g != nil
	})
	if !someCached {
		return false, nil
	}
	e := cache.Done(ctx)
	err = errors.Join(err, e)
	return e == nil, err
}

func flatten[T any](slices [][]T) []T {
	length := 0
	for _, slice := range slices {
		length += len(slice)
	}
	results := make([]T, 0, length)
	for _, slice := range slices {
		results = append(results, slice...)
	}
	return results
}

func loading(msg string, f func() error) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + msg
	s.Start()
	defer s.Stop()
	return f()
}

func run(ctx context.Context) error {
	if err := cmd.CheckRequiredCommands(); err != nil {
		return err
	}

	ghqRoot, err := cmd.GetGhqRoot(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ghq root: %w", err)
	}

	md, err := cache.LoadMetadata(ctx)
	if err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}

	if !md.Initialized() {
		// 同期的なキャッシュ更新
		var updated bool
		err := loading("Fetching repositories...", func() error {
			u, err := updateCache(ctx)
			updated = u
			return err
		})
		if !updated {
			return err
		}
		// 少なくとも１つキャッシュが更新されたなら続行する。
	}

	if md.Initialized() && md.IsStale() {
		// 非同期的なキャッシュ更新
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		p := pool.New().WithErrors().WithContext(ctx)
		p.Go(func(ctx context.Context) error {
			_, err := updateCache(ctx)
			return err
		})
		defer func() {
			cancel()
			if err := p.Wait(); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}

	// Load and process repositories
	repositoryGroups, err := cache.FetchRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch repositories: %w", err)
	}
	allRepos := make([]models.Repository, 0)
	for _, repos := range repositoryGroups {
		allRepos = append(allRepos, repos.Repositories()...)
	}

	// Add local repositories from ghq
	ghqRepos, err := cmd.FetchGHQRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch ghq repositories: %w", err)
	}
	allRepos = append(allRepos, ghqRepos...)

	// Remove duplicates
	allRepos = deduplicateRepositories(allRepos)

	allRepos = checkCloneStatus(allRepos, ghqRoot)

	selected, err := cmd.Select(ctx, allRepos)
	if err != nil {
		if errors.Is(err, cmd.ErrRepositoryNotSelected) {
			return nil
		}
		return fmt.Errorf("error selecting repository: %w", err)
	}

	if !selected.Cloned {
		err := loading(
			fmt.Sprintf("Cloning %s/%s...", selected.OrgName, selected.Name),
			func() error {
				return cmd.CloneRepository(ctx, selected.GetGitURL())
			})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	clonePath, err := selected.GetClonePath(ghqRoot)
	if err != nil {
		return fmt.Errorf("failed to get repository path: %w", err)
	}

	fmt.Println(clonePath)
	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()
	if err := run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		cancel()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
