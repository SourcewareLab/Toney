package github

import (
	"fmt"
	"strings"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type GitHubSetupModel struct {
	Width   int
	Height  int
	Focused bool
	Step    int // 0: token, 1: owner, 2: repo, 3: confirmation
	Inputs  []textinput.Model
	Error   string
	Success bool
}

type GitHubSetupCompleteMsg struct {
	Success bool
	Error   error
}

func NewGitHubSetupModel() *GitHubSetupModel {
	inputs := make([]textinput.Model, 3)

	// Token input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter your GitHub Personal Access Token"
	inputs[0].Focus()
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].EchoCharacter = '*'

	// Owner input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Enter repository owner (username or organization)"

	// Repo input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Enter repository name"

	return &GitHubSetupModel{
		Step:    0,
		Inputs:  inputs,
		Focused: true,
	}
}

func (m *GitHubSetupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *GitHubSetupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case GitHubSetupCompleteMsg:
		if msg.Error != nil {
			m.Error = msg.Error.Error()
			m.Success = false
		} else {
			m.Success = true
			m.Error = ""
		}

	case tea.KeyMsg:
		if !m.Focused {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.Step < 2 {
				// Move to next step
				m.Inputs[m.Step].Blur()
				m.Step++
				m.Inputs[m.Step].Focus()
				return m, textinput.Blink
			} else if m.Step == 2 {
				// Final step - save configuration
				return m, m.saveConfiguration()
			} else if m.Success {
				// Setup complete, return to menu
				return m, tea.Quit
			}

		case "tab", "shift+tab":
			// Navigate between inputs
			if m.Step > 0 {
				m.Inputs[m.Step].Blur()
				if msg.String() == "tab" {
					m.Step = (m.Step + 1) % 3
				} else {
					m.Step = (m.Step - 1 + 3) % 3
				}
				m.Inputs[m.Step].Focus()
				return m, textinput.Blink
			}
		}
	}

	// Update current input
	if m.Step < 3 && !m.Success {
		m.Inputs[m.Step], cmd = m.Inputs[m.Step].Update(msg)
	}

	return m, cmd
}

func (m *GitHubSetupModel) View() string {
	if m.Success {
		return m.renderSuccessView()
	}

	if m.Error != "" {
		return m.renderErrorView()
	}

	return m.renderSetupView()
}

func (m *GitHubSetupModel) renderSetupView() string {
	var content strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#b4befe")).
		Bold(true).
		Render("🔧 GitHub Integration Setup")

	content.WriteString(title + "\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render("Configure your GitHub integration to access issues as structured notes.\n\n")

	content.WriteString(instructions)

	// Step indicators
	stepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086"))
	activeStepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#b4befe")).Bold(true)

	steps := []string{"Token", "Owner", "Repository"}
	stepIndicators := ""
	for i, step := range steps {
		if i == m.Step {
			stepIndicators += activeStepStyle.Render(fmt.Sprintf("[%d] %s", i+1, step))
		} else if i < m.Step {
			stepIndicators += stepStyle.Render(fmt.Sprintf("[✓] %s", step))
		} else {
			stepIndicators += stepStyle.Render(fmt.Sprintf("[ ] %s", step))
		}
		if i < len(steps)-1 {
			stepIndicators += "  →  "
		}
	}
	content.WriteString(stepIndicators + "\n\n")

	// Current input
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		Padding(1, 2)

	if m.Step < 3 {
		// Input labels
		labels := []string{
			"GitHub Personal Access Token:",
			"Repository Owner:",
			"Repository Name:",
		}

		content.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#cdd6f4")).
			Bold(true).
			Render(labels[m.Step]) + "\n")

		content.WriteString(inputStyle.Render(m.Inputs[m.Step].View()) + "\n\n")

		// Help text
		helpTexts := []string{
			"Create a Personal Access Token at: https://github.com/settings/tokens\nRequired scopes: repo (for private repos) or public_repo (for public repos)",
			"Enter the username or organization that owns the repository",
			"Enter the name of the repository you want to sync issues from",
		}

		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6c7086")).
			Italic(true)

		content.WriteString(helpStyle.Render(helpTexts[m.Step]) + "\n\n")
	}

	// Navigation help
	navHelp := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6c7086")).
		Render("Press Enter to continue • Tab/Shift+Tab to navigate • Esc to cancel")

	content.WriteString(navHelp)

	// Wrap in container
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		Padding(2, 4).
		Width(min(80, m.Width-4))

	return containerStyle.Render(content.String())
}

func (m *GitHubSetupModel) renderSuccessView() string {
	successStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#a6e3a1")).
		Padding(2, 4).
		Width(min(60, m.Width-4))

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a6e3a1")).
		Bold(true).
		Render("✅ GitHub Integration Configured Successfully!") + "\n\n"

	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render(fmt.Sprintf("Repository: %s/%s\n", m.Inputs[1].Value(), m.Inputs[2].Value())) +
		"Status: Enabled\n\n" +
		"The 'GitHub Issues' option will now appear in the main menu.\n" +
		"You can browse and convert issues to structured notes.\n\n" +
		"Press Enter to return to the main menu."

	return successStyle.Render(content)
}

func (m *GitHubSetupModel) renderErrorView() string {
	errorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f38ba8")).
		Padding(2, 4).
		Width(min(60, m.Width-4))

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f38ba8")).
		Bold(true).
		Render("❌ Setup Failed") + "\n\n"

	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("#cdd6f4")).
		Render(fmt.Sprintf("Error: %s\n\n", m.Error)) +
		"Please check your credentials and try again.\n\n" +
		"Press Enter to retry or Esc to cancel."

	return errorStyle.Render(content)
}

func (m *GitHubSetupModel) saveConfiguration() tea.Cmd {
	return func() tea.Msg {
		token := strings.TrimSpace(m.Inputs[0].Value())
		owner := strings.TrimSpace(m.Inputs[1].Value())
		repo := strings.TrimSpace(m.Inputs[2].Value())

		if token == "" || owner == "" || repo == "" {
			return GitHubSetupCompleteMsg{
				Success: false,
				Error:   fmt.Errorf("all fields are required"),
			}
		}

		// Test the credentials by making a simple API call
		api := &GitHubAPI{
			token: token,
			owner: owner,
			repo:  repo,
		}

		_, err := api.FetchIssues()
		if err != nil {
			return GitHubSetupCompleteMsg{
				Success: false,
				Error:   fmt.Errorf("failed to validate credentials: %w", err),
			}
		}

		// Save to configuration
		viper.Set("github.enabled", true)
		viper.Set("github.token", token)
		viper.Set("github.owner", owner)
		viper.Set("github.repo", repo)

		if err := viper.WriteConfig(); err != nil {
			return GitHubSetupCompleteMsg{
				Success: false,
				Error:   fmt.Errorf("failed to save configuration: %w", err),
			}
		}

		// Reload config
		if err := config.SetConfig(); err != nil {
			return GitHubSetupCompleteMsg{
				Success: false,
				Error:   fmt.Errorf("failed to reload configuration: %w", err),
			}
		}

		return GitHubSetupCompleteMsg{
			Success: true,
			Error:   nil,
		}
	}
}

func (m *GitHubSetupModel) SetFocus(focused bool) {
	m.Focused = focused
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
