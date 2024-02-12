package tandemrepeat

import (
	"speciale/suffixtree"
)

func addLeafList(st suffixtree.SuffixTree) {
	root := st.GetRoot()

	//make a depth first traversal and collect all leafs that are acceisble through a node and add them to the set in the node:
	var dfs func(node *suffixtree.SuffixTreeNode) []int
	dfs = func(node *suffixtree.SuffixTreeNode) []int {
		// if leaf node
		if node.IsLeaf() {
			node.LeafList = []int{node.Label}
			return node.LeafList
		} else {
			var longest int

			for _, child := range node.Children {
				if child != nil {
					childLeafList := dfs(child)
					node.LeafList = append(node.LeafList, childLeafList...)

					if len(childLeafList) > longest {
						longest = len(childLeafList)
						node.BiggestChild = child
					}
				}
			}

		}
		return node.LeafList
	}

	dfs(root)
}
