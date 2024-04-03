package lce

import (
	"speciale/suffixtree"
)

// Create L, E, R arrays
func createLERArrays(st suffixtree.SuffixTreeInterface) ([]int, []int, []int) {

	//do an euler tour of the suffix tree
	L, E, R := eulerTour(st)

	return L, E, R
}

// traverse the tree in a depth first manner
func eulerTour(st suffixtree.SuffixTreeInterface) ([]int, []int, []int) {
	//euler labels
	nextEulerLabel := 0
	nextEulerStep := 0

	//tables
	L := make([]int, 2*st.GetSize()-1)
	E := make([]int, 2*st.GetSize()-1)
	R := make([]int, st.GetSize())

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
			dfs(child, depth+1)
		}

		//process self again
		L[nextEulerStep] = depth           //note the depth of current eulerStep
		E[nextEulerStep] = node.EulerLabel //map eulerStep to eulerLabel
		nextEulerStep++

	}
	dfs(st.GetRoot(), 0)

	return L, E, R

}
