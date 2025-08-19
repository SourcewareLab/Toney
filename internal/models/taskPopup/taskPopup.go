package taskpopup

import (
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/keymap"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskPopup struct {
	Width      int
	Height     int
	Type       enums.TaskPopup
	Keymap     keymap.TaskPopupMap
	Help       help.Model
	Form       *Form
	Select     *SelectStatus
	DeleteForm *DeleteForm
	ShowSelect bool
}

func NewPopup(w int, h int, typ enums.TaskPopup) *TaskPopup {
	return &TaskPopup{
		Width:      w,
		Height:     h,
		Type:       typ,
		Keymap:     keymap.NewTaskPopupMap(),
		Help:       help.New(),
		Form:       NewForm(w, h),
		Select:     NewSelect(w/3, h),
		DeleteForm: NewDeleteForm(w/3, h),
		ShowSelect: false,
	}
}

func (m *TaskPopup) Init() tea.Cmd {
	return nil
}

func (m *TaskPopup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Enter):
			if m.Type == enums.CreateTask && !m.ShowSelect {
				m.ShowSelect = true
				return m, nil
			}
			return m, func() tea.Msg {
				return messages.TaskPopupMessage{ // returning all values, daily.go will handle the logic
					Title:     m.Form.TitleInput.Value(),
					Type:      m.Type,
					Desc:      m.Form.DescInput.Value(),
					Status:    m.Select.Opts[m.Select.Selected],
					IsDeleted: m.DeleteForm.isDeleting,
				}
			}
		case key.Matches(msg, m.Keymap.Exit):
			return m, func() tea.Msg {
				return messages.TaskPopupMessage{Type: enums.ClosePopup} // closes the popup
			}
		}
	}

	switch m.Type {
	case enums.CreateTask:
		if m.ShowSelect {
			updated, _ := m.Select.Update(msg)
			if sel, ok := updated.(*SelectStatus); ok { // Type matching, cause I cant assign it straightaway
				m.Select = sel
				return m, nil
			}
		}
		updated, _ := m.Form.Update(msg)
		if form, ok := updated.(*Form); ok { // Type matching, cause I cant assign it straightaway
			m.Form = form
			return m, nil
		}
	case enums.ChangeStatus:
		updated, _ := m.Select.Update(msg)
		if sel, ok := updated.(*SelectStatus); ok { // Type matching, cause I cant assign it straightaway
			m.Select = sel
			return m, nil
		}
	case enums.Delete:
		updated, _ := m.DeleteForm.Update(msg)
		if del, ok := updated.(*DeleteForm); ok { // Type matching, cause I cant assign it straightaway
			m.DeleteForm = del
			return m, nil
		}
	case enums.Edit:
		updated, _ := m.Form.Update(msg)
		if form, ok := updated.(*Form); ok { // Type matching, cause I cant assign it straightaway
			m.Form = form
			return m, nil
		}
	}

	return m, nil
}

func (m *TaskPopup) View() string {
	view := ""

	switch m.Type {
	case enums.CreateTask:
		if m.ShowSelect {
			view = lipgloss.JoinVertical(lipgloss.Center,
				styles.GetSelectStatus(m.Width, m.Height/2),
				lipgloss.Place(m.Width, m.Height/2, lipgloss.Center, lipgloss.Top, m.Select.View()))
		} else {
			view = m.Form.View()
		}
	case enums.ChangeStatus:
		view = lipgloss.JoinVertical(lipgloss.Center,
			styles.GetSelectStatus(m.Width, m.Height/2),
			lipgloss.Place(m.Width, m.Height/2, lipgloss.Center, lipgloss.Top, m.Select.View()))
	case enums.Delete:
		view = lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.Place(m.Width, m.Height/2, lipgloss.Center, lipgloss.Center, styles.GetSelectStatus(m.Width, m.Height/3)),
			lipgloss.Place(m.Width, m.Height/2, lipgloss.Center, lipgloss.Top, m.DeleteForm.View()))
	case enums.Edit:
		view = m.Form.View()
	}

	binds := m.Keymap.Bindings()
	if m.Type == enums.CreateTask || m.Type == enums.Edit {
		binds = append(binds, m.Form.Keymap.Bindings()...)
	}

	return lipgloss.JoinVertical(lipgloss.Center, view, m.Help.View(keymap.NewDynamic(binds)))
}
