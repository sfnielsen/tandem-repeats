package tandemrepeat

import (
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
