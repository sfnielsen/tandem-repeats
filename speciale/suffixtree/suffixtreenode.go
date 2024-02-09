package suffixtree

// SuffixTreeNode represents a node in the suffix tree.
type SuffixTreeNode struct {
	//standard fields
	Label    int
	Parent   *SuffixTreeNode
	StartIdx int
	EndIdx   int
	Children [256]*SuffixTreeNode // assume ASCII characters

	// Fields required for tandem repeats O(nlogn) tandem repeats
	Marked      bool // probably not necessary
	DfsInterval DfsInterval
	LeafList  map[int]*SuffixTreeNode
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
	//if all children are nil, the node is a leaf
	for _, child := range node.Children {
		if child != nil {
			return false
		}
	}
	return true
}
