package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/n3xem/gh-otui/cache"
	"github.com/n3xem/gh-otui/cmd"
	"github.com/n3xem/gh-otui/github"
	"github.com/n3xem/gh-otui/models"

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

func run(ctx context.Context) error {
	if err := cmd.CheckRequiredCommands(); err != nil {
		return err
	}

	ghqRoot, err := cmd.GetGhqRoot(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ghq root: %w", err)
	}

	// Handle cache creation
	if len(os.Args) > 1 && os.Args[1] == "--cache" {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Fetching repositories..."
		s.Start()
		hosts := auth.KnownHosts()
		gihubClients := make([]*github.Client, 0, len(hosts))
		for _, host := range hosts {
			client, err := github.NewClient(api.ClientOptions{
				Host: host,
			})
			if err != nil {
				return err
			}
			gihubClients = append(gihubClients, client)
		}
		bufErrors := make([]error, 0, 8)
		for _, client := range gihubClients {
			g, err := github.FetchUserRepositories(ctx, client)
			if err != nil {
				bufErrors = append(bufErrors, err)
				continue
			}
			if err := cache.Save(ctx, g); err != nil {
				bufErrors = append(bufErrors, err)
				continue
			}
		}
		for _, client := range gihubClients {
			orgs, err := github.NewOrganizations(ctx, client)
			if err != nil {
				return err
			}
			for _, org := range orgs {
				g, err := org.FetchRepositories(ctx)
				if err != nil {
					bufErrors = append(bufErrors, err)
					continue
				}
				if err := cache.Save(ctx, g); err != nil {
					bufErrors = append(bufErrors, err)
					continue
				}
			}
		}
		for _, client := range gihubClients {
			gs, err := github.FetchCollaboratingRepositories(ctx, client)
			if err != nil {
				bufErrors = append(bufErrors, err)
				continue
			}
			for g := range gs {
				if err := cache.Save(ctx, g); err != nil {
					bufErrors = append(bufErrors, err)
					continue
				}
			}
		}
		s.Stop()

		if len(bufErrors) > 0 {
			return errors.Join(bufErrors...)
		}

		fmt.Fprintln(os.Stderr, "Cache saved successfully")
		return nil
	}

	allRepos := make([]models.Repository, 0)

	// Load and process repositories
	repositoryGroups, err := cache.FetchRepositories(ctx)
	if err != nil {
		return fmt.Errorf("cache not found. Please create cache with: %s --cache", os.Args[0])
	}
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
		return fmt.Errorf("error selecting repository: %w", err)
	}

	if !selected.Cloned {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Cloning %s/%s...", selected.OrgName, selected.Name)
		s.Start()
		if err := cmd.CloneRepository(ctx, selected.GetGitURL()); err != nil {
			s.Stop()
			return err
		}
		s.Stop()
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
