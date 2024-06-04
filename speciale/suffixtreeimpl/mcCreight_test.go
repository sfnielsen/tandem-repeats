package suffixtreeimpl

// Test that size of the suffix tree has correct number of leaves
import (
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"testing"
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
			Alphabet: stringgenerators.AlphabetAB,
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

func TestSizeMc(t *testing.T) {
	st := ConstructMcCreightSuffixTree("abab$")
	if st.GetSize() != 8 {
		t.Errorf("Expected size to be 8, got %d", st.GetSize())
	}
}

// test that nodes in naive and mcCreight are the same
func TestMCSizeSameAsNaiveOnBigStrings(t *testing.T) {

	for i := 0; i < 1; i++ {
		//str := randomGenerator.GenerateString(1000)
		str := "abaabbabababbababaaababababbababbbbbabababaaaabaab$"

		st := ConstructMcCreightSuffixTree(str)
		st2 := ConstructNaiveSuffixTree(str)
		//cast st to mcCreight tree
		//PrintSuffixes(*st.(*McCreightSuffixTree))

		if st.GetSize() != st2.GetSize() {
			t.Errorf("Expected size to be %d, got %d", st2.GetSize(), st.GetSize())
		}

	}

}

// test that amount of leaves is correct (string size N produces N leaves)
func TestMcCreightSuffixTreeNLeaves(t *testing.T) {

	//check that we have n leaves on string size n
	for i := 0; i < 5; i++ {

		str := randomGenerator.GenerateString(1000)
		st := ConstructMcCreightSuffixTree(str)

		//count leaves
		leaves := 0
		var dfs func(node *suffixtree.SuffixTreeNode)
		dfs = func(node *suffixtree.SuffixTreeNode) {
			if node.IsLeaf() {
				leaves++
			} else {
				for _, child := range node.Children {
					if child != nil {
						dfs(child)
					}
				}
			}
		}
		dfs(st.GetRoot())

		if leaves != len(str) {
			t.Errorf("Expected %d leaves, got %d", len(str), leaves)
		}
	}

}

// test that we have the same tree as naive suffix tree
func TestMcCreightSuffixTreeSameAsNaive(t *testing.T) {

	str := randomGenerator.GenerateString(1000)
	st := ConstructMcCreightSuffixTree(str)
	st2 := ConstructNaiveSuffixTree(str)

	//compare trees by looking at topological structure
	//as well as start/end indices and labels
	var dfs func(node *suffixtree.SuffixTreeNode, node2 *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode, node2 *suffixtree.SuffixTreeNode) {
		if node.IsLeaf() {
			if !node2.IsLeaf() {
				t.Errorf("Expected leaf, got internal node")
			}
		} else {
			for _, child := range node.Children {
				if child != nil {
					if node2.Children[rune(st.GetInternalString()[child.StartIdx])] == nil {
						t.Errorf("Expected child, got nil")
					}
					if node2.Children[rune(st.GetInternalString()[child.StartIdx])].Label != child.Label {
						t.Errorf("Expected label %d, got %d", child.Label, node2.Children[rune(st.GetInternalString()[child.StartIdx])].Label)
					}
					if node2.Children[rune(st.GetInternalString()[child.StartIdx])].StartIdx != child.StartIdx {
						t.Errorf("Expected startIdx %d, got %d", child.StartIdx, node2.Children[rune(st.GetInternalString()[child.StartIdx])].StartIdx)
					}
					if node2.Children[rune(st.GetInternalString()[child.StartIdx])].EndIdx != child.EndIdx {
						t.Errorf("Expected endIdx %d, got %d", child.EndIdx, node2.Children[rune(st.GetInternalString()[child.StartIdx])].EndIdx)
					}
					dfs(child, node2.Children[rune(st.GetInternalString()[child.StartIdx])])
				}
			}
		}
	}
	dfs(st.GetRoot(), st2.GetRoot())
	dfs(st2.GetRoot(), st.GetRoot())

}

// verify that suffix links exist
func TestMcCreightSuffixLinksExist(t *testing.T) {
	str := randomGenerator.GenerateString(1000)
	st := ConstructMcCreightSuffixTree(str)

	//verify that all internal nodes have a suffix link
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		if !node.IsLeaf() {
			if node.SuffixLink == nil {
				t.Errorf("Expected suffix link, got nil")
			}
			for _, child := range node.Children {
				if child != nil {
					dfs(child)
				}
			}
		}
	}
	dfs(st.GetRoot())

}

// test that suffix links are correct, so that s(av) = v
func TestMcCreightSuffixLinksCorrect(t *testing.T) {
	str := randomGenerator.GenerateString(1000)
	st := ConstructMcCreightSuffixTree(str)

	//save string going down to nodes
	var dfs func(node *suffixtree.SuffixTreeNode, suffix string)
	dfs = func(node *suffixtree.SuffixTreeNode, suffix string) {
		if !node.IsLeaf() {

			//check if there is a suffix link
			if node.SuffixLink == nil {
				t.Errorf("Expected suffix link, got nil")
			}

			for _, child := range node.Children {
				if child != nil {
					//now check that string on suffix link is the same as the string on the node

					sl := node.SuffixLink
					if sl != st.GetRoot() {
						suffixL := st.GetInputString()[sl.StartIdx : sl.EndIdx+1]
						parent := sl.Parent

						//backtrack till root
						for parent != st.GetRoot() {
							suffixL = st.GetInputString()[parent.StartIdx:parent.EndIdx+1] + suffixL
							parent = parent.Parent
						}
						if suffix[1:] != suffixL {
							t.Errorf("Expected suffix %s, got %s", suffix[1:], suffixL)
						}

					}
					dfs(child, suffix+st.GetInputString()[child.StartIdx:child.EndIdx+1])

				}

			}
		} else {
			if node.SuffixLink != nil {
				if st.GetInputString()[node.StartIdx] != st.GetInputString()[node.SuffixLink.StartIdx] {
					t.Errorf("Leaf should not have a suffix link")
				}
			}
		}
	}
	dfs(st.GetRoot(), "")

}
