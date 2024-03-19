package tandemrepeat

import (
	"fmt"
	"speciale/suffixtree"
)

//Phase 2
func Algorithm2() {

}

// Algorithm 1
// Combines algorithm 1a and 1b to find tandem repeats
func Algorithm1(tree suffixtree.SuffixTreeInterface) [][]TandemRepeat {
	s := tree.GetInputString()
	leftMostCoveringRepeats := make([][]TandemRepeat, len(s))

	// intialize nested slice
	for i := range leftMostCoveringRepeats {
		leftMostCoveringRepeats[i] = make([]TandemRepeat, 0)
	}

	// Compute the blocks and Z-values
	li, si := LZDecomposition(tree)
	blocks := CreateLZBlocks(li, si)

	println("blocks:")
	fmt.Println(blocks)

	// add idx to dfs table
	idxToDfsTable := getIdxtoDfsTable(tree)

	// Process block B for tandem repeats that satisfy condition 1
	Algorithm1a(s, blocks, idxToDfsTable, &leftMostCoveringRepeats)

	// Process block B for tandem repeats that satisfy condition 2
	Algorithm1b(s, blocks, idxToDfsTable, &leftMostCoveringRepeats)

	for _, k := range leftMostCoveringRepeats {
		fmt.Println(k)
	}
	return leftMostCoveringRepeats
}

func Algorithm1a(s string, blocks []int, idxToDfsTable []int, leftMostCoveringRepeats *[][]TandemRepeat) {
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
				addToLeftMostCoveringRepeats(leftMostCoveringRepeats, start, k)
			}
		}
	}
}

func Algorithm1b(s string, blocks []int, idxToDfsTable []int, leftMostCoveringRepeats *[][]TandemRepeat) {
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
			if k1+k2 >= k && k1 > 0 && start+k < h1 && k2 > 0 {
				addToLeftMostCoveringRepeats(leftMostCoveringRepeats, start, k)
			}
		}
	}

}

// get all tandem rpeeats by RIGHT rotating on the branching repeats
func rightRotation(allBranchingRepeats []TandemRepeat, st suffixtree.SuffixTreeInterface) []TandemRepeat {
	var allTandemRepeats = make([]TandemRepeat, 0)

	for _, k := range allBranchingRepeats {
		// add tandem repeat until length is 0
		i := 0
		// left rotate until we no longer have a tandem repeat (or we reach the start of the string)
		for k.Start+i+2*(k.Length) < len(st.GetInputString()) {
			if st.GetInputString()[k.Start+i] == st.GetInputString()[(k.Start+i)+2*(k.Length)] {
				i += 1
				allTandemRepeats = append(allTandemRepeats, TandemRepeat{k.Start + i, k.Length, 2})
			} else {
				break
			}

		}

	}
	allTandemRepeats = append(allTandemRepeats, allBranchingRepeats...)
	return allTandemRepeats
}

// add tandemrepeat to leftMostCoveringRepeats at index start if the last inserted tandem repeat at index start is not of same length
func addToLeftMostCoveringRepeats(leftMostCoveringRepeats *[][]TandemRepeat, start int, k int) {
	if len((*leftMostCoveringRepeats)[start]) == 0 {
		(*leftMostCoveringRepeats)[start] = append((*leftMostCoveringRepeats)[start], TandemRepeat{start, k, 2})
	} else if (*leftMostCoveringRepeats)[start][len((*leftMostCoveringRepeats)[start])-1].Length != k {
		(*leftMostCoveringRepeats)[start] = append((*leftMostCoveringRepeats)[start], TandemRepeat{start, k, 2})
	}
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
