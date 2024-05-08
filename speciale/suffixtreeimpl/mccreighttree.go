package suffixtreeimpl

import (
	"speciale/suffixtree"
)

// McCreightSuffixTree implements the SuffixTree interface using the O(n) McCreight construction algorithm.
type McCreightSuffixTree struct {
	suffixtree.SuffixTree
}

func (st *McCreightSuffixTree) ConstructSuffixTree() {

	//insert first suffix manually
	//just the entire length on root
	newNode := &suffixtree.SuffixTreeNode{
		Label:    0,
		Parent:   st.Root,
		StartIdx: 0,
		EndIdx:   len(st.InputString) - 1,
	}
	st.Root.Children[rune(st.InputString[0])] = newNode
	previousI := newNode

	// now we insert the rest of the suffixes iteratively
	for i := 1; i < len(st.InputString); i++ {
		// Insert all suffixes of inputString into the suffix tree
		previousI = st.insertSuffix(i, previousI)

	}
	// Add DFS labels
	st.AddDFSLabelsAndLeafBools()
}

// InsertSuffix inserts the suffix starting at the given index into the suffix tree.
func (st *McCreightSuffixTree) insertSuffix(suffixStartIdx int, previousI *suffixtree.SuffixTreeNode) *suffixtree.SuffixTreeNode {
	suffix := st.InputString[suffixStartIdx:]

	// variables to keep track of the current node and the depth of the fast scan
	var currentNode *suffixtree.SuffixTreeNode
	var depth, fsEndDepth int

	//Case A: if h(i-1) is root, then s(h(i-1)) is root
	if previousI.Parent == st.Root {
		currentNode = st.Root
		depth = 0
		fsEndDepth = 0

	} else if previousI.Parent.Parent == st.Root && previousI.Parent.EdgeLength() == 1 {
		//Case B: if h(i-1) has length 1 meaning that s(h(i-1)) is root, then s(p(h(i-1))) is also root
		currentNode = st.Root
		previousI.Parent.SuffixLink = st.Root
		depth = 0
		fsEndDepth = 0

	} else if previousI.Parent.Parent == st.Root {
		//Case C: if h(i-1) has length > 1, then we can find s(h(i-1)) by following the suffix link of p(h(i-1))
		currentNode = previousI.Parent.Parent.SuffixLink
		depth = 0
		fsEndDepth = previousI.Parent.EdgeLength() - 1
	} else {
		//Case D: general case, we can find s(h(i-1)) by following the suffix link of h(i-1) and fastscan
		// suffix link is not the root
		currentNode = previousI.Parent.Parent.SuffixLink

		//we need to know the depth from going up from previousI to parent,parent and then jumping the suffixlink
		depth = (len(st.InputString) - suffixStartIdx) - previousI.EdgeLength() - previousI.Parent.EdgeLength()
		fsEndDepth = depth + previousI.Parent.EdgeLength()
	}

	//we do not need to fastscan if we no substring to fastscan through
	if depth == fsEndDepth {
		goto slowscan
	}
	// fast scan
	for {

		currentNode = currentNode.Children[rune(suffix[depth])]

		// Case A: the edge is just as long as the remaining fast scan length
		if currentNode.EdgeLength() == fsEndDepth-depth {
			// we update suffix link of h(i-1) to be currentNode (which is s(h(i-1)))
			previousI.Parent.SuffixLink = currentNode
			depth = fsEndDepth

			goto slowscan //we end in a node, and can continue slow scan

		} else if
		// Case B: the edge is longer than the remaining fast scan length
		currentNode.EdgeLength() > fsEndDepth-depth {
			// we split the edge and insert the suffix
			newI := st.splitEdge(currentNode, suffixStartIdx+depth, fsEndDepth-depth, len(st.InputString)-1, suffixStartIdx)
			//new internal node created is s(h(i-1)), so we update suffix link of h(i-1) to be newChild.Parent
			previousI.Parent.SuffixLink = newI.Parent
			return newI
		} else {
			// Case C: the edge is shorter than the remaining fast scan length
			// we continue fast scan
			depth += currentNode.EdgeLength()
		}
	}

	// if we did not split edge and return, we need to continue slow scanning

slowscan:

	// slow scan
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
					return st.splitEdge(child, suffixStartIdx+depth, j, len(st.InputString)-1, suffixStartIdx)
				}
			}
			currentNode = child
			depth += currentEdgeSize

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
			return newNode
		}
	}
}

func (st *McCreightSuffixTree) createNewChildAndInternalNode(originalChild *suffixtree.SuffixTreeNode, startIdx, splitIdx, endIdx, suffixOffset int) (*suffixtree.SuffixTreeNode, *suffixtree.SuffixTreeNode) {
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

	return newChild, internalNode
}

// splitEdge splits the edge of the given child node into two edges on inertion of a suffix. This McCreight variant returns the new child node.
func (st *McCreightSuffixTree) splitEdge(originalChild *suffixtree.SuffixTreeNode, startIdx, splitIdx, endIdx, suffixOffset int) *suffixtree.SuffixTreeNode {
	// Create a new child
	newChild, internalNode := st.createNewChildAndInternalNode(originalChild, startIdx, splitIdx, endIdx, suffixOffset)

	// Update parent by removing original child and adding internal node
	// This is done by overwriting the original child with the internal node
	originalChild.Parent.Children[rune(st.InputString[internalNode.StartIdx])] = internalNode

	// Update original child
	originalChild.Parent = internalNode
	originalChild.StartIdx += splitIdx

	// Add original child and new child to internal node
	internalNode.Children[rune(st.InputString[originalChild.StartIdx])] = originalChild
	internalNode.Children[rune(st.InputString[newChild.StartIdx])] = newChild

	return newChild
}

// ConstructMcCreightSuffixTree creates a new mcCreight Suffixtree instance with the given input string.
func ConstructMcCreightSuffixTree(inputString string) suffixtree.SuffixTreeInterface {

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
	root.SuffixLink = root

	// Create a McCreightSuffixTree instance
	var st suffixtree.SuffixTreeInterface = &McCreightSuffixTree{suffixtree.SuffixTree{Root: root, InputString: inputString, Size: 0}}

	// Construct the suffix tree
	st.ConstructSuffixTree()

	// Return the interface value
	return st
}
