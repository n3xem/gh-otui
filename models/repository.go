package models

import (
	"fmt"
	"path/filepath"
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

func (r Repository) GetClonePath(ghqRoot string) (string, error) {
	return filepath.Join(ghqRoot, r.Host, r.OrgName, r.Name), nil
}

func (r Repository) GetGitURL() string {
	return fmt.Sprintf("git@%s:%s/%s", r.Host, r.OrgName, r.Name)
}

func (r Repository) FormattedLine() string {
	cloneStatus := " "
	if r.Cloned {
		cloneStatus = "✓"
	}
	return fmt.Sprintf("%s %s/%s/%s", cloneStatus, r.Host, r.OrgName, r.Name)
}

type RepositoryGroup struct {
	host         string
	organization string
	repositories []Repository
}

func NewRepositoryGroup(repositories ...Repository) (*RepositoryGroup, error) {
	if len(repositories) == 0 {
		return nil, fmt.Errorf("no repositories provided")
	}

	g := &RepositoryGroup{
		host:         repositories[0].Host,
		organization: repositories[0].OrgName,
		repositories: make([]Repository, 0, len(repositories)),
	}

	for _, repo := range repositories {
		if err := g.Add(repo); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func (g *RepositoryGroup) Host() string {
	return g.host
}

func (g *RepositoryGroup) Organization() string {
	return g.organization
}

func (g *RepositoryGroup) Add(repo Repository) error {
	if repo.Host != g.host {
		return fmt.Errorf("repository host %s does not match group host %s", repo.Host, g.host)
	}
	if repo.OrgName != g.organization {
		return fmt.Errorf("repository organization %s does not match group organization %s", repo.OrgName, g.organization)
	}
	g.repositories = append(g.repositories, repo)
	return nil
}

func (g *RepositoryGroup) Repositories() []Repository {
	return g.repositories
}
