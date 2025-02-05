package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Repository struct {
	Name    string `json:"name"`
	OrgName string
	HtmlUrl string `json:"html_url"`
	Host    string
	Cloned  bool
}

type Organization struct {
	Login string `json:"login"`
}

// pecoで選択するための文字列を生成する関数
func formatRepoLine(repo Repository) string {
	cloneStatus := " "
	if repo.Cloned {
		cloneStatus = "✓"
	}
	return fmt.Sprintf("%s %s/%s/%s", cloneStatus, repo.Host, repo.OrgName, repo.Name)
}

// リポジトリに関連するメソッドを追加
func (r Repository) GetClonePath() string {
	return filepath.Join(os.Getenv("HOME"), "ghq", r.Host, r.OrgName, r.Name)
}

func (r Repository) GetGitURL() string {
	return fmt.Sprintf("git@%s:%s/%s", r.Host, r.OrgName, r.Name)
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}

func fetchOrganizations(client *api.RESTClient) []Organization {
	var orgs []Organization
	err := client.Get("user/orgs", &orgs)
	handleError(err, "組織の取得に失敗")
	return orgs
}

func fetchRepositories(client *api.RESTClient, orgs []Organization) []Repository {
	var allRepos []Repository
	for _, org := range orgs {
		var repos []Repository
		err := client.Get(fmt.Sprintf("orgs/%s/repos?per_page=100", org.Login), &repos)
		if err != nil {
			fmt.Printf("リポジトリの取得に失敗 (%s): %v\n", org.Login, err)
			continue
		}
		for i := range repos {
			repos[i].OrgName = org.Login
			hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
			repos[i].Host = strings.Split(hostWithPath, "/")[0]
		}
		allRepos = append(allRepos, repos...)
	}
	return allRepos
}

func checkCloneStatus(repos []Repository) []Repository {
	for i, repo := range repos {
		if _, err := os.Stat(repo.GetClonePath()); err == nil {
			repos[i].Cloned = true
		}
	}
	return repos
}

func processSelectedRepository(repos []Repository, selected string) {
	for _, repo := range repos {
		repoLine := formatRepoLine(repo)
		trimmedRepoLine := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(repoLine), "✓"))
		trimmedSelected := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(selected), "✓"))

		if trimmedRepoLine == trimmedSelected {
			if !repo.Cloned {
				cmd := exec.Command("ghq", "get", repo.GetGitURL())
				if output, err := cmd.CombinedOutput(); err != nil {
					handleError(err, fmt.Sprintf("リポジトリのクローンに失敗\nOutput: %s", string(output)))
				}
			}
			fmt.Println(repo.GetClonePath())
			return
		}
	}
}

func main() {
	loadCache := func() ([]Repository, error) {
		cacheData, err := os.ReadFile("cache")
		if err != nil {
			return nil, err
		}
		var repos []Repository
		if err := json.Unmarshal(cacheData, &repos); err != nil {
			return nil, err
		}
		return repos, nil
	}

	saveCache := func(repos []Repository) {
		cacheData, err := json.Marshal(repos)
		if err != nil {
			fmt.Printf("キャッシュの作成に失敗: %v\n", err)
			return
		}
		if err := os.WriteFile("cache", cacheData, 0644); err != nil {
			fmt.Printf("キャッシュの保存に失敗: %v\n", err)
			return
		}
		fmt.Println("キャッシュを保存しました")
	}

	// --cache フラグが指定された場合
	if len(os.Args) > 1 && os.Args[1] == "--cache" {
		client, err := api.DefaultRESTClient()
		handleError(err, "GitHub APIクライアントの初期化に失敗")

		orgs := fetchOrganizations(client)
		allRepos := fetchRepositories(client, orgs)
		saveCache(allRepos)
		return
	}

	// キャッシュからデータを読み込む
	allRepos, err := loadCache()
	if err != nil {
		fmt.Println("キャッシュが見つかりません。以下のコマンドでキャッシュを作成してください：")
		fmt.Printf("%s --cache\n", os.Args[0])
		os.Exit(1)
	}

	allRepos = checkCloneStatus(allRepos)

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
	handleError(err, "pecoの実行に失敗")

	selected := strings.TrimSpace(string(out))
	if selected == "" {
		fmt.Println("選択されていません")
		return
	}

	processSelectedRepository(allRepos, selected)
}
