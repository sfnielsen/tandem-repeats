package stringgenerators

import "math/rand"

type StringGenerator interface {
	GenerateString(n int) string
	SetSeed(providedSeed int)
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
	//Contains 128 ASCII characters, except $ which is used as a sentinel character
	AlphabetByte string = "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f" +
		"\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
		" !\"#%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f"
)
