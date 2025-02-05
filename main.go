package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gh-otui/cache"
	"gh-otui/cmd"
	"gh-otui/github"
	"gh-otui/models"

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
		s.Suffix = " Connecting to GitHub..."
		s.Start()

		client, err := github.NewClient()
		if err != nil {
			s.Stop()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
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
		repos := client.FetchRepositories(orgs)
		s.Stop()

		s.Suffix = " Saving cache..."
		s.Start()
		if err := cache.SaveCache(repos); err != nil {
			s.Stop()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		s.Stop()
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

	// Convert github.Repository to models.Repository
	allRepos := make([]models.Repository, len(githubRepos))
	for i, repo := range githubRepos {
		allRepos[i] = models.Repository{
			Name:    repo.Name,
			OrgName: repo.OrgName,
			HtmlUrl: repo.HtmlUrl,
			Host:    repo.Host,
		}
	}

	allRepos = checkCloneStatus(allRepos, ghqRoot)

	var lines []string
	for _, repo := range allRepos {
		lines = append(lines, repo.FormattedLine())
	}

	selected, err := cmd.RunPeco(lines)
	if err != nil {
		fmt.Printf("Failed to run peco: %v\n", err)
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
