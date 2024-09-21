package trie

type Node struct {
	Data     rune
	Children []*Node
}

func (n *Node) New(data rune) *Node {
	return &Node{Data: data}
}

func (n *Node) Insert(t *Trie, data rune) {
	if n.Children == nil {
		n.Children = make([]*Node, 0)
	}

	newNode := n.New(data)
	n.Children = append(n.Children, newNode)
	t.CurrentNode = newNode
}

func (n *Node) InsertString(t *Trie, s string) {
	// Base case: if the string is empty, stop the recursion
	if len(s) == 0 {
		return
	}

	// Insert the first character of the string
	n.Insert(t, rune(s[0]))

	// Recursively insert the remaining substring
	t.CurrentNode.InsertString(t, s[1:])
}


type Trie struct {
	RootNode    *Node
	CurrentNode *Node
}

func (t *Trie) New(s string) {
	t.RootNode = &Node{Data: ' '}
	t.CurrentNode = t.RootNode

	for _, c := range s {
		t.CurrentNode.Insert(t, c)
	}
}

func (t *Trie) LoadString(s string) {
	// start comparing from the Roodnode's children
	// if match is found then look for similarity in the
	// children until the end of the string
	// wherever the similarity diverges insert the remaining
	// part of hte string into the last node's children

	// if no match is found in any of the children then
	// insert a new node in the rootNode's children and load hte string there

	// if a similar string is already loaded then return

	// 1. if a substring exists in any pathway then
	// only start adding from the last common node

	// 2. if no substring exists then add the whole string
	// from the root node

	// 3. if the string is already loaded then return

	nodes := t.FindMatch(s)
	if len(nodes) == len(s) {
		return
	} else if len(nodes) == 0 {
		t.RootNode.InsertString(t, s)
	} else {
		lastNode := nodes[len(nodes)-1]
		lastNode.InsertString(t, s[len(nodes):])

	}

}

// Wrapper function for the first call
func (t *Trie) FindMatch(s string) []*Node {
	// Start recursion from the root node with an empty pathway
	return t.findMatch(t.RootNode, s, []*Node{})
}

func (t *Trie) findMatch(n *Node, s string, pathway []*Node) []*Node {
	// Base condition: if the string is fully consumed, return the pathway
	if len(s) == 0 {
		return pathway
	}

	// Get the current character to match
	c := rune(s[0])

	// Loop through the children of the current node to find the match
	for _, child := range n.Children {
		if child.Data == c {
			// Append the matching child to the pathway
			pathway = append(pathway, child)

			// Recursively call the function with the rest of the string
			result := t.findMatch(child, s[1:], pathway)

			// If the result is non-nil, return it (indicating a successful match)
			if result != nil {
				return result
			}
		}
	}

	// If no match is found, return the pathway up to the last matched node
	return pathway
}

