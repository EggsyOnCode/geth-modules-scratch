package trie

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



