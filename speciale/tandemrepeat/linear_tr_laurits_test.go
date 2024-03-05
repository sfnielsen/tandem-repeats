package tandemrepeat

import (
	"speciale/suffixtreeimpl"
	"testing"
)

// test that naive_tr works on small example
func TestLauritsTR(t *testing.T) {
	//define a simple string
	//generate some big strings from the stringgenerators
	s := randomGenerator_ab.GenerateString(100)

	lcaLookup(3, 4, suffixtreeimpl.ConstructNaiveSuffixTree(s))
	if 1 != 2 {
		t.Errorf("Expected %d tandem repeats, got %d", 1, 2)
	}
}
