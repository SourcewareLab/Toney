package filetree

func FlattenVisibleTree(root *Node) []*Node {
	if root == nil {
		return []*Node{}
	}

	var result []*Node
	var walk func(n *Node) // Removed unusued 'depth int', 0 references or implementations found for this across the project.

	walk = func(n *Node) {
		result = append(result, n)
		if n.IsExpanded {
			for _, child := range n.Children {
				walk(child)
			}
		}
	}
	walk(root)
	return result
}
