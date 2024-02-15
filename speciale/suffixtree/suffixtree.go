package suffixtree

// SuffixTree defines the interface for a suffix tree.
type SuffixTreeInterface interface {
	// Getters for the fields of the suffix tree
	ConstructSuffixTree()
	GetRoot() *SuffixTreeNode
	GetInputString() string
	GetSize() int
	IncrementSize()
	AddBiggestChildToNodes()
	//TODO
	//PrintTree()

	// Function to search for a substring in the suffix tree
	// Should write the indices of the occurrences of the substring in the input string to the standard output
	// in standard cigar format
	//SearchSubstring(substring string) []int

}

type SuffixTree struct {
	Root        *SuffixTreeNode
	InputString string
	Size        int
}

// GetRoot returns the root node of the suffix tree.
func (n *SuffixTree) GetRoot() *SuffixTreeNode {
	return n.Root
}

// GetInputString returns the input string used to construct the suffix tree.
func (n *SuffixTree) GetInputString() string {
	return n.InputString
}

// GetSize returns the size of the suffix tree.
func (n *SuffixTree) GetSize() int {
	return n.Size
}

// GetSize returns the size of the suffix tree.
func (n *SuffixTree) IncrementSize() {
	n.Size++
}

func (n *SuffixTree) AddBiggestChildToNodes() {
	return
}
