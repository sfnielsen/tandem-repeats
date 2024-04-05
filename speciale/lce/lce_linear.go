package lce

import (
	"math"
	"speciale/suffixtree"
)

// LCELinear holds the preprocessed data for the LCE linear time algorithm
type LCELinear struct {
	//forward LCE queries
	suffixTree                  *suffixtree.SuffixTreeInterface
	L                           []int
	E                           []int
	R                           []int
	blocks                      [][]int
	LPrime                      []int
	BPrime                      []int
	LPrimeST                    *sparseTable
	NormalizedBlockSparseTables []*sparseTable

	//backward LCE queries
	//TBD
}
type sparseTable = [][]int

// main function for the LCE linear time preprocessing
func PreProcessLCE(st suffixtree.SuffixTreeInterface) *LCELinear {
	//create L,E,R arrays
	L, E, R := createLERArrays(st)
	// divide L into blocks
	blocks := createLBlocks(L)
	// compute L' and B'
	LPrime, BPrime := computeLPrimeandBPrime(blocks)
	// compute sparse table for L'
	LPrimeSparseTable := computeSparseTable(LPrime)
	// precompute all possible normalized blocks
	NormalizedBlockSparseTables := computeNormalizedBlockSparseTables(blocks)

	return &LCELinear{&st, L, E, R, blocks, LPrime, BPrime, LPrimeSparseTable, NormalizedBlockSparseTables}
}

// create the three arrays: L,E,R by doing an euler tour of the suffix tree
func createLERArrays(st suffixtree.SuffixTreeInterface) ([]int, []int, []int) {
	//euler labels
	nextEulerLabel := 0
	nextEulerStep := 0

	//tables
	L := make([]int, 2*st.GetSize()-1)
	E := make([]int, 2*st.GetSize()-1)
	R := make([]int, st.GetSize())

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

	return L, E, R

}

// create blocks from L array
func createLBlocks(L []int) [][]int {
	//create blocks
	n := len(L)
	blockSize := int(math.Ceil(math.Log2(float64(n)) / 2))
	numBlocks := int(math.Ceil(float64(n) / float64(blockSize))) //number of blocks
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
	table := make([][]int, n)
	for i := 0; i < len(LPrime); i++ {
		table[i] = make([]int, int(math.Ceil(math.Log2(float64(n))))+1)

	}

	//first compute first row (trivial)
	for i, v := range LPrime {
		table[i][0] = v
	}

	//compute the rest of the sparse table
	for j := 1; 1<<j <= n; j++ {
		for i := 0; i+(1<<j) <= n; i++ {
			table[i][j] = intMin(table[i][j-1], table[i+(1<<(j-1))][j-1])
		}
	}

	return &table
}

// function to compute the LCE between two indices i and j
func RMQLookup(i, j int, sparseTable *sparseTable) int {

	k := int(math.Floor((math.Log2(float64(j - i)))))

	range1 := (*sparseTable)[i][k]
	range2 := (*sparseTable)[j-(1<<k)+1][k]

	return intMin(range1, range2)

}

// helper function to find the minimum of two integers
func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// compute all normalized blocks
func computeNormalizedBlockSparseTables(blocks [][]int) []*sparseTable {
	// Create a 2D slice to store all blocks
	totalNormalizedBlocks := 1 << uint(len(blocks[0]))
	normalizedBlocks := make([]*sparseTable, totalNormalizedBlocks)

	// Total number of possible blocks (each element can be + or -)
	for _, block := range blocks {
		blockAsInt := convertBlockToInt(block)
		if normalizedBlocks[blockAsInt] == nil {
			// subtract the first element from the rest of the block
			normalizedBlock := normalizeBlock(block)
			normalizedBlocks[blockAsInt] = computeSparseTable(normalizedBlock)
		}
	}
	return normalizedBlocks
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
