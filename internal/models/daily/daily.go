package daily

import (
	"fmt"

	"github.com/SourcewareLab/Toney/internal/colors"
	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	taskpopup "github.com/SourcewareLab/Toney/internal/models/taskPopup"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Daily struct {
	Width      int
	Height     int
	List       list.Model
	Tasks      Tasks
	Popup      *taskpopup.TaskPopup
	ShowPopup  bool
	Keymap     keymap.DailyTaskMap
	Help       help.Model
	CurrentTab int
	Tabs       []enums.TaskTabs
}

func NewDaily(w int, h int) *Daily {
	tasks := GetItems()

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
		Keymap:     keymap.NewDailyTaskMap(),
		Help:       help.New(),
		CurrentTab: 0,
		Tabs:       []enums.TaskTabs{enums.All, enums.Unique, enums.Recurring, enums.Github},
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
		case key.Matches(msg, m.Keymap.TabRight):
			m.CurrentTab++
			if m.CurrentTab >= len(m.Tabs) {
				m.CurrentTab = m.CurrentTab % len(m.Tabs)
			}

			m.Refresh()
			return m, nil
		case key.Matches(msg, m.Keymap.TabLeft):
			m.CurrentTab--
			if m.CurrentTab < 0 {
				m.CurrentTab = len(m.Tabs) + m.CurrentTab
			}

			m.Refresh()
			return m, nil
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

	main := lipgloss.JoinVertical(lipgloss.Left,
		styles.GetDailyText(m.Width, m.Height/3),
		m.GetTabs(),
		lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Left, m.List.View()))

	if len(m.List.Items()) == 0 {
		main = lipgloss.JoinVertical(lipgloss.Left,
			styles.GetDailyText(m.Width, m.Height/3),
			lipgloss.Place(m.Width, 2*m.Height/3, lipgloss.Center, lipgloss.Top,
				lipgloss.NewStyle().Foreground(colors.ColorPalette().Text).Render("You have no Tasks!")))
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

	curr := m.Tasks.ItemsAsList()
	switch m.Tabs[m.CurrentTab] {
	case enums.Unique:
		curr = TaskToItems(m.Tasks.Unique)
	case enums.Recurring:
		curr = TaskToItems(m.Tasks.Recurring)
		// TODO: Github
	}

	lst := list.New(curr, TaskDelegate{}, m.Width/2, 2*m.Height/3)
	km := list.DefaultKeyMap()
	km.Quit.Unbind()
	lst.KeyMap = km
	lst.SetShowHelp(false)
	lst.SetShowTitle(false)

	m.List = lst
}
