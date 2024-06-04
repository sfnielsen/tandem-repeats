package lce

import (
	"math"
	"math/bits"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"time"
)

// #######################################################################################
// #######################################################################################
// Types
// #######################################################################################
// #######################################################################################

// LCELinear holds the preprocessed data for the LCE linear time algorithm
// Preprocessing time: O(n),Space complexity: O(n), Query time: O(1)
type LCELinear struct {
	//forward LCE queries
	suffixTree                   *suffixtree.SuffixTreeInterface
	E                            []int
	R                            []int
	blocks                       [][]int
	LPrime                       []int
	BPrime                       []int
	LPrimeST                     *sparseTable
	NormalizedBlockSparseTables  []*sparseTable
	blockIdxToNormalizedBlockIdx []int
	Leafs                        []*suffixtree.SuffixTreeNode
	EulerindexToNode             []*suffixtree.SuffixTreeNode

	//backward LCE queries
	//TBD
}

// LCELinearTwoWays holds the preprocessed data for the LCE linear time algorithm
type LCELinearTwoWays struct {
	forward  *LCELinear //forward LCE queries
	backward *LCELinear //backward LCE queries
}

type stTuple struct {
	level int
	index int
}
type sparseTable = [][]stTuple

// #######################################################################################
// #######################################################################################
// Query functions
// #######################################################################################
// #######################################################################################

// LCELookupForward returns the longest common extension of the nodes at index i and j in the forward direction
func (lceTwoWays *LCELinearTwoWays) LCELookupForward(i, j int) int {
	stringLength := len((*lceTwoWays.forward.suffixTree).GetInputString())
	if i >= stringLength || j >= stringLength {
		return 0
	}

	return lceTwoWays.forward.LCELookup(i, j)
}

// LCELookupBackward returns the longest common extension of the nodes at index i and j in the BACKWARD direction
func (lceTwoWays *LCELinearTwoWays) LCELookupBackward(i, j int) int {
	stringLength := len((*lceTwoWays.backward.suffixTree).GetInputString())
	revI, revJ := stringLength-(j+2), stringLength-(i+2)
	if revI < 0 || revJ < 0 || revI >= stringLength || revJ >= stringLength {
		return 0
	}

	return lceTwoWays.backward.LCELookup(revI, revJ)

}

// find the lowest common extension (lowest common ancestor) of i and j
func (lce *LCELinear) LCELookup(i, j int) int {
	//get eulerindexes for i and j entries
	leaf_i := lce.Leafs[i]
	leaf_j := lce.Leafs[j]
	eulerindex_i := lce.R[leaf_i.EulerLabel]
	eulerindex_j := lce.R[leaf_j.EulerLabel]

	//ensure i < j
	if eulerindex_i > eulerindex_j {
		eulerindex_i, eulerindex_j = eulerindex_j, eulerindex_i
	}

	blockSize := len(lce.blocks[0])

	block_i := eulerindex_i / blockSize
	block_j := eulerindex_j / blockSize

	// case 1:   They are on the same block
	if block_i == block_j {
		//find the normalized block idx
		normalizedBlockIdx := lce.blockIdxToNormalizedBlockIdx[block_i]
		//find the normalized block sparse table
		normalizedBlockSparseTable := lce.NormalizedBlockSparseTables[normalizedBlockIdx]
		//find the LCA in the normalized block
		lca := RMQLookup(eulerindex_i%blockSize, eulerindex_j%blockSize, normalizedBlockSparseTable)
		lcaEulerIdx := block_i*blockSize + lca.index
		return lce.EulerindexToNode[lce.E[lcaEulerIdx]].StringDepth
	} else {
		// case 2:   They are on different blocks (i < j   always)
		// find the LCA in the first block
		lcaBlocki := lce.blockIdxToNormalizedBlockIdx[block_i]
		lca1 := RMQLookup(eulerindex_i%blockSize, blockSize-1, lce.NormalizedBlockSparseTables[lcaBlocki])
		lca1.level += lce.blocks[block_i][0]
		// find the LCA in the last block
		lcaBlockj := lce.blockIdxToNormalizedBlockIdx[block_j]
		lca2 := RMQLookup(0, eulerindex_j%blockSize, lce.NormalizedBlockSparseTables[lcaBlockj])
		lca2.level += lce.blocks[block_j][0]
		// find the LCA between the two blocks

		//edgecase when block i and j are adjacent
		if block_j == block_i+1 {
			if lca1.level < lca2.level {
				return lce.EulerindexToNode[lce.E[block_i*blockSize+lca1.index]].StringDepth // lca1 smallest
			}
			return lce.EulerindexToNode[lce.E[block_j*blockSize+lca2.index]].StringDepth // lca2 smallest
		}

		//i and j are not adjacent
		lca3 := RMQLookup(block_i+1, block_j-1, lce.LPrimeST)
		// find the LCA of the three LCA's
		lca1EulerIdx := block_i*blockSize + lca1.index
		lca2EulerIdx := block_j*blockSize + lca2.index
		lca3EulerIdx := lce.BPrime[lca3.index]

		if lca1.level <= lca2.level && lca1.level <= lca3.level {
			return lce.EulerindexToNode[lce.E[lca1EulerIdx]].StringDepth // lca1 smallest
		} else if lca2.level <= lca3.level {
			return lce.EulerindexToNode[lce.E[lca2EulerIdx]].StringDepth // lca2 smallest
		}
		return lce.EulerindexToNode[lce.E[lca3EulerIdx]].StringDepth // lca3 smallest
	}
}

// #######################################################################################
// #######################################################################################
// Preprocessing functions
// #######################################################################################
// #######################################################################################

// main function for linear time LCE preprocessing in both directions
func PreProcessLCEBothDirections(st suffixtree.SuffixTreeInterface) (*LCELinearTwoWays, int) {
	//create forward LCE
	total := 0
	forwardLCE, time_taken := PreProcessLCE(st)
	total += time_taken
	//create backward LCE
	reversedString := reverseStringWithSentinel(st.GetInputString())
	reversedSuffixTree := suffixtreeimpl.ConstructMcCreightSuffixTree(reversedString)
	backwardLCE, time_taken := PreProcessLCE(reversedSuffixTree)
	total += time_taken

	return &LCELinearTwoWays{forward: forwardLCE, backward: backwardLCE}, total

}

// main function for the LCE linear time preprocessing (in one direction)
func PreProcessLCE(st suffixtree.SuffixTreeInterface) (*LCELinear, int) {
	timer := time.Now()
	//create L,E,R arrays
	L, E, R, EulerindexToNode, leafSlice := createLERArraysAndLeafArrayStack(st)
	total := int(time.Since(timer).Milliseconds())
	// divide L into blocks
	blocks := createLBlocks(L)
	// compute L' and B'
	LPrime, BPrime := computeLPrimeandBPrime(blocks)
	// compute sparse table for L'
	LPrimeSparseTable := computeSparseTable(LPrime)
	// precompute all possible normalized blocks
	NormalizedBlockSparseTables, blockIdxToNormalizedBlockIdx := computeNormalizedBlockSparseTables(blocks)

	return &LCELinear{&st, E, R, blocks, LPrime, BPrime, LPrimeSparseTable, NormalizedBlockSparseTables, blockIdxToNormalizedBlockIdx, leafSlice, EulerindexToNode}, total
}

// create the three arrays: L,E,R by doing an euler tour of the suffix tree
func createLERArrays(st suffixtree.SuffixTreeInterface) ([]int, []int, []int, []*suffixtree.SuffixTreeNode) {
	//euler labels
	nextEulerLabel := 0
	nextEulerStep := 0

	//tables
	L := make([]int, 2*st.GetSize()-1)
	E := make([]int, 2*st.GetSize()-1)
	R := make([]int, st.GetSize())

	EulerindexToNode := make([]*suffixtree.SuffixTreeNode, 2*st.GetSize()-1)

	//perform an euler tour of the suffix tree
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {

		//process self
		node.EulerLabel = nextEulerLabel
		R[nextEulerLabel] = nextEulerStep //make mapping from eulerLabel to the eulertour

		L[nextEulerStep] = depth           //note the depth of current eulerStep
		E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel

		nextEulerLabel++
		nextEulerStep++

		EulerindexToNode[node.EulerLabel] = node

		//process children
		for _, child := range node.Children {
			if child != nil {
				dfs(child, depth+1)

				//process self again (each time we return from a subtree)
				L[nextEulerStep] = depth           //note the depth of current eulerStep
				E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel
				nextEulerStep++

			}
		}

	}
	dfs(st.GetRoot(), 0)

	return L, E, R, EulerindexToNode

}

// create the three arrays: L,E,R by doing an euler tour of the suffix tree using a stack
// also create EulerindexToNode mapping and an array of all leaf nodes
func createLERArraysAndLeafArrayStack(st suffixtree.SuffixTreeInterface) ([]int, []int, []int, []*suffixtree.SuffixTreeNode, []*suffixtree.SuffixTreeNode) {
	stack := suffixtree.Stack{}

	// Push the root node with start flag onto the stack
	stack.Push(&suffixtree.StackItem{Node: st.GetRoot(), IsStart: true})

	//euler labels
	nextEulerLabel := 0
	nextEulerStep := 0

	//tables
	L := make([]int, 2*st.GetSize()-1)
	E := make([]int, 2*st.GetSize()-1)
	R := make([]int, st.GetSize())
	EulerindexToNode := make([]*suffixtree.SuffixTreeNode, 2*st.GetSize()-1)
	leafs := make([]*suffixtree.SuffixTreeNode, len(st.GetInputString()))

	// since first iteration is the root and depth needs to start at 0
	depth := -1

	for len(stack) > 0 {
		item := stack.PopOrNil()
		node := item.Node

		if node.IsLeaf() {
			//We're now entering a new depth of the tree
			depth += 1

			//process self
			node.EulerLabel = nextEulerLabel
			R[nextEulerLabel] = nextEulerStep //make mapping from eulerLabel to the eulertour

			L[nextEulerStep] = depth           //note the depth of current eulerStep
			E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel

			nextEulerLabel++
			nextEulerStep++

			EulerindexToNode[node.EulerLabel] = node

			leafs[node.Label] = node //save leaf node

		} else if item.IsStart {
			//We're now entering a new depth of the tree
			depth += 1

			item.IsStart = false

			//process self
			node.EulerLabel = nextEulerLabel
			R[nextEulerLabel] = nextEulerStep //make mapping from eulerLabel to the eulertour

			L[nextEulerStep] = depth           //note the depth of current eulerStep
			E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel

			nextEulerLabel++
			nextEulerStep++

			EulerindexToNode[node.EulerLabel] = node

			//process children
			for i := len(node.Children) - 1; i >= 0; i-- {
				if node.Children[i] != nil {
					stack.Push(item)
					stack.Push(&suffixtree.StackItem{Node: node.Children[i], IsStart: true})
				}
			}

		} else {
			// We're back up again
			depth -= 1

			//process self again (each time we return from a subtree)
			L[nextEulerStep] = depth           //note the depth of current eulerStep
			E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel
			nextEulerStep++
		}

	}
	return L, E, R, EulerindexToNode, leafs
}

// function to calculate the block size and number of blocks
// avoiding casting, ceil,log2 etc..
func calculateBlocksSizeAndAmount(n int) (int, int) {

	log2n := bits.Len(uint(n))

	blockSize := log2n / 2

	numBlocks := (n + blockSize - 1) / blockSize //blocksize-1 to ensure rounding up

	return blockSize, numBlocks
}

// create blocks from L array
func createLBlocks(L []int) [][]int {
	//create blocks
	n := len(L)
	//calculate block size
	blockSize, numBlocks := calculateBlocksSizeAndAmount(n)

	blocks := make([][]int, numBlocks)

	for i := 0; i < len(blocks); i++ {
		if i == len(blocks)-1 {
			//special case for last block
			blocks[i] = append(blocks[i], L[i*blockSize:]...)
			break
		}

		blocks[i] = append(blocks[i], L[i*blockSize:(i+1)*blockSize]...)

	}

	return blocks
}

// compute L' and B' array
func computeLPrimeandBPrime(blocks [][]int) ([]int, []int) {

	LPrime := make([]int, len(blocks))
	BPrime := make([]int, len(blocks))

	//find smallest element in each block
	idxL := 0
	for i, block := range blocks {
		min := math.MaxInt32
		for j, v := range block {
			if v < min {
				min = v
				//save index in B
				BPrime[i] = idxL + j
			}
		}
		LPrime[i] = min
		idxL += len(block)
	}

	return LPrime, BPrime
}

// function to compute sparse tables (ST) using dynamic programming
func computeSparseTable(LPrime []int) *sparseTable {
	//compute sparse table
	n := len(LPrime)
	log2n := bits.Len(uint(n))
	table := make([][]stTuple, n)

	for i := 0; i < len(LPrime); i++ {
		table[i] = make([]stTuple, log2n)
	}

	//first compute first row (trivial)
	for i, v := range LPrime {
		tup := stTuple{level: v, index: i}
		table[i][0] = tup
	}

	// Compute the rest of the sparse table
	for j := 1; 1<<j <= n; j++ {
		// precompute 2^j and 2^(j-1)
		pow2j := 1 << j
		pow2jMinus1 := pow2j >> 1
		//now insert precomputed values
		for i := 0; i+pow2j <= n; i++ {
			table[i][j] = stTupMin(table[i][j-1], table[i+pow2jMinus1][j-1])
		}
	}

	return &table
}

// function to compute the LCE between two indices i and j
func RMQLookup(i, j int, sparseTable *sparseTable) stTuple {
	diff := uint(j - i)
	if diff == 0 {
		return (*sparseTable)[i][0]
	}
	var k int
	if diff == 1 {
		k = 0
	} else {
		k = bits.Len(diff) - 1
	}

	range1 := (*sparseTable)[i][k]
	range2 := (*sparseTable)[j-(1<<k)+1][k]

	return stTupMin(range1, range2)

}

func stTupMin(a, b stTuple) stTuple {
	if a.level < b.level {
		return a
	}
	return b
}

// compute all normalized blocks
func computeNormalizedBlockSparseTables(blocks [][]int) ([]*sparseTable, []int) {
	// Create a 2D slice to store all blocks
	totalNormalizedBlocks := 1 << uint(len(blocks[0]))
	normalizedBlocks := make([]*sparseTable, totalNormalizedBlocks)

	blockIdxToNormalizedBlockIdx := make([]int, len(blocks)) //we need a mapping from block index to normalized block index

	// Total number of possible blocks (each element can be + or -)
	for idx, block := range blocks {
		blockAsInt := convertBlockToInt(block)
		blockIdxToNormalizedBlockIdx[idx] = blockAsInt
		if normalizedBlocks[blockAsInt] == nil {
			// subtract the first element from the rest of the block
			normalizedBlock := normalizeBlock(block)
			normalizedBlocks[blockAsInt] = computeSparseTable(normalizedBlock)
		}
	}
	return normalizedBlocks, blockIdxToNormalizedBlockIdx
}

func normalizeBlock(block []int) []int {
	// subtract the first element from the rest of the block
	normalizedBlock := make([]int, len(block))
	for i, v := range block {
		if i == 0 {
			continue
		}
		normalizedBlock[i] = v - block[0]
	}
	return normalizedBlock

}

// convertBinaryToValues converts a binary number to a sequence of +1 and -1 values
func convertBlockToInt(block []int) int {
	// Iterate over each bit position
	number := 1
	for i, v := range block {
		if i == 0 {
			continue
		}
		number = number << 1
		if block[i-1]-v == -1 {
			number += 1
		}

	}

	return number
}

// Function that reverses a string ending with a sentinel character
// The sentinel is excluded from the reversal, and added back at the end
// after the reversal
func reverseStringWithSentinel(s string) string {
	runes := []rune(s)
	//remove sentinel
	runes = runes[:len(runes)-1]
	//reverse
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	//add sentinel
	runes = append(runes, '$')
	return string(runes)
}
