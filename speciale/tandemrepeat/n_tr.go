package tandemrepeat

import (
	"speciale/lce"
	"speciale/suffixtree"
)

// Decorate tree with vocabulary in O(n) time and return all tandem repeats in O(z) time
func DecorateTreeAndReturnTandemRepeats(tree suffixtree.SuffixTreeInterface) []TandemRepeat {
	DecorateTreeWithVocabulary(tree)
	return GetAllTandemRepeatsFromDecoratedTree(tree)
}

// Function that runs algorithm 1a,1b,2 and 3 on a suffix tree and decorates it with the tandem repeat vocabulary
func DecorateTreeWithVocabulary(tree suffixtree.SuffixTreeInterface) {

	// Phase 1
	// get leftmost covering repeats
	leftMostCoveringRepeats := Algorithm1(tree)

	//hacky way to get rid of tandem repeats and have ints instead.
	//Could be improved at a later point
	leftMostCoveringRepeatsInts := make([][]int, len(leftMostCoveringRepeats))
	for idx, k := range leftMostCoveringRepeats {
		leftMostCoveringRepeatsInts[idx] = make([]int, 0)
		for _, j := range k {
			leftMostCoveringRepeatsInts[idx] = append(leftMostCoveringRepeatsInts[idx], j.Length)

		}
	}

	// Phase 2
	// Decorate tree with subset of leftmost covering repeats
	Algorithm2StackMethod(tree, leftMostCoveringRepeatsInts)

	// Phase 3
	// Decorate tree with entire vocabulary
	Algorithm3StackMethod(tree)

}

// #######################################################################################
// #######################################################################################
// Algorithm 1
// #######################################################################################
// #######################################################################################

// Combines algorithm 1a and 1b to find tandem repeats
func Algorithm1(tree suffixtree.SuffixTreeInterface) [][]TandemRepeat {
	s := tree.GetInputString()
	leftMostCoveringRepeats := make([][]TandemRepeat, len(s))

	// intialize nested slice
	for i := range leftMostCoveringRepeats {
		leftMostCoveringRepeats[i] = make([]TandemRepeat, 0)
	}

	// Compute the blocks and Z-values
	li := LZDecompositionStackMethod(tree)
	blocks := CreateLZBlocks(li)

	// add idx to dfs table
	idxToDfsTable := getIdxtoDfsTableStackMethod(tree)

	// Do preprecessing for constant time LCE queries
	//lceObject := lce.PreProcessLCEBothDirections(tree)

	IterateBlocksAndExecuteAlgorithm1aAnd1b(tree, blocks, idxToDfsTable, &leftMostCoveringRepeats)

	return leftMostCoveringRepeats

}

// Phase 1, pt 1
// Compute the LZ decomposition of a string using a suffix tree
func LZDecomposition(tree suffixtree.SuffixTreeInterface) []int {

	// Initialize arrays to store the lengths of blocks and their starting positions
	n := len(tree.GetInputString())
	li := make([]int, n)

	// Perform a depth-first traversal of the suffix tree to compute the LZ decomposition
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {

		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			if node.Label < child.Label {
				li[child.Label] = depth
			}

			dfs(child, depth+child.EdgeLength())

		}
	}

	// Perform depth-first traversal starting from the root of the suffix tree
	dfs(tree.GetRoot(), 0)

	return li
}

func LZDecompositionStackMethod(tree suffixtree.SuffixTreeInterface) []int {

	// Initialize arrays to store the lengths of blocks and their starting positions
	n := len(tree.GetInputString())
	li := make([]int, n)

	// Perform a depth-first traversal of the suffix tree to compute the LZ decomposition
	stack := suffixtree.TreeStack{tree.GetRoot()}
	for len(stack) > 0 {
		node := stack.PopOrNil()
		if node == nil {
			continue
		}

		for _, child := range node.Children {
			if child == nil {
				continue
			}
			if node.Label < child.Label {
				li[child.Label] = node.StringDepth
			}
			if child.IsLeaf() {
				continue
			}
			stack.Push(child)
		}

	}

	return li

}

// create blocks from the LZ decomposition
func CreateLZBlocks(li []int) []int {
	n := len(li)

	//first block
	blocks := []int{0}
	iB := 0

	//recursively define remaining blocks
	for iB < n-1 {

		//next start
		iB += intMax(1, li[iB])
		blocks = append(blocks, iB)

	}
	return blocks
}

func IterateBlocksAndExecuteAlgorithm1aAnd1b(tree suffixtree.SuffixTreeInterface, blocks []int, idxToDfsTable []int, leftMostCoveringRepeats *[][]TandemRepeat) {
	s := tree.GetInputString()
	for i := 0; i < len(blocks); i++ {
		h := blocks[i]
		h1 := len(s)
		h2 := len(s)
		if i < len(blocks)-1 {
			h1 = blocks[i+1]
		}
		if i < len(blocks)-2 {
			h2 = blocks[i+2]
		}

		// Process block B for tandem repeats that satisfy condition 2
		Algorithm1b(s, h, h1, h2, leftMostCoveringRepeats)
		// Process block B for tandem repeats that satisfy condition 1
		Algorithm1a(s, h, h1, leftMostCoveringRepeats)

	}
}
func longestCommonExtension(s string, i, j int) int {
	n := len(s) - 1
	for i < n && j < n && s[i] == s[j] {
		i++
		j++
	}
	return i
}
func longestCommonExtensionBackwards(s string, i, j int) int {
	length := 0

	// Backward comparison
	for i >= 0 && j >= 0 && s[i] == s[j] {
		length++
		i--
		j--
	}

	return length
}
func Algorithm1a(s string, h int, h1 int, leftMostCoveringRepeats *[][]TandemRepeat) {
	for k := 1; k <= h1-h; k++ {
		q := h1 - k
		k1 := lce.FindLCEForwardSlow(s, q, h1)
		k2 := lce.FindLCEBackwardSlow(s, q-1, h1-1)
		start := intMax(q-k2, q-k+1)
		if k1+k2 >= k && k1 > 0 {
			addToLeftMostCoveringRepeats(leftMostCoveringRepeats, start, k)
		}
	}

}

func Algorithm1b(s string, h int, h1 int, h2 int, leftMostCoveringRepeats *[][]TandemRepeat) {
	for k := 1; k <= h2-h; k++ {
		q := h + k
		k1 := lce.FindLCEForwardSlow(s, q, h)
		k2 := lce.FindLCEBackwardSlow(s, q-1, h-1)
		start := intMax(h-k2, h-k+1)
		if k1+k2 >= k && k1 > 0 && start+k <= h1 && k2 > 0 {
			addToLeftMostCoveringRepeats(leftMostCoveringRepeats, start, k)
		}
	}

}

// add tandemrepeat to leftMostCoveringRepeats at index start if the last inserted tandem repeat at index start is not of same length
func addToLeftMostCoveringRepeats(leftMostCoveringRepeats *[][]TandemRepeat, start int, k int) {
	if len((*leftMostCoveringRepeats)[start]) == 0 {
		(*leftMostCoveringRepeats)[start] = append((*leftMostCoveringRepeats)[start], TandemRepeat{start, k, 2})
	} else if (*leftMostCoveringRepeats)[start][len((*leftMostCoveringRepeats)[start])-1].Length != k {
		(*leftMostCoveringRepeats)[start] = append((*leftMostCoveringRepeats)[start], TandemRepeat{start, k, 2})
	}
}

// return the maximum of two integers
func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// #######################################################################################
// #######################################################################################
// Algorithm 2
// #######################################################################################
// #######################################################################################
func Algorithm2(tree suffixtree.SuffixTreeInterface, leftMostCoveringRepeatsInts [][]int) {

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

		// Process the current node
		if !(tree.GetRoot() == node) {
			ProcessNodeAlg2(node, leftMostCoveringRepeatsInts, depth)
		}

	}

	// Perform depth-first traversal starting from the root of the suffix tree
	dfs(tree.GetRoot(), 0)

}

func Algorithm2StackMethod(tree suffixtree.SuffixTreeInterface, leftMostCoveringRepeatsInts [][]int) {

	//need two-way stack for bottom up traversal
	stack := suffixtree.Stack{&suffixtree.StackItem{Node: tree.GetRoot(), IsStart: true}}
	for len(stack) > 0 {
		item := stack.PopOrNil()
		node := item.Node

		// Traverse the children of the current node
		// need to traverse in reverse order
		if item.IsStart {
			item.IsStart = false
			stack.Push(item)

			for i := len(node.Children) - 1; i >= 0; i-- {
				if node.Children[i] == nil {
					continue
				}
				stack.Push(&suffixtree.StackItem{Node: node.Children[i], IsStart: true})
			}
		} else {
			//now we are going bottom up!
			//here we can process the node
			if !(tree.GetRoot() == node) {
				ProcessNodeAlg2(node, leftMostCoveringRepeatsInts, node.StringDepth)
			}
		}
	}
}

// Check if we can decorate this node with a tandem repeat from the leftmost covering set
func ProcessNodeAlg2(node *suffixtree.SuffixTreeNode, leftMostCoveringRepeatsInt [][]int, depth int) {
	//fmt.Println(node.Label, "hattemand", depth, depth-node.EdgeLength(), leftMostCoveringRepeatsInt[node.Label])

	//the label of any node (internal or leaf) is the smallest index in the subtree
	pList := leftMostCoveringRepeatsInt[node.Label]

	//if list is empty break
	if len(pList) == 0 {
		return
	}
	l := (pList)[len(pList)-1] * 2 // multiply by 2 to get entire repeat length \alpha \alpha

	parentDepth := depth - node.EdgeLength()

	for l > parentDepth {
		node.TandemRepeatDeco = append(node.TandemRepeatDeco, l-parentDepth)
		//remove the last element from the list
		pList = pList[:len(pList)-1]
		leftMostCoveringRepeatsInt[node.Label] = leftMostCoveringRepeatsInt[node.Label][:len(leftMostCoveringRepeatsInt[node.Label])-1]
		//if the list is empty, break
		if len(pList) == 0 {
			break
		}
		//update l
		l = pList[len(pList)-1] * 2

	}

}

// #######################################################################################
// #######################################################################################
// Algorithm 3
// #######################################################################################
// #######################################################################################
func Algorithm3(tree suffixtree.SuffixTreeInterface) {
	// Bottom-up traversal of the suffix tree
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			dfs(child, depth+child.EdgeLength())
		}

		// Process the current node
		if !(tree.GetRoot() == node) {
			if node.TandemRepeatDeco != nil {
				//attempt suffix walk
				for _, v := range node.TandemRepeatDeco {
					attemptSuffixWalk(tree, node, v)

				}
			}

		}

	}

	// Perform depth-first traversal starting from the root of the suffix tree
	dfs(tree.GetRoot(), 0)
}

func Algorithm3StackMethod(tree suffixtree.SuffixTreeInterface) {

	stack := suffixtree.Stack{&suffixtree.StackItem{Node: tree.GetRoot(), IsStart: true}}

	for len(stack) > 0 {
		item := stack.PopOrNil()
		node := item.Node

		// Traverse the children of the current node
		// need to traverse in reverse order
		if item.IsStart {
			item.IsStart = false
			stack.Push(item)

			for i := len(node.Children) - 1; i >= 0; i-- {
				if node.Children[i] == nil {
					continue
				}
				stack.Push(&suffixtree.StackItem{Node: node.Children[i], IsStart: true})
			}
		} else {
			//now we are going bottom up!
			//here we can process the node
			if !(tree.GetRoot() == node) {
				if node.TandemRepeatDeco != nil {
					//attempt suffix walk
					for _, v := range node.TandemRepeatDeco {
						attemptSuffixWalk(tree, node, v)

					}
				}

			}
		}
	}

}

func attemptSuffixWalk(st suffixtree.SuffixTreeInterface, node *suffixtree.SuffixTreeNode, tandemRepeatLengthOnEdge int) {
	//put this tandem repeat into the complete list
	if node.TandemRepeatDecoComplete == nil {
		node.TandemRepeatDecoComplete = make(map[int]bool)
	} else {
		_, ok := node.TandemRepeatDecoComplete[tandemRepeatLengthOnEdge]
		if ok {
			//case where tandem repeat is already in the list
			return
		}
	}
	node.TandemRepeatDecoComplete[tandemRepeatLengthOnEdge] = true

	//naming conventions according to the paper
	v := node
	u := v.Parent
	uMark := u.SuffixLink

	beta := tandemRepeatLengthOnEdge
	betaSum := 0 //used to keep track of full length of beta for case where we traverse multiple nodes
	char := st.GetInternalString()[v.StartIdx]
	vMark := uMark.Children[char]

	for beta > 0 {
		//case where beta ends in a node
		if vMark.EdgeLength() == beta {
			//check if child with "alpha" exists
			char := st.GetInternalString()[v.Label]
			if vMark.Children[char] != nil {

				beta = 1 // We just pass vMark, so beta is 1, as we just take a single step down next edge of vMarkMark. (v'')

				//check if tr already is marked
				if vMark.Children[char].TandemRepeatDecoComplete == nil {
					vMark.Children[char].TandemRepeatDecoComplete = make(map[int]bool)
				}
				_, ok := vMark.Children[char].TandemRepeatDecoComplete[beta]
				if ok {
					//case where tandem repeat is already in the list
					return
				}
				//add to complete list
				vMark.Children[char].TandemRepeatDecoComplete[beta] = true

				v = vMark.Children[char]
				u = v.Parent
				uMark = u.SuffixLink
				vMark = uMark.Children[char]
				betaSum = 0
				if uMark == st.GetRoot() {
					//if we reach the root, we are done
					return
				}

			} else {
				//failure - 'alpha' was not present in the subtree
				return
			}
		}
		if vMark.EdgeLength() > beta {
			//check if alpha is present in extension of beta
			if st.GetInternalString()[vMark.StartIdx+beta] == st.GetInternalString()[v.Label] {
				//success - continue suffix walk

				if vMark.TandemRepeatDecoComplete == nil {
					vMark.TandemRepeatDecoComplete = make(map[int]bool)
				}
				//check if tr already is marked
				_, ok := vMark.TandemRepeatDecoComplete[beta+1]
				if ok {
					//case where tandem repeat is already in the list
					return
				}
				//add to complete list
				vMark.TandemRepeatDecoComplete[beta+1] = true

				beta++ // Beta has now grown by one, as we have added \alpha
				v = vMark
				u = v.Parent
				uMark = u.SuffixLink
				vMark = uMark.Children[st.GetInternalString()[v.StartIdx]]
				betaSum = 0

			} else {
				//failure - alpha was not present in the extension of beta
				return
			}
		}
		if vMark.EdgeLength() < beta {
			// we need to fastscan further
			betaSum += vMark.EdgeLength()
			beta -= vMark.EdgeLength()
			vMark = vMark.Children[st.GetInternalString()[v.StartIdx+betaSum]] //should exist by construction
		}
	}
}

// #######################################################################################
// #######################################################################################
// Additional functions
// #######################################################################################
// #######################################################################################

// function to output all tandem repeats by traverseing subtrees and outputting labels
func GetAllTandemRepeatsFromDecoratedTree(tree suffixtree.SuffixTreeInterface) []TandemRepeat {

	tandemRepeats := make([]TandemRepeat, 0)

	//make a dfs label to idx mapping
	idxToDfsTable := getIdxtoDfsTableStackMethod(tree)
	//now reverse it
	dfsToIdxTable := make([]int, len(idxToDfsTable))
	for i, k := range idxToDfsTable {
		dfsToIdxTable[k] = i
	}

	//do a dfs
	var dfs func(node *suffixtree.SuffixTreeNode, depth int)
	dfs = func(node *suffixtree.SuffixTreeNode, depth int) {
		// Traverse the children of the current node
		for _, child := range node.Children {
			if child == nil {
				continue
			}
			dfs(child, depth+child.EdgeLength())
		}

		// Process the current node - should work for both internal and leaf nodes
		for k := range node.TandemRepeatDecoComplete {
			for leafDfs := node.DfsInterval.Start; leafDfs <= node.DfsInterval.End; leafDfs++ {
				leafIdx := dfsToIdxTable[leafDfs]
				tandemRepeats = append(tandemRepeats, TandemRepeat{leafIdx, (node.Parent.StringDepth + k) / 2, 2})
			}

		}

	}

	dfs(tree.GetRoot(), 0)

	return tandemRepeats

}
