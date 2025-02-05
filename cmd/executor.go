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
	requiredCommands := []string{"gh", "peco", "ghq"}
	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("%s command not found", cmd)
		}
	}
	return nil
}

func RunPeco(lines []string) (string, error) {
	cmd := exec.Command("peco")
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
