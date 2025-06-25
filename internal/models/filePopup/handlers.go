package filepopup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SourcewareLab/Toney/internal/enums"
	filetree "github.com/SourcewareLab/Toney/internal/fileTree"
	"github.com/SourcewareLab/Toney/internal/messages"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleEnter(m *FilePopup) (tea.Model, tea.Cmd) {
	path := GetPath(m.Node)

	switch m.Type {
	case enums.FileCreate:
		Create(path, m.TextInput.Value(), m.Node)
	case enums.FileDelete:
		Delete(path[0:len(path)-1], m.TextInput.Value())
	case enums.FileRename:
		Rename(path[0:len(path)-1], m.TextInput.Value())
	case enums.FileMove:
		Move(path[0:len(path)-1], m.TextInput.Value())
	default:
		fmt.Println("default")
	}

	return m, tea.Batch(func() tea.Msg {
		return messages.HidePopupMessage{}
	},
		func() tea.Msg {
			return messages.RefreshFileExplorerMsg{}
		})
}

func Move(path string, value string) {
	err := os.Rename(path, value)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Rename(path string, value string) {
	newpathArr := strings.Split(path, "/")
	newpath := strings.Join(newpathArr[0:len(newpathArr)-1], "/") + "/" + value
	fmt.Println(newpath)
	err := os.Rename(path, newpath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Delete(path string, value string) {
	if value != "y" {
		return
	}

	err := os.Remove(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Create(path string, value string, n *filetree.Node) {
	if !n.IsDirectory {
		pathArr := strings.Split(path, "/")
		pathArr = pathArr[0 : len(pathArr)-2]
		path = strings.Join(pathArr, "/") + "/"
	}

	if strings.HasSuffix(value, "/") {
		err := os.MkdirAll(path+value, 0o755)
		if err != nil {
			fmt.Println(err.Error())
			// os.Exit(1)
		}
	} else {
		_, err := os.Create(path + value)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
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

func GetPath(node *filetree.Node) string {
	if node == nil {
		return ""
	}

	segments := []string{}
	for n := node; n != nil; n = n.Parent {
		segments = append([]string{n.Name}, segments...)
	}

	return filepath.Join(segments...)
}
