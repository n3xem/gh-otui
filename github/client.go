package github

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
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
	clients map[string]*api.RESTClient
}

func NewClient() (*Client, error) {
	hosts := auth.KnownHosts()
	clients := make(map[string]*api.RESTClient)
	for _, host := range hosts {
		client, err := api.NewRESTClient(api.ClientOptions{
			Host: host,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize client for %s: %w", host, err)
		}
		clients[host] = client
	}

	return &Client{clients: clients}, nil
}

func (c *Client) FetchOrganizations() ([]Organization, error) {
	var allOrgs []Organization
	for host, client := range c.clients {
		var orgs []Organization
		if err := client.Get("user/orgs", &orgs); err != nil {
			return nil, fmt.Errorf("failed to fetch organizations from %s: %w", host, err)
		}
		allOrgs = append(allOrgs, orgs...)
	}
	return allOrgs, nil
}

func (c *Client) FetchRepositories(orgs []Organization, page int) (repos []Repository, nextPage int, err error) {
	var allRepos []Repository

	for _, client := range c.clients {
		for _, org := range orgs {
			resp, err := client.Request("GET", fmt.Sprintf("orgs/%s/repos?per_page=100&page=%d", org.Login, page), nil)
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
	}
	return allRepos, nextPage, nil
}
