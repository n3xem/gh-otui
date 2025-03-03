package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/n3xem/gh-otui/models"
)

func GetGhqRoot() (string, error) {
	cmd := exec.Command("ghq", "root")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ghq root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func CheckRequiredCommands() error {
	requiredCommands := []string{"gh", "ghq"}
	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("%s command not found", cmd)
		}
	}

	// Check for peco or fzf
	if _, err := exec.LookPath("peco"); err != nil {
		if _, err := exec.LookPath("fzf"); err != nil {
			return fmt.Errorf("neither peco nor fzf command found")
		}
	}
	return nil
}

func RunSelector(lines []string) (string, error) {
	selector := os.Getenv("GH_OTUI_SELECTOR")
	if selector == "" {
		if _, err := exec.LookPath("peco"); err == nil {
			selector = "peco"
		} else if _, err := exec.LookPath("fzf"); err == nil {
			selector = "fzf"
		} else {
			return "", fmt.Errorf("neither peco nor fzf command found")
		}
	}

	cmd := exec.Command(selector)
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func CloneRepository(gitURL string) error {
	cmd := exec.Command("ghq", "get", gitURL)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %s: %w", string(output), err)
	}
	return nil
}

// ClonedGhqRepository represents a git repository managed by ghq
type ClonedGhqRepository struct {
	FullPath string
}

// ListGhqRepositories returns a list of all repositories managed by ghq
func ListGhqRepositories() ([]ClonedGhqRepository, error) {
	cmd := exec.Command("ghq", "list", "--full-path")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	var repositories []ClonedGhqRepository
	for _, path := range paths {
		if path != "" {
			repositories = append(repositories, ClonedGhqRepository{FullPath: path})
		}
	}
	return repositories, nil
}

func (c ClonedGhqRepository) ToRepository() (models.Repository, error) {
	parts := strings.Split(c.FullPath, "/")
	if len(parts) < 4 {
		return models.Repository{}, fmt.Errorf("invalid repository path: %s", c.FullPath)
	}

	repoName := parts[len(parts)-1]
	orgName := parts[len(parts)-2]
	host := parts[len(parts)-3]

	return models.Repository{
		Name:    repoName,
		OrgName: orgName,
		Host:    host,
		HtmlUrl: fmt.Sprintf("https://%s/%s/%s", host, orgName, repoName),
		Cloned:  true,
	}, nil
}
