package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TESTS pending (WIP!!!)

func TestNewBranchNode(t *testing.T) {
	bn := NewBranchNode()
	assert.NotNil(t, bn)
	assert.Equal(t, 17, len(bn.children))
	assert.Nil(t, bn.value)
}

func TestEncodeBranchNode(t *testing.T) {
	bn := NewBranchNode()
	encoded, err := bn.Encode()
	assert.NoError(t, err)
	assert.NotNil(t, encoded)
}

func TestBranchNodeHash(t *testing.T) {
	bn := NewBranchNode()
	hash := bn.Hash()
	assert.NotNil(t, hash)
	assert.Equal(t, 32, len(hash)) // Assuming keccak256 returns a 32-byte hash
}

// func TestAddAccountToTrie(t *testing.T) {
// 	db, err := leveldb.OpenFile("testdb", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	trie := NewTrie(db)
// 	account := &account{pubKey: ecdsa.PublicKey{}, balance: 100}
// 	account.AddAccountToTrie(trie, db)

// 	// Retrieve account from DB to verify it was stored correctly
// 	accountHash := keccak256.New().Hash(rlp.EncodeToBytes(account)) // Assuming a method to encode account exists
// 	data, err := db.Get(accountHash, nil)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, data)
// }

// func TestTraverseLeafNode(t *testing.T) {
// 	// Setup
// 	db, err := leveldb.OpenFile("testdb", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	trie := NewTrie(db)
// 	// Add a sample account
// 	account := &account{pubKey: ecdsa.PublicKey{}, balance: 100}
// 	account.AddAccountToTrie(trie, db)

// 	// Test traversal
// 	nibbles := []byte{...} // Example nibbles corresponding to the account key
// 	trie.decode(nibbles, nil) // Assuming this would add the account
// 	assert.NotNil(t, trie.nodes)
// }
