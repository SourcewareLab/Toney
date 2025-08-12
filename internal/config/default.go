package config

import "github.com/charmbracelet/lipgloss"

func DefaultConfig() Config {
	return Config{
		General: GeneralConfig{
			Editor:      []string{"nvim"},
			NotesDir:    ".toney", // From $HOME Directory
			StartScript: []string{},
			StopScript:  []string{},
			Script:      []string{},
		},
		Keybinds: KeybindsConfig{
			Global: GlobalKeybinds{
				Up:     "up",
				Down:   "down",
				Script: "ctrl+s",
			},
			Fuzz: FuzzyKeybinds{
				Up:           "up",
				Down:         "down",
				Enter:        "enter",
				StartWriting: "/",
				Exit:         "esc",
			},
			Diary: DiaryKeybinds{
				ScrollUp:   "up",
				ScrollDown: "down",
				Edit:       "e",
				Finder:     "f",
				BackToMenu: "esc",
			},
			Home: HomeKeybinds{
				FocusViewer:   "V",
				FocusExplorer: "F",
				Create:        "c",
				Rename:        "r",
				ScrollUp:      "up",
				ScrollDown:    "down",
				Move:          "m",
				Delete:        "d",
				Edit:          "enter",
				BackToMenu:    "esc",
				Finder:        "f",
			},
			Daily: DailyKeybinds{
				Create:          "c",
				CreateRecurring: "r",
				Edit:            "e",
				StatusChange:    "s",
				Delete:          "d",
				ExitPopup:       "esc",
				Enter:           "enter",
				FormUp:          "ctrl+up",
				FormDown:        "ctrl+down",
				BackToMenu:      "esc",
			},
		},
		Styles: StylesConfig{
			Text:             "#cdd6f4",
			Background:       "#1e1e2e",
			Border:           "#45475a",
			FocusedBorder:    "#b4befe",
			MenuSelectedBg:   "#b4befe",
			MenuSelectedText: "#1e1e2e",
			ErrorBg:          "#11111b",
			ErrorText:        "#f38ba8",
			Icons: IconsConfig{
				FolderIcon: "󰷏",
				FileIcon:   "",
				TaskIcons: TaskIcons{
					CompletedIcon: "✓",
					AbandonedIcon: "×",
					PendingIcon:   "~",
					StartedIcon:   "○",
				},
			},
			TaskStyles: TaskStylesConfig{
				FocusedBar:   "#b4befe",
				UnfocusedBar: "#45475a",
				CompletedStyle: TaskStateStyle{
					Title: "#a6e3a1",
					Desc:  "#5a7a57",
				},
				AbandonedStyle: TaskStateStyle{
					Title: "#f38ba8",
					Desc:  "#894454",
				},
				PendingStyle: TaskStateStyle{
					Title: "#6c7086",
					Desc:  "#313244",
				},
				StartedStyle: TaskStateStyle{
					Title: "#f9e2af",
					Desc:  "#a38e65",
				},
			},
			Renderer: RendererConfig{
				ItemPrefix:        lipgloss.NewStyle().Foreground(lipgloss.Color("#45475a")).Render("• "),
				EnumerationPrefix: ". ",
				Document: BlockStyle{
					BlockPrefix: "\n",
					BlockSuffix: "\n",
					Margin:      1,
					Color:       "#cdd6f4",
				},
				BlockQuote: BlockStyle{
					Indent:      1,
					IndentToken: "│ ",
				},
				Heading: HeadingStyle{
					Base: BlockStyle{
						BlockSuffix: "\n",
						Color:       "#94e2d5", // Pink
						Bold:        true,
					},
					Levels: []BlockStyle{
						{Prefix: " ", Suffix: " ", Color: "#74c7ec", Background: "#313244", Bold: true}, // Sapphire bg
						{Prefix: "## "},
						{Prefix: "### "},
						{Prefix: "#### "},
						{Prefix: "##### "},
						{Prefix: "###### ", Color: "#9399b2", Bold: false}, // Overlay2
					},
				},
				HorizontalRule: InlineStyle{
					Color:  "#585b70", // Surface2
					Format: "\n────────\n",
				},
				List: ListStyle{
					LevelIndent: 2,
					Styles: InlineStyle{
						Color:  "#cdd6f4", // Green
						Bold:   false,
						Suffix: "\n",
					},
					Task: &TaskStyle{
						Ticked:   lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Render("• [✓] "),
						Unticked: lipgloss.NewStyle().Foreground(lipgloss.Color("#585b70")).Render("• [ ] "),
					},
				},
				Link:          InlineStyle{Color: "#89b4fa", Underline: true}, // Blue
				Image:         InlineStyle{Color: "#fab387", Underline: true}, // Peach
				Emph:          InlineStyle{Italic: true, Color: "#89b4fa"},    // Pink
				Strong:        InlineStyle{Bold: true, Color: "#74c7ec"},
				Strikethrough: InlineStyle{Strikethrough: true, Color: "#f38ba8"}, // Red

				Code: BlockStyle{
					Prefix:     " ",
					Suffix:     " ",
					Color:      "#89dceb", // Sky
					Background: "#11111b", // Base bg
				},

				CodeBlock: CodeBlockStyle{
					BlockStyle: BlockStyle{Background: "#11111b", Color: "#89dceb", Margin: 1},
					Text:       InlineStyle{Color: "#cdd6f4"}, // Text
				},

				Table: TableStyle{
					BlockStyle:      BlockStyle{Color: "#cdd6f4"},
					Header:          BlockStyle{Color: "#f5c2e7", Background: "#1e1e2e", Bold: true},
					Cell:            BlockStyle{Color: "#cdd6f4", Background: "#1e1e2e"},
					CenterSeparator: "|",
					ColumnSeparator: "|",
					RowSeparator:    "-",
				},
			},
		},
	}
}
