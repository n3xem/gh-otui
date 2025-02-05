package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/n3xem/gh-otui/cache"
	"github.com/n3xem/gh-otui/cmd"
	"github.com/n3xem/gh-otui/github"
	"github.com/n3xem/gh-otui/models"
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
				if err := cmd.CloneRepository(repo.GetGitURL()); err != nil {
					return err
				}
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
		client, err := github.NewClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		orgs, err := client.FetchOrganizations()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Convert github.Organization to models.Organization
		modelOrgs := make([]models.Organization, len(orgs))
		for i, org := range orgs {
			modelOrgs[i] = models.Organization{Login: org.Login}
		}

		githubRepos := client.FetchRepositories(orgs)

		// Convert github.Repository to models.Repository
		modelRepos := make([]models.Repository, len(githubRepos))
		for i, repo := range githubRepos {
			modelRepos[i] = models.Repository{
				Name:    repo.Name,
				OrgName: repo.OrgName,
				HtmlUrl: repo.HtmlUrl,
				Host:    repo.Host,
			}
		}

		if err := cache.SaveCache(modelRepos); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cache saved successfully")
		return
	}

	// Load and process repositories
	allRepos, err := cache.LoadCache()
	if err != nil {
		fmt.Println("Cache not found. Please create cache with:")
		fmt.Printf("%s --cache\n", os.Args[0])
		os.Exit(1)
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
