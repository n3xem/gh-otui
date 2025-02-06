package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
