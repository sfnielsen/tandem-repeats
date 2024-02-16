package suffixtree

// SuffixTree defines the interface for a suffix tree.
type SuffixTreeInterface interface {
	// Getters for the fields of the suffix tree
	ConstructSuffixTree()
	GetRoot() *SuffixTreeNode
	GetInputString() string
	GetSize() int

	// AddDFSLabels adds DFS labels to the nodes in the suffix tree.
	AddDFSLabels()

	// AddBiggestChildToNodes adds the biggest child to each node in the suffix tree.
	AddBiggestChildToNodes()
}

type SuffixTree struct {
	Root        *SuffixTreeNode
	InputString string
	Size        int
}

// GetRoot returns the root node of the suffix tree.
func (n *SuffixTree) GetRoot() *SuffixTreeNode {
	return n.Root
}

// GetInputString returns the input string used to construct the suffix tree.
func (n *SuffixTree) GetInputString() string {
	return n.InputString
}

// GetSize returns the size of the suffix tree.
func (n *SuffixTree) GetSize() int {
	return n.Size
}

func (st *SuffixTree) AddBiggestChildToNodes() {
	var dfs func(node *SuffixTreeNode)
	dfs = func(node *SuffixTreeNode) {
		if node.IsLeaf() {
			node.BiggestChild = nil
		} else {
			var longest int
			var biggestFoundChild *SuffixTreeNode = nil
			for _, child := range node.Children {
				if child != nil {
					dfs(child)

					if child.DfsInterval.End-child.DfsInterval.Start+1 > longest {
						longest = child.DfsInterval.End - child.DfsInterval.Start + 1
						biggestFoundChild = child
					}
				}
			}
			node.BiggestChild = biggestFoundChild
		}
	}
	dfs(st.Root)
}

// Adds DFS labels.
// Leaves are assigned a single number, and internal nodes are assigned a range of numbers
// corresponding to the leaves in their subtree.
func (st *SuffixTree) AddDFSLabels() {
	// assign dfs intervals and count up the size of the tree
	// this can easily be done during construction, but this is just a naive implementation
	dfsNumber := 0
	var dfs func(node *SuffixTreeNode) int
	dfs = func(node *SuffixTreeNode) int {
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
		st.incrementSize()
		return 0
	}
	dfs(st.Root)
}

// incrementSize increments the size of the suffix tree by one.
//only called internally so no need for a public function
func (n *SuffixTree) incrementSize() {
	n.Size++
}
