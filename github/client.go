package github

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

type Organization struct {
	Login string `json:"login"`
}

type Repository struct {
	Name    string `json:"name"`
	HtmlUrl string `json:"html_url"`
	OrgName string
	Host    string
}

type Client struct {
	client *api.RESTClient
}

func NewClient() (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GitHub API client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) FetchOrganizations() ([]Organization, error) {
	var orgs []Organization
	err := c.client.Get("user/orgs", &orgs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch organizations: %w", err)
	}
	return orgs, nil
}

func (c *Client) FetchRepositories(orgs []Organization) []Repository {
	var allRepos []Repository
	for _, org := range orgs {
		var repos []Repository
		err := c.client.Get(fmt.Sprintf("orgs/%s/repos?per_page=100", org.Login), &repos)
		if err != nil {
			fmt.Printf("Failed to fetch repositories (%s): %v\n", org.Login, err)
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
