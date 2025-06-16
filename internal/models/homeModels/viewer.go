package homemodels

import (
	"fmt"
	"os"
	"strings"

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
	Path      string
	Renderer  *glamour.TermRenderer
	isReading bool
	CurrRead  string
}

func NewViewer() *Viewer {
	vp := viewport.New(0, 0)
	vp.YOffset = 0
	vp.Style = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		MarginTop(1).
		Padding(1, 1).
		BorderForeground(lipgloss.Color("#45475a"))
	vp.SetContent("Select a file to view its contents")

	return &Viewer{
		Viewport:  vp,
		isReading: false,
	}
}

type ChangeFileMessage struct {
	Path string
}

type RefreshView struct {
	Content string
	Path    string
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ChangeFileMessage:
		m.Path = msg.Path
		content := m.ReadFile()
		m.Viewport.SetContent(content)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

		m.Renderer, _ = glamour.NewTermRenderer(
			glamour.WithWordWrap(m.Width-2),
			glamour.WithAutoStyle(),
		)

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
}

func (m *Viewer) Header() string {
	return ""
}

func (m *Viewer) ReadFile() string {
	path := strings.TrimSuffix(m.Path, "/")

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		content = ([]byte)(fmt.Sprintf("An error occured while reading the file\n%s", err.Error()))
	}

	rendered := m.RenderMarkdown(string(content), m.Width)

	return rendered
}

func (m *Viewer) RenderMarkdown(md string, width int) string {
	out, _ := m.Renderer.Render(md)

	return out
}
