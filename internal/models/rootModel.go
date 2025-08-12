package models

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/SourcewareLab/Toney/internal/messages"
	"github.com/SourcewareLab/Toney/internal/models/daily"
	"github.com/SourcewareLab/Toney/internal/models/diary"
	errorpopup "github.com/SourcewareLab/Toney/internal/models/error"
	filepopup "github.com/SourcewareLab/Toney/internal/models/filePopup"
	"github.com/SourcewareLab/Toney/internal/models/github"
	homemodel "github.com/SourcewareLab/Toney/internal/models/homeModel"
	"github.com/SourcewareLab/Toney/internal/models/menu"
	"github.com/SourcewareLab/Toney/internal/utils"
	overlay "github.com/rmhubbert/bubbletea-overlay"

	"github.com/SourcewareLab/Toney/internal/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	Width          int
	Height         int
	Page           enums.Page
	Home           *homemodel.HomeModel
	Menu           *menu.Menu
	Daily          *daily.Daily
	Diary          *diary.Diary
  GitHub        *github.GitHubModel
	ErrorPopup     *errorpopup.ErrorPopup
	Overlay        *overlay.Model
	CurrentPage    enums.Page
	ShowPopup      bool
	FilePopupType  enums.PopupType
	FilePopup      *filepopup.FilePopup
	isLoading      bool
	IsShowingError bool
}

func NewRoot() *RootModel {
	return &RootModel{
		Page:      enums.MenuPage,
		ShowPopup: false,
		isLoading: true,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.CloseError:
		m.IsShowingError = false
		return m, nil
	case messages.ErrorMsg:
		m.IsShowingError = true
		m.ErrorPopup = errorpopup.NewErrorPopup(m.Width, m.Height, msg.Msg, msg.Title, msg.Locn)

		var curr tea.Model
		switch m.CurrentPage {
		case enums.DailyPage:
			curr = m.Daily
		case enums.HomePage:
			curr = m.Home
		case enums.DiaryPage:
			curr = m.Diary
		}

		m.Overlay = overlay.New(m.ErrorPopup, curr, overlay.Right, overlay.Top, 2, 2)
		return m, tea.Tick(3*time.Second, func(_ time.Time) tea.Msg { return messages.CloseError{} })
	case messages.FzfSelection:
		switch m.CurrentPage {
		case enums.MenuPage:
			pg, cmd := m.Menu.Update(msg)
			if d, ok := pg.(*menu.Menu); ok { // Type matching, cause I cant assign it straightaway
				m.Menu = d
				return m, cmd
			}
		case enums.HomePage:
			pg, cmd := m.Home.Update(msg)
			if d, ok := pg.(*homemodel.HomeModel); ok { // Type matching, cause I cant assign it straightaway
				m.Home = d
				return m, cmd
			}
		case enums.DailyPage:
			pg, cmd := m.Daily.Update(msg)
			if d, ok := pg.(*daily.Daily); ok { // Type matching, cause I cant assign it straightaway
				m.Daily = d
				return m, cmd
			}
		case enums.DiaryPage:
			pg, cmd := m.Diary.Update(msg)
			if d, ok := pg.(*diary.Diary); ok { // Type matching, cause I cant assign it straightaway
				m.Diary = d
				return m, cmd
			}
		case enums.GitHubPage:
			if m.GitHub != nil {
				pg, cmd := m.GitHub.Update(msg)
				if d, ok := pg.(*github.GitHubModel); ok { // Type matching, cause I cant assign it straightaway
					m.GitHub = d
					return m, cmd
				}
			}
		}
	case messages.TaskPopupMessage:
		if m.CurrentPage == enums.DailyPage {
			dailypage, cmd := m.Daily.Update(msg)
			if d, ok := dailypage.(*daily.Daily); ok { // Type matching, cause I cant assign it straightaway
				m.Daily = d
				return m, cmd
			}
			return m, cmd
		}
	case messages.ChangePage:
		switch msg.Page {
		case enums.MenuPage:
			m.CurrentPage = enums.MenuPage
		case enums.HomePage:
			m.CurrentPage = enums.HomePage
		case enums.DailyPage:
			m.CurrentPage = enums.DailyPage
		case enums.DiaryPage:
			m.CurrentPage = enums.DiaryPage
		case enums.GitHubPage:
			m.CurrentPage = enums.GitHubPage
			// Always create a new model to handle setup vs issues dynamically
			githubModel := github.NewGitHubSetupOrIssuesModel(m.Width, m.Height)
			if gitHubIssuesModel, ok := githubModel.(*github.GitHubModel); ok {
				m.GitHub = gitHubIssuesModel
				// Auto-sync issues when entering the page if configured
				return m, m.GitHub.SyncIssues()
			} else {
				// Handle setup model - we need to store it differently
				// For now, we'll handle this in the View method
				m.GitHub = nil
			}
		case enums.Quit:
			return m, tea.Quit
		}
	case messages.ShowLoader:
		m.isLoading = true
		return m, nil
	case messages.HideLoader:
		m.isLoading = false
		return m, nil
	case messages.ShowPopupMessage:
		m.FilePopup = filepopup.NewPopup(msg.Type, msg.Curr)
		m.ShowPopup = true
	case messages.HidePopupMessage:
		m.ShowPopup = false
	case messages.RefreshFileExplorerMsg:
		m.Home.FileExplorer.Update(msg)
		return m, nil
	case messages.EditorClose:
		if m.CurrentPage == enums.DiaryPage {
			m.Diary.Update(msg)
			return m, nil
		}
		if msg.Err != nil {
			return m, utils.ReturnError("Root", "Error Closing Editor", msg.Err)
		}
		m.Home.FileExplorer.Update(msg)
		m.Home.Viewer.Update(msg)

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case config.AppConfig.Keybinds.Global.Script:
			// Script
			script := strings.Join(config.AppConfig.General.Script, " ")
			command := exec.Command("bash", "-c", script)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				return m, utils.ReturnError("Root", "Error Closing Editor", err)
			}
			return m, nil
		case "ctrl+c":
			// Stop script
			script := strings.Join(config.AppConfig.General.StopScript, " ")
			command := exec.Command("bash", "-c", script)

			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				return m, utils.ReturnError("Root", "Error Closing Editor", err)
			}
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		if m.Home != nil { // Checking whether this is an app resize or app open
			m.Menu.Update(msg)
			m.Home.Update(msg)
			m.Diary.Update(msg)
			m.Daily.Update(msg)
		} else {
			m.Home = homemodel.NewHome(msg.Width, msg.Height)
			m.Menu = menu.NewMenu(msg.Width, msg.Height)
			m.Daily = daily.NewDaily(msg.Width, msg.Height)
			m.Diary = diary.NewDiary(msg.Width, msg.Height)
		}

		m.isLoading = false

		return m, nil
	}

	var cmd tea.Cmd

	if m.ShowPopup {
		_, cmd = m.FilePopup.Update(msg)
	} else {
		switch m.CurrentPage {
		case enums.MenuPage:
			_, cmd = m.Menu.Update(msg)
		case enums.HomePage:
			_, cmd = m.Home.Update(msg)
		case enums.DailyPage:
			_, cmd = m.Daily.Update(msg)
		case enums.DiaryPage:
			_, cmd = m.Diary.Update(msg)
		case enums.GitHubPage:
			if m.GitHub != nil {
				_, cmd = m.GitHub.Update(msg)
			}
		default:
			fmt.Printf("UNHANDLED MSG: %#v\n", msg)
		}
	}

	return m, cmd
}

func (m *RootModel) View() string {
	if m.isLoading {
		return lipgloss.NewStyle().Render("Loading...")
	}

	if m.IsShowingError {
		return m.Overlay.View()
	}

	if m.ShowPopup && m.FilePopup != nil {
		return lipgloss.Place(m.Width+2, m.Height+2, lipgloss.Center, lipgloss.Center, m.FilePopup.View())
	}

	switch m.CurrentPage {
	case enums.HomePage:
		return lipgloss.NewStyle().Background(colors.ColorPalette().Background).Render(m.Home.View())
	case enums.MenuPage:
		return m.Menu.View()
	case enums.DailyPage:
		return m.Daily.View()
	case enums.DiaryPage:
		return m.Diary.View()
	case enums.GitHubPage:
		if m.GitHub != nil {
			return m.GitHub.View()
		}
		return "GitHub model not initialized"
	default:
		return lipgloss.NewStyle().Background(colors.ColorPalette().Background).Render(m.Home.View())
	}
}
