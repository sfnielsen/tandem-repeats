package suffixtree

// SuffixTreeNode represents a node in the suffix tree.
type SuffixTreeNode struct {
	//standard fields
	Label    int // for leafs it is the index of the suffix, for internal nodes it is the smallest index of the suffix in the subtree
	Parent   *SuffixTreeNode
	StartIdx int                   // start index of the substring in the input string (inclusive)
	EndIdx   int                   // end index of the substring in the input string (inclusive)
	Children [1400]*SuffixTreeNode // assume ASCII characters

	// Fields required for McCreight's algorithm
	SuffixLink *SuffixTreeNode

	// Fields required for algorithm: O(nlogn) tandem repeats
	DfsInterval  DfsInterval
	BiggestChild *SuffixTreeNode

	// Fields required for algorithm: O(n) tandem repeats
	TandemRepeatDeco         []int
	TandemRepeatDecoComplete map[int]bool
	StringDepth              int

	L_k           *SuffixTreeNode
	TrailingZeros int
	A_v           []bool
	A_v_int       int

	// Fields required for algorithm: LCE linear
	EulerLabel int
}
type DfsInterval struct {
	Start int
	End   int
}

// EdgeLength returns the length of the edge represented by the node.
func (node *SuffixTreeNode) EdgeLength() int {
	return node.EndIdx - node.StartIdx + 1
}

// check if the node is a leaf
func (node *SuffixTreeNode) IsLeaf() bool {
	// if its the root
	if node.Label == -1 {
		return false
	}
	//if all children are nil, the node is a leaf
	for _, child := range node.Children {
		if child != nil {
			return false
		}
	}
	return true
}
