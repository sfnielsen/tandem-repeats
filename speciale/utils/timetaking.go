package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"speciale/stringgenerators"
	"strings"
	"time"
)

func SaveResults(results []TimingResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"InputSize", "Algorithm", "RunningTime", "Complexity", "Alphabet"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, result := range results {
		durationString := strings.TrimSuffix(result.RunningTime.String(), "ms")

		row := []string{
			fmt.Sprintf("%d", result.InputSize),
			result.Algorithm,
			durationString,
			result.ExpectedComplexity,
			result.Alphabet,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func TakeTimeAndSave(functions []AlgorithmInterface, maxSize int, steps int, alphabet string) {
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	fmt.Println(filename)
	var randomGenerator stringgenerators.StringGenerator = &stringgenerators.RandomStringGenerator{Alphabet: alphabet}

	var results []TimingResult

	for i := maxSize / steps; i <= maxSize; i += int(maxSize / steps) {
		fmt.Println(i)
		// Run each type 10 times
		for range [10]int{} {
			// Construct suffix tree
			inputString := randomGenerator.GenerateString(i)
			for _, function := range functions {
				time := function.GetTime(inputString)

				results = append(results,
					TimingResult{
						InputSize:          i,
						Algorithm:          function.GetName(),
						RunningTime:        time,
						ExpectedComplexity: function.GetExpectedComplexity(),
						Alphabet:           alphabet,
					})
			}

		}

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}

func TakeTimeAllAlphabets(functions []AlgorithmInterface, maxSize int, steps int) {
	//iterate all alphabettypes
	alphabets := []string{stringgenerators.AlphabetA, stringgenerators.AlphabetAB, stringgenerators.AlphabetDNA, stringgenerators.AlphabetProtein, stringgenerators.AlphabetByte}
	var results []TimingResult
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	fmt.Println(filename)

	itr := 0
	//iterate all alphabets
	for _, alphabet := range alphabets {
		itr++
		var randomGenerator stringgenerators.StringGenerator = nil
		if alphabet == "fib" {
			randomGenerator = &stringgenerators.FibonacciStringGenerator{}
		} else {
			randomGenerator = &stringgenerators.RandomStringGenerator{Alphabet: alphabet}
		}

		for i := maxSize / steps; i < maxSize; i += int(maxSize / steps) {
			fmt.Println(i)
			// Run each type 10 times
			for range [10]int{} {
				// Construct suffix tree
				inputString := randomGenerator.GenerateString(i)
				for _, function := range functions {
					time := function.GetTime(inputString)

					if itr == len(alphabets) {
						alphabet = "Byte"
					}
					results = append(results,
						TimingResult{
							InputSize:          i,
							Algorithm:          function.GetName(),
							RunningTime:        time,
							ExpectedComplexity: function.GetExpectedComplexity(),
							Alphabet:           alphabet,
						})
				}

			}

		}

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}
