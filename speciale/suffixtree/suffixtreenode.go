package suffixtree

// SuffixTreeNode represents a node in the suffix tree.
type SuffixTreeNode struct {
	Label    int
	Parent   *SuffixTreeNode
	StartIdx int
	EndIdx   int
	Children map[rune]*SuffixTreeNode
}

// EdgeLength returns the length of the edge represented by the node.
func (node *SuffixTreeNode) EdgeLength() int {
	return node.EndIdx - node.StartIdx + 1
}
