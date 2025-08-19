package daily

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/models/github"
	"github.com/charmbracelet/bubbles/list"
)

type Tasks struct {
	All    []Task       `json:"all"`
	Github []GithubTask `json:"github"`
}

type GithubTask struct {
	TaskTitle string           `json:"title"`
	TaskDesc  string           `json:"desc"`
	Status    enums.TaskStatus `json:"status"`
	Ref       string           `json:"ref"`
	Repo      string           `json:"repo"`
	Owner     string           `json:"owner"`
	// We will not store these data in the file
	Link     string   `json:"-"`
	Labels   []string `json:"-"`
	Assignee []string `json:"-"`
}

func (m GithubTask) Title() string       { return m.TaskTitle }
func (m GithubTask) Description() string { return m.TaskDesc }
func (m GithubTask) FilterValue() string { return m.TaskTitle }

type Task struct {
	TaskTitle string           `json:"title"`
	TaskDesc  string           `json:"desc"`
	Status    enums.TaskStatus `json:"status"`
	// We will not store these data in the file
	Index int `json:"-"` // Point to index in the array
}

func (m Task) Title() string       { return m.TaskTitle }
func (m Task) Description() string { return m.TaskDesc }
func (m Task) FilterValue() string { return m.TaskTitle }

func (m Tasks) ItemsAsList() []list.Item {
	return TasksToItems(m.All)
}

func (m Tasks) ItemsAsListWithGithub() []list.Item {
	lst1 := TasksToItems(m.All)
	lst2 := GithubTaskToItems(m.Github)

	return append(lst1, lst2...)
}

func TasksToItems(tasks []Task) []list.Item {
	list := make([]list.Item, 0)
	for i, v := range tasks {
		v.Index = i
		list = append(list, v)
	}
	return list
}

func GithubTaskToItems(tasks []GithubTask) []list.Item {
	list := make([]list.Item, 0)
	for _, v := range tasks {
		list = append(list, v)
	}
	return list
}

func GetItems() Tasks {
	path := GetPath()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		tasks := Tasks{
			All:    make([]Task, 0),
			Github: make([]GithubTask, 0),
		}

		data, err2 := json.Marshal(tasks)
		if err2 != nil {
			fmt.Println(err2.Error())
			return Tasks{}
		}

		err2 = os.WriteFile(path, data, 0o644)
		if err2 != nil {
			fmt.Println(err2.Error())
		}
	} else if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	content, _ := os.ReadFile(path)

	tasks := Tasks{}

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		fmt.Println(err.Error())
	}

	return tasks
}

func WriteItems(tasks Tasks) {
	path := GetPath()

	data, _ := json.Marshal(tasks)

	os.WriteFile(path, data, 0o644)
}

func GetPath() string {
	home, _ := os.UserHomeDir()
	date := time.Now().Format("2006-01-02")

	return filepath.Join(home, config.AppConfig.General.NotesDir, ".daily", date)
}

func ConvertGitHubIssuestoGithubTasks(issues []github.GitHubIssue) []GithubTask {
	tasks := make([]GithubTask, len(issues))
	for i, issue := range issues {
		// Extract owner and repo from the repo field (format: "owner/repo")
		repoParts := strings.Split(issue.Repo, "/")
		owner := ""
		repo := ""
		if len(repoParts) == 2 {
			owner = repoParts[0]
			repo = repoParts[1]
		}

		tasks[i] = GithubTask{
			TaskTitle: issue.IssueTitle,
			TaskDesc:  issue.Body,
			Status:    enums.Pending, // GitHub issues are typically "pending" in task context
			Ref:       fmt.Sprintf("#%d", issue.Number),
			Repo:      repo,
			Owner:     owner,
			Link:      issue.HTMLURL,
			Labels:    convertLabelsToStrings(issue.Labels),
			Assignee:  issue.Assignees,
		}
	}
	return tasks
}

func convertLabelsToStrings(labels []github.Label) []string {
	result := make([]string, len(labels))
	for i, label := range labels {
		result[i] = label.Name
	}
	return result
}
