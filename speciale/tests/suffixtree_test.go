// suffixtree_test.go
package suffixtreeimpl

import (
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
// Run tests using the `go test tests` command in your terminal

// TestSuffixTreeNodeCreation tests the creation of SuffixTreeNode instances.
func TestSuffixTreeNodeCreation(t *testing.T) {
	// Create a SuffixTreeNode instance
	node := &suffixtree.SuffixTreeNode{
		Label:    42,
		Parent:   nil,
		StartIdx: 0,
		EndIdx:   10,
		Children: make(map[rune]*suffixtree.SuffixTreeNode),
	}

	if node.Label != 42 {
		t.Errorf("Expected Label to be 42, got %d", node.Label)
	}
}

// Test creation of a NaiveSuffixTree instance
func TestNaiveSuffixTreeCreation(t *testing.T) {
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree("banana$")

	if st.GetRoot() == nil {
		t.Errorf("Expected root node to be non-nil")
	}

	if st.GetInputString() != "banana$" {
		t.Errorf("Expected input string to be 'banana', got %s", st.GetInputString())
	}
}

// Test that size of the suffix tree is correct
func TestNaiveSuffixTreeSize(t *testing.T) {
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree("abab$")

	if st.GetSize() != 8 {
		t.Errorf("Expected size to be 7, got %d", st.GetSize())
	}
}

// verify that we have n leaves
func TestNaiveSuffixTreeLeaves(t *testing.T) {
	// Create a NaiveSuffixTree instance
	const teststring string = "babababababababababababababbbbbbbbbbbbbbbbbbbbbb$"
	st := suffixtreeimpl.ConstructNaiveSuffixTree(teststring)

	leaves := 0
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if len(node.Children) == 0 {
			leaves++
		}
		for _, child := range node.Children {
			dfs(child)
		}
	}
	dfs(st.GetRoot())

	if leaves != len(teststring) {
		t.Errorf("Expected 4 leaves, got %d", leaves)
	}
}

// test that labels of leaves are 0,1,...,n-1
func TestNaiveSuffixTreeLeafLabels(t *testing.T) {
	// Create a NaiveSuffixTree instance
	const teststring string = "badsfkdsnfjkdsnvjkfndsvjkfnsjdvnfjkdsnvfkjsnfkjdsb$"
	st := suffixtreeimpl.ConstructNaiveSuffixTree(teststring)
	//create a set and add all the labels to it
	leafLabels := make(map[int]bool)
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if len(node.Children) == 0 {
			leafLabels[node.Label] = true
		}
		for _, child := range node.Children {
			dfs(child)
		}
	}
	dfs(st.GetRoot())
	//verify that all labels are in the set
	for i := 0; i < len(teststring); i++ {
		if !leafLabels[i] {
			t.Errorf("Expected label %d to be present", i)
		}
	}
}
