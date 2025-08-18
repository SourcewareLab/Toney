package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/spf13/viper"
)

type GitHubAPI struct {
	client *http.Client
	token  string
	user   string // authenticated user login
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
	Assignees   []GitHubUser  `json:"assignees"`
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
		labels[i] = Label(apiLabel)
	}
	return labels
}

func (api *GitHubAPI) FetchAllIssuesForUser() ([]GitHubIssue, error) {
	if api.token == "" {
		return nil, fmt.Errorf("GitHub token is missing")
	}

	// Determine authenticated user login for assignment checks
	if api.user == "" {
		ureq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err == nil {
			ureq.Header.Set("Authorization", "Bearer "+api.token)
			ureq.Header.Set("Accept", "application/vnd.github.v3+json")
			ureq.Header.Set("User-Agent", "Toney-GitHub-Integration")
			if uresp, err := api.client.Do(ureq); err == nil {
				func() {
					defer uresp.Body.Close()
					if uresp.StatusCode == http.StatusOK {
						var me GitHubUser
						if json.NewDecoder(uresp.Body).Decode(&me) == nil {
							api.user = me.Login
						} else {
							io.Copy(io.Discard, uresp.Body)
						}
					} else {
						io.Copy(io.Discard, uresp.Body)
					}
				}()
			}
		}
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
	var repos []GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to parse repos JSON: %w", err)
	}

	// Bounded parallel fetching of repo issues
	workerCount := 8
	if len(repos) < workerCount {
		workerCount = len(repos)
	}
	issuesCh := make(chan []GitHubIssue, len(repos))
	repoCh := make(chan GitHubRepo)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for r := range repoCh {
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
			// No "since" filter to ensure we see all open issues
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
				var apiIssues []GitHubAPIIssue
				if err := json.NewDecoder(resp.Body).Decode(&apiIssues); err != nil {
					return
				}
				batch := make([]GitHubIssue, 0, len(apiIssues))
				for _, apiIssue := range apiIssues {
					if apiIssue.PullRequest != nil {
						continue
					}
					// collect assignees and check if assigned to me
					assignees := make([]string, 0, len(apiIssue.Assignees))
					assignedToMe := false
					for _, a := range apiIssue.Assignees {
						assignees = append(assignees, a.Login)
						if api.user != "" && a.Login == api.user {
							assignedToMe = true
						}
					}
					batch = append(batch, GitHubIssue{
						Number:       apiIssue.Number,
						IssueTitle:   apiIssue.Title,
						Body:         apiIssue.Body,
						State:        apiIssue.State,
						Author:       apiIssue.User.Login,
						Labels:       convertLabels(apiIssue.Labels),
						CreatedAt:    apiIssue.CreatedAt,
						UpdatedAt:    apiIssue.UpdatedAt,
						HTMLURL:      apiIssue.HTMLURL,
						Repo:         fmt.Sprintf("%s/%s", r.Owner.Login, r.Name),
						Assignees:    assignees,
						AssignedToMe: assignedToMe,
					})
				}
				if len(batch) > 0 {
					issuesCh <- batch
				}
			}()
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}
	for _, r := range repos {
		repoCh <- r
	}
	close(repoCh)
	go func() {
		wg.Wait()
		close(issuesCh)
	}()

	all := make([]GitHubIssue, 0)
	for batch := range issuesCh {
		all = append(all, batch...)
	}

	// Update last_sync to now (RFC3339) and persist
	viper.Set("github.last_sync", time.Now().UTC().Format(time.RFC3339))
	_ = viper.WriteConfig() // best-effort; non-fatal on error

	return all, nil
}
