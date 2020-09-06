package model

type Node struct {
	Element interface {} `json:"element"`
	Parent *Node         `json:"parent"`
	Children []*Node     `json:"children"`
}

func (node *Node) IsRoot() bool {
	return node.Parent == nil
}

func (node *Node) IsLeaf() bool {
	return len(node.Children) == 0
}