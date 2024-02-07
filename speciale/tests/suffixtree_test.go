// suffixtree_test.go
package suffixtreeimpl

import (
	"speciale/suffixtree"
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
		Children: []*suffixtree.SuffixTreeNode{},
	}

	if node.Label != 42 {
		t.Errorf("Expected Label to be 42, got %d", node.Label)
	}
}
