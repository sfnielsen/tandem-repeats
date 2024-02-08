package stringgenerators

// Different alphabet types for generating strings

type AlphabetStringGenerator struct {
	Alphabet string
}

func (g *AlphabetStringGenerator) GenerateString(n int) string {
	// Create a byte slice of length n
	b := make([]byte, n+1)

	// Iterate over the byte slice and fill it with random characters from the alphabet
	for i := range b[:n] {
		b[i] = g.Alphabet[i%len(g.Alphabet)]
	}
	// Add the sentinel character at the end
	b[n] = '$'

	// Return the byte slice as a string
	return string(b)
}
