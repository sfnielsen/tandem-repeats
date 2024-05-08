package suffixtree

// TreeStack is a stack of SuffixTreeNodes only for top down traversal
type TreeStack []*SuffixTreeNode

func (t *TreeStack) Push(node *SuffixTreeNode) {
	*t = append(*t, node)
}

func (t *TreeStack) PopOrNil() *SuffixTreeNode {
	if len(*t) == 0 {
		return nil
	}
	node := (*t)[len(*t)-1]
	*t = (*t)[:len(*t)-1]
	return node
}

// StackItem is a struct that holds a node and a flag in order to go back during dfs run
type StackItem struct {
	Node    *SuffixTreeNode
	IsStart bool // Flag indicating whether it's the start index of DFS numbering
}

type Stack []*StackItem

func (s *Stack) Push(item *StackItem) {
	*s = append(*s, item)
}

func (s *Stack) PopOrNil() *StackItem {
	if len(*s) == 0 {
		return nil
	}
	lastIndex := len(*s) - 1
	popped := (*s)[lastIndex]
	*s = (*s)[:lastIndex]
	return popped
}
