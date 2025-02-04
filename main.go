package main

import (
	"fmt"
	"os"
	"path/filepath"

	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
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

// Model はアプリケーションの状態を保持する構造体
type Model struct {
	repos    []Repository
	cursor   int
	selected int
}

// Init は初期化時に実行される
func (m *Model) Init() tea.Cmd {
	// リポジトリごとにクローン済みかチェック
	for i, repo := range m.repos {
		repoPath := fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), ".ghq/github.com", repo.OrgName)
		if _, err := os.Stat(filepath.Join(repoPath, repo.Name)); err == nil {
			m.repos[i].Cloned = true
		}
	}
	return nil
}

// Update はイベントに応じて状態を更新する
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.cursor
			if !m.repos[m.cursor].Cloned {
				err := exec.Command("ghq", "get", fmt.Sprintf("https://github.com/%s/%s",
					m.repos[m.cursor].OrgName,
					m.repos[m.cursor].Name,
				)).Run()
				if err != nil {
					return m, nil
				}
				m.repos[m.cursor].Cloned = true
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View は画面の表示を定義する
func (m *Model) View() string {
	if m.selected >= 0 {
		return fmt.Sprintf("%s/.ghq/github.com/%s/%s",
			os.Getenv("HOME"),
			m.repos[m.selected].OrgName,
			m.repos[m.selected].Name,
		)
	}

	s := "リポジトリ一覧:\n\n"

	for i, repo := range m.repos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		cloneStatus := " "
		if repo.Cloned {
			cloneStatus = "✓"
		}
		s += fmt.Sprintf("%s %s %s/%s\n", cursor, cloneStatus, repo.OrgName, repo.Name)
	}

	s += "\n(↑/↓ または j/k で移動, Enter で選択, q で終了)\n"
	return s
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

	initialModel := Model{
		repos:    allRepos,
		cursor:   0,
		selected: -1,
	}

	p := tea.NewProgram(&initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
