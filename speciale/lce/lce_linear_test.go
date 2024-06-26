package lce

import (
	"math"
	"speciale/stringgenerators"
	"speciale/suffixtreeimpl"
	"strconv"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
var (
	setupCompleted          bool
	randomGenerator_protein stringgenerators.StringGenerator
	randomGenerator_ab      stringgenerators.StringGenerator
	randomGenerator_dna     stringgenerators.StringGenerator
	randomGenerator_byte    stringgenerators.StringGenerator
	randomGenerator_a       stringgenerators.StringGenerator
)

// addLeafList adds leaflists to the suffix tree
func init() {
	if !setupCompleted {
		randomGenerator_protein = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetProtein}
		randomGenerator_ab = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetAB}
		randomGenerator_dna = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetDNA}
		randomGenerator_byte = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetByte}
		randomGenerator_a = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetA}
		setupCompleted = true
	}
}

func TestLCEArrays(t *testing.T) {

	s := randomGenerator_ab.GenerateString(3000)

	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)

	L, E, R, _, _ := createLERArraysAndLeafArrayStack(st)
	n := st.GetSize()

	//check that L,E,R sizes are correct
	if len(L) != 2*n-1 {
		t.Errorf("Expected L size to be %d, got %d", 2*n-1, len(L))
	}
	if len(E) != 2*n-1 {
		t.Errorf("Expected E size to be %d, got %d", 2*n-1, len(E))
	}
	if len(R) != n {
		t.Errorf("Expected R size to be %d, got %d", n, len(R))
	}

	//Test that levels/depths are +-1 of each other
	for i := 0; i < len(L)-1; i++ {
		if L[i] != L[i+1]+1 && L[i] != L[i+1]-1 {
			t.Errorf("Expected L[i] to be L[i+1]+-1")
		}
	}
	//first and last depth is 0
	if L[0] != 0 || L[len(L)-1] != 0 {
		t.Errorf("Expected first and last depth to be 0 (ALWAYS root)")
	}

	//setup for additional testing
	allInts := make(map[int]int)
	for _, v := range R {
		allInts[E[v]] = L[v]
	}

	//check that R maps to all 0,1,...,n-1 different ints in L
	if len(allInts) != n {
		t.Errorf("Expected R to map to all 0,1,...,n-1 different ints")
	}
	//check that all same E values also have same L values
	for _, v := range E {
		if allInts[E[v]] != L[v] {
			t.Errorf("Expected L to be the same for all same E values")
		}

	}

}

// Test that the Blocks from L are created correctly
func TestLBlocks(t *testing.T) {
	randomGenerator_ab.SetSeed(1)
	s := randomGenerator_ab.GenerateString(721)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	L, _, _, _, _ := createLERArraysAndLeafArrayStack(st)
	blocks := createLBlocks(L)
	n := len(L)
	blockSize := int(math.Ceil(math.Log2(float64(n)) / 2))
	numBlocks := int(math.Ceil(float64(n) / float64(blockSize)))

	//check that blocks are correct size
	if len(blocks) != numBlocks {
		t.Errorf("Expected %d blocks, got %d", numBlocks, len(blocks))
	}

	LFromBlocks := 0
	//check that blocks are correct size
	for i := 0; i < len(blocks); i++ {
		LFromBlocks += len(blocks[i])

		if len(blocks[i]) != blockSize && i != len(blocks)-1 {
			t.Errorf("Expected block size to be %d, got %d", blockSize, len(blocks[i]))
		}

		if i == len(blocks)-1 {
			if len(blocks[i]) != n%blockSize || len(blocks[i]) == 0 {
				t.Errorf("Expected block size to be %d, got %d", n%blockSize, len(blocks[i]))
			}
		}
	}

	//check that total length of blocks is equal to L
	if LFromBlocks != len(L) {
		t.Errorf("Expected total length of blocks to be %d, got %d", len(L), LFromBlocks)
	}

}

// Test that L' and B' are computed correctly
func TestLPrimeandBPrime(t *testing.T) {
	randomGenerator_ab.SetSeed(1)
	s := randomGenerator_ab.GenerateString(721)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	L, _, _, _, _ := createLERArraysAndLeafArrayStack(st)
	blocks := createLBlocks(L)
	LPrime, BPrime := computeLPrimeandBPrime(blocks)
	sizeNormalBlock := len(blocks[0])

	//check that L' and B' are correct size
	if len(LPrime) != len(blocks) {
		t.Errorf("Expected L' to be of size %d, got %d", len(blocks), len(LPrime))
	}
	if len(BPrime) != len(blocks) {
		t.Errorf("Expected B' to be of size %d, got %d", len(blocks), len(BPrime))
	}

	//check that L' and B' are correct
	for i, block := range blocks {
		min := math.MaxInt32
		for _, v := range block {
			if v < min {
				min = v
			}
		}
		if LPrime[i] != min {
			t.Errorf("Expected L'[i] to be the smallest element in the block")
		}

		minIndex := -1
		for j, v := range block {
			if v == min {
				minIndex = j
				break
			}
		}

		if BPrime[i] != i*sizeNormalBlock+minIndex {
			t.Errorf("Expected B'[i] to be the index of the smallest element in the block")
		}

		if L[BPrime[i]] != min {
			t.Errorf("Expected L[B'[i]] (%d) to be the smallest element in the block (%d)", L[BPrime[i]], min)
		}
	}

}

// Test that the sparse table is computed correctly
func TestSparseTable(t *testing.T) {
	//setup
	randomGenerator_ab.SetSeed(1)
	s := randomGenerator_ab.GenerateString((1 << 10) + 1) //1025
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	L, _, _, _, _ := createLERArraysAndLeafArrayStack(st)
	blocks := createLBlocks(L)
	LPrime, _ := computeLPrimeandBPrime(blocks)
	sparseTable := computeSparseTable(LPrime)

	//iterate all combinations of i,j (i<=j) and check that RMQ is correct
	for i := 0; i < len(LPrime)-1; i++ {
		for j := i + 1; j < len(LPrime); j++ {
			//naively find the smallest element in the range
			min := math.MaxInt32
			for k := i; k <= j; k++ {
				if LPrime[k] < min {
					min = LPrime[k]
				}
			}
			//check that RMQ on the sparse table achieves the same result as the naive approach
			if RMQLookup(i, j, sparseTable).level != min {
				t.Errorf("Expected RMQ to be %d, got %d for indexes i=%d,j=%d", min, RMQLookup(i, j, sparseTable), i, j)
			}

		}
	}

}

func TestConvertBinaryToValues(t *testing.T) {
	// test whether the testarray produces the same int as the binary string parsed to int
	binaryString := "111010"

	// Parse binary string to integer
	correct, _ := strconv.ParseInt(binaryString, 2, 64)

	testArray := []int{3, 4, 5, 4, 5, 4}
	resultingInt := convertBlockToInt(testArray)
	if correct != int64(resultingInt) {
		// test failed
		t.Errorf("Expected %d, got %d", correct, resultingInt)
	}

	// test longer binary string
	binaryString = "111010100101011110101011"
	correct, _ = strconv.ParseInt(binaryString, 2, 64)

	//create testarray
	testArray = []int{3, 4, 5, 4, 5, 4, 5, 4, 3, 4, 3, 4, 3, 4, 5, 6, 7, 6, 7, 6, 7, 6, 7, 8}
	resultingInt = convertBlockToInt(testArray)
	if correct != int64(resultingInt) {
		// test failed
		t.Errorf("Expected %d, got %d", correct, resultingInt)
	}

}

func TestComputeNormalizedBlock(t *testing.T) {
	randomGenerator_ab.SetSeed(1)
	s := randomGenerator_ab.GenerateString((1 << 10) + 1) //1025
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	L, _, _, _, _ := createLERArraysAndLeafArrayStack(st)
	blocks := createLBlocks(L)

	normalizedBlocks, blockIdxToNormalizedBlockIdx := computeNormalizedBlockSparseTables(blocks)

	for i, block := range blocks {
		blockStart := block[0]
		blockAsBinary := blockIdxToNormalizedBlockIdx[i]
		blockST := normalizedBlocks[blockAsBinary]
		for i := 0; i < len(block)-1; i++ {
			for j := i + 1; j < len(block); j++ {
				//naively find the smallest element in the range
				min := math.MaxInt32
				for k := i; k <= j; k++ {
					if block[k] < min {
						min = block[k]
					}
				}
				//check that RMQ on the sparse table achieves the same result as the naive approach
				if RMQLookup(i, j, blockST).level+blockStart != min {
					t.Errorf("Expected RMQ to be %d, got %d for indexes i=%d,j=%d", min, RMQLookup(i, j, blockST), i, j)
				}

			}
		}
	}

}
