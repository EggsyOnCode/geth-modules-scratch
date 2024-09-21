package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieInsertSingleChar(t *testing.T) {
	trie := &Trie{}
	trie.New("a")

	// Assert that the root node exists and has one child with the correct data
	assert.NotNil(t, trie.RootNode, "Root node should not be nil")
	assert.Equal(t, ' ', trie.RootNode.Data, "Root node data should be a space")

	assert.Len(t, trie.RootNode.Children, 1, "Root node should have 1 child")
	assert.Equal(t, 'a', trie.RootNode.Children[0].Data, "First child node should contain 'a'")
}

func TestTrieInsertString(t *testing.T) {
	trie := &Trie{}
	trie.New("hello")

	// root node is NOT included in the pathway as of yet

	// Check if the Trie root node has the correct structure after insertion
	expected := "hello"
	fmt.Printf("len of hello is %d\n", len(expected))
	nodes := trie.FindMatch(expected)
	fmt.Printf("len of nodes is %d\n", len(nodes))


	assert.Equal(t, len(nodes), len(expected), "The number of nodes should be equal to the length of the string")

	for i, node := range nodes { 
		assert.Equal(t, rune(expected[i]), node.Data, "Node data should match the string character")
	}
}

func TestLoadStringAndMatch(t *testing.T) {
	trie := &Trie{}
	trie.New("hello")
	trie.LoadString("world")

	// Check if both strings 'hello' and 'world' are inserted correctly
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5}, 
		{"world", 5}, 
		{"hell", 4},  
		{"wor", 3},   
	}

	for _, test := range tests {
		nodes := trie.FindMatch(test.input)
		assert.Equal(t, len(nodes), test.expected, "Expected %d nodes for string %s", test.expected, test.input)
	}
}

func TestPartialMatchAndInsert(t *testing.T) {
	trie := &Trie{}
	trie.New("hello")
	trie.LoadString("helium")

	// Check if both 'hello' and 'helium' share a common prefix and diverge properly
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},  
		{"helium", 6}, 
	}

	for _, test := range tests {
		nodes := trie.FindMatch(test.input)
		assert.Equal(t, len(nodes), test.expected, "Expected %d nodes for string %s", test.expected, test.input)
	}
}

func TestTrieInsertEmptyString(t *testing.T) {
	trie := &Trie{}
	trie.New("")

	// Assert that the root node exists and has no children since the string is empty
	assert.NotNil(t, trie.RootNode, "Root node should not be nil")
	assert.Equal(t, ' ', trie.RootNode.Data, "Root node data should be a space")
	assert.Empty(t, trie.RootNode.Children, "Root node should have no children")
}
