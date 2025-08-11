package github

import (
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/SourcewareLab/Toney/internal/ui/widgets"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitHubModel struct {
	Width       int
	Height      int
	Issues      []GitHubIssue
	List        list.Model
	Loading     bool
	Error       string
	Focused     bool
	filterState string // "all" | "open" | "closed"
	sortBy      string // "title" | "updated"
	showHelp    bool
}

type GitHubIssue struct {
	Number     int     `json:"number"`
	IssueTitle string  `json:"title"`
	Body       string  `json:"body"`
	State      string  `json:"state"`
	Author     string  `json:"user"`
	Labels     []Label `json:"labels"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	HTMLURL    string  `json:"html_url"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// renderHelpView shows an overlay-style help screen with keybindings
func (m *GitHubModel) renderHelpView() string {
	pal := theme.DefaultDarkPalette()
	header := widgets.Header{Palette: pal, Width: m.Width, Title: "GitHub · Help", Right: "? to close"}.View()
	overlay := widgets.HelpOverlay{
		Palette: pal,
		Width:   m.Width,
		Height:  m.Height - lipgloss.Height(header) - 1,
		Title:   "Keybindings",
		Sections: []widgets.HelpSection{
			{Title: "Navigation", Rows: []string{"↑/↓: Move selection", "enter/c: Convert to note", "esc: Back"}},
			{Title: "Actions", Rows: []string{"r: Refresh", "f: Cycle filter (all/open/closed)", "s: Toggle sort (title/updated)", "?: Toggle help"}},
		},
	}.View()
	footer := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{{Key: "?", Desc: "close"}, {Key: "esc", Desc: "back"}}}.View()
	return lipgloss.JoinVertical(lipgloss.Left, header, overlay, footer)
}

// issueDelegate renders GitHub issues with status and label chips.
type issueDelegate struct{ pal theme.Palette }

func (d issueDelegate) Height() int                               { return 2 }
func (d issueDelegate) Spacing() int                              { return 0 }
func (d issueDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d issueDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	issue, ok := listItem.(GitHubIssue)
	if !ok {
		return
	}
	selected := index == m.Index()

	// Title line: left title, right updated time
	title := fmt.Sprintf("#%d %s", issue.Number, issue.IssueTitle)
	updated := timeAgo(issue.UpdatedAt)
	titleStyle := lipgloss.NewStyle().Foreground(d.pal.Fg)
	if selected {
		titleStyle = titleStyle.Foreground(d.pal.SelectedFg).Background(d.pal.SelectedBg).Bold(true)
	}
	rightStyle := lipgloss.NewStyle().Foreground(d.pal.Muted)
	if selected {
		rightStyle = rightStyle.Foreground(d.pal.SelectedFg)
	}
	// Space-fill to align right timestamp
	totalW := m.Width()
	left := titleStyle.Render(title)
	right := rightStyle.Render(updated)
	pad := totalW - lipgloss.Width(left) - lipgloss.Width(right)
	if pad < 1 {
		pad = 1
	}
	titleRow := lipgloss.JoinHorizontal(lipgloss.Left, left, strings.Repeat(" ", pad), right)

	// Status chip
	statusText := "OPEN"
	statusBg := d.pal.Success
	if strings.ToLower(issue.State) == "closed" {
		statusText = "CLOSED"
		statusBg = d.pal.Error
	}
	chip := lipgloss.NewStyle().Foreground(d.pal.SelectedFg).Background(statusBg).Padding(0, 1).Bold(true)
	statusChip := chip.Render(statusText)

	// Label chips (use GitHub label colors)
	labelChips := make([]string, 0, len(issue.Labels))
	for _, lb := range issue.Labels {
		clr := lipgloss.Color("#" + lb.Color)
		lc := lipgloss.NewStyle().Foreground(d.pal.Fg).Background(clr).Padding(0, 1)
		if selected {
			lc = lc.Foreground(d.pal.SelectedFg)
		}
		labelChips = append(labelChips, lc.Render(lb.Name))
	}
	chipsRow := lipgloss.JoinHorizontal(lipgloss.Left, append([]string{statusChip}, labelChips...)...)

	// Author and updated summary from Description()
	desc := issue.Description()
	descStyle := lipgloss.NewStyle().Foreground(d.pal.Muted)
	if selected {
		descStyle = descStyle.Foreground(d.pal.SelectedFg)
	}
	infoRow := lipgloss.JoinHorizontal(lipgloss.Left, chipsRow, "  ", descStyle.Render(desc))

	line := lipgloss.JoinVertical(lipgloss.Left, titleRow, infoRow)
	_, _ = io.WriteString(w, line)
}

// timeAgo returns a short relative time like "2h", "3d" given an ISO8601 timestamp.
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

// Implement list.Item interface for GitHubIssue
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
	return fmt.Sprintf("[%s] #%d %s", status, i.Number, i.IssueTitle)
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

	pal := theme.DefaultDarkPalette()
	delegate := issueDelegate{pal: pal}

	l := list.New(items, delegate, w/2, 2*h/3)
	l.Title = ""
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(true)

	l.Styles.Title = lipgloss.NewStyle().
		Foreground(pal.Fg).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pal.Border).
		Padding(0, 1).
		Margin(0, 0, 1, 0)
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(pal.Fg)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(pal.Accent)

	return &GitHubModel{
		Width:       w,
		Height:      h,
		List:        l,
		Loading:     false,
		Issues:      []GitHubIssue{},
		Focused:     true,
		filterState: "all",
		sortBy:      "title",
		showHelp:    false,
	}
}

func NewGitHubSetupOrIssuesModel(w int, h int) tea.Model {
	if !config.AppConfig.GitHub.Enabled ||
		config.AppConfig.GitHub.Token == "" ||
		config.AppConfig.GitHub.Owner == "" ||
		config.AppConfig.GitHub.Repo == "" {
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

	case OpenEditorMsg:
		return m, m.openFileInEditor(msg.FilePath, msg.IssueTitle)

	case tea.KeyMsg:
		if !m.Focused {
			return m, nil
		}

		switch msg.String() {
		case "?":
			m.showHelp = !m.showHelp
			return m, nil
		case "r", "ctrl+r":
			m.Loading = true
			m.Error = ""
			return m, m.SyncIssues()

		case "enter":
			if len(m.Issues) > 0 && m.List.Index() < len(m.Issues) {
				selectedIssue := m.Issues[m.List.Index()]
				return m, m.convertIssueToNote(selectedIssue)
			}

		case "c":
			if len(m.Issues) > 0 && m.List.Index() < len(m.Issues) {
				selectedIssue := m.Issues[m.List.Index()]
				return m, m.convertIssueToNote(selectedIssue)
			}
		case "f":
			switch m.filterState {
			case "all":
				m.filterState = "open"
			case "open":
				m.filterState = "closed"
			default:
				m.filterState = "all"
			}
			m.applyListItems()
			return m, nil
		case "s":
			if m.sortBy == "title" {
				m.sortBy = "updated"
			} else {
				m.sortBy = "title"
			}
			m.applyListItems()
			return m, nil
		case "esc":
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

// applyListItems filters, sorts, and loads issues into the list
func (m *GitHubModel) applyListItems() {
	filtered := make([]GitHubIssue, 0, len(m.Issues))
	for _, is := range m.Issues {
		switch m.filterState {
		case "open":
			if is.State == "open" {
				filtered = append(filtered, is)
			}
		case "closed":
			if is.State == "closed" {
				filtered = append(filtered, is)
			}
		default:
			filtered = append(filtered, is)
		}
	}

	// Sort
	switch m.sortBy {
	case "updated":
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].UpdatedAt > filtered[j].UpdatedAt
		})
	default:
		sort.Slice(filtered, func(i, j int) bool {
			return strings.ToLower(filtered[i].IssueTitle) < strings.ToLower(filtered[j].IssueTitle)
		})
	}

	// Load into list
	items := make([]list.Item, len(filtered))
	for idx, is := range filtered {
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

	if m.showHelp {
		return m.renderHelpView()
	}

	return m.renderIssuesView()
}

func (m *GitHubModel) renderDisabledView() string {
	pal := theme.DefaultDarkPalette()
	header := widgets.Header{Palette: pal, Width: m.Width, Title: "GitHub", Right: "disabled"}.View()
	box := lipgloss.NewStyle().
		Border(theme.Borders()).
		BorderForeground(pal.Border).
		Padding(1, 2).
		Width(m.Width)
	content := lipgloss.NewStyle().
		Foreground(pal.Fg).
		Render("GitHub integration is disabled.\n\nRun 'toney github setup' to enable it.")
	footer := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{{Key: "esc", Desc: "back"}}}.View()
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		box.Render(content),
		footer,
	)
}

func (m *GitHubModel) renderLoadingView() string {
	pal := theme.DefaultDarkPalette()
	repoInfo := fmt.Sprintf("%s/%s", config.AppConfig.GitHub.Owner, config.AppConfig.GitHub.Repo)

	// Helper function for responsive width
	maxWidth := 80
	if m.Width-4 < maxWidth {
		maxWidth = m.Width - 4
	}

	header := widgets.Header{Palette: pal, Width: m.Width, Title: "GitHub · Loading", Right: ""}.View()
	containerStyle := lipgloss.NewStyle().
		Width(maxWidth).
		Height(10).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(theme.Borders()).
		BorderForeground(pal.Border)

	titleStyle := lipgloss.NewStyle().
		Foreground(pal.Fg).
		Bold(true).
		Margin(0, 0, 1, 0)

	loadingStyle := lipgloss.NewStyle().
		Foreground(pal.Fg).
		Bold(true).
		Margin(0, 0, 1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(pal.Muted).
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
	navigation := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{{Key: "esc", Desc: "back"}}}.View()

	// Combine content with navigation at bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		centeredContent,
		navigation,
	)
}

func (m *GitHubModel) renderErrorView() string {
	pal := theme.DefaultDarkPalette()
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
		Border(theme.Borders()).
		BorderForeground(pal.Border)

	titleStyle := lipgloss.NewStyle().
		Foreground(pal.Fg).
		Bold(true).
		Margin(0, 0, 1, 0)

	errorStyle := lipgloss.NewStyle().
		Foreground(pal.Error).
		Bold(true).
		Margin(0, 0, 1, 0)

	descStyle := lipgloss.NewStyle().
		Foreground(pal.Muted).
		Margin(0, 0, 1, 0)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(" Error"),
		errorStyle.Render(m.Error),
		descStyle.Render("Press 'r' to retry or 'esc' to go back"),
	)

	contentContainer := lipgloss.NewStyle().
		Width(m.Width).
		Height(availableHeight).
		Align(lipgloss.Center, lipgloss.Center)

	centeredContent := contentContainer.Render(containerStyle.Render(content))

	navigation := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{{Key: "r", Desc: "retry"}, {Key: "esc", Desc: "back"}}}.View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		centeredContent,
		navigation,
	)
}

func (m *GitHubModel) renderIssuesView() string {
	pal := theme.DefaultDarkPalette()
	repoInfo := fmt.Sprintf("%s/%s", config.AppConfig.GitHub.Owner, config.AppConfig.GitHub.Repo)
	header := widgets.Header{Palette: pal, Width: m.Width, Title: "GitHub Issues", Right: repoInfo}.View()
	footer := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{
		{Key: "↑↓", Desc: "navigate"},
		{Key: "enter", Desc: "convert"},
		{Key: "r", Desc: "refresh"},
		{Key: "f", Desc: "filter"},
		{Key: "s", Desc: "sort"},
		{Key: "?", Desc: "help"},
		{Key: "esc", Desc: "back"},
	}}.View()

	bodyH := m.Height - lipgloss.Height(header) - lipgloss.Height(footer)
	if bodyH < 3 {
		bodyH = 3
	}

	// Compose top summary text and list inside a bordered box
	// repoInfo already defined above for header
	openCount, closedCount := 0, 0
	for _, issue := range m.Issues {
		if issue.State == "open" {
			openCount++
		} else {
			closedCount++
		}
	}
	summaryText := fmt.Sprintf("%s\n%d open • %d closed   |   filter: %s   |   sort: %s",
		repoInfo, openCount, closedCount, m.filterState, m.sortBy)
	summary := lipgloss.NewStyle().Foreground(pal.Fg).Width(m.Width).Height(bodyH/3).Align(lipgloss.Center, lipgloss.Center).Render(summaryText)
	listArea := lipgloss.Place(m.Width, 2*bodyH/3, lipgloss.Center, lipgloss.Center, m.List.View())

	if len(m.List.Items()) == 0 {
		listArea = lipgloss.Place(m.Width, 2*bodyH/3, lipgloss.Center, lipgloss.Top,
			lipgloss.NewStyle().Foreground(pal.Fg).Render("You have no Issues!"))
	}

	container := lipgloss.NewStyle().
		Border(theme.Borders()).
		BorderForeground(pal.Border).
		Width(m.Width).
		Height(bodyH).
		Padding(0, 1)

	body := container.Render(lipgloss.JoinVertical(lipgloss.Left, summary, listArea))

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *GitHubModel) SetFocus(focused bool) {
	m.Focused = focused
}

func (m *GitHubModel) SyncIssues() tea.Cmd {
	return func() tea.Msg {
		api := NewGitHubAPI()
		issues, err := api.FetchIssues()
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

func (m *GitHubModel) convertIssueToNote(issue GitHubIssue) tea.Cmd {
	return func() tea.Msg {
		converter := NewNoteConverter()
		filePath, err := converter.ConvertIssueToNote(issue)
		if err != nil {
			return GitHubSyncMsg{
				Issues: m.Issues,
				Error:  fmt.Errorf("failed to convert issue to note: %w", err),
			}
		}

		return OpenEditorMsg{
			FilePath:   filePath,
			IssueTitle: issue.IssueTitle,
		}
	}
}

type OpenEditorMsg struct {
	FilePath   string
	IssueTitle string
}

func (m *GitHubModel) openFileInEditor(filePath, issueTitle string) tea.Cmd {
	editorArgs := append(config.AppConfig.General.Editor, filePath)
	cmd := exec.Command(editorArgs[0], editorArgs[1:]...)
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return GitHubSyncMsg{
				Issues: m.Issues,
				Error:  fmt.Errorf("failed to open editor: %w", err),
			}
		}
		return GitHubSyncMsg{
			Issues: m.Issues,
			Error:  nil,
		}
	})
}

func (m *GitHubModel) GetCurrentPage() enums.Page {
	return enums.Page(10)
}
