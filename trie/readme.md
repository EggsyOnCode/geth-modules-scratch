Explanation about Decode func of MPT:

    Branch Node:
        If a child exists for the current nibble, it traverses that child.
        If not, it creates a new leaf node, assigns it the encoded account, and updates the branch node.

    Extension Node:
        If the nibble matches the shared nibble, it continues to traverse down the trie.
        If it doesn’t match, it creates a new branch node and an extension node to maintain the prefix.

    Leaf Node:
        If the key matches, it updates the existing leaf node with the new encoded account.
        If the key doesn’t match, it creates a new branch node to handle the conflict and updates the existing leaf node.