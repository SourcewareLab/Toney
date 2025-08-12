package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
)

type GitHubAPI struct {
	client *http.Client
	token  string
}

func NewGitHubAPI() *GitHubAPI {
	return &GitHubAPI{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: config.AppConfig.GitHub.Token,
	}
}

type GitHubRepo struct {
	Name  string     `json:"name"`
	Owner GitHubUser `json:"owner"`
}

type GitHubAPIIssue struct {
	Number      int           `json:"number"`
	Title       string        `json:"title"`
	Body        string        `json:"body"`
	State       string        `json:"state"`
	User        GitHubUser    `json:"user"`
	Labels      []GitHubLabel `json:"labels"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	HTMLURL     string        `json:"html_url"`
	PullRequest *struct{}     `json:"pull_request,omitempty"`
}

type GitHubUser struct {
	Login string `json:"login"`
}

type GitHubLabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func convertLabels(apiLabels []GitHubLabel) []Label {
	labels := make([]Label, len(apiLabels))
	for i, apiLabel := range apiLabels {
		labels[i] = Label{
			Name:  apiLabel.Name,
			Color: apiLabel.Color,
		}
	}
	return labels
}

func (api *GitHubAPI) FetchIssues() ([]GitHubIssue, error) {
	if api.token == "" {
		return nil, fmt.Errorf("GitHub token is missing")
	}
	return api.FetchAllIssuesForUser()
}

func (api *GitHubAPI) FetchAllIssuesForUser() ([]GitHubIssue, error) {
	if api.token == "" {
		return nil, fmt.Errorf("GitHub token is missing")
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create repos request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Toney-GitHub-Integration")

	q := req.URL.Query()
	q.Add("per_page", "100")
	q.Add("sort", "updated")
	q.Add("affiliation", "owner,collaborator,organization_member")
	req.URL.RawQuery = q.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error listing repos (status %d): %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read repos response: %w", err)
	}
	var repos []GitHubRepo
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repos JSON: %w", err)
	}

	all := make([]GitHubIssue, 0)
	for _, r := range repos {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", r.Owner.Login, r.Name)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("Authorization", "Bearer "+api.token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "Toney-GitHub-Integration")

		q := req.URL.Query()
		q.Add("state", "open")
		q.Add("per_page", "100")
		q.Add("sort", "updated")
		q.Add("direction", "desc")
		req.URL.RawQuery = q.Encode()

		resp, err := api.client.Do(req)
		if err != nil {
			continue
		}
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				io.Copy(io.Discard, resp.Body)
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			var apiIssues []GitHubAPIIssue
			if err := json.Unmarshal(body, &apiIssues); err != nil {
				return
			}
			for _, apiIssue := range apiIssues {
				if apiIssue.PullRequest != nil {
					continue
				}
				all = append(all, GitHubIssue{
					Number:     apiIssue.Number,
					IssueTitle: apiIssue.Title,
					Body:       apiIssue.Body,
					State:      apiIssue.State,
					Author:     apiIssue.User.Login,
					Labels:     convertLabels(apiIssue.Labels),
					CreatedAt:  apiIssue.CreatedAt,
					UpdatedAt:  apiIssue.UpdatedAt,
					HTMLURL:    apiIssue.HTMLURL,
					Repo:       fmt.Sprintf("%s/%s", r.Owner.Login, r.Name),
				})
			}
		}()
	}
	return all, nil
}
