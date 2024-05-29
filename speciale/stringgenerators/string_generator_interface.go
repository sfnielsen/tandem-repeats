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
	AlphabetByte string = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f" +
		"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
		" !\"#%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f" +
		"\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f" +
		"\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f" +
		"\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf" +
		"\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf" +
		"\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf" +
		"\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf" +
		"\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef" +
		"\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
)
