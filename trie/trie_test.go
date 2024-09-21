package trie

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTrieInsert(t *testing.T) {
	trie := Trie{}
	trie.New("hello")
	trie.LoadString("world")

	fmt.Printf("trie is %v\n", trie)
	assert.NotNil(t, trie.RootNode)
	assert.NotNil(t, trie.CurrentNode)
}


