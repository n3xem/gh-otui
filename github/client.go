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
	"github.com/sourcegraph/conc/pool"
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

func (r Repository) ToDomain() models.Repository {
	return models.Repository{
		Name:    r.Name,
		OrgName: r.OrgName,
		Host:    r.Host,
		HtmlUrl: r.HtmlUrl,
	}
}

type Client struct {
	client *api.RESTClient
	host   string
}

func (c *Client) fetchOrgRepositories(ctx context.Context, org string, page int) (repos []Repository, nextPage int, lastPage int, err error) {
	resp, err := c.client.RequestWithContext(ctx, "GET", fmt.Sprintf("orgs/%s/repos?per_page=100&page=%d", org, page), nil)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch organization repositories for %s: %w", org, err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to unmarshal organization repositories for %s: %w", org, err)
	}

	// Linkヘッダーの処理
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		next, last, err := parseLinkHeader(linkHeader)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to parse link header: %w", err)
		}
		nextPage = next
		lastPage = last
	}

	for i := range repos {
		repos[i].OrgName = org
		hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
		repos[i].Host = strings.Split(hostWithPath, "/")[0]
	}
	return repos, nextPage, lastPage, nil
}

func parseLinkHeader(linkHeader string) (nextPage, lastPage int, err error) {
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		if strings.Contains(link, `rel="next"`) {
			parts := strings.Split(link, ";")
			urlPart := strings.Trim(parts[0], " <>")
			parsedURL, err := url.Parse(urlPart)
			if err != nil {
				return 0, 0, err
			}
			query := parsedURL.Query()
			if pageStr := query.Get("page"); pageStr != "" {
				nextPage, err = strconv.Atoi(pageStr)
				if err != nil {
					return 0, 0, err
				}
			}
		} else if strings.Contains(link, `rel="last"`) {
			parts := strings.Split(link, ";")
			urlPart := strings.Trim(parts[0], " <>")
			parsedURL, err := url.Parse(urlPart)
			if err != nil {
				return 0, 0, err
			}
			query := parsedURL.Query()
			if pageStr := query.Get("page"); pageStr != "" {
				lastPage, err = strconv.Atoi(pageStr)
				if err != nil {
					return 0, 0, err
				}
			}
		}
	}
	return nextPage, lastPage, nil
}

func FetchCollaboratingRepositories(ctx context.Context, client *Client) (iter.Seq[*models.RepositoryGroup], error) {
	type key struct {
		host string
		org  string
	}
	groups := make(map[key]*models.RepositoryGroup)
	ghRepos, err := fetchCollaboratingRepositories(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch collaborating repositories: %w", err)
	}
	repos := mapValues(ghRepos, Repository.ToDomain)
	for _, repo := range repos {
		key := key{
			host: repo.Host,
			org:  repo.OrgName,
		}
		if group, ok := groups[key]; ok {
			err := group.Add(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to add repository to group: %w", err)
			}
		} else {
			group, err := models.NewRepositoryGroup(repo)
			if err != nil {
				return nil, fmt.Errorf("failed to create repository group: %w", err)
			}
			groups[key] = group
		}
	}
	return maps.Values(groups), nil
}

func fetchCollaboratingRepositories(ctx context.Context, client *Client) ([]Repository, error) {
	maxAttempts := 100
	allRepos := make([]Repository, 0, 10000)
	repos, nextPage, lastPage, err := client.fetchUserRepositories(ctx, affiliationCollaborator, 1)
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, repos...)
	if lastPage == 0 || lastPage == 1 {
		return allRepos, nil
	}
	lastPage = min(lastPage, maxAttempts)
	p := pool.NewWithResults[[]Repository]().WithContext(ctx).WithMaxGoroutines(5)
	for page := nextPage; page <= lastPage; page++ {
		p.Go(func(ctx context.Context) ([]Repository, error) {
			repos, _, _, err := client.fetchUserRepositories(ctx, affiliationCollaborator, page)
			if err != nil {
				return nil, err
			}
			return repos, nil
		})
	}
	repoLists, err := p.Wait()
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, flatten(repoLists)...)
	return allRepos, nil
}

func FetchUserRepositories(ctx context.Context, client *Client) (*models.RepositoryGroup, error) {
	ghRepos, err := fetchUserRepositories(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories for user: %w", err)
	}

	repos := mapValues(ghRepos, Repository.ToDomain)

	g, err := models.NewRepositoryGroup(repos...)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository group: %w", err)
	}
	return g, nil
}

func fetchUserRepositories(ctx context.Context, client *Client) ([]Repository, error) {
	maxAttempts := 100
	allRepos := make([]Repository, 0, 10000)

	repos, nextPage, lastPage, err := client.fetchUserRepositories(ctx, affiliationOwner, 1)
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, repos...)
	if lastPage == 0 || lastPage == 1 {
		return repos, nil
	}
	lastPage = min(lastPage, maxAttempts)

	p := pool.NewWithResults[[]Repository]().WithContext(ctx).WithMaxGoroutines(5)
	for page := nextPage; page <= lastPage; page++ {
		p.Go(func(ctx context.Context) ([]Repository, error) {
			repos, _, _, err := client.fetchUserRepositories(ctx, affiliationOwner, page)
			if err != nil {
				return nil, err
			}
			return repos, nil
		})
	}

	repoLists, err := p.Wait()
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, flatten(repoLists)...)
	return allRepos, nil
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

func (o *OwnerOrganization) fetchRepositories(ctx context.Context) ([]Repository, error) {
	maxAttempts := 100 // 安全のための最大ページ数
	allRepos := make([]Repository, 0, 10000)

	repos, nextPage, lastPage, err := o.client.fetchOrgRepositories(ctx, o.name, 1)
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, repos...)
	if lastPage == 0 || lastPage == 1 {
		return allRepos, nil
	}
	lastPage = min(lastPage, maxAttempts)

	p := pool.NewWithResults[[]Repository]().WithContext(ctx).WithMaxGoroutines(5)
	for page := nextPage; page <= lastPage; page++ {
		p.Go(func(ctx context.Context) ([]Repository, error) {
			repos, _, _, err := o.client.fetchOrgRepositories(ctx, o.name, page)
			if err != nil {
				return nil, err
			}
			return repos, nil
		})
	}

	repoLists, err := p.Wait()
	if err != nil {
		return nil, err
	}
	allRepos = append(allRepos, flatten(repoLists)...)
	return allRepos, nil
}

func flatten[T any](slices [][]T) []T {
	length := 0
	for _, slice := range slices {
		length += len(slice)
	}
	results := make([]T, 0, length)
	for _, slice := range slices {
		results = append(results, slice...)
	}
	return results
}

func mapValues[U any, S ~[]U, V any](values S, f func(U) V) []V {
	results := make([]V, 0, len(values))
	for _, value := range values {
		results = append(results, f(value))
	}
	return results
}

func (o *OwnerOrganization) FetchRepositories(ctx context.Context) (*models.RepositoryGroup, error) {
	ghRepos, err := o.fetchRepositories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories for organization %s: %w", o.name, err)
	}

	repos := mapValues(ghRepos, Repository.ToDomain)

	g, err := models.NewRepositoryGroup(repos...)
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
func (c *Client) fetchUserRepositories(ctx context.Context, a affiliation, page int) (repos []Repository, nextPage int, lastPage int, err error) {
	resp, err := c.client.RequestWithContext(ctx, "GET", fmt.Sprintf("user/repos?per_page=100&page=%d&affiliation=%s", page, a), nil)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch user repositories: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to unmarshal user repositories: %w", err)
	}

	// Linkヘッダーの処理
	linkHeader := resp.Header.Get("Link")
	if linkHeader != "" {
		next, last, err := parseLinkHeader(linkHeader)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to parse link header: %w", err)
		}
		nextPage = next
		lastPage = last
	}

	// リポジトリ情報の補完
	for i := range repos {
		hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
		repos[i].Host = strings.Split(hostWithPath, "/")[0]
		// user/reposの場合、ownerがリポジトリのオーナー
		repos[i].OrgName = strings.Split(strings.TrimPrefix(repos[i].HtmlUrl, "https://"+repos[i].Host+"/"), "/")[0]
	}

	return repos, nextPage, lastPage, nil
}
