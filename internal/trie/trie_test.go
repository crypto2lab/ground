package trie

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrieInsertion(t *testing.T) {
	key := []byte("account::eclesio::TNT")
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, 1000000000000)

	trie := NewTrie()
	require.Nil(t, trie.root)

	err := trie.Insert(key, value)
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	expectedTrie := &Trie{
		root: &LeafNode{
			Path:  KeyToNibbles(key),
			Value: value,
		},
	}
	require.Equal(t, trie, expectedTrie)
}

func TestTrieInsertionAndUpdate_KeepAsLeafNode(t *testing.T) {
	key := []byte("account::eclesio::TNT")
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, 1000000000000)

	trie := NewTrie()
	require.Nil(t, trie.root)

	err := trie.Insert(key, value)
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	expectedTrie := &Trie{
		root: &LeafNode{
			Path:  KeyToNibbles(key),
			Value: value,
		},
	}
	require.Equal(t, trie, expectedTrie)

	sameKey := []byte("account::eclesio::TNT")
	differentValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, 10)

	require.NotNil(t, trie.root)
	err = trie.Insert(sameKey, differentValue)
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	expectedTrie = &Trie{
		root: &LeafNode{
			Path:  KeyToNibbles(sameKey),
			Value: differentValue,
		},
	}
	require.Equal(t, trie, expectedTrie)
}

func TestTrieInsertionAndUpdate_ChangeToBranchNode(t *testing.T) {
	key := []byte("account::eclesio::TNT")
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, 1000000000000)

	trie := NewTrie()
	require.Nil(t, trie.root)

	err := trie.Insert(key, value)
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	expectedTrie := &Trie{
		root: &LeafNode{
			Path:  KeyToNibbles(key),
			Value: value,
		},
	}
	require.Equal(t, trie, expectedTrie)

	anotherKey := []byte("account::flavia::TNT")
	differentValue := make([]byte, 8)
	binary.LittleEndian.PutUint64(differentValue, 10)

	require.NotNil(t, trie.root)
	err = trie.Insert(anotherKey, differentValue)
	require.NoError(t, err)
	require.IsType(t, &BranchNode{}, trie.root)

	keyNibbles := KeyToNibbles(key)
	anotherKeyNibbles := KeyToNibbles(anotherKey)

	matched := nibblesEqualsUntil(keyNibbles, anotherKeyNibbles)

	expectedTrie = &Trie{
		root: &BranchNode{
			Path:  keyNibbles[:matched],
			Value: nil,
			Children: [16]Node{
				5: &LeafNode{
					Path:  keyNibbles[matched+1:],
					Value: value,
				},
				6: &LeafNode{
					Path:  anotherKeyNibbles[matched+1:],
					Value: differentValue,
				},
			},
		},
	}
	require.Equal(t, trie, expectedTrie)
}
