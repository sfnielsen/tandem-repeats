package suffixtreeimpl

import "speciale/suffixtree"

// NaiveTree implements the SuffixTree interface using a naive construction algorithm.
type NaiveSuffixTree struct {
	Root        *suffixtree.SuffixTreeNode
	InputString string
	Size        int
}

// GetRoot returns the root node of the suffix tree.
func (n NaiveSuffixTree) GetRoot() *suffixtree.SuffixTreeNode {
	return n.Root
}

// GetInputString returns the input string used to construct the suffix tree.
func (n NaiveSuffixTree) GetInputString() string {
	return n.InputString
}

// GetSize returns the size of the suffix tree.
func (n NaiveSuffixTree) GetSize() int {
	return n.Size
}

func (n NaiveSuffixTree) PrintTree() {
	println("Printing tree")
}

// NewNaiveSuffixTree creates a new NaiveSuffixTree instance with the given input string.
func ConstructNaiveSuffixTree(inputString string) suffixtree.SuffixTree {

	//ensure that the input string ends with a $ character
	if inputString[len(inputString)-1] != '$' {
		inputString += "$"
	}

	// Create a root node
	root := &suffixtree.SuffixTreeNode{
		Label:    -1,
		StartIdx: -1,
		EndIdx:   -2,
		//parent is nil by default
		//Children is an array of pointers to SuffixTreeNode which is initialized on creation
	}

	// Create a NaiveSuffixTree
	st := NaiveSuffixTree{
		Root:        root,
		InputString: inputString,
	}

	// Construct the suffix tree
	for i := 0; i < len(inputString); i++ {
		// Insert all suffixes of inputString into the suffix tree
		st.InsertSuffix(i)
	}

	// Add DFS labels
	st.AddDFSLabels(st.Root)

	// Return the interface value
	return st
}

// Adds DFS labels.
// Leaves are assigned a single number, and internal nodes are assigned a range of numbers
// corresponding to the leaves in their subtree.
func (st *NaiveSuffixTree) AddDFSLabels(node *suffixtree.SuffixTreeNode) {
	// assign dfs intervals and count up the size of the tree
	// this can easily be done during construction, but this is just a naive implementation
	dfsNumber := 0
	var dfs func(node *suffixtree.SuffixTreeNode) int
	dfs = func(node *suffixtree.SuffixTreeNode) int {
		// if leaf node
		if node.IsLeaf() {
			node.DfsInterval.Start = dfsNumber
			node.DfsInterval.End = dfsNumber
			dfsNumber++
		} else {
			//if NOT leaf node
			node.DfsInterval.Start = dfsNumber
			for _, child := range node.Children {
				if child != nil {
					dfs(child)
				}
			}
			node.DfsInterval.End = dfsNumber - 1 // -1 because we have already incremented dfsNumber for the next leaf
		}
		st.Size++ // increment size of tree
		return 0
	}
	dfs(st.Root)
}

// InsertSuffix inserts the suffix starting at the given index into the suffix tree.
func (st *NaiveSuffixTree) InsertSuffix(suffixStartIdx int) {
	suffix := st.InputString[suffixStartIdx:]

	// Start at the root
	currentNode := st.Root

	depth := 0
	// infinite loop
	for {

		// Check if the current node has a child with the first character of the suffix
		letter := rune(suffix[depth])
		// Check if there is an edge to follow
		child := currentNode.Children[letter]
		if child != nil {
			// If there is, slow scan through the edge
			// If the edge is longer than our string, we are guaranteed to mismatch on $ character anyways.
			currentEdgeSize := child.EdgeLength()
			for j := 0; j < currentEdgeSize; j++ {
				if suffix[depth+j] != st.InputString[child.StartIdx+j] {
					// If the characters do not match, split the edge and insert the suffix
					st.splitEdge(child, suffixStartIdx+depth, j, len(st.InputString)-1, suffixStartIdx)
					return
				}
			}
			currentNode = child
			depth += currentEdgeSize
			// Check if current node exists
		} else {
			// If it does not, create a new node and insert it as a child of the current node
			// Note that we will always end here if we match completely (as we have $ character)
			newNode := &suffixtree.SuffixTreeNode{
				Label:    suffixStartIdx,
				Parent:   currentNode,
				StartIdx: suffixStartIdx + depth,
				EndIdx:   len(st.InputString) - 1,
			}
			currentNode.Children[rune(suffix[depth])] = newNode
			return
		}
	}
}

func (st *NaiveSuffixTree) splitEdge(originalChild *suffixtree.SuffixTreeNode, startIdx, splitIdx, endIdx, suffixOffset int) {
	// Create a new child
	newChild := &suffixtree.SuffixTreeNode{
		Label:    suffixOffset,
		Parent:   nil,
		StartIdx: startIdx + splitIdx,
		EndIdx:   endIdx,
	}

	// Create a new internal node
	internalNode := &suffixtree.SuffixTreeNode{
		Label:    originalChild.Label,
		Parent:   originalChild.Parent,
		StartIdx: originalChild.StartIdx,
		EndIdx:   originalChild.StartIdx + splitIdx - 1,
	}

	// Add internal node as parent to new child
	newChild.Parent = internalNode

	// Update parent by removing original child and adding internal node
	// This is done by overwriting the original child with the internal node
	originalChild.Parent.Children[rune(st.InputString[internalNode.StartIdx])] = internalNode

	// Update original child
	originalChild.Parent = internalNode
	originalChild.StartIdx += splitIdx

	// Check if they have the same starting character
	if st.InputString[originalChild.StartIdx] == st.InputString[newChild.StartIdx] {
		println("problems :D")
	}

	// Add original child and new child to internal node
	internalNode.Children[rune(st.InputString[originalChild.StartIdx])] = originalChild
	internalNode.Children[rune(st.InputString[newChild.StartIdx])] = newChild
}
