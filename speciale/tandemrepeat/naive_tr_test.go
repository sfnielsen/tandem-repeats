package tandemrepeat

import (
	"testing"
)

// test that naive_tr works on small example
func TestNaiveTrSimpleTest(t *testing.T) {
	//define a simple string
	s := "abab$"
	// Create a NaiveTandemRepeat instance
	tr := FindTandemRepeatsNaive(s)
	if len(tr) != 1 {
		t.Errorf("Expected 1 tandem repeat, got %d", len(tr))
	}

	//extract the string from the tandem repeat and verify that it is correct ab$
	if GetTandemRepeatSubstring(tr[0], s) != "ab" {
		t.Errorf("Expected tandem repeat to be 'ab', got %s", GetTandemRepeatSubstring(tr[0], s))
	}

}

// now test that we get 3 tandem repeats in the string ababab$
func TestNaiveTrThreeRepeats(t *testing.T) {
	//define a simple string
	s := "ababab$"
	// Create a NaiveTandemRepeat instance
	tr := FindTandemRepeatsNaive(s)
	if len(tr) != 3 {
		t.Errorf("Expected 3 tandem repeats, got %d", len(tr))
	}
	//extract the string from the tandem repeat and verify that it is correct ab$
	if GetTandemRepeatSubstring(tr[0], s) != "ab" {
		t.Errorf("Expected tandem repeat to be 'ab', got %s", GetTandemRepeatSubstring(tr[0], s))
	}
	if GetTandemRepeatSubstring(tr[1], s) != "ba" {
		t.Errorf("Expected tandem repeat to be 'ba', got %s", GetTandemRepeatSubstring(tr[1], s))
	}
	if GetTandemRepeatSubstring(tr[2], s) != "ab" {
		t.Errorf("Expected tandem repeat to be 'ab', got %s", GetTandemRepeatSubstring(tr[2], s))
	}
}

// test that we get "amount" of tandem repeats
func TestNaiveTrAmountIsX(t *testing.T) {
	s := "aaaaaaaa$"
	// Create a NaiveTandemRepeat instance
	tr := FindTandemRepeatsNaive(s)
	
	if len(tr) != 16 {
		t.Errorf("Expected 4 tandem repeats, got %d", len(tr))
	}

}
