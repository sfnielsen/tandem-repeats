package stringgenerators

import (
	"math/rand"
	"time"
)

type RandomStringGenerator struct {
	Alphabet string
	seed     int
}

func (g *RandomStringGenerator) SetSeed(providedSeed int) {
	g.seed = providedSeed
}

func (g *RandomStringGenerator) setRandSeed() {
	if g.seed != 0 {
		rand.Seed(int64(g.seed))
	} else {
		// If seed is not set, use the current time as a seed
		rand.Seed(time.Now().UnixNano())
	}

}

func (g *RandomStringGenerator) GenerateString(n int) string {
	g.setRandSeed()

	// Create a byte slice of length n
	b := make([]byte, n+1)

	// Iterate over the byte slice and fill it with random characters from the alphabet
	for i := range b[:n] {
		// Generate a random index within the range of the alphabet
		randomIndex := rand.Intn(len(g.Alphabet))
		// Select a random character from the alphabet and assign it to the byte slice
		b[i] = g.Alphabet[randomIndex]
	}
	// Add the sentinel character at the end
	b[n] = '$'

	// Return the byte slice as a string
	return string(b)
}
