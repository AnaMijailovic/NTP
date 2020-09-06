package model

type Tree struct {
	Root *Node `json:"root"`
}

func (tree *Tree) IsEmpty() bool {
	return tree.Root == nil
}

func (tree *Tree) Depth(node *Node) int {
	if node.Parent == nil {
		return 0
	} else {
		return 1 + tree.Depth(node.Parent)
	}
}
