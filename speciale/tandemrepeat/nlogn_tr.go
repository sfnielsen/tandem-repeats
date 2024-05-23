package tandemrepeat

import (
	"speciale/suffixtree"
)

// get all tandem repeats by left rotating on the branching repeats
func GetAllTandemRepeats(allBranchingRepeats []TandemRepeat, st suffixtree.SuffixTreeInterface) []TandemRepeat {
	var allTandemRepeats = make([]TandemRepeat, 0)

	for _, k := range allBranchingRepeats {
		// add tandem repeat until length is 0
		i := 0
		// left rotate until we no longer have a tandem repeat (or we reach the start of the string)
		for k.Start-i-1 >= 0 {
			i += 1
			if st.GetInternalString()[k.Start-i] == st.GetInternalString()[(k.Start-i)+2*(k.Length)] {
				allTandemRepeats = append(allTandemRepeats, TandemRepeat{k.Start - i, k.Length, 2})
			} else {
				break
			}
		}

	}
	allTandemRepeats = append(allTandemRepeats, allBranchingRepeats...)
	return allTandemRepeats
}

// GetTandemRepeatSubstring returns the substring of the tandem repeat
func getIdxtoDfsTable(st suffixtree.SuffixTreeInterface) []int {
	//create table
	var idxToDfsTable []int = make([]int, len(st.GetInternalString()))

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

func getIdxtoDfsTableStackMethod(st suffixtree.SuffixTreeInterface) []int {
	var idxToDfsTable []int = make([]int, len(st.GetInternalString()))

	stack := suffixtree.TreeStack{st.GetRoot()}
	for len(stack) > 0 {
		node := stack.PopOrNil()
		if node.IsLeaf() {
			idxToDfsTable[node.Label] = node.DfsInterval.Start
		} else {
			for _, child := range node.Children {
				if child != nil {
					stack.Push(child)
				}
			}
		}
	}

	return idxToDfsTable
}

// FindAllTandemRepeatsLogarithmic finds tandem repeats in a suffix tree in O(nlogn + z) time
func FindAllTandemRepeatsLogarithmic(st suffixtree.SuffixTreeInterface) []TandemRepeat {
	//find all branching repeats in O(nlogn) time
	trBranching := FindAllBranchingTandemRepeatsLogarithmic(st)

	//get all tandem repeats by left rotating on the branching repeats
	//this is O(z) time (up to O(n^2))
	allTandemRepeats := GetAllTandemRepeats(trBranching, st)
	return allTandemRepeats

}

func reverseSlice(slice []int) []int {
	reverse := make([]int, len(slice))
	for i, v := range slice {
		reverse[v] = i
	}
	return reverse
}

// FindTandemRepeatsLogarithmic finds tandem repeats in a suffix tree in O(nlogn) time
func FindAllBranchingTandemRepeatsLogarithmic(st suffixtree.SuffixTreeInterface) []TandemRepeat {

	//we create the a idx to dfs mapping
	idxToDfsTable := getIdxtoDfsTableStackMethod(st)

	//create Dfs to idx mapping, this is an alternative to leaf lists
	dfsToIdxTable := reverseSlice(idxToDfsTable)

	//add biggest child to each node
	st.AddBiggestChildToNodes()

	//store all branching repeats slice
	allBranchingRepeats := make([]TandemRepeat, 0)

	// now we run stoye and gusfield 'optimized algorithm'
	var dfs func(node *suffixtree.SuffixTreeNode, depth int) []int
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) []int {
		depth = depth + node.EdgeLength()

		leafList := []int{} // leaflist dynamically added to the node

		for _, child := range node.Children {
			if child == nil {
				continue
			}

			// iterate over elements from leaflistPrime
			//step 2a is performed implicitly by traversal of the children (minus the biggest child)
			if node.BiggestChild != child {

				//iterate over all leafs in leaflistPrime
				for dfsNumber := child.DfsInterval.Start; dfsNumber <= child.DfsInterval.End; dfsNumber++ {

					//step 2b
					i := dfsToIdxTable[dfsNumber]
					j := i + depth
					if j < len(st.GetInternalString()) {
						dfsVal := idxToDfsTable[j]
						//check if j is in dfs interval
						if dfsVal >= node.DfsInterval.Start && dfsVal <= node.DfsInterval.End {

							//now check if we are branching
							if i+2*depth < len(st.GetInternalString()) && st.GetInternalString()[i] != st.GetInternalString()[i+2*depth] {
								//we have a branching repeat, add to slice
								allBranchingRepeats = append(allBranchingRepeats, TandemRepeat{i, depth, 2})
							}
						}
					}

					//step 2c
					j = dfsToIdxTable[dfsNumber]
					i = j - depth
					if i >= 0 && i < len(st.GetInternalString()) {
						dfsVal := idxToDfsTable[i]
						//check if i is in dfs interval
						//this check is simplified compared to the paper
						//in the paper we check if i is in LL(v), but we only need to check that it is in LL(v')
						if dfsVal >= node.BiggestChild.DfsInterval.Start && dfsVal <= node.BiggestChild.DfsInterval.End {

							//now check if we are branching
							if i+2*depth < len(st.GetInternalString()) && st.GetInternalString()[i] != st.GetInternalString()[i+2*depth] {
								//we have a branching repeat
								allBranchingRepeats = append(allBranchingRepeats, TandemRepeat{i, depth, 2})
							}
						}

					}
				}
			}

			// step 1, marking internal nodes is done implicitly by a depth-first traversal
			if !child.IsLeaf() {
				dfs(child, depth)
			}
		}

		//case for internal nodes
		return leafList

	}

	dfs(st.GetRoot(), 0)

	return allBranchingRepeats
}
