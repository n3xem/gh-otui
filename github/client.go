package github

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"net/url"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/n3xem/gh-otui/models"
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
	host   string
}

func (c *Client) fetchOrgRepositories(ctx context.Context, org string, page int) (repos []Repository, nextPage int, err error) {
	var allRepos []Repository

	resp, err := c.client.RequestWithContext(ctx, "GET", fmt.Sprintf("orgs/%s/repos?per_page=100&page=%d", org, page), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch organization repositories for %s: %w", org, err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal organization repositories for %s: %w", org, err)
	}

	// Linkヘッダーの処理
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := strings.Split(linkHeader, ",")
		for _, link := range links {
			if strings.Contains(link, `rel="next"`) {
				parts := strings.Split(link, ";")
				urlPart := strings.Trim(parts[0], " <>")
				parsedURL, err := url.Parse(urlPart)
				if err != nil {
					continue
				}
				query := parsedURL.Query()
				if pageStr := query.Get("page"); pageStr != "" {
					nextPage, err = strconv.Atoi(pageStr)
					if err != nil {
						continue
					}
				}
			}
		}
	}

	for i := range repos {
		repos[i].OrgName = org
		hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
		repos[i].Host = strings.Split(hostWithPath, "/")[0]
	}
	allRepos = append(allRepos, repos...)
	return allRepos, nextPage, nil
}

func FetchCollaboratingRepositories(ctx context.Context, client *Client) (iter.Seq[*models.RepositoryGroup], error) {
	page := 1
	maxAttempts := 100
	allRepos := make([]models.Repository, 0, 10000)
	type key struct {
		host string
		org  string
	}
	groups := make(map[key]*models.RepositoryGroup)
	for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
		repos, nextPage, err := client.fetchUserRepositories(ctx, affiliationCollaborator, page)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			key := key{
				host: repo.Host,
				org:  repo.OrgName,
			}
			if group, ok := groups[key]; ok {
				err := group.Add(models.Repository{
					Name:    repo.Name,
					OrgName: repo.OrgName,
					Host:    repo.Host,
					HtmlUrl: repo.HtmlUrl,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to add repository to group: %w", err)
				}
			} else {
				group, err := models.NewRepositoryGroup(models.Repository{
					Name:    repo.Name,
					OrgName: repo.OrgName,
					Host:    repo.Host,
					HtmlUrl: repo.HtmlUrl,
				})
				if err != nil {
					return nil, fmt.Errorf("failed to create repository group: %w", err)
				}
				groups[key] = group
			}
		}
		page = nextPage
		maxAttempts--
	}
	if maxAttempts == 0 {
		return nil, fmt.Errorf("リポジトリの取得が上限に達しました")
	}
	return maps.Values(groups), nil
}

func FetchUserRepositories(ctx context.Context, client *Client) (*models.RepositoryGroup, error) {
	page := 1
	maxAttempts := 100
	allRepos := make([]models.Repository, 0, 10000)
	for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 {
		repos, nextPage, err := client.fetchUserRepositories(ctx, affiliationOwner, page)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			allRepos = append(allRepos, models.Repository{
				Name:    repo.Name,
				OrgName: repo.OrgName,
				Host:    repo.Host,
				HtmlUrl: repo.HtmlUrl,
			})
		}
		page = nextPage
		maxAttempts--
	}
	if maxAttempts == 0 {
		return nil, fmt.Errorf("リポジトリの取得が上限に達しました")
	}
	g, err := models.NewRepositoryGroup(allRepos...)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository group: %w", err)
	}
	return g, nil
}

type OwnerOrganization struct {
	name   string
	client *Client
}

func NewOrganizations(ctx context.Context, client *Client) ([]*OwnerOrganization, error) {
	gitOrgs, err := client.FetchOrganizations(ctx)
	if err != nil {
		return nil, err
	}
	organizations := make([]*OwnerOrganization, 0, len(gitOrgs))
	for _, gitOrg := range gitOrgs {
		org, err := NewOrganization(gitOrg.Login, client)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, org)
	}
	return organizations, nil
}

func NewOrganization(orgName string, client *Client) (*OwnerOrganization, error) {
	if orgName == "" {
		return nil, fmt.Errorf("organization name cannot be empty")
	}
	return &OwnerOrganization{
		name:   orgName,
		client: client,
	}, nil
}

func (o *OwnerOrganization) FetchRepositories(ctx context.Context) (*models.RepositoryGroup, error) {
	page := 1
	maxAttempts := 100 // 安全のための最大ページ数
	allRepos := make([]models.Repository, 0, 10000)

	for page > 0 && len(allRepos) < 10000 && maxAttempts > 0 { // 追加の安全対策
		repos, nextPage, err := o.client.fetchOrgRepositories(ctx, o.name, page)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			allRepos = append(allRepos, models.Repository{
				Name:    repo.Name,
				OrgName: repo.OrgName,
				Host:    repo.Host,
				HtmlUrl: repo.HtmlUrl,
			})
		}
		page = nextPage
		maxAttempts--
	}

	if maxAttempts == 0 {
		return nil, fmt.Errorf("リポジトリの取得が上限に達しました")
	}
	g, err := models.NewRepositoryGroup(allRepos...)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository group: %w", err)
	}
	return g, nil
}

func NewClient(opts api.ClientOptions) (*Client, error) {
	client, err := api.NewRESTClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client for %s: %w", opts.Host, err)
	}

	return &Client{client: client, host: opts.Host}, nil
}

func (c *Client) FetchOrganizations(ctx context.Context) ([]Organization, error) {
	var orgs []Organization
	if err := c.client.DoWithContext(ctx, "GET", "user/orgs", nil, &orgs); err != nil {
		return nil, fmt.Errorf("failed to fetch organizations from %s: %w", c.host, err)
	}
	return orgs, nil
}

type affiliation string

const (
	affiliationOwner        affiliation = "owner"
	affiliationCollaborator affiliation = "collaborator"
)

// fetch login user's repositories
func (c *Client) fetchUserRepositories(ctx context.Context, a affiliation, page int) (repos []Repository, nextPage int, err error) {
	resp, err := c.client.RequestWithContext(ctx, "GET", fmt.Sprintf("user/repos?per_page=100&page=%d&affiliation=%s", page, a), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch user repositories: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal user repositories: %w", err)
	}

	// Linkヘッダーの処理
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		links := strings.Split(linkHeader, ",")
		for _, link := range links {
			if strings.Contains(link, `rel="next"`) {
				parts := strings.Split(link, ";")
				urlPart := strings.Trim(parts[0], " <>")
				parsedURL, err := url.Parse(urlPart)
				if err != nil {
					continue
				}
				query := parsedURL.Query()
				if pageStr := query.Get("page"); pageStr != "" {
					nextPage, err = strconv.Atoi(pageStr)
					if err != nil {
						continue
					}
				}
			}
		}
	}

	// リポジトリ情報の補完
	for i := range repos {
		hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
		repos[i].Host = strings.Split(hostWithPath, "/")[0]
		// user/reposの場合、ownerがリポジトリのオーナー
		repos[i].OrgName = strings.Split(strings.TrimPrefix(repos[i].HtmlUrl, "https://"+repos[i].Host+"/"), "/")[0]
	}

	return repos, nextPage, nil
}
