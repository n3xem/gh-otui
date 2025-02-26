package main

import (
	"fmt"
	"os"
	"strings"
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

func processSelectedRepository(repos []models.Repository, selected string, ghqRoot string) error {
	for _, repo := range repos {
		repoLine := repo.FormattedLine()
		trimmedRepoLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(repoLine), "✓"))
		trimmedSelected := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(selected), "✓"))

		if trimmedRepoLine == trimmedSelected {
			if !repo.Cloned {
				s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				s.Suffix = fmt.Sprintf(" Cloning %s/%s...", repo.OrgName, repo.Name)
				s.Start()
				if err := cmd.CloneRepository(repo.GetGitURL()); err != nil {
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

func main() {
	if err := cmd.CheckRequiredCommands(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	ghqRoot, err := cmd.GetGhqRoot()
	if err != nil {
		fmt.Printf("Failed to get ghq root: %v\n", err)
		os.Exit(1)
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
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			s.Stop()

			// 自分のリポジトリを取得
			s.Suffix = " Fetching user repositories..."
			s.Start()
			page := 1
			maxAttempts := 100

			for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
				repos, nextPage, err := client.FetchUserRepositories(page)
				if err != nil {
					s.Stop()
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				allRepos = append(allRepos, repos...)
				page = nextPage
				maxAttempts--
			}
			s.Stop()

			s.Suffix = " Fetching organizations..."
			s.Start()
			orgs, err := client.FetchOrganizations()
			if err != nil {
				s.Stop()
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			s.Stop()

			s.Suffix = " Fetching repositories..."
			s.Start()
			page = 1
			maxAttempts = 100 // 安全のための最大ページ数

			for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 { // 追加の安全対策
				repos, nextPage, err := client.FetchRepositories(orgs, page)
				if err != nil {
					s.Stop()
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				allRepos = append(allRepos, repos...)
				page = nextPage
				maxAttempts--
			}

			if maxAttempts == 0 {
				s.Stop()
				fmt.Printf("Error: リポジトリの取得が上限に達しました\n")
				os.Exit(1)
			}
			s.Stop()

			s.Suffix = " Saving cache..."
			s.Start()
			if err := cache.SaveCache(allRepos); err != nil {
				s.Stop()
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			s.Stop()
		}
		fmt.Println("Cache saved successfully")
		return
	}

	// Load and process repositories
	githubRepos, err := cache.LoadCache()
	if err != nil {
		fmt.Println("Cache not found. Please create cache with:")
		fmt.Printf("%s --cache\n", os.Args[0])
		os.Exit(1)
	}

	// Get all repositories in ghq root
	ghqPaths, err := cmd.ListGhqRepositories()
	if err != nil {
		fmt.Printf("Failed to get ghq repositories: %v\n", err)
		os.Exit(1)
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

	selected, err := cmd.RunSelector(lines)
	if err != nil {
		fmt.Printf("Failed to run selector: %v\n", err)
		os.Exit(1)
	}

	if selected == "" {
		fmt.Println("No repository selected")
		return
	}

	if err := processSelectedRepository(allRepos, selected, ghqRoot); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
