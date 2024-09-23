package trie

import (
	"bytes"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wealdtech/go-merkletree/keccak256"
)

// 1. structs for branch, ext and leaf node
// 2. algo to differentiate between leaf and ext node using prefix
// a util func to pre-process the data keys into compact hex encoding for mpt
// a util for hashing a node (keccak256(rlp.encode(rand struct)))
// levelDB to store KV (address <-> account data)
// methods for nodes
// node interface
// a method to traverse the trie and and return if a leaf node already exists tehre
// a method to cmp two paths and where they diverge
// a method to get the value of a node (a call to the db)
// a method to set the value of a node (a call to the db)
// {calls to the DB are not thread safe}
// ext node:- we traverse the trie whilst maintaining a
// temp [] of struct {nodeHash, key} (fro backtracking)
// and when we find a divergence from a common path (which)
// is being recorded in the temp array
// adn len(temp) > 0 we know where to
// create an ext node
// revised criterion for ext node:-  common prefix + commonlaity in
// key-end of LN
// revised criterion for branch node:-  common prefix & NO commonlaity in
// key-end of LN & divergence count of 1 nibble

type node interface {
	Hash() []byte
	Encode() ([]byte, error)
}

type (
	branchNode struct {
		// 17 children, each references a node via its hash (to make it modular we are treating it as []byte)
		children [17][]byte
		// rlp encoded
		value []byte
	}
	extensionNode struct {
		// the common prefix
		prefix       int
		sharedNibble []byte
		nextNode     []byte
	}
	leafNode struct {
		prefix int
		keyEnd []byte
		value  []byte
	}
	account struct {
		pubKey  ecdsa.PublicKey
		balance int
	}
	trie struct {
		// hash -> node
		nodes    map[string]node
		rootHash []byte
		dB       *leveldb.DB
	}
)

func NewBranchNode() *branchNode {
	return &branchNode{
		children: [17][]byte{},
		value:    nil,
	}
}

func (bn *branchNode) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(bn)
}

func (bn *branchNode) Hash() []byte {
	encoded, _ := bn.Encode()
	return keccak256.New().Hash(encoded)
}

func NewExtNode() *extensionNode {
	return &extensionNode{
		prefix:       -1, // -1 -> uninit
		sharedNibble: nil,
		nextNode:     nil,
	}
}

func (en *extensionNode) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(en)
}

func (en *extensionNode) Hash() []byte {
	encoded, _ := en.Encode()
	return keccak256.New().Hash(encoded)
}

func NewLeafNode() *leafNode {
	return &leafNode{
		prefix: -1, // -1 -> uninit
		keyEnd: nil,
		value:  nil,
	}
}

func (ln *leafNode) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(ln)
}

func (ln *leafNode) Hash() []byte {
	encoded, _ := ln.Encode()
	return keccak256.New().Hash(encoded)
}

func NewTrie() *trie {
	newT := &trie{
		nodes: map[string]node{},
	}

	// by default, the root node is a branch node
	bn := NewBranchNode()
	newT.nodes[string(bn.Hash())] = bn
	newT.rootHash = bn.Hash()

	return newT
}

// func receives account struct, we encode it , hash it
// store it in DB and add it to the trie
// AddAccount(account)
// addToTrie(keyHash, encodedAccount)
// prePrc() -> preProcessing of the keyHash -> nibbles path (depending on the prefix)
// decode() -> prefix is used to distinguish between leaf and ext node

func (ac *account) AddAccountToTrie(t *trie, db *leveldb.DB) {
	// encode the account
	var w bytes.Buffer
	if err := rlp.Encode(&w, ac); err != nil {
		panic("mpt: failed to encode the account ")
	}
	// hash it
	accountHash := keccak256.New().Hash(w.Bytes())
	// store it in DB
	if err := db.Put(accountHash, w.Bytes(), nil); err != nil {
		panic("mpt: failed to store the account in the DB")
	}
	// add it to the trie

}

func (t *trie) addToTrie(keyHash []byte, encodedAccount []byte) {
	// pre-process the keyHash
	nibbles := preProcess(keyHash)
	// decode the nibbles
	t.decode(nibbles, encodedAccount)
}

// traverses the trie
// checks if a leaf node already exists; if yes updates it
// if not, creates a new leaf node and manges hte
// creation of intermediary nodes (branch and ext nodes)
// once each node is processed , its hash is recaluclaetd
// and re-persisted to the db
func (t *trie) decode(nibbles []byte, encodedAccount []byte) {
	n := t.nodes[string(t.rootHash)]
	t.traverse(n.Hash(), nibbles, encodedAccount, nil)
}

// traverses and updates the trie
// TODO: persist nodes to DB (go routines)
func (t *trie) traverse(n []byte, nibbles []byte, encodedAccount []byte, temp []any) {
	currNibble := nibbles[0]
	switch node := t.nodes[string(n)].(type) {
	case *branchNode:
		if node.children[currNibble] != nil {
			// If the child exists, we traverse it
			// Add common prefix to temp
			tempStruct := struct {
				nodeHash []byte
				key      byte
			}{
				node.Hash(),
				currNibble,
			}

			temp = append(temp, tempStruct)
			// Recur to the next node with the remaining nibbles
			t.traverse(node.children[currNibble], nibbles[1:], encodedAccount, temp)
		} else {
			// If the child does not exist, we create a new leaf node
			ln := NewLeafNode()
			ln.value = encodedAccount
			ln.keyEnd = nibbles[1:] // Remaining nibbles for the leaf
			ln.prefix = len(nibbles) - 1

			// Add pointer to the new leaf node in the branch node
			node.children[currNibble] = ln.Hash()

			// Calculate and persist the branch node hash
			// (Assuming you have a method for hashing)
			t.nodes[string(node.Hash())] = node // Persist the updated branch node
		}
	case *extensionNode:
		// Check if the shared nibble matches
		if len(nibbles) > 1 && nibbles[1] == node.sharedNibble[0] {
			// If it matches, we continue traversing through the next node
			t.traverse(node.nextNode, nibbles[1:], encodedAccount, temp)
		} else {
			// If it doesn't match, we need to create a branch node
			branch := NewBranchNode()
			// Create the new extension node with the shared prefix
			newExtNode := NewExtNode()
			newExtNode.prefix = node.prefix
			newExtNode.sharedNibble = node.sharedNibble
			newExtNode.nextNode = node.nextNode

			// Set the new extension node as a child of the branch
			branch.children[currNibble] = newExtNode.Hash()

			// Create the new leaf node and link it
			ln := NewLeafNode()
			ln.value = encodedAccount
			ln.keyEnd = nibbles[1:]
			ln.prefix = len(nibbles) - 1
			branch.children[nibbles[1]] = ln.Hash()

			// Persist the new nodes
			t.nodes[string(branch.Hash())] = branch
			t.nodes[string(newExtNode.Hash())] = newExtNode
			t.nodes[string(ln.Hash())] = ln
		}
	case *leafNode:
		// If we encounter a leaf node
		if bytes.Equal(node.keyEnd, nibbles[1:]) {
			// If the key matches, update the existing leaf node
			node.value = encodedAccount
			// Persist the updated leaf node
			t.nodes[string(node.Hash())] = node
		} else {
			// If the key does not match, we create a new branch node
			branch := NewBranchNode()
			// Move existing leaf node to the branch
			branch.children[node.keyEnd[0]] = node.Hash()
			// Create the new leaf node
			ln := NewLeafNode()
			ln.value = encodedAccount
			ln.keyEnd = nibbles[1:]
			ln.prefix = len(nibbles) - 1
			branch.children[currNibble] = ln.Hash()

			// Persist the new branch and updated leaf node
			t.nodes[string(branch.Hash())] = branch
			t.nodes[string(ln.Hash())] = ln
		}
	}

	// Calculate and persist the root hash
	go t.UpdateRootHash()

}

func (t *trie) UpdateRootHash() {
	if len(t.nodes) == 0 {
		t.rootHash = nil // No nodes in the trie
		return
	}
	// Start from the root node
	rootNode := t.nodes[string(t.rootHash)] // You would need to define how to access the root node
	t.rootHash = t.calculateMerkleRoot(rootNode)
}

func (t *trie) calculateMerkleRoot(node node) []byte {
	switch n := node.(type) {
	case *leafNode:
		// If it's a leaf node, return its hash
		return hash(n.value) // Assuming value is stored in the leaf node
	case *branchNode:
		var childHashes [][]byte
		for _, child := range n.children {
			if child != nil {
				childHashes = append(childHashes, t.calculateMerkleRoot(t.nodes[string(child)]))
			}
		}
		// Combine child hashes to create the branch node hash
		if len(childHashes) > 0 {
			combined := []byte{}
			for _, h := range childHashes {
				combined = append(combined, h...)
			}
			return hash(combined)
		}
	case *extensionNode:
		// Handle extension node (similar to branch nodes)
		return t.calculateMerkleRoot(t.nodes[string(n.nextNode)]) // assuming nextNode points to the child
	}
	return nil // In case of empty or unrecognized node type
}

// PreProcess function to convert key hash to nibbles
func preProcess(keyHash []byte) []byte {
	var nibbles []byte
	for _, b := range keyHash {
		// High nibble
		nibbles = append(nibbles, b>>4)
		// Low nibble
		nibbles = append(nibbles, b&0x0F)
	}
	return nibbles
}

func hash(data []byte) []byte {
	enc, _ := rlp.EncodeToBytes(data)
	return keccak256.New().Hash(enc)
}
