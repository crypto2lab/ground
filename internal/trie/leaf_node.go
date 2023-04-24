package trie

type LeafNode struct {
	Path  []byte
	Value []byte
}

func (l *LeafNode) Copy() Node {
	pathCopy := make([]byte, len(l.Path))
	valeCopy := make([]byte, len(l.Value))

	copy(pathCopy, l.Path[:])
	copy(valeCopy, l.Value[:])
	return &LeafNode{
		Path:  pathCopy,
		Value: valeCopy,
	}
}

func NewLeafNode(path, value []byte) *LeafNode {
	leafNodeValue := make([]byte, len(value))
	copy(leafNodeValue, value[:])

	return &LeafNode{
		Value: leafNodeValue,
		Path:  path,
	}
}
