package suffixtree

import (
	"speciale/stringgenerators"
	"testing"
)

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

// test that mapping works
func TestInputStringToInternalStringMapping(t *testing.T) {
	str := "abcdefghijklmnopqrstuvwxyz..,abcabccccccccca.,.-$"
	str2, alphabetSize := InputStringToInternalString(str)

	if alphabetSize != 30 {
		t.Errorf("Expected 30, got %d", alphabetSize)
	}

	//make array of alphabetsize and ensure that each value fits within the alphabet size
	array := make([]bool, alphabetSize)
	for _, c := range str2 {
		//ensure all values are within the alphabet size
		if int(c) >= alphabetSize || int(c) < 0 {
			t.Errorf("Expected %d, got %d", alphabetSize, c)
		}
		array[int(c)] = true
	}
	// check if array is all true
	for _, b := range array {
		if !b {
			t.Errorf("Expected true, got false")
		}
	}

}

// test mapping on byte-alphabet
func TestInputStringToInternalStringMappingByte(t *testing.T) {
	randomGenerator_byte.SetSeed(1)
	str := randomGenerator_byte.GenerateString(100000)

	str2, alphabetSize := InputStringToInternalString(str)

	if alphabetSize != 128 {
		t.Errorf("Expected 256, got %d", alphabetSize)
	}

	//make array of alphabetsize and ensure that each value fits within the alphabet size
	array := make([]bool, alphabetSize)
	for _, c := range str2 {
		//ensure all values are within the alphabet size
		if byte(c) > byte(alphabetSize-1) {
			t.Errorf("Expected to be less than %d, got %d", byte(alphabetSize), byte(c))
		} else {
			array[byte(c)] = true
		}
	}
	// check if array is all true
	for _, b := range array {
		if !b {
			t.Errorf("Expected true, got false")
		}
	}
}
