package trie

import (
	"github.com/crypto2lab/ground/lib/comparable"
)

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

	switch node := t.root.(type) {
	case *BranchNode:
		return node.Insert(keyNibbles, value)

	case *LeafNode:
		matched := nibblesEqualsUntil(node.Path, keyNibbles)
		if matched == len(keyNibbles) && matched == len(node.Path) {
			t.root = NewLeafNode(node.Path, value)
			return nil
		}

		if matched == len(node.Path) {
			branchNode := NewBranchNodeWithValue(node.Path, node.Value)
			branchNode.SetBranch(keyNibbles[matched], NewLeafNode(keyNibbles[matched+1:], value))
			t.root = branchNode
			return nil
		}

		if matched == len(keyNibbles) {
			branchNode := NewBranchNodeWithValue(keyNibbles, value)
			branchNode.SetBranch(node.Path[matched], NewLeafNode(node.Path[matched+1:], node.Value))
			t.root = branchNode
			return nil
		}

		branchNodePath := make([]byte, matched)
		copy(branchNodePath, node.Path[:matched])
		branchNode := NewBranchNode(branchNodePath)

		newLeafNodePath := make([]byte, len(node.Path[matched+1:]))
		copy(newLeafNodePath, node.Path[matched+1:])

		branchNode.SetBranch(
			node.Path[matched], NewLeafNode(newLeafNodePath, node.Value))

		branchNode.SetBranch(
			keyNibbles[matched], NewLeafNode(keyNibbles[matched+1:], value))

		t.root = branchNode
	}

	return nil
}

func nibblesEqualsUntil(nibblesA, nibblesB []byte) int {
	maxIterations := comparable.Min(len(nibblesA), len(nibblesB))
	for idx := 0; idx < maxIterations; idx++ {
		if nibblesA[idx] == nibblesB[idx] {
			continue
		}

		return idx
	}

	return maxIterations
}

func NewTrie() *Trie {
	return &Trie{}
}
