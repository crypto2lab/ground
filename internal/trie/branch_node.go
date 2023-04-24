package trie

type BranchNode struct {
	Path     []byte
	Value    []byte
	Children [16]Node
}

func (b *BranchNode) Insert(pathNibbles, value []byte) error {
	return nil
}

func (b *BranchNode) SetBranch(branchNibble byte, node Node) {
	b.Children[branchNibble] = node
}

func (b *BranchNode) Copy() Node {
	childrenCopy := [16]Node{}
	for idx, children := range b.Children {
		childrenCopy[idx] = children.Copy()
	}

	pathCopy := make([]byte, len(b.Path))
	valeCopy := make([]byte, len(b.Value))

	copy(pathCopy, b.Path[:])
	copy(valeCopy, b.Value[:])
	return &BranchNode{
		Path:     pathCopy,
		Value:    valeCopy,
		Children: childrenCopy,
	}
}

func NewBranchNodeWithValue(path, value []byte) *BranchNode {
	return &BranchNode{
		Path:     path,
		Value:    value,
		Children: [16]Node{},
	}
}

func NewBranchNode(path []byte) *BranchNode {
	return &BranchNode{
		Path:     path,
		Value:    nil,
		Children: [16]Node{},
	}
}
