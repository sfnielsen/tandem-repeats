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

	//backward LCE queries
	//TBD
}

// main function for the LCE linear time preprocessing
func PreProcessLCE(st suffixtree.SuffixTreeInterface) *LCELinear {
	//create L,E,R arrays
	L, E, R := createLERArrays(st)

	blocks := createLBlocks(L)

	LPrime, BPrime := computeLPrimeandBPrime(blocks)

	return &LCELinear{&st, L, E, R, blocks, LPrime, BPrime}
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
