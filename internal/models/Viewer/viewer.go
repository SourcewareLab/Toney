package viewer

import (
	"fmt"
	"os"
	"strings"

	"toney/internal/messages"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Viewer struct {
	IsFocused bool
	Height    int
	Width     int
	Viewport  viewport.Model
	Ready     bool
	Path      string
	Renderer  *glamour.TermRenderer
}

func NewViewer(w int, h int) *Viewer {
	vp := viewport.New(w*3/4, h)
	vp.YOffset = 0
	vp.Style = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		MarginTop(1).
		Padding(1, 1).
		BorderForeground(lipgloss.Color("#45475a"))
	vp.SetContent("Select a file to view its contents")

	return &Viewer{
		Viewport: vp,
		Height:   h,
		Width:    w,
		Renderer: nil,
	}
}

func InitRenderer(w int) tea.Cmd {
	return func() tea.Msg {
		rnd, _ := glamour.NewTermRenderer(
			glamour.WithWordWrap(w*3/4-2),
			glamour.WithAutoStyle(),
		)

		return messages.RendererCreated{
			Renderer: rnd,
		}
	}
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.RendererCreated:
		m.Renderer = msg.Renderer
		return m, nil
	case messages.ChangeFileMessage:
		m.Path = msg.Path
		content := m.ReadFile()
		m.Viewport.SetContent(content)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
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
