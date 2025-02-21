package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	host   string
}

func NewClient(opts api.ClientOptions) (*Client, error) {
	client, err := api.NewRESTClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client for %s: %w", opts.Host, err)
	}

	return &Client{client: client, host: opts.Host}, nil
}

func (c *Client) FetchOrganizations() ([]Organization, error) {
	var orgs []Organization
	if err := c.client.Get("user/orgs", &orgs); err != nil {
		return nil, fmt.Errorf("failed to fetch organizations from %s: %w", c.host, err)
	}
	return orgs, nil
}

func (c *Client) FetchRepositories(orgs []Organization, page int) (repos []Repository, nextPage int, err error) {
	var allRepos []Repository

	for _, org := range orgs {
		resp, err := c.client.Request("GET", fmt.Sprintf("orgs/%s/repos?per_page=100&page=%d", org.Login, page), nil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch repositories for %s: %w", org.Login, err)
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal repositories for %s: %w", org.Login, err)
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
			repos[i].OrgName = org.Login
			hostWithPath := strings.TrimPrefix(repos[i].HtmlUrl, "https://")
			repos[i].Host = strings.Split(hostWithPath, "/")[0]
		}
		allRepos = append(allRepos, repos...)
	}
	return allRepos, nextPage, nil
}
