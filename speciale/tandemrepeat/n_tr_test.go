package tandemrepeat

import (
	"fmt"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
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
			expectedS: []int{-1, -1, 0, 0, 1, 2, 0, 1, 1, 2, 2, 0, 1, 2, 0, 1, 0, -1},
			lzBlocks:  []int{0, 1, 2, 3, 8, 11, 17},
		},

		{
			//just another simple test
			input:     "abbbaabbbb$",
			expectedL: []int{0, 0, 2, 1, 1, 4, 3, 3, 2, 1, 0},
			expectedS: []int{-1, -1, 1, 1, 0, 0, 1, 1, 1, 1, -1},
			lzBlocks:  []int{0, 1, 2, 4, 5, 9, 10},
		},
	}

	for _, tc := range LZTestCases {
		// Compute the LZ decomposition
		tree := suffixtreeimpl.ConstructMcCreightSuffixTree(tc.input)

		li, si := LZDecomposition(tree)
		lzB := CreateLZBlocks(li, si)
		// Compare the computed values with the expected values
		for i := range li {
			if li[i] != tc.expectedL[i] {
				t.Errorf("Test case failed for input '%s': Expected li[%d] = %d, but got %d", tc.input, i, tc.expectedL[i], li[i])
			}
			if si[i] != tc.expectedS[i] {
				t.Errorf("Test case failed for input '%s': Expected si[%d] = %d, but got %d", tc.input, i, tc.expectedS[i], si[i])
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

// test alg 1
func TestAlgorithm1(t *testing.T) {
	// Test on example from paper https://doi.org/10.1016/j.jcss.2004.03.004
	input := "abaabaabbaaabaaba$"
	//input := "aaaaaaaaaaaaaaaaaaaaaaaaaababababaaabbbabbabaabbaaaabababaabababab$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	Algorithm1(st)
}

func TestAllRepeatTypesOfLinearAlgoPhase1(t *testing.T) {
	// intialize maps for all repeat types for both algorithm 1 and nlogn tandem repeat algorithm
	phase1Repeats := make(map[string]TandemRepeat)
	allRepeats := make(map[string]TandemRepeat)

	randomGenerator_ab.SetSeed(45)
	input := randomGenerator_ab.GenerateString(2324)
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
		for _, v2 := range v {
			BranchingTandemRepeatTypesPhase1 = append(BranchingTandemRepeatTypesPhase1, v2)
		}
	}
	allTandemRepeatsTypesPhase1 := rightRotation(BranchingTandemRepeatTypesPhase1, st)

	for _, v := range allTandemRepeatsTypesPhase1 {
		phase1Repeats[GetTandemRepeatSubstring(v, input)] = v
	}

	fmt.Println(st.GetInputString())
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

// test that all the sets are sorted and contains no duplicates
func TestAlgorithm1SetsAreSortedAndNoDuplicates(t *testing.T) {
	for i := 0; i < 30; i++ {
		randomGenerator_ab.SetSeed(113)
		input := randomGenerator_ab.GenerateString(1213)
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
}

func TestAlgorithm2(t *testing.T) {
	input := "abaabaabbaaabaaba$"
	//input := "aaaaaaaaaaaaaaaaaaaaaaaaaababababaaabbbabbabaabbaaaabababaabababab$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	Algorithm1(st)
}

func TestDecorateTreeIsCorrectOnTestString(t *testing.T) {
	input := "abaabaabbaaabaaba$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	DecorateTreeWithVocabulary(st)

	//Bottom-up traversal of the suffix tree
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			dfs(child, depth+child.EdgeLength())
		}

		//if tree has a decoration print int
		if node.TandemRepeatDeco != nil {

			fmt.Println("new dayz")
			fmt.Println(GetTandemRepeatSubstring(TandemRepeat{node.Label, (depth - node.EdgeLength() + node.TandemRepeatDeco[0]) / 2, 2}, input))

			fmt.Printf("Node: %v, Decoration: %v\n", node.Label, node.TandemRepeatDeco)
		}

	}
	dfs(st.GetRoot(), 0)
}

func TestDecorateTreeIsCorrectOnRandomString(t *testing.T) {
	input := "abaabaabbaaabaaba$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	DecorateTreeAndReturnTandemRepeats(st)

	//Bottom-up traversal of the suffix tree
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			dfs(child, depth+child.EdgeLength())
		}

		//if tree has a decoration print int
		if node.TandemRepeatDecoComplete != nil {

			fmt.Println("new dayz")
			for k := range node.TandemRepeatDecoComplete {
				fmt.Println(GetTandemRepeatSubstring(TandemRepeat{node.Label, (depth - node.EdgeLength() + k) / 2, 2}, input))
			}
			fmt.Printf("Node: %v, Decoration: %v\n", node.Label, node.TandemRepeatDecoComplete)
		}

	}
	dfs(st.GetRoot(), 0)
}

func TestThatWeReturnAllTandemRepeats(t *testing.T) {

	randomGenerator_ab.SetSeed(1)
	for i := 0; i < 10; i++ {
		input := randomGenerator_ab.GenerateString(20)

		// get all tandem repeats
		tandemRepeats := FindTandemRepeatsNaive(input)
		for _, v := range tandemRepeats {
			fmt.Println(GetTandemRepeatSubstring(v, input), "at index:", v.Start)
		}

		//input := "abaabaabbaaabaaba$"
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)

		tandemRepeats2 := DecorateTreeAndReturnTandemRepeats(st)
		for _, v := range tandemRepeats2 {
			fmt.Println(v)
			fmt.Println(GetTandemRepeatSubstring(v, input), "at index:", v.Start)
		}

		println()

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
