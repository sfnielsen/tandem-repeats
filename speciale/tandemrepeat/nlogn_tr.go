package tandemrepeat

import (
	"speciale/suffixtree"
)

// addLeafList adds leaflists to the suffix tree
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
			// case for internal node
			var longest int
			for _, child := range node.Children {
				if child != nil {
					//get the leaflist of the child
					childLeafList := dfs(child)
					node.LeafList = append(node.LeafList, childLeafList...)
					//keep track of the biggest child found
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

// get all tandem repeats
func getAllTandemRepeats(allBranchingRepeats []TandemRepeat, st suffixtree.SuffixTree) []TandemRepeat {
	var allTandemRepeats []TandemRepeat = make([]TandemRepeat, 0)

	for _, v := range allBranchingRepeats {
		// add tandem repeat until length is 0
		i := 0
		// left rotate until we no longer have a tandem repeat (or we reach the start of the string)
		for v.Start-i-1 >= 0 {
			i += 1
			if st.GetInputString()[v.Start-i] == st.GetInputString()[(v.Start-i)+2*(v.length)] {
				allTandemRepeats = append(allTandemRepeats, TandemRepeat{v.Start - i, v.length, 2})
			} else {
				break
			}
		}

	}
	allTandemRepeats = append(allTandemRepeats, allBranchingRepeats...)
	return allTandemRepeats
}

// GetTandemRepeatSubstring returns the substring of the tandem repeat
func getIdxtoDfsTable(st suffixtree.SuffixTree) []int {
	//create table
	var idxToDfsTable []int = make([]int, len(st.GetInputString()))

	//fill table with another dfs...
	//this can be done during construction or in another dfs...
	//but we can optimize it later
	var dfs func(node *suffixtree.SuffixTreeNode) int
	dfs = func(node *suffixtree.SuffixTreeNode) int {
		if node.IsLeaf() {
			idxToDfsTable[node.Label] = node.DfsInterval.Start
			return node.DfsInterval.Start
		} else {
			for _, child := range node.Children {
				if child != nil {
					dfs(child)
				}
			}
		}
		return 0
	}
	dfs(st.GetRoot())

	return idxToDfsTable
}

// FindTandemRepeatsLogarithmic finds tandem repeats in a suffix tree in O(nlogn) time
func FindTandemRepeatsLogarithmic(st suffixtree.SuffixTree) []TandemRepeat {
	//first we need to add the leaflist to the suffix tree
	addLeafList(st)

	//we create the a idx to dfs mapping
	idxToDfsTable := getIdxtoDfsTable(st)

	//store all branching repeats
	var allBranchingRepeats []TandemRepeat = make([]TandemRepeat, 0)

	// now we run stoye and gusfield 'optimized algorithm'
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		depth = depth + node.EdgeLength()

		for _, child := range node.Children {
			if child == nil {
				continue
			}
			// iterate over elements from leaflistPrime
			//step 2a is performed implicitly by traversal of the children (minus the biggest child)
			if node.BiggestChild != child {

				for _, leaf := range child.LeafList {

					//step 2b
					i := leaf
					j := i + depth
					if j < len(st.GetInputString()) {
						dfsVal := idxToDfsTable[j]
						//check if j is in dfs interval
						if dfsVal >= node.DfsInterval.Start && dfsVal <= node.DfsInterval.End {

							//now check if we are branching
							if i+2*depth < len(st.GetInputString()) && st.GetInputString()[i] != st.GetInputString()[i+2*depth] {
								//we have a branching repeat
								allBranchingRepeats = append(allBranchingRepeats, TandemRepeat{leaf, depth, 2})
							}
						}
					}
					//step 2c
					j = leaf
					i = j - depth
					if i >= 0 && i < len(st.GetInputString()) {
						dfsVal := idxToDfsTable[i]
						//check if i is in dfs interval
						if dfsVal >= node.DfsInterval.Start && dfsVal <= node.DfsInterval.End {

							//now check if we are branching
							if i+2*depth < len(st.GetInputString()) && st.GetInputString()[i] != st.GetInputString()[i+2*depth] {
								//we have a branching repeat
								allBranchingRepeats = append(allBranchingRepeats, TandemRepeat{i, depth, 2})
							}
						}
					}
				}
			}
			// step 1, marking internal nodes is done implicitly by a depth-first traversal
			dfs(child, depth)
		}
	}
	dfs(st.GetRoot(), 0)

	return getAllTandemRepeats(allBranchingRepeats, st)
}
