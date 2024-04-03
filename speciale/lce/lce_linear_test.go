package lce

import (
	"speciale/stringgenerators"
	"speciale/suffixtreeimpl"
	"testing"
)

// Test functions must start with the word "Test" and take a *testing.T parameter.
var (
	setupCompleted          bool
	randomGenerator_protein stringgenerators.StringGenerator
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
		setupCompleted = true
	}
}

func TestLCEArrays(t *testing.T) {

	s := randomGenerator_ab.GenerateString(3000)

	st := suffixtreeimpl.ConstructMcCreightSuffixTree(s)

	L, E, R := createLERArrays(st)
	n := st.GetSize()

	//check that L,E,R sizes are correct
	if len(L) != 2*n-1 {
		t.Errorf("Expected L size to be %d, got %d", 2*n-1, len(L))
	}
	if len(E) != 2*n-1 {
		t.Errorf("Expected E size to be %d, got %d", 2*n-1, len(E))
	}
	if len(R) != n {
		t.Errorf("Expected R size to be %d, got %d", n, len(R))
	}

	//Test that levels/depths are +-1 of each other
	for i := 0; i < len(L)-1; i++ {
		if L[i] != L[i+1]+1 && L[i] != L[i+1]-1 {
			t.Errorf("Expected L[i] to be L[i+1]+-1")
		}
	}
	//first and last depth is 0
	if L[0] != 0 || L[len(L)-1] != 0 {
		t.Errorf("Expected first and last depth to be 0 (ALWAYS root)")
	}

	//setup for additional testing
	allInts := make(map[int]int)
	for _, v := range R {
		allInts[E[v]] = L[v]
	}

	//check that R maps to all 0,1,...,n-1 different ints in L
	if len(allInts) != n {
		t.Errorf("Expected R to map to all 0,1,...,n-1 different ints")
	}
	//check that all same E values also have same L values
	for _, v := range E {
		if allInts[E[v]] != L[v] {
			t.Errorf("Expected L to be the same for all same E values")
		}

	}

}
