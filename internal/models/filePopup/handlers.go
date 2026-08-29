package filepopup

import (
	"fmt"
	"os"
	"strings"

	"github.com/SourcewareLab/Toney/v2/internal/enums"
	filetree "github.com/SourcewareLab/Toney/v2/internal/fileTree"
	"github.com/SourcewareLab/Toney/v2/internal/messages"
	"github.com/SourcewareLab/Toney/v2/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleEnter(m *FilePopup) (tea.Model, tea.Cmd) {
	path := GetPath(m.Node)
	var err error

	switch m.Type {
	case enums.FileCreate:
		err = Create(path, m.TextInput.Value(), m.Node)
	case enums.FileDelete:
		err = Delete(path[0:len(path)-1], m.TextInput.Value())
	case enums.FileRename:
		err = Rename(path[0:len(path)-1], m.TextInput.Value())
	case enums.FileMove:
		err = Move(path[0:len(path)-1], m.TextInput.Value())
	default:
		fmt.Println("default")
	}

	if err != nil {
		return m, utils.ReturnError("FilePopup", "File Operation Error", err)
	}

	return m, tea.Batch(func() tea.Msg {
		return messages.HidePopupMessage{}
	},
		func() tea.Msg {
			return messages.RefreshFileExplorerMsg{}
		})
}

func Move(path string, value string) error {
	err := os.Rename(path, value)
	if err != nil {
		return err
	}
	return nil
}

func Rename(path string, value string) error {
	newpathArr := strings.Split(path, "/")
	newpath := strings.Join(newpathArr[0:len(newpathArr)-1], "/") + "/" + value
	err := os.Rename(path, newpath)
	if err != nil {
		return err
	}
	return nil
}

func Delete(path string, value string) error {
	if value != "y" {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func Create(path string, value string, n *filetree.Node) error {
	if !n.IsDirectory {
		pathArr := strings.Split(path, "/")
		pathArr = pathArr[0 : len(pathArr)-2]
		path = strings.Join(pathArr, "/") + "/"
	}

	if strings.HasSuffix(value, "/") {
		err := os.MkdirAll(path+value, 0o755)
		if err != nil {
			return err
		}
	} else {
		_, err := os.Create(path + value)
		if err != nil {
			return err
		}
	}
	return nil
}

func MapExpanded(new *filetree.Node, old *filetree.Node) {
	if old.IsExpanded {
		new.IsExpanded = true
	}

	if len(old.Children) != len(new.Children) {
		return
	}

	for idx, val := range old.Children {
		if val.Name != new.Children[idx].Name {
			continue
		}

		if !val.IsDirectory || !new.Children[idx].IsDirectory {
			continue
		}

		MapExpanded(new.Children[idx], val)
	}
}

func GetPath(n *filetree.Node) string {
	c := n

	home, _ := os.UserHomeDir()

	path := ""

	for {
		if c == nil {
			return home + "/" + path
		}

		path = c.Name + "/" + path
		c = c.Parent
	}
}
