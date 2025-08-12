package github

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitHubModel struct {
	Width         int
	Height        int
	Issues        []GitHubIssue
	List          list.Model
	Loading       bool
	Error         string
	Focused       bool
	sortBy        string       // "title" | "updated"
	showingIssue  bool         // true when showing issue overlay
	selectedIssue *GitHubIssue // issue being viewed in overlay
	showAssigned  bool         // filter: show only issues assigned to me
}

type GitHubIssue struct {
	Number       int      `json:"number"`
	IssueTitle   string   `json:"title"`
	Body         string   `json:"body"`
	State        string   `json:"state"`
	Author       string   `json:"user"`
	Labels       []Label  `json:"labels"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	HTMLURL      string   `json:"html_url"`
	Repo         string   `json:"-"`
	Assignees    []string `json:"-"`
	AssignedToMe bool     `json:"-"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type keyMap struct {
	Retry          key.Binding
	Back           key.Binding
	ToggleAssigned key.Binding
}

var keys = keyMap{
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	ToggleAssigned: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "assigned only"),
	),
}

type issueDelegate struct{}

func (d issueDelegate) Height() int                               { return 2 }
func (d issueDelegate) Spacing() int                              { return 0 }
func (d issueDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d issueDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	issue, ok := listItem.(GitHubIssue)
	if !ok {
		return
	}
	selected := index == m.Index()
	colors := colors.ColorPalette()

	// Truncate title if too long
	title := fmt.Sprintf("#%d %s", issue.Number, issue.IssueTitle)
	maxTitleWidth := m.Width() - 20 // Leave space for timestamp
	if len(title) > maxTitleWidth {
		title = title[:maxTitleWidth-3] + "..."
	}

	updated := timeAgo(issue.UpdatedAt)
	titleStyle := lipgloss.NewStyle().Foreground(colors.Text)
	if selected {
		titleStyle = titleStyle.Foreground(colors.MenuSelectedText).Background(colors.MenuSelectedBg).Bold(true)
	}
	rightStyle := lipgloss.NewStyle().Foreground(colors.Text)
	if selected {
		rightStyle = rightStyle.Foreground(colors.MenuSelectedText)
	}

	totalW := m.Width()
	left := titleStyle.Render(title)
	right := rightStyle.Render(updated)
	pad := totalW - lipgloss.Width(left) - lipgloss.Width(right)
	if pad < 1 {
		pad = 1
	}
	titleRow := lipgloss.JoinHorizontal(lipgloss.Left, left, strings.Repeat(" ", pad), right)

	desc := issue.Description()
	maxDescWidth := m.Width() - 10
	if len(desc) > maxDescWidth {
		desc = desc[:maxDescWidth-3] + "..."
	}

	descStyle := lipgloss.NewStyle().Foreground(colors.Text)
	if selected {
		descStyle = descStyle.Foreground(colors.MenuSelectedText)
	}

	infoRow := descStyle.Render(desc)
	line := lipgloss.JoinVertical(lipgloss.Left, titleRow, infoRow)
	_, _ = io.WriteString(w, line)
}

func timeAgo(ts string) string {
	if ts == "" {
		return ""
	}
	// Try a few common GitHub formats
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05Z07:00"}
	var t time.Time
	var err error
	for _, l := range layouts {
		t, err = time.Parse(l, ts)
		if err == nil {
			break
		}
	}
	if err != nil {
		return ""
	}
	d := time.Since(t)
	if d < time.Minute {
		return "now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	if d < 30*24*time.Hour {
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
	if d < 365*24*time.Hour {
		return fmt.Sprintf("%dmo", int(d.Hours()/(24*30)))
	}
	return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
}

func (i GitHubIssue) FilterValue() string {
	return i.IssueTitle
}

func (i GitHubIssue) Title() string {
	if i.IssueTitle == "" {
		return fmt.Sprintf("Issue #%d", i.Number)
	}
	status := "OPEN"
	if i.State == "closed" {
		status = "CLOSED"
	}
	repoPrefix := ""
	if i.Repo != "" {
		repoPrefix = fmt.Sprintf("%s ", i.Repo)
	}
	return fmt.Sprintf("%s[%s] #%d %s", repoPrefix, status, i.Number, i.IssueTitle)
}

func (i GitHubIssue) Description() string {
	desc := strings.ReplaceAll(i.Body, "\n", " ")
	desc = strings.TrimSpace(desc)
	if len(desc) > 80 {
		desc = desc[:77] + "..."
	}
	labels := ""
	if len(i.Labels) > 0 {
		labelNames := make([]string, len(i.Labels))
		for j, label := range i.Labels {
			labelNames[j] = label.Name
		}
		labels = fmt.Sprintf(" [%s]", strings.Join(labelNames, ", "))
	}
	author := fmt.Sprintf("@%s", i.Author)
	return fmt.Sprintf("%s%s - %s", author, labels, desc)
}

type GitHubSyncMsg struct {
	Issues []GitHubIssue
	Error  error
}

func NewGitHubModel(w int, h int) *GitHubModel {
	items := []list.Item{}

	delegate := issueDelegate{}

	l := list.New(items, delegate, w/2, 2*h/3)
	l.Title = ""
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(true)

	colors := colors.ColorPalette()
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(colors.Text).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Border).
		Padding(0, 1).
		Margin(0, 0, 1, 0)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(colors.Text)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(colors.Text)

	return &GitHubModel{
		Width:         w,
		Height:        h,
		List:          l,
		Loading:       false,
		Issues:        []GitHubIssue{},
		Focused:       true,
		sortBy:        "title",
		showingIssue:  false,
		selectedIssue: nil,
		showAssigned:  false,
	}
}

func NewGitHubSetupOrIssuesModel(w int, h int) tea.Model {
	if !config.AppConfig.GitHub.Enabled ||
		config.AppConfig.GitHub.Token == "" {
		return NewGitHubSetupModel()
	}

	model := NewGitHubModel(w, h)
	model.Loading = true
	return model
}

func (m *GitHubModel) Init() tea.Cmd {
	return nil
}

func (m *GitHubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case GitHubSyncMsg:
		m.Loading = false
		if msg.Error != nil {
			m.Error = msg.Error.Error()
			return m, nil
		}

		m.Issues = msg.Issues
		m.Error = ""
		m.applyListItems()

	case tea.KeyMsg:
		if !m.Focused {
			return m, nil
		}

		switch msg.String() {
		case "r", "ctrl+r":
			m.Loading = true
			m.Error = ""
			return m, m.SyncIssues()

		case "enter":
			if len(m.List.Items()) > 0 && m.List.Index() < len(m.List.Items()) {
				if item, ok := m.List.Items()[m.List.Index()].(GitHubIssue); ok {
					m.selectedIssue = &item
					m.showingIssue = true
					return m, nil
				}
			}

		case "s":
			if m.sortBy == "title" {
				m.sortBy = "updated"
			} else {
				m.sortBy = "title"
			}
			m.applyListItems()
			return m, nil
		case "a":
			m.showAssigned = !m.showAssigned
			m.applyListItems()
			return m, nil
		case "esc":
			if m.showingIssue {
				m.showingIssue = false
				m.selectedIssue = nil
				return m, nil
			}
			return m, func() tea.Msg {
				return messages.ChangePage{
					Page: enums.MenuPage,
				}
			}
		}
	}

	if m.Focused {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}

func (m *GitHubModel) applyListItems() {
	base := make([]GitHubIssue, len(m.Issues))
	copy(base, m.Issues)

	issues := base
	if m.showAssigned {
		filtered := make([]GitHubIssue, 0, len(base))
		for _, is := range base {
			if is.AssignedToMe {
				filtered = append(filtered, is)
			}
		}
		issues = filtered
	}

	// Sort
	switch m.sortBy {
	case "updated":
		sort.Slice(issues, func(i, j int) bool {
			return issues[i].UpdatedAt > issues[j].UpdatedAt
		})
	default:
		sort.Slice(issues, func(i, j int) bool {
			return strings.ToLower(issues[i].IssueTitle) < strings.ToLower(issues[j].IssueTitle)
		})
	}

	items := make([]list.Item, len(issues))
	for idx, is := range issues {
		items[idx] = is
	}
	m.List.SetItems(items)
}

func (m *GitHubModel) View() string {
	if !config.AppConfig.GitHub.Enabled {
		return m.renderDisabledView()
	}

	if m.Loading {
		return m.renderLoadingView()
	}

	if m.Error != "" {
		return m.renderErrorView()
	}

	if m.showingIssue && m.selectedIssue != nil {
		return m.renderIssueOverlay()
	}
	return m.renderIssuesView()
}

func (m *GitHubModel) renderDisabledView() string {
	h := help.New()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.ColorPalette().Border).
		Padding(1, 2).
		Width(m.Width)
	content := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Render("GitHub integration is disabled.\n\nRun 'toney github setup' to enable it.")
	help := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		PaddingLeft(2).
		Render(h.ShortHelpView([]key.Binding{keys.Back}))
	return lipgloss.JoinVertical(lipgloss.Left,
		box.Render(content),
		help,
	)
}

func (m *GitHubModel) renderLoadingView() string {
	repoInfo := "all repositories"

	// Helper function for responsive width
	maxWidth := 80
	if m.Width-4 < maxWidth {
		maxWidth = m.Width - 4
	}

	containerStyle := lipgloss.NewStyle().
		Width(maxWidth).
		Height(10).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.ColorPalette().Border)

	titleStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Bold(true).
		Margin(0, 0, 1, 0)

	loadingStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Bold(true).
		Margin(0, 0, 1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Margin(0, 0, 1, 0)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(" GitHub Issues"),
		loadingStyle.Render(" Loading issues..."),
		descStyle.Render(fmt.Sprintf("Fetching from %s", repoInfo)),
	)

	// Center the main content in available space
	contentContainer := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center)

	centeredContent := contentContainer.Render(containerStyle.Render(content))

	// Footer
	h := help.New()
	navigation := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		PaddingLeft(2).
		Render(h.ShortHelpView([]key.Binding{keys.Back}))

	// Combine content with navigation at bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredContent,
		navigation,
	)
}

func (m *GitHubModel) renderErrorView() string {
	// Calculate responsive dimensions
	availableHeight := m.Height - 3

	maxWidth := 80
	if m.Width-4 < maxWidth {
		maxWidth = m.Width - 4
	}

	containerStyle := lipgloss.NewStyle().
		Width(maxWidth).
		Height(14).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.ColorPalette().Border)

	titleStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Bold(true).
		Margin(0, 0, 1, 0)

	errorStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Bold(true).
		Margin(0, 0, 1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Margin(0, 0, 1, 0)

	h := help.New()
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(" Error"),
		errorStyle.Render(m.Error),
		descStyle.Render(h.ShortHelpView([]key.Binding{keys.Retry, keys.Back})),
	)

	contentContainer := lipgloss.NewStyle().
		Width(m.Width).
		Height(availableHeight).
		Align(lipgloss.Center, lipgloss.Center)

	centeredContent := contentContainer.Render(containerStyle.Render(content))

	navigation := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		PaddingLeft(2).
		Render(h.ShortHelpView([]key.Binding{keys.Retry, keys.Back}))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredContent,
		navigation,
	)
}

func (m *GitHubModel) renderIssuesView() string {
	repoInfo := "All repositories"

	// Since we only fetch open issues from API, all issues are open
	shownCount := len(m.List.Items())
	total := len(m.Issues)
	mode := "all"
	if m.showAssigned {
		mode = "assigned"
	}
	summaryText := fmt.Sprintf("GitHub Issues - %s (%s)\n%d shown / %d total", repoInfo, mode, shownCount, total)
	summary := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		Width(m.Width).
		Height(m.Height/3).
		Align(lipgloss.Center, lipgloss.Center).
		Render(summaryText)

	listArea := lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Center, m.List.View())

	if len(m.List.Items()) == 0 {
		listArea = lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Top,
			lipgloss.NewStyle().Foreground(colors.ColorPalette().Text).Render("No open issues found!"))
	}

	h := help.New()
	main := lipgloss.JoinVertical(lipgloss.Left, summary, listArea)
	help := lipgloss.NewStyle().
		Foreground(colors.ColorPalette().Text).
		PaddingLeft(2).
		Render(h.ShortHelpView([]key.Binding{keys.Retry, keys.ToggleAssigned, keys.Back}))

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}

func (m *GitHubModel) renderIssueOverlay() string {
	colors := colors.ColorPalette()
	issue := m.selectedIssue

	// Create overlay content
	title := fmt.Sprintf("#%d %s", issue.Number, issue.IssueTitle)

	// Wrap long titles
	titleStyle := lipgloss.NewStyle().
		Foreground(colors.Text).
		Bold(true).
		Width(m.Width - 8).
		Align(lipgloss.Left)

	// Author and state info
	infoText := fmt.Sprintf("Author: %s | State: %s | Updated: %s",
		issue.Author, strings.ToUpper(issue.State), timeAgo(issue.UpdatedAt))
	infoStyle := lipgloss.NewStyle().
		Foreground(colors.Text).
		Width(m.Width - 8).
		Align(lipgloss.Left)

	// Labels if any
	labelsText := ""
	if len(issue.Labels) > 0 {
		labelNames := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labelNames[i] = label.Name
		}
		labelsText = "Labels: " + strings.Join(labelNames, ", ")
	}

	// Compose content starting with title and info
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		"",
		infoStyle.Render(infoText),
	)

	if labelsText != "" {
		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			infoStyle.Render(labelsText),
		)
	}

	// Add description section with prominent styling if body exists
	if issue.Body != "" {
		// Description header
		descHeaderStyle := lipgloss.NewStyle().
			Foreground(colors.Text).
			Bold(true).
			Width(m.Width - 8).
			Align(lipgloss.Left)

		// Description body with better visibility
		bodyStyle := lipgloss.NewStyle().
			Foreground(colors.Text).
			Width(m.Width-8).
			Align(lipgloss.Left).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(colors.Border).
			PaddingLeft(2).
			MarginTop(1)

		content = lipgloss.JoinVertical(lipgloss.Left,
			content,
			"",
			descHeaderStyle.Render("Description:"),
			bodyStyle.Render(issue.Body),
		)
	}

	// Create bordered overlay
	overlay := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Border).
		Padding(1, 2).
		Width(m.Width - 4).
		Height(m.Height - 6).
		Render(content)

	// Help text
	helpText := "Press ESC to close"
	help := lipgloss.NewStyle().
		Foreground(colors.Text).
		Width(m.Width).
		Align(lipgloss.Center).
		PaddingLeft(2).
		Render(helpText)

	// Center the overlay on screen
	centeredOverlay := lipgloss.Place(m.Width, m.Height-2, lipgloss.Center, lipgloss.Center, overlay)

	return lipgloss.JoinVertical(lipgloss.Left, centeredOverlay, help)
}

func (m *GitHubModel) SetFocus(focused bool) {
	m.Focused = focused
}

func (m *GitHubModel) SyncIssues() tea.Cmd {
	return func() tea.Msg {
		api := NewGitHubAPI()
		issues, err := api.FetchAllIssuesForUser()
		// Sort by title as soon as we fetch
		if err == nil {
			sort.Slice(issues, func(i, j int) bool {
				return strings.ToLower(issues[i].IssueTitle) < strings.ToLower(issues[j].IssueTitle)
			})
		}
		return GitHubSyncMsg{
			Issues: issues,
			Error:  err,
		}
	}
}

func (m *GitHubModel) GetCurrentPage() enums.Page {
	return enums.Page(10)
}
