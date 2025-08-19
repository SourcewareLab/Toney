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
	vp := viewport.New(w*3/4-8, h-6)
	vp.YOffset = 0
	vp.SetContent(
		lipgloss.Place(w*3/4-8, h-6, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(colors.ColorPalette().Text).Render("Select a file to view its contents"),
		))

	r, _ := glamour.NewTermRenderer(glamour.WithStyles(config.ToGlamourStyle(config.AppConfig.Styles.Renderer)),
		glamour.WithWordWrap(w*3/4-10))

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

		m.Viewport.Height = msg.Height - 6
		m.Viewport.Width = msg.Width*3/4 - 8

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
	borderColor := colors.ColorPalette().Border
	if m.IsFocused {
		borderColor = colors.ColorPalette().FocusedBorder
	}

	// Use the updated width calculation from homeModel
	viewerWidth := m.Width - 4
	m.Viewport.Width = viewerWidth - 6
	m.Viewport.Height = m.Height - 6

	// Explicit border definition to ensure right border shows
	viewerStyle := lipgloss.NewStyle().
		Width(viewerWidth).
		Height(m.Height-3).
		Padding(1, 2).
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}).
		BorderForeground(borderColor)

	return viewerStyle.Render(m.Viewport.View())
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
