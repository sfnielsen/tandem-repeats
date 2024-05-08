package suffixtree

// SuffixTree defines the interface for a suffix tree.
type SuffixTreeInterface interface {
	// Getters for the fields of the suffix tree
	ConstructSuffixTree()
	GetRoot() *SuffixTreeNode
	GetInputString() string
	GetInternalString() string
	GetSize() int
	GetAlphabetSize() int

	// AddDFSLabels adds DFS labels to the nodes in the suffix tree.
	AddDFSLabelsAndLeafBools()

	// AddStringDepth adds string depth to the nodes in the suffix tree.
	AddStringDepth()

	// AddBiggestChildToNodes adds the biggest child to each node in the suffix tree.
	AddBiggestChildToNodes()

	// Compute leafs of the suffix tree and return them as a slice of nodes
	ComputeLeafsStackMethod() []*SuffixTreeNode
}

type SuffixTree struct {
	Root           *SuffixTreeNode
	InputString    string
	InternalString string
	AlphabetSize   int
	Size           int
}

// GetRoot returns the root node of the suffix tree.
func (n *SuffixTree) GetRoot() *SuffixTreeNode {
	return n.Root
}

// GetInputString returns the input string used to construct the suffix tree.
func (n *SuffixTree) GetInputString() string {
	return n.InputString
}

// GetInternalString returns the internal string used to construct the suffix tree.
func (n *SuffixTree) GetInternalString() string {
	return n.InternalString
}

// GetSize returns the size of the suffix tree.
func (n *SuffixTree) GetSize() int {
	return n.Size
}

// GetAlphabetSize returns the size of the alphabet used to construct the suffix tree.
func (n *SuffixTree) GetAlphabetSize() int {
	return n.AlphabetSize
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
		node.NodeIsLeaf = true
		if node.Label == -1 {
			node.NodeIsLeaf = false
		}
		//if all children are nil, the node is a leaf
		for _, child := range node.Children {
			if child != nil {
				node.NodeIsLeaf = false
				break
			}
		}
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

/*
This function add the following variables to each node:

	- Dfs Labels start and end
	- NodeIsLeaf
	- StringDepth
*/
func (st *SuffixTree) AddDFSLabelsAndLeafBools() {
	stack := Stack{}
	dfsNumber := 0 // Initialize DFS index

	// Push the root node with start flag onto the stack
	stack.Push(&StackItem{Node: st.Root, IsStart: true})

	for len(stack) > 0 {
		item := stack.PopOrNil()
		node := item.Node

		//check if the node is a leaf and save it as a bool
		node.NodeIsLeaf = true
		if node.Label == -1 {
			node.NodeIsLeaf = false
		}
		for _, child := range node.Children {
			if child != nil {
				node.NodeIsLeaf = false
				break
			}
		}

		// set the dfs interval
		if node.IsLeaf() {
			// Set the stringdepth of the node
			node.StringDepth = node.Parent.StringDepth + node.EdgeLength()

			node.DfsInterval.Start = dfsNumber
			node.DfsInterval.End = dfsNumber
			dfsNumber++
			st.incrementSize()
		} else if item.IsStart {
			// Set the stringdepth of the node
			if node != st.Root {
				node.StringDepth = node.Parent.StringDepth + node.EdgeLength()
			} else {
				node.StringDepth = 0
			}

			// We need top push the item again since we will encounter it again when traversing back up
			item.IsStart = false
			stack.Push(item)

			//if NOT leaf node
			node.DfsInterval.Start = dfsNumber

			// add children to the queue in reverse order
			// this is done to ensure that the children are popped in the correct order
			for i := len(node.Children) - 1; i >= 0; i-- {
				if node.Children[i] != nil {
					stack.Push(&StackItem{Node: node.Children[i], IsStart: true})
				}
			}
			st.incrementSize()

		} else {
			node.DfsInterval.End = dfsNumber - 1 // -1 because we have already incremented dfsNumber for the next leaf
		}

	}

}

//add string depth
func (st *SuffixTree) AddStringDepth() {
	var dfs func(node *SuffixTreeNode, depth int)
	dfs = func(node *SuffixTreeNode, depth int) {
		depth = depth + node.EdgeLength()
		node.StringDepth = depth
		for _, child := range node.Children {
			if child != nil {
				dfs(child, depth)
			}
		}
	}
	dfs(st.Root, 0)
}

// incrementSize increments the size of the suffix tree by one.
//only called internally so no need for a public function
func (n *SuffixTree) incrementSize() {
	n.Size++
}

func (n *SuffixTree) ComputeLeafs() []*SuffixTreeNode {
	leafs := make([]*SuffixTreeNode, len(n.GetInputString()))
	var dfs func(node *SuffixTreeNode)
	dfs = func(node *SuffixTreeNode) {
		if node.IsLeaf() {
			leafs[node.Label] = node
		} else {
			for _, child := range node.Children {
				if child != nil {
					dfs(child)
				}
			}
		}
	}
	dfs(n.Root)

	return leafs
}

func (n *SuffixTree) ComputeLeafsStackMethod() []*SuffixTreeNode {
	leafs := make([]*SuffixTreeNode, len(n.GetInputString()))
	stack := TreeStack{n.GetRoot()}

	for len(stack) > 0 {
		node := stack.PopOrNil()
		if node.IsLeaf() {
			leafs[node.Label] = node
		} else {
			for _, child := range node.Children {
				if child != nil {
					stack.Push(child)
				}
			}
		}
	}

	return leafs
}
