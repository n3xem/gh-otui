package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
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

func processSelectedRepository(ctx context.Context, repos []models.Repository, selected string, ghqRoot string) error {
	for _, repo := range repos {
		repoLine := repo.FormattedLine()
		trimmedRepoLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(repoLine), "✓"))
		trimmedSelected := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(selected), "✓"))

		if trimmedRepoLine == trimmedSelected {
			if !repo.Cloned {
				s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				s.Suffix = fmt.Sprintf(" Cloning %s/%s...", repo.OrgName, repo.Name)
				s.Start()
				if err := cmd.CloneRepository(ctx, repo.GetGitURL()); err != nil {
					s.Stop()
					return err
				}
				s.Stop()
			}
			clonePath, err := repo.GetClonePath(ghqRoot)
			if err != nil {
				return fmt.Errorf("failed to get repository path: %w", err)
			}
			fmt.Println(clonePath)
			return nil
		}
	}
	return fmt.Errorf("repository not found")
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
		hosts := auth.KnownHosts()
		var allRepos []github.Repository
		for _, host := range hosts {
			s.Suffix = fmt.Sprintf(" Connecting to %s...", host)
			s.Start()
			client, err := github.NewClient(api.ClientOptions{
				Host: host,
			})
			if err != nil {
				s.Stop()
				return err
			}
			s.Stop()

			// 自分のリポジトリを取得
			s.Suffix = " Fetching user repositories..."
			s.Start()
			page := 1
			maxAttempts := 100

			for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
				repos, nextPage, err := client.FetchUserRepositories(ctx, page)
				if err != nil {
					s.Stop()
					return err
				}
				allRepos = append(allRepos, repos...)
				page = nextPage
				maxAttempts--
			}
			s.Stop()

			s.Suffix = " Fetching organizations..."
			s.Start()
			orgs, err := client.FetchOrganizations(ctx)
			if err != nil {
				s.Stop()
				return err
			}
			s.Stop()

			s.Suffix = " Fetching repositories..."
			s.Start()
			page = 1
			maxAttempts = 100 // 安全のための最大ページ数

			for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 { // 追加の安全対策
				repos, nextPage, err := client.FetchRepositories(ctx, orgs, page)
				if err != nil {
					s.Stop()
					return err
				}
				allRepos = append(allRepos, repos...)
				page = nextPage
				maxAttempts--
			}

			if maxAttempts == 0 {
				s.Stop()
				return fmt.Errorf("リポジトリの取得が上限に達しました")
			}
			s.Stop()

			s.Suffix = " Saving cache..."
			s.Start()
			if err := cache.SaveCache(allRepos); err != nil {
				s.Stop()
				return err
			}
			s.Stop()
		}
		fmt.Println("Cache saved successfully")
		return nil
	}

	// Load and process repositories
	githubRepos, err := cache.LoadCache()
	if err != nil {
		return fmt.Errorf("cache not found. Please create cache with: %s --cache", os.Args[0])
	}

	// Get all repositories in ghq root
	ghqPaths, err := cmd.ListGhqRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ghq repositories: %w", err)
	}

	// Convert ghq paths directly to repositories
	var allRepos []models.Repository
	// Convert github.Repository to models.Repository
	for _, repo := range githubRepos {
		allRepos = append(allRepos, models.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			Host:    repo.Host,
			HtmlUrl: repo.HtmlUrl,
		})
	}

	// Add local repositories from ghq
	for _, ghqRepo := range ghqPaths {
		repo, err := ghqRepo.ToRepository()
		if err != nil {
			continue // Skip invalid repository paths
		}
		allRepos = append(allRepos, repo)
	}

	// Remove duplicates
	allRepos = deduplicateRepositories(allRepos)

	allRepos = checkCloneStatus(allRepos, ghqRoot)

	var lines []string
	for _, repo := range allRepos {
		lines = append(lines, repo.FormattedLine())
	}

	selected, err := cmd.RunSelector(ctx, lines)
	if err != nil {
		return fmt.Errorf("failed to run selector: %w", err)
	}

	if selected == "" {
		return fmt.Errorf("no repository selected")
	}

	if err := processSelectedRepository(ctx, allRepos, selected, ghqRoot); err != nil {
		return fmt.Errorf("error processing selected repository: %w", err)
	}
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
