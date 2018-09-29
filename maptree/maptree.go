package maptree

type MapTreeNode struct {
	INode  uint64
	Parent *MapTreeNode
}

func NewMapTreeNodeWithParent(inode uint64, parent *MapTreeNode) *MapTreeNode {
	res := new(MapTreeNode)
	res.INode = inode
	res.Parent = parent
	return res
}

func NewMapTreeNode(inode uint64) *MapTreeNode {
	res := new(MapTreeNode)
	res.INode = inode
	return res
}

func (a *MapTreeNode) ToRoot() *MapTreeNode {
	for a.Parent != nil {
		a = a.Parent
	}
	return a
}
