package filetree

func FlattenVisibleTree(root *Node) []*Node {
	var result []*Node
	var walk func(n *Node, depth int)

	walk = func(n *Node, depth int) {
		result = append(result, n)
		if n.IsExpanded {
			for _, child := range n.Children {
				walk(child, depth+1)
			}
		}
	}
	walk(root, 0)
	return result
}
