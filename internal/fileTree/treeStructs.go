package filetree

type Node struct {
	Name        string
	Depth       int
	Parent      *Node
	Children    []*Node
	IsDirectory bool
	IsExpanded  bool
}
