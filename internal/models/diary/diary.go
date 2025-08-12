package diary

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/models/fzf"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Diary struct {
	Width        int
	Height       int
	Keymap       keymap.DiaryMap
	DirPath      string
	CurrFileName string
	ShowFinder   bool
	CurrDate     time.Time
	Vp           viewport.Model
	Help         help.Model
	Finder       fzf.FuzzyFinder
	Renderer     *glamour.TermRenderer
}

func NewDiary(w int, h int) *Diary {
	today := time.Now().Format("2006-01-02") + ".md"
	dirpath := config.AppConfig.General.NotesDir

	r, _ := glamour.NewTermRenderer(glamour.WithStyles(config.ToGlamourStyle(config.AppConfig.Styles.Renderer)),
		glamour.WithWordWrap(w))
	content, _ := r.Render(ReadDiary(dirpath, today))

	vp := viewport.New(w, h-1)
	vp.Style = styles.BorderStyle().
		BorderForeground(colors.ColorPalette().Border).
		Foreground(colors.ColorPalette().Text)
	vp.SetContent(content)

	files, _ := AllFiles(dirpath)

	return &Diary{
		Width:        w,
		Height:       h,
		Keymap:       keymap.NewDiaryMap(),
		Vp:           vp,
		Help:         help.New(),
		Finder:       fzf.NewFzf(files, w, h),
		CurrFileName: today,
		DirPath:      dirpath,
		Renderer:     r,
	}
}

func (m *Diary) Init() tea.Cmd {
	return nil
}

func (m *Diary) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.FzfSelection:
		if msg.Exited {
			m.ShowFinder = false
			return m, nil
		}
		m.CurrFileName = msg.Selection
		m.Refresh()
		m.ShowFinder = false
		return m, nil
	case messages.EditorClose:
		m.Refresh()
		return m, nil
	case tea.KeyMsg:
		if m.ShowFinder {
			var cmd tea.Cmd
			m.Finder, cmd = m.Finder.Update(msg)
			return m, cmd
		}
		switch {
		case key.Matches(msg, m.Keymap.Edit):
			home, _ := os.UserHomeDir()

			args := config.AppConfig.General.Editor
			args = append(args, filepath.Join(home, m.DirPath, ".diary", m.CurrFileName))

			c := exec.Command(args[0], args[1:]...)
			cmd := tea.ExecProcess(c, func(err error) tea.Msg {
				return messages.EditorClose{
					Err: err,
				}
			})

			return m, cmd
		case key.Matches(msg, m.Keymap.OpenFinder):
			m.ShowFinder = true
			return m, nil
		case key.Matches(msg, m.Keymap.BackToMenu):
			return m, func() tea.Msg {
				return messages.ChangePage{Page: enums.MenuPage}
			}
		default:
			var cmd tea.Cmd
			m.Vp, cmd = m.Vp.Update(msg)
			return m, cmd
		}
	}
	var cmd tea.Cmd
	if m.ShowFinder {
		m.Finder, cmd = m.Finder.Update(msg)
	} else {
		m.Vp, cmd = m.Vp.Update(msg)
	}
	return m, cmd
}

func (m *Diary) View() string {
	if m.ShowFinder {
		return m.Finder.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, m.Vp.View(),
		lipgloss.NewStyle().PaddingLeft(2).Render(m.Help.View(keymap.NewDynamic(m.Keymap.Bindings()))))
}

func (m *Diary) Refresh() {
	data := ReadDiary(m.DirPath, m.CurrFileName)
	content, _ := m.Renderer.Render(string(data))
	m.Vp.SetContent(content)
}

func ReadDiary(dirpath string, today string) string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, dirpath, ".diary", today)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, _ := os.Create(path)

		now := time.Now()
		day := now.Day()
		suffix := getOrdinalSuffix(day)
		formatted := fmt.Sprintf("%d%s %s %d", day, suffix, now.Format("January"), now.Year())

		f.WriteString(fmt.Sprintf("# %s \n", formatted))
		f.Close()
	} else if err != nil {
		return fmt.Sprintf("Creating: %t", os.IsNotExist(err))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err.Error()
	}

	return string(data)
}

func AllFiles(dir string) ([]string, error) {
	home, _ := os.UserHomeDir()
	entries, err := os.ReadDir(filepath.Join(home, dir, ".diary"))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func getOrdinalSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
