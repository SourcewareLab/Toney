package homemodels

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Integrate Loader/Spinner
type Viewer struct {
	IsFocused bool
	Height    int
	Width     int
	Viewport  viewport.Model
	Ready     bool
	Content   string
	Path      string
}

func NewViewer() *Viewer {
	vp := viewport.New(0, 0)
	vp.YOffset = 0
	vp.Style = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		MarginTop(1).
		BorderForeground(lipgloss.Color("#45475a"))
	vp.SetContent("Select a file to view its contents")

	return &Viewer{
		Content:  "Select a file to view its contents",
		Viewport: vp,
	}
}

type ChangeFileMessage struct {
	Path string
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ChangeFileMessage:
		m.Path = msg.Path
		m.Viewport.SetContent(RenderMarkdown(ReadFile(msg.Path)))
		m.Viewport.YOffset = 0
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

		m.Viewport.Width = m.Width*3/4 - 1
		m.Viewport.Height = m.Height + 1
		m.Viewport.YOffset = 0
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Viewer) View() string {
	if m.IsFocused {
		m.Viewport.Style = m.Viewport.Style.BorderForeground(lipgloss.Color("#bb9af7"))
	} else {
		m.Viewport.Style = m.Viewport.Style.BorderForeground(lipgloss.Color("#45475a"))
	}

	return m.Viewport.View()
	// return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#45475a")).Render(m.Viewport.View())
	// return style.Width(m.Width*3/4 - 2).Height(m.Height - 2).MarginTop(1).Render(RenderMarkdown(m.Viewport.View()))
}

func (m Viewer) Header() string {
	return ""
}

func ReadFile(path string) string {
	data, err := os.ReadFile(path[0 : len(path)-1])
	if err != nil {
		return "An error occurred when reading file!"
	}

	return string(data)
}

func RenderMarkdown(md string) string {
	out, _ := glamour.Render(md, "dark")

	return out
}
