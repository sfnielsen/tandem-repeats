package lce

import (
	"math"
	"speciale/suffixtree"
)

// LCELinear holds the preprocessed data for the LCE linear time algorithm
type LCELinear struct {
	//forward LCE queries
	suffixTree *suffixtree.SuffixTreeInterface
	L          []int
	E          []int
	R          []int
	blocks     [][]int
	LPrime     []int
	BPrime     []int
	LPrimeST   [][]int

	//backward LCE queries
	//TBD
}

// main function for the LCE linear time preprocessing
func PreProcessLCE(st suffixtree.SuffixTreeInterface) *LCELinear {
	//create L,E,R arrays
	L, E, R := createLERArrays(st)

	blocks := createLBlocks(L)

	LPrime, BPrime := computeLPrimeandBPrime(blocks)

	LPrimeSparseTable := computeSparseTable(LPrime)

	return &LCELinear{&st, L, E, R, blocks, LPrime, BPrime, LPrimeSparseTable}
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
func computeSparseTable(LPrime []int) [][]int {
	//compute sparse table
	n := len(LPrime)
	sparseTable := make([][]int, n)
	for i := 0; i < len(LPrime); i++ {
		sparseTable[i] = make([]int, int(math.Ceil(math.Log2(float64(n))))+1)

	}

	//first compute first row (trivial)
	for i, v := range LPrime {
		sparseTable[i][0] = v
	}

	//compute the rest of the sparse table
	for j := 1; 1<<j <= n; j++ {
		for i := 0; i+(1<<j) <= n; i++ {
			sparseTable[i][j] = minInt(sparseTable[i][j-1], sparseTable[i+(1<<(j-1))][j-1])
		}
	}

	return sparseTable
}

// function to compute the LCE between two indices i and j
func RMQLookup(i, j int, sparseTable [][]int) int {

	k := int(math.Floor((math.Log2(float64(j - i)))))

	range1 := sparseTable[i][k]
	range2 := sparseTable[j-(1<<k)+1][k]

	return minInt(range1, range2)

}

// helper function to find the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
