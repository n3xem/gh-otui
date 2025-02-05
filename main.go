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

// ghq rootのパスを取得する関数
func getGhqRoot() (string, error) {
	cmd := exec.Command("ghq", "root")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ghq root: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
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
func (r Repository) GetClonePath(ghqRoot string) (string, error) {
	return filepath.Join(ghqRoot, r.Host, r.OrgName, r.Name), nil
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

func checkCloneStatus(repos []Repository, ghqRoot string) []Repository {
	for i, repo := range repos {
		path, _ := repo.GetClonePath(ghqRoot) // エラーは無視して続行
		if _, err := os.Stat(path); err == nil {
			repos[i].Cloned = true
		}
	}
	return repos
}

func processSelectedRepository(repos []Repository, selected string, ghqRoot string) {
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
			clonePath, err := repo.GetClonePath(ghqRoot)
			if err != nil {
				handleError(err, "リポジトリパスの取得に失敗")
			}
			fmt.Println(clonePath)
			return
		}
	}
}

func getCachePath() string {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "gh", "extensions", "gh-otui")
	return filepath.Join(configDir, "cache.json")
}

func main() {
	// 必要なコマンドの存在確認
	requiredCommands := []string{"gh", "peco", "ghq"}
	for _, cmd := range requiredCommands {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Printf("%sコマンドが見つかりません。インストールしてください。\n", cmd)
			os.Exit(1)
		}
	}

	ghqRoot, err := getGhqRoot()
	if err != nil {
		fmt.Printf("ghq rootの取得に失敗: %v\n", err)
		os.Exit(1)
	}

	loadCache := func() ([]Repository, error) {
		cacheData, err := os.ReadFile(getCachePath())
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

		// キャッシュディレクトリを作成
		cacheDir := filepath.Dir(getCachePath())
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Printf("キャッシュディレクトリの作成に失敗: %v\n", err)
			return
		}

		if err := os.WriteFile(getCachePath(), cacheData, 0644); err != nil {
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

	allRepos = checkCloneStatus(allRepos, ghqRoot)

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

	processSelectedRepository(allRepos, selected, ghqRoot)
}
