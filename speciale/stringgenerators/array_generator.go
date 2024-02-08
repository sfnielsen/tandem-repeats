package stringgenerators

import (
	"math/rand"
	"speciale/suffixtree"
)

func GenerateStringArray(numberOfStrings int, stringLength int, generators []suffixtree.StringGenerator) []string {
	// Create a slice of strings of length numberOfStrings
	strings := make([]string, numberOfStrings)

	// Iterate over the slice and fill it with random strings of length stringLength
	for i := range strings {
		// select random generator from generator slice
		randomIndex := rand.Intn(len(generators))
		generator := generators[randomIndex]
		strings[i] = generator.GenerateString(stringLength)
	}

	// Return the slice of strings
	return strings
}
