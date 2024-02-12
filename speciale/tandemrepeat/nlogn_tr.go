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

func FindTandemRepeatsLogarithmic(st suffixtree.SuffixTree) []TandemRepeat {
	addLeafList(st)
	var allBranchingRepeats []TandemRepeat = make([]TandemRepeat, 0)

	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		depth = depth + node.EdgeLength()

		for _, child := range node.Children {
			if child != nil || !child.IsLeaf() {
				// step 1, marking internal nodes is done implicitly by a depth-first traversal
				dfs(child, depth)

				// step 2a, find the set LL'(v) which is leaf list minus the leaf list of the biggest child
				var leafListPrime []int
				for _, child := range node.Children {
					if child == nil || child.Label == node.BiggestChild.Label {
						continue
					}
					leafListPrime = append(leafListPrime, child.LeafList...)
				}

				// step 2b
				//for _, child := range node.Children {
				//j := child.Label + depth

				//}

			}
		}
	}
	dfs(st.GetRoot(), 0)

	return allBranchingRepeats
}
