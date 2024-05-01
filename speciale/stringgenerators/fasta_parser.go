package stringgenerators

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// ExtractSequenceFromFasta extracts a random sequence of the provided length from a FASTA file
func ExtractSequenceFromFasta(filename string, length int) (string, error) {

	// Read FASTA file
	data, err := ioutil.ReadFile("../../fasta/" + filename)
	if err != nil {
		return "", err
	}

	// Split the data into lines
	lines := strings.Split(string(data), "\n")

	var sb strings.Builder
	// Concatenate the sequence lines (skip the header line starting with ">")
	for _, line := range lines {
		if !strings.HasPrefix(line, ">") {
			sb.WriteString(line)
		}
	}
	str := sb.String()
	// Generate a random starting index within the valid range
	rand.Seed(time.Now().UnixNano())
	startIndex := rand.Intn(len(str) - length)

	// Extract the sequence starting from the random index with the provided length
	extractedSequence := str[startIndex : startIndex+length]

	return extractedSequence, nil
}
