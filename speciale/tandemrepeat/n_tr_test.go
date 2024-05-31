package tandemrepeat

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"speciale/lce"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
	"unsafe"
)

// test LZDecomposition
func TestLZDecompositionOnSimpleStrings(t *testing.T) {

	LZTestCases := []struct {
		input     string
		expectedL []int
		expectedS []int
		lzBlocks  []int
	}{

		{
			// test on example from paper https://doi.org/10.1016/j.jcss.2004.03.004
			input:     "abaabaabbaaabaaba$",
			expectedL: []int{0, 0, 1, 5, 4, 3, 2, 1, 3, 2, 6, 6, 5, 4, 3, 2, 1, 0},
			lzBlocks:  []int{0, 1, 2, 3, 8, 11, 17},
		},

		{
			//just another simple test
			input:     "abbbaabbbb$",
			expectedL: []int{0, 0, 2, 1, 1, 4, 3, 3, 2, 1, 0},
			lzBlocks:  []int{0, 1, 2, 4, 5, 9, 10},
		},
	}

	for _, tc := range LZTestCases {
		// Compute the LZ decomposition
		tree := suffixtreeimpl.ConstructMcCreightSuffixTree(tc.input)

		li := LZDecompositionStackMethod(tree)
		lzB := CreateLZBlocks(li)
		// Compare the computed values with the expected values
		for i := range li {
			if li[i] != tc.expectedL[i] {
				t.Errorf("Test case failed for input '%s': Expected li[%d] = %d, but got %d", tc.input, i, tc.expectedL[i], li[i])
			}

		}
		for i := range lzB {
			if lzB[i] != tc.lzBlocks[i] {
				t.Errorf("Test case failed for input '%s': Expected lzBlocks[%d] = %d, but got %d", tc.input, i, tc.lzBlocks[i], lzB[i])
			}
		}
		if len(lzB) != len(tc.lzBlocks) {
			t.Errorf("Test case failed for input '%s': Expected %d lzBlocks, but got %d", tc.input, len(tc.lzBlocks), len(lzB))
		}

	}

}

// Test that all sets from algorithm 1 are sorted
func TestAlgorithm1SetsAreSorted(t *testing.T) {
	randomGenerator_ab.SetSeed(77)
	input := randomGenerator_ab.GenerateString(1531)
	//input := "abaabaabbaaabaaba$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	leftMostCoveringSet := Algorithm1(st)
	for idx_v, v := range leftMostCoveringSet {
		for i := 0; i < len(v)-1; i++ {
			if !(v[i].Length < v[i+1].Length) {
				t.Errorf("Expected leftMostCoveringSet to be sorted and to not have duplicates,%d, %d, %d", idx_v, v[i].Length, v[i+1].Length)
			}
		}
	}
}

// test that all algorithm 1 results are tandem repeats
func TestAlg1OnlyFindsTandemRepeats(t *testing.T) {
	randomGenerator_ab.SetSeed(1)
	for i := 0; i < 5; i++ {
		input := randomGenerator_ab.GenerateString(500)
		//input := "abaabaabbaaabaaba$"
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

		leftMostCoveringSet := Algorithm1(st)
		for _, v := range leftMostCoveringSet {
			for _, v2 := range v {
				if input[v2.Start:v2.Start+v2.Length] != input[v2.Start+v2.Length:v2.Start+2*v2.Length] {
					t.Errorf("Expected %s to be a tandem repeat", GetTandemRepeatSubstring(v2, input))
				}
			}
		}
	}
}

func TestAllRepeatTypesOfLinearAlgoPhase1(t *testing.T) {
	// intialize maps for all repeat types for both algorithm 1 and nlogn tandem repeat algorithm
	phase1Repeats := make(map[string]TandemRepeat)
	allRepeats := make(map[string]TandemRepeat)

	randomGenerator_ab.SetSeed(40)
	input := randomGenerator_ab.GenerateString(5000)
	//input := "abaabaabbaaabaaba$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

	tandemRepeats := FindAllTandemRepeatsLogarithmic(st)
	for _, v := range tandemRepeats {
		allRepeats[GetTandemRepeatSubstring(v, input)] = v
	}

	leftMostCoveringSet := Algorithm1(st)
	// add all repeats from nested slice "leftmostcoveringset" to a single slice of tandemrepeats so we can input to getALlTandemRepeats
	var BranchingTandemRepeatTypesPhase1 []TandemRepeat
	for _, v := range leftMostCoveringSet {
		BranchingTandemRepeatTypesPhase1 = append(BranchingTandemRepeatTypesPhase1, v...)
	}
	allTandemRepeatsTypesPhase1 := rightRotation(BranchingTandemRepeatTypesPhase1, st)

	for _, v := range allTandemRepeatsTypesPhase1 {
		phase1Repeats[GetTandemRepeatSubstring(v, input)] = v
	}

	// check that all repeats from phase 1 is in all repeats
	for k, v := range phase1Repeats {
		if _, ok := allRepeats[k]; !ok {
			t.Errorf("Expected %s to be in all repeats at index %d", k, v.Start)
		}
	}

	//check that all repeats from all repeats is in phase 1
	for k, v := range allRepeats {
		if _, ok := phase1Repeats[k]; !ok {
			t.Errorf("Expected %s to be in phase 1 repeats at index %d", k, v.Start)
		}
	}
}

// Test that algorithm 2 only decorates tree with actual tandem repeats
func TestAlg2OnlyDecoratesTreeWithTandemRepeats(t *testing.T) {

	//Run alg 2
	randomGenerator_ab.SetSeed(1)
	for i := 0; i < 10; i++ {
		input := randomGenerator_ab.GenerateString(100)

		tree := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

		// Phase 1
		// get leftmost covering repeats
		leftMostCoveringRepeats := Algorithm1(tree)

		//hacky way to get rid of tandem repeats and have ints instead.
		//Could be improved at a later point
		leftMostCoveringRepeatsInts := make([][]int, len(leftMostCoveringRepeats))
		for idx, k := range leftMostCoveringRepeats {
			leftMostCoveringRepeatsInts[idx] = make([]int, 0)
			for _, j := range k {
				leftMostCoveringRepeatsInts[idx] = append(leftMostCoveringRepeatsInts[idx], j.Length)

			}
		}

		// Phase 2
		// Decorate tree with subset of leftmost covering repeats
		Algorithm2StackMethod(tree, leftMostCoveringRepeatsInts)

		//Bottom-up traversal of the suffix tree
		var dfs func(node *suffixtree.SuffixTreeNode)
		dfs = func(node *suffixtree.SuffixTreeNode) {
			// Traverse the children of the current node
			for _, child := range node.Children {
				if child == nil {
					continue
				}
				dfs(child)
			}

			//if tree has a decoration print int
			if node.TandemRepeatDeco != nil {
				for _, k := range node.TandemRepeatDeco {
					if (node.Parent.StringDepth+k)%2 != 0 {
						t.Errorf("Expected tandem repeat but has odd length: %d", node.StringDepth+k)
					}
					lenOfRepeat := node.Parent.StringDepth + k
					if input[node.Label:node.Label+lenOfRepeat/2] != input[node.Label+lenOfRepeat/2:node.Label+lenOfRepeat] {
						t.Errorf("Expected %s to be a tandem repeat", GetTandemRepeatSubstring(TandemRepeat{node.Label, (node.Parent.StringDepth + k) / 2, 2}, input))
					}
				}
			}

		}
		dfs(tree.GetRoot())
	}

}

func TestDecorateTreeOnlyDecoratesTreeWithTandemRepeats(t *testing.T) {

	randomGenerator_ab.SetSeed(4)
	for i := 0; i < 10; i++ {
		input := randomGenerator_ab.GenerateString(3000)
		tree := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

		DecorateTreeWithVocabulary(tree)

		//Bottom-up traversal of the suffix tree
		var dfs func(node *suffixtree.SuffixTreeNode)
		dfs = func(node *suffixtree.SuffixTreeNode) {
			// Traverse the children of the current node
			for _, child := range node.Children {
				if child == nil {
					continue
				}
				dfs(child)
			}

			//if tree has a decoration print int
			if node.TandemRepeatDecoComplete != nil {
				for k := range node.TandemRepeatDecoComplete {
					if (node.Parent.StringDepth+k)%2 != 0 {
						t.Errorf("Expected tandem repeat but has odd length: %d", node.StringDepth+k)
					}
					lenOfRepeat := node.Parent.StringDepth + k
					if input[node.Label:node.Label+lenOfRepeat/2] != input[node.Label+lenOfRepeat/2:node.Label+lenOfRepeat] {
						t.Errorf("Expected %s to be a tandem repeat", GetTandemRepeatSubstring(TandemRepeat{node.Label, (node.Parent.StringDepth + k) / 2, 2}, input))
					}
				}
			}

		}
		dfs(tree.GetRoot())
	}
}

func TestThatWeReturnAllTandemRepeats(t *testing.T) {

	randomGenerator_ab.SetSeed(1)
	for i := 0; i < 10; i++ {
		input := randomGenerator_ab.GenerateString(2518)

		// get all tandem repeats
		tandemRepeats := FindTandemRepeatsNaive(input)

		st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

		tandemRepeats2 := DecorateTreeAndReturnTandemRepeats(st)

		//Cgecj that both algorithm finds the same tandem repeats
		for i := range tandemRepeats {
			for j := range tandemRepeats2 {
				if tandemRepeats[i].Start == tandemRepeats2[j].Start {
					break
				}
				if j == len(tandemRepeats2)-1 {
					t.Errorf("Expected tandem repeats to be the same. Not the case for string: %s", input)
				}
			}
		}

		//check that no duplicates are found in n time alg
		for i := range tandemRepeats2 {
			for j := range tandemRepeats2 {
				if tandemRepeats2[i] == tandemRepeats2[j] && i != j {
					t.Errorf("Expected no duplicates in tandem repeats. Not the case for string: %s", input)
				}
			}
		}

		//check that all found tandem repeats are actual repeats
		for _, v := range tandemRepeats2 {
			if input[v.Start:v.Start+v.Length] != input[v.Start+v.Length:v.Start+2*v.Length] {
				t.Errorf("Expected %s to be a tandem repeat. Error was in string %s", GetTandemRepeatSubstring(v, input), input)
			}
		}

		//check that the length of the two slices are the same
		if len(tandemRepeats) != len(tandemRepeats2) {
			t.Errorf("Expected tandem repeats to be the same length: true amount is %d, but O(n) alg found %d", len(tandemRepeats), len(tandemRepeats2))
		}

	}
}

// #####################################################################################
// #####################################################################################
// Helper functions used for testing
// #####################################################################################
// #####################################################################################

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

func TestForwardLookup(t *testing.T) {
	randomGenerator_ab.SetSeed(410)
	s := randomGenerator_ab.GenerateString(1000)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	lceObject := lce.PreProcessLCE(st)

	// run through all pairs of i and j
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			// find the LCE using the slow method
			realLength := lce.FindLCEForwardSlow(s, i, j)
			// find the LCE using the LCELookup method
			lce := lceObject.LCELookup(i, j)

			if realLength != lce {
				t.Errorf("Expected %d, got %d", realLength, lce)
			}
		}
	}

}

// Test that LCE backward (and forward) works
func TestBackwardAndForwardLookup(t *testing.T) {
	randomGenerator_ab.SetSeed(40)
	s := randomGenerator_ab.GenerateString(1634)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	//stringLength := len(st.GetInputString())
	lceObject := lce.PreProcessLCEBothDirections(st)

	// run through all pairs of i and j
	for i := 1; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {

			// find the LCE using the slow method
			realLengthFW := lce.FindLCEForwardSlow(s, i, j)
			realLengthBW := lce.FindLCEBackwardSlow(s, i-1, j-1)
			// find the LCE using the LCELookup method
			//lca := lceObject.backward.LCELookup(stringLength-j-1, stringLength-i-1)

			// check if the LCE is correct
			lceFW := lceObject.LCELookupForward(i, j)
			lceBW := lceObject.LCELookupBackward(i-1, j-1)

			if realLengthFW != lceFW {
				t.Errorf("Expected %d, got %d", realLengthFW, lceFW)
			}
			if realLengthBW != lceBW {
				t.Errorf("Expected %d, got %d", realLengthBW, lceBW)
			}

		}
	}

}

// #####################################################################################
// #####################################################################################
// Benchmarking
// #####################################################################################
// #####################################################################################
func BenchmarkExample(b *testing.B) {
	debug.SetGCPercent(5000)
	randomGenerator_ab.SetSeed(42)
	str := randomGenerator_ab.GenerateString(3000000)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(str)

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	FindAllBranchingTandemRepeatsLogarithmic(st)

	defer pprof.StopCPUProfile()

}

func BenchmarkExample2(b *testing.B) {
	var e suffixtree.SuffixTreeNode
	fmt.Printf("Size of %T struct: %d bytes", e, unsafe.Sizeof(e))
	fmt.Println()

}
