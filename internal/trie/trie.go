package trie

import "fmt"

type NodeType byte

const (
	Branch NodeType = iota
	Leaf
)

type Node interface {
	Copy() Node
}

type Trie struct {
	root Node
}

func (t *Trie) Insert(key, value []byte) error {
	keyNibbles := KeyToNibbles(key)

	if t.root == nil {
		t.root = NewLeafNode(keyNibbles, value)
		return nil
	}

	node := &t.root
	for {

		if branch, ok := (*node).(*BranchNode); ok {
			if len(keyNibbles) == 0 {
				branch.Value = make([]byte, len(value))
				copy(branch.Value, value[:])
				return nil
			}

			matched := nibblesEqualsUntil(branch.Path, keyNibbles)
			if matched < len(branch.Path) {
				return fmt.Errorf("invalid path, expected min match: %d, got: %d",
					len(branch.Path), matched)
			}

			if matched == len(branch.Path) {
				branch.Value = make([]byte, len(value))
				copy(branch.Value, value[:])
				return nil
			}

			forkNibble, remaining := keyNibbles[matched], keyNibbles[matched+1:]
			child := branch.Children[forkNibble]
			if child == nil {
				return fmt.Errorf("%w: nil child at %d", ErrChildrenNotFound, forkNibble)
			}

			keyNibbles = remaining
			node = &child
			continue
		}

		if leaf, ok := (*node).(*LeafNode); ok {
			matched := nibblesEqualsUntil(leaf.Path, keyNibbles)
			if matched == len(keyNibbles) && matched == len(leaf.Path) {
				t.root = NewLeafNode(leaf.Path, value)
				return nil
			}

			if matched == len(leaf.Path) {
				branchNode := NewBranchNodeWithValue(leaf.Path, leaf.Value)
				branchNode.SetBranch(keyNibbles[matched], NewLeafNode(keyNibbles[matched+1:], value))
				t.root = branchNode
				return nil
			}

			if matched == len(keyNibbles) {
				branchNode := NewBranchNodeWithValue(keyNibbles, value)
				branchNode.SetBranch(leaf.Path[matched], NewLeafNode(leaf.Path[matched+1:], leaf.Value))
				t.root = branchNode
				return nil
			}

			branchNodePath := make([]byte, matched)
			copy(branchNodePath, leaf.Path[:matched])
			branchNode := NewBranchNode(branchNodePath)

			newLeafNodePath := make([]byte, len(leaf.Path[matched+1:]))
			copy(newLeafNodePath, leaf.Path[matched+1:])

			branchNode.SetBranch(
				leaf.Path[matched], NewLeafNode(newLeafNodePath, leaf.Value))

			branchNode.SetBranch(
				keyNibbles[matched], NewLeafNode(keyNibbles[matched+1:], value))

			t.root = branchNode
			return nil
		}
	}
}

func NewTrie() *Trie {
	return &Trie{}
}
