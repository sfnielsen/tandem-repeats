package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"speciale/tandemrepeat"
	"strings"
	"time"
)

type Algorithm func(inputSize int)

type TimingResult struct {
	InputSize   int
	Algorithm   string
	RunningTime time.Duration
}

func SaveResults(results []TimingResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"InputSize", "Algorithm", "RunningTime"}
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
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

type SuffixTreeConstructionType func(string) suffixtree.SuffixTree
type TandemRepeatFinderType func(suffixtree.SuffixTree) []tandemrepeat.TandemRepeat

type AlgorithmSetup struct {
	SuffixTreeConstructor SuffixTreeConstructionType
	TandemRepeatFinder    TandemRepeatFinderType
}

func TakeTimeAndSave(setup AlgorithmSetup, maxSize int, steps int) {
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	fmt.Println(filename)
	var randomGenerator stringgenerators.StringGenerator = &stringgenerators.RandomStringGenerator{Alphabet: stringgenerators.AlphabetAB}

	var results []TimingResult

	for i := maxSize / steps; i < maxSize; i += int(maxSize / steps) {
		fmt.Println(i)
		// Run each type 10 times
		for range [10]int{} {
			// Construct suffix tree
			inputString := randomGenerator.GenerateString(i)

			//Take time on suffix tree construction
			startST := time.Now()
			sa := setup.SuffixTreeConstructor(inputString)
			elapsedST := time.Since(startST)

			// Find tandem repeats and take time
			startTR := time.Now()
			setup.TandemRepeatFinder(sa)
			elapsedTR := time.Since(startTR)

			results = append(results,
				TimingResult{
					InputSize:   i,
					Algorithm:   "SuffixTree",
					RunningTime: elapsedST,
				},
				TimingResult{
					InputSize:   i,
					Algorithm:   "TandemRepeat",
					RunningTime: elapsedTR,
				})

		}

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}
