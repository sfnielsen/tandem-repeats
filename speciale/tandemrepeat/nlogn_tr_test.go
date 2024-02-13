package tandemrepeat

import (
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
var (
	setupCompleted     bool
	alphabetGenerator  stringgenerators.StringGenerator
	randomGenerator    stringgenerators.StringGenerator
	fibonacciGenerator stringgenerators.StringGenerator
)

// addLeafList adds leaflists to the suffix tree
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

// test that we find the correct tandem repeats in a simple string
func TestFindTandemRepeatsLogarithmicVerySimpleExample(t *testing.T) {
	//define a simple string
	s := "babakk$"
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree(s)
	//find the tandem repeats
	tr := FindTandemRepeatsLogarithmic(st)

	if len(tr) != 2 {
		t.Errorf("Expected 2 tandem repeat, got %d", len(tr))
	}
	//extract the string from the tandem repeat and verify that it is correct ab$
	if GetTandemRepeatSubstring(tr[0], s) != "ba" && GetTandemRepeatSubstring(tr[1], s) != "ba" {
		t.Errorf("Expected tandem repeat to be 'ba'")
	}
	if GetTandemRepeatSubstring(tr[0], s) != "k" && GetTandemRepeatSubstring(tr[1], s) != "k" {
		t.Errorf("Expected tandem repeat to be 'k'")

	}
}

// test that we can find ALL tandem repeats
func TestFindTandemRepeatsLogarithmicSimpleExample(t *testing.T) {
	//generate some big strings from the stringgenerators
	s := randomGenerator.GenerateString(5000)

	// find tandem repeats with the naive_tr
	tr1 := FindTandemRepeatsNaive(s)

	tr2 := FindTandemRepeatsLogarithmic(suffixtreeimpl.ConstructNaiveSuffixTree(s))

	//check length are the same
	if len(tr1) != len(tr2) {
		t.Errorf("Expected %d tandem repeats, got %d", len(tr1), len(tr2))
	}

	//make 2 sets and check that they are equal
	//if not print the difference
	set1 := make(map[TandemRepeat]bool)
	for _, v := range tr1 {
		set1[v] = true
	}
	set2 := make(map[TandemRepeat]bool)
	for _, v := range tr2 {
		set2[v] = true
	}
	//not print difference if they are not equal
	if len(set1) != len(set2) {
		t.Errorf("Sets are not equal. Set1: %d, Set2: %d", len(set1), len(set2))
	}
	mismatches := 0
	for k := range set1 {
		if !set2[k] {
			t.Errorf("Set2 does not contain %v", k)
			mismatches++
		}
	}
	for k := range set2 {
		if !set1[k] {
			t.Errorf("Set1 does not contain %v", k)
			mismatches++
		}
	}
	if mismatches > 0 {
		t.Errorf("Total mismatches: %d", mismatches)

	}
}
