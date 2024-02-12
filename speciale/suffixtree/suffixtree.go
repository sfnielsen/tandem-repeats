package suffixtree

// SuffixTree defines the interface for a suffix tree.
type SuffixTree interface {
	// Getters for the fields of the suffix tree
	GetRoot() *SuffixTreeNode
	GetInputString() string
	GetSize() int
	//TODO
	//PrintTree()

	// Function to search for a substring in the suffix tree
	// Should write the indices of the occurrences of the substring in the input string to the standard output
	// in standard cigar format
	//SearchSubstring(substring string) []int

}
