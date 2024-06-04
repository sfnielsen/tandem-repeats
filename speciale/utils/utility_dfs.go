package utils

import (
	"speciale/suffixtree"
)

// performan depth first run of a suffix tree and for each node encountered, write down in an int-slice how many children it has
func DfsCountChildren(st suffixtree.SuffixTreeInterface) []int {
	childrenCount := make([]int, 128)
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		no_children := 0
		for _, child := range node.Children {
			if child != nil {
				no_children++
				dfs(child)
			}
		}
		childrenCount[no_children]++
	}
	dfs(st.GetRoot())
	return childrenCount
}
