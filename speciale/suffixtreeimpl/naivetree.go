package suffixtreeimpl

import (
	"speciale/suffixtree"
)

// NaiveTree implements the SuffixTree interface using a naive construction algorithm.
type NaiveSuffixTree struct {
	suffixtree.SuffixTree
}

func (st *NaiveSuffixTree) ConstructSuffixTree() {
	for i := 0; i < len(st.InputString); i++ {
		// Insert all suffixes of inputString into the suffix tree
		st.insertSuffix(i)
	}
	// Add DFS labels
	st.AddDFSLabels()
}

// InsertSuffix inserts the suffix starting at the given index into the suffix tree.
func (st *NaiveSuffixTree) insertSuffix(suffixStartIdx int) {
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

// splitEdge splits the edge of the given child node into two edges on inertion of a suffix.
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

	// Add original child and new child to internal node
	internalNode.Children[rune(st.InputString[originalChild.StartIdx])] = originalChild
	internalNode.Children[rune(st.InputString[newChild.StartIdx])] = newChild
}

// NewNaiveSuffixTree creates a new NaiveSuffixTree instance with the given input string.
func ConstructNaiveSuffixTree(inputString string) suffixtree.SuffixTreeInterface {

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
	var st suffixtree.SuffixTreeInterface = &NaiveSuffixTree{suffixtree.SuffixTree{Root: root, InputString: inputString, Size: 0}}

	// Construct the suffix tree
	st.ConstructSuffixTree()

	// Return the interface value
	return st
}
