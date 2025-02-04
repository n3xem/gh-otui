package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/cli/go-gh/v2/pkg/api"
)

// Repository は表示するリポジトリの情報を保持する構造体
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Stars       int    `json:"stargazers_count"`
	OrgName     string
	Cloned      bool
}

// pecoで選択するための文字列を生成する関数
func formatRepoLine(repo Repository) string {
	cloneStatus := " "
	if repo.Cloned {
		cloneStatus = "✓"
	}
	return fmt.Sprintf("%s %s/%s", cloneStatus, repo.OrgName, repo.Name)
}

func main() {
	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	var orgs []struct {
		Login string `json:"login"`
	}

	err = client.Get("user/orgs", &orgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	var allRepos []Repository
	for _, org := range orgs {
		var repos []Repository
		err = client.Get(fmt.Sprintf("orgs/%s/repos", org.Login), &repos)
		if err != nil {
			fmt.Printf("リポジトリの取得に失敗: %v\n", err)
			continue
		}
		for i := range repos {
			repos[i].OrgName = org.Login
		}
		allRepos = append(allRepos, repos...)
	}

	// リポジトリごとにクローン済みかチェック
	for i, repo := range allRepos {
		repoPath := fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), "ghq/github.com", repo.OrgName)
		if _, err := os.Stat(filepath.Join(repoPath, repo.Name)); err == nil {
			allRepos[i].Cloned = true
		}
	}

	// pecoに渡す文字列を準備
	var lines []string
	for _, repo := range allRepos {
		lines = append(lines, formatRepoLine(repo))
	}

	// pecoコマンドを実行
	cmd := exec.Command("peco")
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running peco:", err)
		return
	}

	// 選択された行を処理
	selected := strings.TrimSpace(string(out))
	if selected == "" {
		fmt.Println("No selection made")
		return
	}

	// 選択されたリポジトリを特定
	for _, repo := range allRepos {
		repoLine := formatRepoLine(repo)
		// 先頭の空白、チェックマーク、その後の空白を確実に除去
		trimmedRepoLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(repoLine), "✓"))
		trimmedSelected := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(selected), "✓"))

		if trimmedRepoLine == trimmedSelected {
			repoPath := fmt.Sprintf("%s/ghq/github.com/%s/%s",
				os.Getenv("HOME"),
				repo.OrgName,
				repo.Name,
			)

			if !repo.Cloned {
				cmd := exec.Command("ghq", "get", fmt.Sprintf("https://github.com/%s/%s",
					repo.OrgName,
					repo.Name,
				))
				if output, err := cmd.CombinedOutput(); err != nil {
					fmt.Printf("Error cloning repository: %v\nOutput: %s\n", err, string(output))
					return
				}
			}

			// パスを出力して終了
			fmt.Println(repoPath)
			return
		}
	}
}
