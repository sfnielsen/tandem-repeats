package tandemrepeat

import (
	"fmt"
	"math/bits"
	"speciale/suffixtree"
)

func getPositionOfLeftmostBitWithValue1(n int) int {
	//return the position of the leftmost bit with value 1
	//if no bit is set, return -1
	//if the input is 0, return 0
	if n == 0 {
		return 0
	}
	//find the position of the leftmost bit with value 1
	//use the fact that -n is the two's complement of n
	//and that n & -n is the rightmost bit with value 1
	return n & -n
}
func trailingZerosCount(num int) int {
	if num == 0 {
		return 0
	}
	return bits.TrailingZeros16(uint16(num))
}
func leadingZerosCount(num int) int {
	if num == 0 {
		return 0
	}
	return bits.LeadingZeros8(uint8(num))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func boolsToInts(slice []bool) int {
	result := 0
	for _, b := range slice {
		result = (result << 1) | boolToInt(b)
	}
	return result
}
func addIVvalues(st suffixtree.SuffixTreeInterface) {
	var dfs func(node *suffixtree.SuffixTreeNode) *suffixtree.SuffixTreeNode
	dfs = func(node *suffixtree.SuffixTreeNode) *suffixtree.SuffixTreeNode {
		if node.IsLeaf() {
			node.L_k = node
			node.TrailingZeros = trailingZerosCount(node.DfsInterval.Start)
			return node
		} else {
			node.TrailingZeros = trailingZerosCount(node.DfsInterval.Start)
			node.L_k = node
			for _, child := range node.Children {
				if child != nil {
					childL_k := dfs(child)
					if childL_k.TrailingZeros > node.TrailingZeros {
						node.L_k = childL_k
						node.TrailingZeros = childL_k.TrailingZeros
					}
				}
			}
			return node.L_k
		}
	}
	dfs(st.GetRoot())
}

func preprocessMappingFromTreeToB(st suffixtree.SuffixTreeInterface) {
	//bottom up traversal of the suffix tree
	//for each node we calculate i(v)
	//i(v) is the node in the subtree for which the dfs label has the largest number of consecutive zeros at its right end
	//iterate slice of leafs and do bottom up traversal:
	// dfs traversal of the tree

	inputStringBitLength := bits.Len(uint(st.GetSize()))
	bool_array := make([]bool, inputStringBitLength)
	for i := range bool_array {
		bool_array[i] = false
	}
	addIVvalues(st)

	var dfs func(node *suffixtree.SuffixTreeNode, a_v []bool)
	dfs = func(node *suffixtree.SuffixTreeNode, a_v []bool) {
		a_v_copy := make([]bool, len(a_v))
		copy(a_v_copy, a_v)
		if node.Label != -1 {
			a_v_copy[node.TrailingZeros] = true
			node.A_v_int = boolsToInts(a_v_copy)
		}
		if !node.IsLeaf() {
			for _, child := range node.Children {
				if child != nil {
					dfs(child, a_v_copy)
				}
			}
		}

	}
	dfs(st.GetRoot(), bool_array)

}
func xorSlices(slice1, slice2 []bool) int {
	if len(slice1) != len(slice2) {
		return 0 // or handle the mismatch in a suitable way
	}

	intValue1 := boolsToInts(slice1)
	intValue2 := boolsToInts(slice2)

	xorResult := intValue1 ^ intValue2
	return xorResult
}
func intsToBools(value, size int) []bool {
	result := make([]bool, size)
	for i := size - 1; i >= 0; i-- {
		result[i] = value&1 == 1
		value >>= 1
	}
	return result
}
func leftmostOneToRight(n, position int, bitsize int) int {
	// Shift the bits to the right so that the target position becomes the rightmost bit
	bitLength := bits.Len(uint(n))
	mask := (1 << (bitsize - position - 1)) - 1
	// Apply the mask to get the bits to the right of the specified position
	rightPart := n & mask
	bitLength = bits.Len(uint(rightPart))
	bitLengthmask := bits.Len(uint(mask))
	return position + (bitLengthmask - bitLength) + 1
}
func CreateMask(N int, totalBits int) int {
	// Create a mask with the first N bits set to 1
	mask := (1 << N) - 1
	if N > totalBits {
		return mask
	}
	mask = mask << (totalBits - N)

	return mask
}
func ReplaceBits(num, bits, N int, totalbits int) int {
	// Create a mask with the first N bits set to 1
	mask := CreateMask(N+1, totalbits)

	// Clear the first N bits in num
	num &= mask

	return num
}
func lcaLookup(i int, j int, st suffixtree.SuffixTreeInterface) {
	inputStringBitLength := bits.Len(uint(st.GetSize()))
	preprocessMappingFromTreeToB(st)
	leafs := st.ComputeLeafs()
	leaf_x := leafs[i]
	leaf_y := leafs[j]

	// STEP 1-2
	j_value := trailingZerosCount((leaf_x.L_k.DfsInterval.Start & leaf_y.L_k.DfsInterval.Start))
	fmt.Printf("x in bit:   %d: %b\n", leaf_x.L_k.DfsInterval.Start, leaf_x.L_k.DfsInterval.Start)
	fmt.Printf("y in bit:   %d: %b\n", leaf_y.L_k.DfsInterval.Start, leaf_y.L_k.DfsInterval.Start)
	fmt.Println("J-value / trailing zeros match / height", j_value)
	// STEP 3;  FIND X BAR
	// 3a
	fmt.Printf("A_x;  %d, %b\n", leaf_x.A_v_int, leaf_x.A_v_int)

	l_value_x := trailingZerosCount(leaf_x.A_v_int)
	var x_bar int

	fmt.Println("X     l-value", inputStringBitLength-l_value_x-1)
	if j_value == l_value_x {
		// 3b
		x_bar = i
	} else {
		// 3c
		posOf1bitFromRight := leftmostOneToRight(leaf_x.A_v_int, j_value, inputStringBitLength)
		i_w_x := ReplaceBits(leaf_x.A_v_int, 0, posOf1bitFromRight, inputStringBitLength)
		fmt.Println(i_w_x)
		x_bar = leafs[i_w_x].L_k.Parent.Label
	}
	// Clear the rightmost X bits in the original number
	fmt.Printf("x-bar:   %d: %b\n", x_bar, x_bar)

	// STEP 4;  FIND Y BAR
	l_value_y := trailingZerosCount(leaf_y.A_v_int)
	var y_bar int
	if j_value == l_value_y {
		y_bar = i
	} else {
		posOf1bitFromRight := leftmostOneToRight(leaf_y.A_v_int, j_value, inputStringBitLength)
		i_w_y := ReplaceBits(leaf_x.A_v_int, 0, posOf1bitFromRight, inputStringBitLength)
		y_bar = leafs[i_w_y].L_k.Parent.Label
	}
	// Clear the rightmost X bits in the original number
	fmt.Printf("ybar:   %d: %b\n", y_bar, y_bar)

}
