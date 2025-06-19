package messages

import (
	"toney/internal/enums"
	filetree "toney/internal/fileTree"

	"github.com/charmbracelet/glamour"
)

type (
	ShowLoader struct{}

	HideLoader struct{}

	EditorClose struct {
		Err error
	}

	RendererCreated struct {
		Renderer *glamour.TermRenderer
	}

	ShowPopupMessage struct {
		Type enums.PopupType
		Curr *filetree.Node
	}

	RefreshFileExplorerMsg struct{}

	NvimDoneMsg struct {
		Err string
	}

	ChangeFileMessage struct {
		Path string
	}

	RefreshView struct {
		Content string
		Path    string
	}
	HidePopupMessage struct{}
)
