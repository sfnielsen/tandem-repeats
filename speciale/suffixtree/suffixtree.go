package suffixtree

type SuffixTreeNode struct {
	Label    int
	Parent   *SuffixTreeNode
	StartIdx int
	EndIdx   int
	Children []*SuffixTreeNode
}

// SuffixTree defines the interface for a suffix tree.
type SuffixTree interface {
	GetRoot() *SuffixTreeNode
	SetRoot(root *SuffixTreeNode)

	// Function to get the length of an edge
	// Parameters: node whose edge length is to be calculated
	// Returns: length of the edge (inclusive of start and end indices)
	EdgeLength(node *SuffixTreeNode) int

	// Function to split an edge
	SplitEdge(originalCihld *SuffixTreeNode, startIdx int, splitIdx int, endIdx int, inputString string, suffixOffest int)

	// Function to insert a suffix into the suffix tree
	InsertSuffix(str string, suffixOffest int, root *SuffixTreeNode, inputString string)

	// Function to create a suffix tree from the given input string.
	// Constructs a suffix tree data structure using Ukkonen's algorithm.
	// Parameters: input string for which the suffix tree is to be constructed.
	// Returns: A pointer to the root node of the constructed suffix tree.
	CreateSuffixTree(inputString string)

	// Function to search for a substring in the suffix tree
	// Should write the indices of the occurrences of the substring in the input string to the standard output
	// in standard cigar format
	SearchSubstring(substring string) []int
}
