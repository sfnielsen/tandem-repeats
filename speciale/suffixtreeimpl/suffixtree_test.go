package suffixtreeimpl

import (
	"sort"
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
// Run tests using the `go test tests` command in your terminal

var (
	setupCompleted     bool
	alphabetGenerator  stringgenerators.StringGenerator
	randomGenerator    stringgenerators.StringGenerator
	fibonacciGenerator stringgenerators.StringGenerator
)

func init() {
	if !setupCompleted {
		// Perform setup steps here
		// ...
		alphaGen := &stringgenerators.AlphabetStringGenerator{
			Alphabet: stringgenerators.AlphabetDNA,
		}
		alphabetGenerator = alphaGen

		randGen := &stringgenerators.RandomStringGenerator{
			Alphabet: stringgenerators.AlphabetDNA,
		}
		randomGenerator = randGen

		fibGen := &stringgenerators.FibonacciStringGenerator{
			First:  "b",
			Second: "a",
		}
		fibonacciGenerator = fibGen

		setupCompleted = true
	}
}

// TestSuffixTreeNodeCreation tests the creation of SuffixTreeNode instances.
func TestSuffixTreeNodeCreation(t *testing.T) {
	// Create a SuffixTreeNode instance
	node := &suffixtree.SuffixTreeNode{
		Label:    42,
		Parent:   nil,
		StartIdx: 0,
		EndIdx:   10,
	}

	if node.Label != 42 {
		t.Errorf("Expected Label to be 42, got %d", node.Label)
	}
}

// Test creation of a NaiveSuffixTree instance
func TestNaiveSuffixTreeCreationDoesntFail(t *testing.T) {
	//generate some random string
	test_str := randomGenerator.GenerateString(100000)

	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree(test_str)

	if st.GetRoot() == nil {
		t.Errorf("Expected root node to be non-nil")
	}

	if st.GetInputString() != test_str {
		t.Errorf("Expected input string to be '"+test_str+"', got %s", st.GetInputString())
	}
}

// Test that size of the suffix tree has correct number of leaves
func TestNaiveSuffixTreeSizeSmall(t *testing.T) {
	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree("abab$")

	if st.GetSize() != 8 {
		t.Errorf("Expected size to be 7, got %d", st.GetSize())
	}
}

func TestNaiveSuffixTreeSizeLarge(t *testing.T) {

	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree("cabacbabbacaccabbababababbaaabababababaabacbaababababababaabacbaabaccbabcbba$")

	//in this example there is 77 leaves and 60 internal nodes
	if st.GetSize() != 77+60 {
		t.Errorf("Expected size to be 100, got %d", st.GetSize())
	}

}

// verify that we have n leaves
func TestNaiveSuffixTreeLeavesSizeIsN(t *testing.T) {
	//generate some random string
	test_str := randomGenerator.GenerateString(1000)

	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree(test_str)

	leaves := 0
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if node.IsLeaf() {
			for _, child := range node.Children {
				dfs(child)
			}
		}
		dfs(st.GetRoot())

		if leaves != len(test_str) {
			t.Errorf("Expected 4 leaves, got %d", leaves)
		}
	}
}

// test that labels of leaves are 0,1,...,n-1
func TestNaiveSuffixTreeLeafLabels(t *testing.T) {
	//generate some random string
	test_str := randomGenerator.GenerateString(1000)

	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree(test_str)
	//create a set and add all the labels to it
	leafLabels := make(map[int]bool)
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if node.IsLeaf() {
			leafLabels[node.Label] = true
		}
		for _, child := range node.Children {
			if child != nil {
				dfs(child)
			}
		}
	}
	dfs(st.GetRoot())
	//verify that all labels are in the set
	for i := 0; i < len(test_str); i++ {
		if !leafLabels[i] {
			t.Errorf("Expected label %d to be present", i)
		}
	}
}

// verify that path down to leaf is the actual suffix
func TestNaiveSuffixTreeSuffixes(t *testing.T) {
	// Create a NaiveSuffixTree instance
	fibonacciGenerator := &stringgenerators.FibonacciStringGenerator{
		First:  "b",
		Second: "a",
	}
	var _ stringgenerators.StringGenerator = fibonacciGenerator
	var teststring string = fibonacciGenerator.GenerateString(20)

	st := ConstructNaiveSuffixTree(teststring)
	//walk a path down to a leaf and verify that it is the suffix
	var dfs func(node *suffixtree.SuffixTreeNode, suffix string)
	dfs = func(node *suffixtree.SuffixTreeNode, suffix string) {
		//guard to not check the root
		if (node.StartIdx != -1) && (node.EndIdx != -1) {
			//if we are at a leaf, verify that the suffix is correct
			if node.IsLeaf() {
				if suffix != teststring[node.Label:] {
					t.Errorf("Expected suffix %s, got %s", teststring[node.StartIdx:], suffix)
				}
			}
		}
		for _, child := range node.Children {
			if child != nil {
				dfs(child, suffix+teststring[child.StartIdx:child.EndIdx+1])
			}
		}
	}
	dfs(st.GetRoot(), "")

}

func TestNaiveSuffixTreeOnMultipleFibonnaciStrings(t *testing.T) {
	fibonacciString := fibonacciGenerator.GenerateString(10)

	fibonacciGenerator.GenerateString(10)

	st := ConstructNaiveSuffixTree(fibonacciString)
	//walk a path down to a leaf and verify that it is the suffix
	var dfs func(node *suffixtree.SuffixTreeNode, suffix string)
	dfs = func(node *suffixtree.SuffixTreeNode, suffix string) {
		//guard to not check the root
		if (node.StartIdx != -1) && (node.EndIdx != -1) {
			//if we are at a leaf, verify that the suffix is correct
			if node.IsLeaf() {
				if suffix != fibonacciString[node.Label:] {
					t.Errorf("Expected suffix %s, got %s", fibonacciString[node.StartIdx:], suffix)
				}
			}
		}
		for _, child := range node.Children {
			if child != nil {
				dfs(child, suffix+fibonacciString[child.StartIdx:child.EndIdx+1])
			}
		}
	}
	dfs(st.GetRoot(), "")
}

func TestNaiveSuffixTreeOnMultipleRandomStringTypes(t *testing.T) {
	var teststring_slice []string = stringgenerators.GenerateStringArray(100, 20, []stringgenerators.StringGenerator{alphabetGenerator, randomGenerator})

	for _, teststring := range teststring_slice {
		st := ConstructNaiveSuffixTree(teststring)
		//walk a path down to a leaf and verify that it is the suffix
		var dfs func(node *suffixtree.SuffixTreeNode, suffix string)
		dfs = func(node *suffixtree.SuffixTreeNode, suffix string) {
			//guard to not check the root
			if (node.StartIdx != -1) && (node.EndIdx != -1) {
				//if we are at a leaf, verify that the suffix is correct
				if node.IsLeaf() {
					if suffix != teststring[node.Label:] {
						t.Errorf("Expected suffix %s, got %s", teststring[node.StartIdx:], suffix)
					}
				}
			}
			for _, child := range node.Children {
				if child != nil {
					dfs(child, suffix+teststring[child.StartIdx:child.EndIdx+1])
				}
			}
		}
		dfs(st.GetRoot(), "")
	}
}

// test That dfs labels on leaves are lexicographically ordered
func TestNaiveSuffixTreeDfsLeafLabelsIsLexicographic(t *testing.T) {
	// Create a NaiveSuffixTree instance
	str := randomGenerator.GenerateString(1000)
	st := ConstructNaiveSuffixTree(str)
	//first create a slice of all suffixes of str
	suffixes := make([]string, len(str))
	for i := 0; i < len(str); i++ {
		suffixes[i] = str[i:]
	}
	//then sort the slice
	sort.Strings(suffixes)

	// then verify that the dfs labels are in the same order
	var dfs func(node *suffixtree.SuffixTreeNode)
	var idx int
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if node.IsLeaf() {
			if suffixes[idx] != str[node.Label:] {
				t.Errorf("Expected suffix %s, got %s", suffixes[idx], str[node.Label:])
			}
			idx += 1
		}
		for _, child := range node.Children {
			if child != nil {
				dfs(child)
			}
		}
	}
	dfs(st.GetRoot())

}

// Test that dfs intervals corrospond to the range of the children of the node
func TestNaiveSuffixTreeDfsIntervals(t *testing.T) {
	// Create a NaiveSuffixTree instance
	st := ConstructNaiveSuffixTree(alphabetGenerator.GenerateString(1000))

	//check that each internal node has a start and end corrosponding to the range of its children
	var dfs func(node *suffixtree.SuffixTreeNode) (int, int)
	dfs = func(node *suffixtree.SuffixTreeNode) (int, int) {

		// case for leaf node
		if node.IsLeaf() {
			if node.DfsInterval.Start != node.DfsInterval.End {
				t.Errorf("Expected start and end to be equal, got %d and %d", node.DfsInterval.Start, node.DfsInterval.End)
			}
			//return value
			return node.DfsInterval.Start, node.DfsInterval.End
		}

		// keep track on smallest and largest value found in children
		smallestDfs := st.GetSize() + 1
		largestDfs := -1
		for _, child := range node.Children {
			if child != nil {
				childMin, childMax := dfs(child)
				if childMin < smallestDfs {
					smallestDfs = childMin
				}
				if childMax > largestDfs {
					largestDfs = childMax
				}
			}
		}

		//verify that the start and end is correct
		if node.DfsInterval.Start != smallestDfs {
			t.Errorf("Expected start to be %d, got %d", smallestDfs, node.DfsInterval.Start)
		}
		if node.DfsInterval.End != largestDfs {
			t.Errorf("Expected end to be %d, got %d", largestDfs, node.DfsInterval.End)
		}

		//return value
		return node.DfsInterval.Start, node.DfsInterval.End
	}
	dfs(st.GetRoot())

}
