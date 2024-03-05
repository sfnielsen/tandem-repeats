package tandemrepeat

import (
	"speciale/suffixtree"
)

// Phase 1, pt 1
// Compute the LZ decomposition of a string using a suffix tree
func LZDecomposition(tree suffixtree.SuffixTreeInterface) ([]int, []int) {

	// Initialize arrays to store the lengths of blocks and their starting positions
	n := len(tree.GetInputString())
	li := make([]int, n)
	si := make([]int, n)

	// Perform a depth-first traversal of the suffix tree to compute the LZ decomposition
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {

		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			if node.Label < child.Label {
				li[child.Label] = depth
				si[child.Label] = node.Label
			}

			dfs(child, depth+child.EdgeLength())

		}
	}

	// Perform depth-first traversal starting from the root of the suffix tree
	dfs(tree.GetRoot(), 0)

	return li, si
}

// create blocks from the LZ decomposition
func CreateLZBlocks(li []int, si []int) []int {
	n := len(li)

	//first block
	blocks := []int{0}
	iB := 0

	//recursively define remaining blocks
	for iB < n-1 {

		//next start
		iB += intMax(1, li[iB])
		blocks = append(blocks, iB)

	}
	return blocks
}

// return the maximum of two integers
func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
