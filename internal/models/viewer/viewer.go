package viewer

import (
	"fmt"
	"os"
	"strings"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
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
	isEditing bool
	Renderer  *glamour.TermRenderer
	Keymap    keymap.ViewerKeyMap
}

func NewViewer(w int, h int) *Viewer {
	vp := viewport.New(w*3/4, h)
	vp.YOffset = 0
	pal := colors.ColorPalette()
	vp.Style = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		MarginTop(0).
		Padding(1, 1).
		BorderForeground(pal.Border).
		Foreground(pal.Text)
	vp.SetContent(
		lipgloss.Place(w*3/4, h-2, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(pal.Text).Render("Select a file to view its contents"),
		))

	r, _ := glamour.NewTermRenderer(glamour.WithStyles(config.ToGlamourStyle(config.AppConfig.Styles.Renderer)),
		glamour.WithWordWrap(w*3/4-2))

	return &Viewer{
		Viewport:  vp,
		Height:    h,
		Width:     w,
		isEditing: false,
		Keymap:    keymap.NewViewerKeyMap(),
		Renderer:  r,
	}
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m *Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.EditorClose:
		m.Viewport.SetContent(m.ReadFile(false))
		return m, nil
	case messages.ChangeFileMessage:
		m.Path = msg.Path
		content := m.ReadFile(false)
		m.Viewport.SetContent(content)
		m.Viewport.YOffset = 0
		return m, nil
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		m.Viewport.Height = msg.Height
		m.Viewport.Width = msg.Width * 3 / 4

		return m, nil
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
	pal := colors.ColorPalette()
	if m.IsFocused {
		m.Viewport.Style = m.Viewport.Style.BorderForeground(pal.FocusedBorder)
	} else {
		m.Viewport.Style = m.Viewport.Style.BorderForeground(pal.Border)
	}
	// Simple help line at the bottom, no header/footer widgets
	help := lipgloss.NewStyle().
		Foreground(pal.Text).
		PaddingLeft(2).
		Render("esc: back")
	// Ensure viewport height fits above help line
	bodyH := m.Height - lipgloss.Height(help)
	if bodyH < 3 {
		bodyH = 3
	}
	m.Viewport.Width = m.Width
	m.Viewport.Height = bodyH
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(pal.Border).
		Width(m.Width).
		Height(bodyH).
		Padding(0, 0)
	return lipgloss.JoinVertical(lipgloss.Left, box.Render(m.Viewport.View()), help)
}

func (m *Viewer) Header() string {
	return ""
}

func (m *Viewer) ReadFile(raw bool) string { // Change to editor type when config done
	path := strings.TrimSuffix(m.Path, "/")

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		content = ([]byte)(fmt.Sprintf("An error occured while reading the file:%s\n%s", m.Path, err.Error()))
	}

	if raw {
		return string(content)
	}

	rendered := m.RenderMarkdown(string(content), m.Width)

	return rendered
}

func (m *Viewer) RenderMarkdown(md string, width int) string {
	out, _ := m.Renderer.Render(md)

	return out
}
