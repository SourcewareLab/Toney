package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
)

type NoteConverter struct {
	notesDir string
}

func NewNoteConverter() *NoteConverter {
	return &NoteConverter{
		notesDir: config.AppConfig.General.NotesDir,
	}
}

func (nc *NoteConverter) ConvertIssueToNote(issue GitHubIssue) (string, error) {
	// Create GitHub issues subdirectory if it doesn't exist
	githubNotesDir := filepath.Join(nc.notesDir, "github-issues")
	if err := os.MkdirAll(githubNotesDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create github-issues directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("issue-%d-%s.md", issue.Number, sanitizeFilename(issue.IssueTitle))
	filePath := filepath.Join(githubNotesDir, filename)

	// Generate note content
	content := nc.generateNoteContent(issue)

	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write note file: %w", err)
	}

	return filePath, nil
}

func (nc *NoteConverter) generateNoteContent(issue GitHubIssue) string {
	var content strings.Builder

	// Title and metadata
	content.WriteString(fmt.Sprintf("# Issue #%d: %s\n\n", issue.Number, issue.IssueTitle))

	// Metadata section
	content.WriteString("## Metadata\n\n")
	content.WriteString(fmt.Sprintf("- **Status**: %s\n", strings.Title(issue.State)))
	content.WriteString(fmt.Sprintf("- **Author**: %s\n", issue.Author))
	content.WriteString(fmt.Sprintf("- **Created**: %s\n", formatDate(issue.CreatedAt)))
	content.WriteString(fmt.Sprintf("- **Updated**: %s\n", formatDate(issue.UpdatedAt)))
	content.WriteString(fmt.Sprintf("- **URL**: [View on GitHub](%s)\n", issue.HTMLURL))

	// Labels section
	if len(issue.Labels) > 0 {
		content.WriteString("- **Labels**: ")
		labelStrings := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labelStrings[i] = fmt.Sprintf("`%s`", label.Name)
		}
		content.WriteString(strings.Join(labelStrings, ", "))
		content.WriteString("\n")
	}

	content.WriteString("\n")

	// Description section
	if issue.Body != "" {
		content.WriteString("## Description\n\n")
		content.WriteString(issue.Body)
		content.WriteString("\n\n")
	}

	// Notes section for user additions
	content.WriteString("## Notes\n\n")
	content.WriteString("<!-- Add your personal notes about this issue here -->\n\n")

	// Tasks section
	content.WriteString("## Tasks\n\n")
	content.WriteString("- [ ] Review issue details\n")
	if issue.State == "open" {
		content.WriteString("- [ ] Work on implementation\n")
		content.WriteString("- [ ] Test solution\n")
		content.WriteString("- [ ] Submit pull request\n")
	} else {
		content.WriteString("- [x] Issue resolved\n")
	}
	content.WriteString("\n")

	// Footer with sync info
	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("*Synced from GitHub on %s*\n", time.Now().Format("2006-01-02 15:04:05")))

	return content.String()
}

func sanitizeFilename(title string) string {
	// Replace invalid filename characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := title

	for _, char := range invalid {
		sanitized = strings.ReplaceAll(sanitized, char, "-")
	}

	// Replace multiple spaces and dashes with single dash
	sanitized = strings.ReplaceAll(sanitized, "  ", " ")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ReplaceAll(sanitized, "--", "-")

	// Trim and lowercase
	sanitized = strings.ToLower(strings.Trim(sanitized, "-"))

	// Limit length
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
		sanitized = strings.TrimSuffix(sanitized, "-")
	}

	return sanitized
}

func formatDate(dateStr string) string {
	// Parse GitHub's ISO 8601 format
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr // Return original if parsing fails
	}

	// Format as human-readable date
	return t.Format("January 2, 2006 at 15:04")
}
