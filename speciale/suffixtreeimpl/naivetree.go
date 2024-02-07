package suffixtreeimpl

import "speciale/suffixtree"

// NaiveTree implements the SuffixTree interface using a naive construction algorithm.
type NaiveTree struct {
	// Implement the data structure for the naive suffix tree
}

func (nt *NaiveTree) EdgeLength(node *suffixtree.SuffixTreeNode) int {
	// Implement naive suffix tree construction algorithm
	return node.Label
}

func (nt *NaiveTree) SplitEdge(pattern string) {
	// Implement splitedge operation using the naive suffix tree
}
