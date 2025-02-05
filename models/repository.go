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
