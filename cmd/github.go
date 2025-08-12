package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/models/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub integration commands",
	Long:  `Manage GitHub integration for syncing issues as structured notes.`,
}

var githubSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup GitHub credentials and configuration",
	Long:  `Interactive setup for GitHub integration including token, owner, and repository configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupGitHub(); err != nil {
			fmt.Printf("Error setting up GitHub: %v\n", err)
			os.Exit(1)
		}
	},
}

var githubSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync GitHub issues to notes",
	Long:  `Fetch GitHub issues and convert them to structured notes in your notes directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncGitHubIssues(); err != nil {
			fmt.Printf("Error syncing GitHub issues: %v\n", err)
			os.Exit(1)
		}
	},
}

var githubStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show GitHub integration status",
	Long:  `Display current GitHub configuration and connection status.`,
	Run: func(cmd *cobra.Command, args []string) {
		showGitHubStatus()
	},
}

func setupGitHub() error {
	fmt.Println("🔧 GitHub Integration Setup")
	fmt.Println("===========================")

	if err := config.SetConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Print("Enter your GitHub Personal Access Token: ")
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	token := strings.TrimSpace(string(tokenBytes))
	fmt.Println()

	if token == "" {
		return fmt.Errorf("GitHub token is required")
	}

	viper.Set("github.enabled", true)
	viper.Set("github.token", token)
	// Clear owner/repo if previously set so we default to all repositories
	viper.Set("github.owner", "")
	viper.Set("github.repo", "")

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ GitHub integration configured successfully!\n")
	fmt.Printf("   Scope: All repositories accessible by your token\n")
	fmt.Printf("   Status: Enabled\n")
	fmt.Println("\nYou can now use 'toney github sync' to fetch issues as notes.")

	return nil
}

func syncGitHubIssues() error {
	if err := config.SetConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !config.AppConfig.GitHub.Enabled {
		return fmt.Errorf("GitHub integration is not enabled. Run 'toney github setup' first")
	}

	if config.AppConfig.GitHub.Token == "" {
		return fmt.Errorf("GitHub token not configured. Run 'toney github setup' first")
	}

    // Sync across all repositories
    fmt.Printf("🔄 Syncing issues from all repositories...\n")

    api := github.NewGitHubAPI()
    issues, err := api.FetchIssues()
    if err != nil {
        return fmt.Errorf("failed to fetch issues: %w", err)
    }

	fmt.Printf("📥 Found %d issues\n", len(issues))

	return nil
}

func showGitHubStatus() {
	if err := config.SetConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	fmt.Println("📊 GitHub Integration Status")
	fmt.Println("============================")

	if config.AppConfig.GitHub.Enabled {
		fmt.Println("Status: ✅ Enabled")
		fmt.Println("Scope: All repositories")

		if config.AppConfig.GitHub.Token != "" {
			fmt.Println("Token: ✅ Configured")
		} else {
			fmt.Println("Token: ❌ Not configured")
		}
	} else {
		fmt.Println("Status: ❌ Disabled")
		fmt.Println("Run 'toney github setup' to enable GitHub integration")
	}
}

func init() {
	rootCmd.AddCommand(githubCmd)
	githubCmd.AddCommand(githubSetupCmd)
	githubCmd.AddCommand(githubSyncCmd)
	githubCmd.AddCommand(githubStatusCmd)
}
