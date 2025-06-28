package messages

import (
	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
)

type (
	ShowLoader struct{}

	HideLoader struct{}

	EditorClose struct {
		Err      error
		IsNative bool
		Value    string
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
