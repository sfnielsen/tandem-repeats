package tandemrepeat

import (
	"log"
	"os"
	"runtime/pprof"
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
var (
	setupCompleted          bool
	randomGenerator_protein stringgenerators.StringGenerator
	randomGenerator_fib     stringgenerators.StringGenerator
	randomGenerator_ab      stringgenerators.StringGenerator
	randomGenerator_dna     stringgenerators.StringGenerator
	randomGenerator_byte    stringgenerators.StringGenerator
	randomGenerator_a       stringgenerators.StringGenerator
)

// addLeafList adds leaflists to the suffix tree
func init() {
	if !setupCompleted {
		randomGenerator_protein = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetProtein}
		randomGenerator_ab = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetAB}
		randomGenerator_dna = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetDNA}
		randomGenerator_byte = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetByte}
		randomGenerator_a = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetA}
		randomGenerator_fib = &stringgenerators.FibonacciStringGenerator{}
		setupCompleted = true
	}
}

// test that we find the correct tandem repeats in a simple string
func TestFindTandemRepeatsLogarithmicVerySimpleExample(t *testing.T) {
	//define a simple string
	s := "babakk$"
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree(s)
	//find the tandem repeats
	tr := FindAllTandemRepeatsLogarithmic(st)

	if len(tr) != 2 {
		t.Errorf("Expected 2 tandem repeat, got %d", len(tr))
	}
	//extract the string from the tandem repeat and verify that it is correct ab$
	if GetTandemRepeatSubstring(tr[0], s) != "baba" && GetTandemRepeatSubstring(tr[1], s) != "ba" {
		t.Errorf("Expected tandem repeat to be 'ba'")
	}
	if GetTandemRepeatSubstring(tr[0], s) != "k" && GetTandemRepeatSubstring(tr[1], s) != "kk" {
		t.Errorf("Expected tandem repeat to be 'k'")

	}
}

// test that we can find ALL tandem repeats
func TestFindTandemRepeatsLogarithmicSimpleExample(t *testing.T) {
	//generate some big strings from the stringgenerators

	s := randomGenerator_a.GenerateString(1000)

	// find tandem repeats with the naive_tr
	tr1 := FindTandemRepeatsNaive(s)

	tr2 := FindAllTandemRepeatsLogarithmic(suffixtreeimpl.ConstructNaiveSuffixTree(s))

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

func TestFindTandemRepeatsLogarithmicMultipleStringtypes(t *testing.T) {
	//generate some big strings from the stringgenerators
	stringArrays := stringgenerators.GenerateStringArray(100, 1200, []stringgenerators.StringGenerator{randomGenerator_protein, randomGenerator_ab, randomGenerator_dna, randomGenerator_byte})

	for _, s := range stringArrays {
		// find tandem repeats with the naive_tr
		tr1 := FindTandemRepeatsNaive(s)

		tr2 := FindAllTandemRepeatsLogarithmic(suffixtreeimpl.ConstructNaiveSuffixTree(s))

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
}

// test that biggest child is assigned correctly
func TestNaiveTrBiggestChild(t *testing.T) {
	s := randomGenerator_ab.GenerateString(100)
	// Create a NaiveTandemRepeat instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree(s)
	st.AddBiggestChildToNodes()
	// check that biggest child is not nil for all nodes
	var dfs func(node *suffixtree.SuffixTreeNode)
	dfs = func(node *suffixtree.SuffixTreeNode) {
		// test children
		if node.IsLeaf() {
			if node.BiggestChild != nil {
				// test leaf nodes is nil
				t.Errorf("Expected biggest child to be nil, got %v", node.BiggestChild)
			}
		} else {
			// test internal nodes is not nil
			if node.BiggestChild == nil {
				t.Errorf("Expected biggest child to be non-nil for internal node, got nil")
			}
			// also verify that biggestchild is the biggest child
			var longest int
			var currentBiggestChild *suffixtree.SuffixTreeNode = nil
			for _, child := range node.Children {
				if child != nil {
					dfs(child)
					if child.DfsInterval.End-child.DfsInterval.Start+1 > longest {
						longest = child.DfsInterval.End - child.DfsInterval.Start + 1
						currentBiggestChild = child
					}
				}
			}
			if currentBiggestChild.Label != node.BiggestChild.Label {
				t.Errorf("Biggest child found differs from biggest child in tree, expected %d but found: %d", node.BiggestChild.Label, currentBiggestChild.Label)
			}
		}
	}
	dfs(st.GetRoot())
}

// test that we find same TR with McCreight and Naive
func TestMcCreightVsNaive(t *testing.T) {
	s := randomGenerator_dna.GenerateString(1000)
	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)
	st2 := suffixtreeimpl.ConstructNaiveSuffixTree(s)

	tr1 := FindAllTandemRepeatsLogarithmic(st)
	tr2 := FindAllTandemRepeatsLogarithmic(st2)

	//check that sets are equal
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

// #####################################################################################
// #####################################################################################
// Benchmarking
// #####################################################################################
// #####################################################################################
func BenchmarkBranchingTR(b *testing.B) {
	trees := make([]suffixtree.SuffixTreeInterface, 0)
	for i := 0; i < 10; i++ {
		str := randomGenerator_ab.GenerateString(100000)
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(str)
		trees = append(trees, st)
	}

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	for i := 0; i < 10; i++ {
		DecorateTreeWithVocabulary(trees[i%20])
	}

	defer pprof.StopCPUProfile()

}
