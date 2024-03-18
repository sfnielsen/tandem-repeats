package tandemrepeat

import (
	"fmt"
	"speciale/suffixtree"
)

// Algorithm 1
// Combines algorithm 1a and 1b to find tandem repeats
func Algorithm1(tree suffixtree.SuffixTreeInterface) []TandemRepeat {
	leftMostCoveringRepeats := make([]TandemRepeat, 0)
	s := tree.GetInputString()

	// Compute the blocks and Z-values
	li, si := LZDecomposition(tree)
	blocks := CreateLZBlocks(li, si)

	println("blocks:")
	fmt.Println(blocks)

	// add idx to dfs table
	idxToDfsTable := getIdxtoDfsTable(tree)

	// Process block B for tandem repeats that satisfy condition 1
	tr1 := Algorithm1a(s, blocks, idxToDfsTable)
	leftMostCoveringRepeats = append(leftMostCoveringRepeats, tr1...)

	// Process block B for tandem repeats that satisfy condition 2
	tr2 := Algorithm1b(s, blocks, idxToDfsTable)
	leftMostCoveringRepeats = append(leftMostCoveringRepeats, tr2...)

	for _, k := range tr1 {
		fmt.Println(k)
	}
	println()
	for _, k := range tr2 {
		fmt.Println(k)
	}

	return leftMostCoveringRepeats
}

func Algorithm1a(s string, blocks []int, idxToDfsTable []int) []TandemRepeat {
	tr := make([]TandemRepeat, 0)

	for i := 0; i < len(blocks); i++ {
		h := blocks[i]
		h1 := len(s)
		if i < len(blocks)-1 {
			h1 = blocks[i+1]
		}

		for k := 1; k <= h1-h; k++ {
			q := h1 - k
			k1 := findLCEForwardSlow(s, h1, q)
			k2 := findLCEBackwardSlow(s, h1-1, q-1)
			start := intMax(q-k2, q-k+1)
			if k1+k2 >= k && k1 > 0 {
				tr = append(tr, TandemRepeat{start, 2 * k, 2})
			}

		}

	}
	return tr
}

func Algorithm1b(s string, blocks []int, idxToDfsTable []int) []TandemRepeat {
	tr := make([]TandemRepeat, 0)

	for i := 0; i < len(blocks); i++ {
		h := blocks[i]
		h1 := len(s)
		h2 := len(s)
		if i < len(blocks)-1 {
			h1 = blocks[i+1]
		}
		if i < len(blocks)-2 {
			h2 = blocks[i+2]
		}

		for k := 1; k <= h2-h; k++ {
			q := h + k
			k1 := findLCEForwardSlow(s, h, q)
			k2 := findLCEBackwardSlow(s, h-1, q-1)

			start := intMax(h-k2, h-k+1)

			if k == 3 && i == 3 {
				fmt.Println(h, h1, h2, k1, k2, start, q, "hugo", h-k2, h-k+1)
			}

			if k1+k2 >= k && k1 > 0 && start+k < h1 && k2 > 0 {
				tr = append(tr, TandemRepeat{start, 2 * k, 2})
			}
		}
	}

	return tr
}

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

// find the longest common extension of two suffixes that starts at i and j
func findLCEForwardSlow(s string, i, j int) int {

	lce := 0

	//match letters until we have a mismatch
	for i < len(s) && j < len(s) {
		if s[i] != s[j] {
			return lce
		} else {
			i++
			j++
			lce++
		}
	}
	return lce

}

// find the longest common extension of two suffixes that ends at i and j
func findLCEBackwardSlow(s string, i, j int) int {
	lce := 0

	//match letters until we have a mismatch
	for i >= 0 && j >= 0 {
		if s[i] != s[j] {
			return lce
		} else {
			i--
			j--
			lce++
		}
	}
	return lce
}
