package cmd

import (
	"bufio"
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

	reader := bufio.NewReader(os.Stdin)

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

	fmt.Print("Enter repository owner (username or organization): ")
	owner, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read owner: %w", err)
	}
	owner = strings.TrimSpace(owner)

	if owner == "" {
		return fmt.Errorf("repository owner is required")
	}

	fmt.Print("Enter repository name: ")
	repo, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read repo: %w", err)
	}
	repo = strings.TrimSpace(repo)

	if repo == "" {
		return fmt.Errorf("repository name is required")
	}

	viper.Set("github.enabled", true)
	viper.Set("github.token", token)
	viper.Set("github.owner", owner)
	viper.Set("github.repo", repo)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ GitHub integration configured successfully!\n")
	fmt.Printf("   Repository: %s/%s\n", owner, repo)
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

	fmt.Printf("🔄 Syncing issues from %s/%s...\n",
		config.AppConfig.GitHub.Owner,
		config.AppConfig.GitHub.Repo)

	api := github.NewGitHubAPI()
	issues, err := api.FetchIssues()
	if err != nil {
		return fmt.Errorf("failed to fetch issues: %w", err)
	}

	fmt.Printf("📥 Found %d issues\n", len(issues))

	converter := github.NewNoteConverter()
	successCount := 0

	for _, issue := range issues {
		_, err := converter.ConvertIssueToNote(issue)
		if err != nil {
			fmt.Printf("⚠️  Failed to convert issue #%d: %v\n", issue.Number, err)
			continue
		}
		successCount++
	}

	fmt.Printf("✅ Successfully converted %d issues to notes\n", successCount)
	fmt.Printf("📁 Notes saved to: %s/github-issues/\n", config.AppConfig.General.NotesDir)

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
		fmt.Printf("Repository: %s/%s\n",
			config.AppConfig.GitHub.Owner,
			config.AppConfig.GitHub.Repo)

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
