package daily

import (
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	taskpopup "github.com/SourcewareLab/Toney/internal/models/taskPopup"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/SourcewareLab/Toney/internal/ui/theme"
	"github.com/SourcewareLab/Toney/internal/ui/widgets"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Daily struct {
	Width     int
	Height    int
	List      list.Model
	Tasks     Tasks
	Popup     *taskpopup.TaskPopup
	ShowPopup bool
	Keymap    keymap.DailyTaskMap
	Help      help.Model
}

func NewDaily(w int, h int) *Daily {
	tasks := GetItems()

	lst := list.New(tasks.ItemsAsList(), TaskDelegate{}, w/2, 2*h/3)
	km := list.DefaultKeyMap()
	km.Quit.Unbind()
	lst.KeyMap = km
	lst.SetShowHelp(false)

	return &Daily{
		Width:  w,
		Height: h,
		List:   lst,
		Tasks:  tasks,
		Keymap: keymap.NewDailyTaskMap(),
		Help:   help.New(),
	}
}

func (m *Daily) Init() tea.Cmd {
	return nil
}

func (m *Daily) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.ChangeStatus)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.DeleteTask):
			m.Popup = taskpopup.NewPopup(m.Width, m.Height, enums.Delete)
			m.ShowPopup = true
			return m, nil
		case key.Matches(msg, m.Keymap.EditTask):
			item := m.List.SelectedItem()

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

	pal := theme.DefaultDarkPalette()
	header := widgets.Header{Palette: pal, Width: m.Width, Title: "Daily Tasks", Right: ""}.View()
	footer := widgets.Footer{Palette: pal, Width: m.Width, Hints: []widgets.KeyHint{
		{Key: "↑↓", Desc: "navigate"},
		{Key: "enter", Desc: "select"},
		{Key: "esc", Desc: "back"},
	}}.View()

	bodyH := m.Height - lipgloss.Height(header) - lipgloss.Height(footer)
	if bodyH < 3 {
		bodyH = 3
	}

	// Compose top summary text and list inside a bordered box
	summary := styles.GetDailyText(m.Width, bodyH/3)
	listArea := lipgloss.Place(m.Width, 2*bodyH/3, lipgloss.Center, lipgloss.Center, m.List.View())

	if len(m.List.Items()) == 0 {
		listArea = lipgloss.Place(m.Width, 2*bodyH/3, lipgloss.Center, lipgloss.Top,
			lipgloss.NewStyle().Foreground(pal.Fg).Render("You have no Tasks!"))
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

func (m *Daily) Refresh() {
	m.Tasks = GetItems()

	lst := list.New(m.Tasks.ItemsAsList(), TaskDelegate{}, m.Width/2, 2*m.Height/3)
	km := list.DefaultKeyMap()
	km.Quit.Unbind()
	lst.KeyMap = km
	lst.SetShowHelp(false)

	m.List = lst
}
