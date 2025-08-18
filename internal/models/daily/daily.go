package daily

import (
	"fmt"
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/models/github"
	taskpopup "github.com/SourcewareLab/Toney/internal/models/taskPopup"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitHubSyncMsg struct {
	Issues []github.GitHubIssue
	Error  error
}

type Daily struct {
	Width         int
	Height        int
	List          list.Model
	Tasks         Tasks
	GitHub        *github.GitHubModel
	Popup         *taskpopup.TaskPopup
	ShowPopup     bool
	Keymap        keymap.DailyTaskMap
	Help          help.Model
	CurrentTab    int
	Tabs          []enums.TaskTabs
	LoadingGithub bool
	GithubError   string
}

func NewDaily(w int, h int) *Daily {
	tasks := GetItems()
	githubModel := github.NewGitHubModel(w, h)

	lst := list.New(tasks.ItemsAsList(), TaskDelegate{}, w/2, 2*h/3)
	km := list.DefaultKeyMap()
	km.Quit.Unbind()
	lst.KeyMap = km
	lst.SetShowHelp(false)
	lst.SetShowTitle(false)

	return &Daily{
		Width:      w,
		Height:     h,
		List:       lst,
		Tasks:      tasks,
		GitHub:     githubModel,
		Keymap:     keymap.NewDailyTaskMap(),
		Help:       help.New(),
		CurrentTab: 0,
		Tabs:       []enums.TaskTabs{enums.Tasks, enums.Github},
	}
}

func (m *Daily) Init() tea.Cmd {
	return m.SyncGitHub()
}

func (m *Daily) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GitHubSyncMsg:
		m.LoadingGithub = false
		if msg.Error != nil {
			m.GithubError = fmt.Sprintf("GitHub Error: %s", msg.Error.Error())
			m.Tasks.Github = []GithubTask{}
		} else {
			m.GithubError = ""
			m.Tasks.Github = convertGitHubIssuestoGithubTasks(msg.Issues)
		}
		m.refreshList()
		return m, nil
	case messages.TaskPopupMessage:
		switch msg.Type {
		case enums.CreateRecurring:
			fallthrough
		case enums.CreateUnique:
			m.CreateTask(msg, msg.Type == enums.CreateUnique)
		case enums.Delete:
			m.DeleteTask(msg)
		case enums.ChangeStatus:
			m.StatusChangeTask(msg)
		case enums.Edit:
			m.EditTask(msg)
		}

		m.Refresh()
		m.ShowPopup = false
		return m, nil
	case tea.KeyMsg:
		if m.ShowPopup {
			updated, cmd := m.Popup.Update(msg)
			if popup, ok := updated.(*taskpopup.TaskPopup); ok { // Type matching, cause I cant assign it straightaway
				m.Popup = popup
				return m, cmd
			}
		}
		switch {
		case key.Matches(msg, m.Keymap.CreateUnique):
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.CreateUnique)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.CreateRecurring):
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.CreateRecurring)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.ChangeStatus):
			// Check if selected item is a GitHub task
			item := m.List.SelectedItem()
			if _, isGithubTask := item.(GithubTask); isGithubTask {
				return m, nil // Don't allow status changes for GitHub tasks
			}
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.ChangeStatus)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.DeleteTask):
			// Check if selected item is a GitHub task
			item := m.List.SelectedItem()
			if _, isGithubTask := item.(GithubTask); isGithubTask {
				return m, nil // Don't allow deletion of GitHub tasks
			}
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.Delete)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.EditTask):
			item := m.List.SelectedItem()

			// Check if it's a GitHub task first
			if _, isGithubTask := item.(GithubTask); isGithubTask {
				return m, nil // Don't allow editing of GitHub tasks
			}

			task, ok := item.(Task)
			if !ok { // Making sure that item is of type Task
				return m, nil
			}

			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.Edit)
			m.Popup.Form.TitleInput.SetValue(task.TaskTitle)
			m.Popup.Form.DescInput.SetValue(task.TaskDesc)

			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.BackToMenu):
			return m, func() tea.Msg {
				return messages.ChangePage{
					Page: enums.MenuPage,
				}
			}
		case key.Matches(msg, m.Keymap.TabRight):
			m.CurrentTab++
			if m.CurrentTab >= len(m.Tabs) {
				m.CurrentTab = m.CurrentTab % len(m.Tabs)
			}

			m.Refresh()
			// Auto-refresh GitHub data when switching to GitHub tab
			if m.Tabs[m.CurrentTab] == enums.Github {
				m.LoadingGithub = true
				m.GithubError = ""
				return m, m.SyncGitHub()
			}
			return m, nil
		case key.Matches(msg, m.Keymap.TabLeft):
			m.CurrentTab--
			if m.CurrentTab < 0 {
				m.CurrentTab = len(m.Tabs) + m.CurrentTab
			}

			m.Refresh()
			// Auto-refresh GitHub data when switching to GitHub tab
			if m.Tabs[m.CurrentTab] == enums.Github {
				m.LoadingGithub = true
				m.GithubError = ""
				return m, m.SyncGitHub()
			}
			return m, nil
		case key.Matches(msg, m.Keymap.RefreshGithub):
			// Force refresh GitHub data asynchronously
			m.LoadingGithub = true
			m.GithubError = ""
			return m, m.SyncGitHub()
		}
	}

	var cmd tea.Cmd

	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

func (m *Daily) View() string {
	if m.ShowPopup {
		return m.Popup.View()
	}

	statusLine := ""
	if m.LoadingGithub {
		statusLine = lipgloss.NewStyle().
			Foreground(colors.ColorPalette().Text).
			Render("🔄 Loading GitHub issues...")
	} else if m.GithubError != "" {
		statusLine = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Render(m.GithubError)
	}

	main := lipgloss.JoinVertical(lipgloss.Left,
		styles.GetDailyText(m.Width, m.Height/3),
		m.GetTabs(),
		statusLine,
		lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Left, m.List.View()))

	if len(m.List.Items()) == 0 {
		emptyMessage := "You have no Tasks!"
		if m.LoadingGithub && m.Tabs[m.CurrentTab] == enums.Github {
			emptyMessage = "Loading GitHub issues..."
		} else {
			switch m.Tabs[m.CurrentTab] {
			case enums.Tasks:
				if m.LoadingGithub {
					emptyMessage = "No local tasks. Loading GitHub issues..."
				} else {
					emptyMessage = "No tasks found. Create unique with 'c' or recurring with 'r', sync GitHub with 'ctrl+r'!"
				}
			case enums.Github:
				if config.AppConfig.GitHub.Enabled && config.AppConfig.GitHub.Token != "" {
					if m.GithubError != "" {
						emptyMessage = "No GitHub issues (check error above)"
					} else {
						emptyMessage = "No GitHub issues found. Press 'ctrl+r' to refresh!"
					}
				} else {
					emptyMessage = "GitHub not configured. Run 'toney github setup'."
				}
			}
		}
		main = lipgloss.JoinVertical(lipgloss.Left,
			styles.GetDailyText(m.Width, m.Height/3),
			m.GetTabs(),
			statusLine,
			lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Top,
				lipgloss.NewStyle().Foreground(colors.ColorPalette().Text).Render(emptyMessage)))
	}

	help := lipgloss.NewStyle().PaddingLeft(2).Render(m.Help.View(keymap.NewDynamic(m.Keymap.Bindings())))

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}

func (m *Daily) GetTabs() string {
	cfg := config.AppConfig.Styles.Renderer.Heading.Levels[0] // Using the same Heading choice as renderer
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Color)).Background(lipgloss.Color(cfg.Background)).Padding(0, 1)

	text := ""
	for _, v := range m.Tabs {
		style := style
		if v == m.Tabs[m.CurrentTab] {
			style = style.Background(colors.ColorPalette().MenuSelectedBg).Foreground(colors.ColorPalette().MenuSelectedText)
		}
		text += fmt.Sprintf(" %s ", style.Render(string(v)))
	}
	return lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, text)
}

func (m *Daily) Refresh() {
	m.Tasks = GetItems()
	m.refreshList()
}

func (m *Daily) refreshList() {
	curr := m.Tasks.ItemsAsListWithGithub()
	switch m.Tabs[m.CurrentTab] {
	case enums.Github:
		curr = GithubTaskToItems(m.Tasks.Github)
	}

	lst := list.New(curr, TaskDelegate{}, m.Width/2, 2*m.Height/3)
	km := list.DefaultKeyMap()
	km.Quit.Unbind()
	lst.KeyMap = km
	lst.SetShowHelp(false)
	lst.SetShowTitle(false)

	m.List = lst
}

func (m *Daily) SyncGitHub() tea.Cmd {
	if !config.AppConfig.GitHub.Enabled || config.AppConfig.GitHub.Token == "" {
		return nil
	}

	return func() tea.Msg {
		api := github.NewGitHubAPI()
		issues, err := api.FetchAllIssuesForUser()
		return GitHubSyncMsg{
			Issues: issues,
			Error:  err,
		}
	}
}

func convertGitHubIssuestoGithubTasks(issues []github.GitHubIssue) []GithubTask {
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
