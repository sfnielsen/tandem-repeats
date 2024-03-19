package tandemrepeat

import (
	"fmt"
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

func TestAlgorithm2(t *testing.T) {
	input := "abaabaabbaaabaaba$"
	//input := "aaaaaaaaaaaaaaaaaaaaaaaaaababababaaabbbabbabaabbaaaabababaabababab$"
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(input)
	Algorithm1(st)
}
