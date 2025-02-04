package main

import (
	"fmt"

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
}

// Model はアプリケーションの状態を保持する構造体
type Model struct {
	repos    []Repository
	cursor   int
	selected int
}

// Init は初期化時に実行される
func (m Model) Init() tea.Cmd {
	return nil
}

// Update はイベントに応じて状態を更新する
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		}
	}
	return m, nil
}

// View は画面の表示を定義する
func (m Model) View() string {
	s := "リポジトリ一覧:\n\n"

	for i, repo := range m.repos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s/%s\n", cursor, repo.OrgName, repo.Name)
	}

	s += "\n"
	if m.selected >= 0 && m.selected < len(m.repos) {
		repo := m.repos[m.selected]
		s += fmt.Sprintf("選択したリポジトリの詳細:\n")
		s += fmt.Sprintf("組織: %s\n", repo.OrgName)
		s += fmt.Sprintf("名前: %s\n", repo.Name)
		s += fmt.Sprintf("説明: %s\n", repo.Description)
		s += fmt.Sprintf("言語: %s\n", repo.Language)
		s += fmt.Sprintf("スター数: %d\n", repo.Stars)
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

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
