package lce

import (
	"math"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
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
	Leafs                       []*suffixtree.SuffixTreeNode
	EulerindexToNode            []*suffixtree.SuffixTreeNode

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

// LCELookupForward returns the longest common extension of the nodes at index i and j in the forward direction
func (lceTwoWays *LCELinearTwoWays) LCELookupForward(i, j int) *suffixtree.SuffixTreeNode {
	return lceTwoWays.forward.LCELookup(i, j)
}

// LCELookupBackward returns the longest common extension of the nodes at index i and j in the BACKWARD direction
func (lceTwoWays *LCELinearTwoWays) LCELookupBackward(i, j int) *suffixtree.SuffixTreeNode {
	stringLength := len((*lceTwoWays.backward.suffixTree).GetInputString())
	return lceTwoWays.backward.LCELookup(stringLength-(j+2), stringLength-(i+2))

}

func (lce *LCELinear) LCELookup(i, j int) *suffixtree.SuffixTreeNode {
	//find the lowest common ancestor of i and j
	leaf_i := lce.Leafs[i]
	leaf_j := lce.Leafs[j]
	eulerindex_i := lce.R[leaf_i.EulerLabel]
	eulerindex_j := lce.R[leaf_j.EulerLabel]
	if eulerindex_i > eulerindex_j {
		eulerindex_i, eulerindex_j = eulerindex_j, eulerindex_i
	}

	blockSize := len(lce.blocks[0])
	block_i := math.Floor(float64(eulerindex_i) / float64(blockSize))
	block_j := math.Floor(float64(eulerindex_j) / float64(blockSize))

	// case 1:   They are on the same block
	if block_i == block_j {
		//find the normalized block
		normalizedBlock := normalizeBlock(lce.blocks[int(block_i)])
		//find the normalized block sparse table
		normalizedBlockSparseTable := lce.NormalizedBlockSparseTables[convertBlockToInt(normalizedBlock)]
		//find the LCA in the normalized block
		lca := RMQLookup(eulerindex_i%blockSize, eulerindex_j%blockSize, normalizedBlockSparseTable)
		lcaEulerIdx := int(block_i)*blockSize + lca.index
		return lce.EulerindexToNode[lce.E[lcaEulerIdx]]
	} else {
		// case 2:   They are on different blocks (i < j   always)
		// find the LCA in the first block
		lcaBlocki := convertBlockToInt(normalizeBlock(lce.blocks[int(block_i)]))
		lca1 := RMQLookup(eulerindex_i%blockSize, blockSize-1, lce.NormalizedBlockSparseTables[lcaBlocki])
		lca1.level += lce.blocks[int(block_i)][0]
		// find the LCA in the last block
		lcaBlockj := convertBlockToInt(normalizeBlock(lce.blocks[int(block_j)]))
		lca2 := RMQLookup(0, eulerindex_j%blockSize, lce.NormalizedBlockSparseTables[lcaBlockj])
		lca2.level += lce.blocks[int(block_j)][0]
		// find the LCA between the two blocks

		//edgecase when block i and j are adjacent
		if block_j == block_i+1 {
			if lca1.level < lca2.level {
				return lce.EulerindexToNode[lce.E[int(block_i)*blockSize+lca1.index]]
			}
			return lce.EulerindexToNode[lce.E[int(block_j)*blockSize+lca2.index]]
		}

		//i and j are not adjacent
		lca3 := RMQLookup(int(block_i)+1, int(block_j)-1, lce.LPrimeST)
		// find the LCA of the three LCA's
		lca1EulerIdx := int(block_i)*blockSize + lca1.index
		lca2EulerIdx := int(block_j)*blockSize + lca2.index
		lca3EulerIdx := lce.BPrime[lca3.index]

		if lca1.level < lca2.level {
			if lca1.level < lca3.level {
				return lce.EulerindexToNode[lce.E[lca1EulerIdx]]
			}
		} else if lca2.level < lca3.level {
			return lce.EulerindexToNode[lce.E[lca2EulerIdx]]
		}
		return lce.EulerindexToNode[lce.E[lca3EulerIdx]]
	}
}

// main function for linear time LCE preprocessing in both directions
func PreProcessLCEBothDirections(st suffixtree.SuffixTreeInterface) *LCELinearTwoWays {
	//create forward LCE
	forwardLCE := PreProcessLCE(st)

	//create backward LCE
	reversedString := reverseStringWithSentinel(st.GetInputString())
	reversedSuffixTree := suffixtreeimpl.ConstructMcCreightSuffixTree(reversedString)
	reversedSuffixTree.AddStringDepth()

	backwardLCE := PreProcessLCE(reversedSuffixTree)

	return &LCELinearTwoWays{forward: forwardLCE, backward: backwardLCE}

}

// main function for the LCE linear time preprocessing (in one direction)
func PreProcessLCE(st suffixtree.SuffixTreeInterface) *LCELinear {

	//create L,E,R arrays
	L, E, R, EulerindexToNode := createLERArrays(st)
	// divide L into blocks
	blocks := createLBlocks(L)
	// compute L' and B'
	LPrime, BPrime := computeLPrimeandBPrime(blocks)
	// compute sparse table for L'
	LPrimeSparseTable := computeSparseTable(LPrime)
	// precompute all possible normalized blocks
	NormalizedBlockSparseTables := computeNormalizedBlockSparseTables(blocks)
	// compute leafs of tree
	leafSlice := st.ComputeLeafs()

	return &LCELinear{&st, L, E, R, blocks, LPrime, BPrime, LPrimeSparseTable, NormalizedBlockSparseTables, leafSlice, EulerindexToNode}
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
	table := make([][]stTuple, n)
	for i := 0; i < len(LPrime); i++ {
		table[i] = make([]stTuple, int(math.Ceil(math.Log2(float64(n))))+1)
	}

	//first compute first row (trivial)
	for i, v := range LPrime {
		tup := stTuple{level: v, index: i}
		table[i][0] = tup
	}

	//compute the rest of the sparse table
	for j := 1; 1<<j <= n; j++ {
		for i := 0; i+(1<<j) <= n; i++ {
			table[i][j] = stTupMin(table[i][j-1], table[i+(1<<(j-1))][j-1])
		}
	}

	return &table
}

// function to compute the LCE between two indices i and j
func RMQLookup(i, j int, sparseTable *sparseTable) stTuple {
	if i == j {
		return (*sparseTable)[i][0]
	}
	var k uint
	if j-i == 1 {
		k = 0
	} else {
		k = uint(math.Floor((math.Log2(float64(uint(j - i))))))

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
