package filetree

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/styles"
)

func CreateTree() (*Node, error) {

	root, err := buildTree(nil, config.NoteRootPath, 0)
	if err != nil {
		return nil, err
	}

	root.IsExpanded = true

	return root, nil
}

func buildTree(parent *Node, path string, depth int) (*Node, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := Node{
		Name:        info.Name(),
		Parent:      parent,
		Depth:       depth,
		IsDirectory: info.IsDir(),
		IsExpanded:  false,
	}

	if node.IsDirectory {
		files, _ := os.ReadDir(path)
		for _, file := range files {
			childPath := filepath.Join(path, file.Name())
			child, _ := buildTree(&node, childPath, depth+1)
			node.Children = append(node.Children, child)
		}
	}

	return &node, nil
}

func BuildNodeTree(n *Node, prefix string, isLast bool, curr *Node) string {
	var sb strings.Builder

	branch := "├─ "
	if isLast {
		branch = "└─ "
	}

	icon := "📄"
	if n.IsDirectory {
		icon = "📁"
	}

	newPrefix := prefix

	// Build the line for this node
	line := ""
	if n.Parent == nil {
		line = icon + " " + n.Name
	} else {
		line = prefix + branch + icon + " " + n.Name
	}

	// Style current node
	if n == curr {
		line = styles.CurrentNodeStyle.Render(line)
	}

	sb.WriteString(line + "\n")

	// Update prefix for children
	if n.Parent != nil {
		if isLast {
			newPrefix += "   "
		} else {
			newPrefix += "│  "
		}
	}

	if n.IsExpanded {
		for i, child := range n.Children {
			sb.WriteString(BuildNodeTree(child, newPrefix, i == len(n.Children)-1, curr))
		}
	}

	return sb.String()
}
