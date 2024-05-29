package stringgenerators

import (
	"math/rand"
	"strings"
)

type StringGenerator interface {
	GenerateString(n int) string
	SetSeed(providedSeed int)
}

func CreateMaxiAlphabet() string {
	var alphabet string = AlphabetByte

	// Add Latin characters, excluding '$'
	for i := 0; i < 128; i++ {
		char := string(rune(i))
		if char != "$" && !strings.Contains(alphabet, char) {
			alphabet += char
		}
	}

	// Add Greek characters (U+0370 to U+03FF), excluding '$'
	for i := 0x370; i <= 0x3FF; i++ {
		char := string(rune(i))
		if char != "$" && !strings.Contains(alphabet, char) {
			alphabet += char
		}
	}

	// Add Cyrillic characters (U+0400 to U+04FF), excluding '$'
	for i := 0x400; i <= 0x4FF; i++ {
		char := string(rune(i))
		if char != "$" && !strings.Contains(alphabet, char) {
			alphabet += char
		}
	}

	return alphabet
}

func GenerateStringFromGivenAlphabet(alphabet string, n int) string {
	// Create a byte slice of length n
	b := make([]byte, n+1)
	length := len(alphabet)
	// Iterate over the byte slice and fill it with random characters from the alphabet
	for i := range b[:n] {
		// Generate a random index within the range of the alphabet
		randomIndex := rand.Intn(length)
		// Select a random character from the alphabet and assign it to the byte slice
		b[i] = alphabet[randomIndex]
	}
	// Add the sentinel character at the end
	b[n] = '$'

	// Return the byte slice as a string
	return string(b)
}

const (
	AlphabetA       string = "A"
	AlphabetAB      string = "AB"
	AlphabetDNA     string = "ACGT"
	AlphabetProtein string = "ACDEFGHIKLMNPQRSTVWY"
	//Contains 256 characters (a byte), except $ which is used as a sentinel character
	AlphabetByte string = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f" +
		"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
		" !\"#%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f"
)
