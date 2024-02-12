package tandemrepeat

import (
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
)

func TestBiggestChildIsBiggestChild(t *testing.T) {

	st := suffixtreeimpl.ConstructNaiveSuffixTree("ababababasdfasdfbasdfasdfawsdfasvaertbnaerfsdfasdf$")
	addLeafList(st)

	//walk a path down to a leaf and verify that it is the suffix
	var dfs func(node *suffixtree.SuffixTreeNode) int
	dfs = func(node *suffixtree.SuffixTreeNode) int {
		if node.IsLeaf() {
			if len(node.LeafList) != 1 {
				t.Errorf("Expected size of leaflist to be 1, got %d", len(node.LeafList))
			}
			return len(node.LeafList)
		} else {
			var longest int
			var biggestFoundChild *suffixtree.SuffixTreeNode = nil
			for _, child := range node.Children {
				if child != nil {
					childLeafList := dfs(child)

					if childLeafList > longest {
						longest = childLeafList
						biggestFoundChild = child
					}
				}
			}
			if longest != len(node.BiggestChild.LeafList) {
				t.Errorf("Leaflist length of biggest child is not the same, expected %d but found: %d", len(node.BiggestChild.LeafList), longest)
			}

			if biggestFoundChild.Label != node.BiggestChild.Label {
				t.Errorf("Biggest child found differs from biggest child in tree, expected %d but found: %d", node.BiggestChild.Label, biggestFoundChild.Label)
			}

			return len(node.LeafList)

		}
	}
	dfs(st.GetRoot())
}
