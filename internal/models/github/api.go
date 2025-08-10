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
	owner  string
	repo   string
}

func NewGitHubAPI() *GitHubAPI {
	return &GitHubAPI{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: config.AppConfig.GitHub.Token,
		owner: config.AppConfig.GitHub.Owner,
		repo:  config.AppConfig.GitHub.Repo,
	}
}

func (api *GitHubAPI) FetchIssues() ([]GitHubIssue, error) {
	if api.token == "" || api.owner == "" || api.repo == "" {
		return nil, fmt.Errorf("GitHub configuration is incomplete")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", api.owner, api.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+api.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Toney-GitHub-Integration")

	// Add query parameters to get both open and closed issues
	q := req.URL.Query()
	q.Add("state", "all")
	q.Add("per_page", "100")
	q.Add("sort", "updated")
	q.Add("direction", "desc")
	req.URL.RawQuery = q.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiIssues []GitHubAPIIssue
	if err := json.Unmarshal(body, &apiIssues); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert API response to our internal format
	issues := make([]GitHubIssue, len(apiIssues))
	for i, apiIssue := range apiIssues {
		issues[i] = GitHubIssue{
			Number:     apiIssue.Number,
			IssueTitle: apiIssue.Title,
			Body:       apiIssue.Body,
			State:      apiIssue.State,
			Author:     apiIssue.User.Login,
			Labels:     convertLabels(apiIssue.Labels),
			CreatedAt:  apiIssue.CreatedAt,
			UpdatedAt:  apiIssue.UpdatedAt,
			HTMLURL:    apiIssue.HTMLURL,
		}
	}

	return issues, nil
}

// GitHubAPIIssue represents the full GitHub API response structure
type GitHubAPIIssue struct {
	Number    int           `json:"number"`
	Title     string        `json:"title"`
	Body      string        `json:"body"`
	State     string        `json:"state"`
	User      GitHubUser    `json:"user"`
	Labels    []GitHubLabel `json:"labels"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
	HTMLURL   string        `json:"html_url"`
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
