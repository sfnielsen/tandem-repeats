package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"speciale/stringgenerators"
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
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
		for range [1]int{} {
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
			for range [2]int{} {
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

func MeasureSizeOfTrees(inputsize int) {
	//iterate all alphabettypes
	maxi_alphabet := stringgenerators.AlphabetByte
	fmt.Println("alphabet length", len(maxi_alphabet))
	var results []TimingResult
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	for range [2]int{} {
		// Construct suffix tree
		inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:2], inputsize)
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
		results = append(results,
			TimingResult{
				InputSize:          2,
				Algorithm:          "Size of tree",
				RunningTime:        time.Duration(st.GetSize()),
				ExpectedComplexity: "n",
				Alphabet:           "Byte",
			})
	}
	//iterate all alphabets
	for idx := 1; idx < len(maxi_alphabet); idx++ {
		fmt.Println(idx)
		for range [10]int{} {
			// Construct suffix tree
			inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:idx], inputsize)
			st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
			results = append(results,
				TimingResult{
					InputSize:          idx,
					Algorithm:          "Size of tree",
					RunningTime:        time.Duration(st.GetSize()),
					ExpectedComplexity: "n",
					Alphabet:           "Byte",
				})

		}
		idx += 4

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}

func DfsAndLookuptime(inputsize int) {
	//iterate all alphabettypes
	maxi_alphabet := stringgenerators.AlphabetByte
	fmt.Println("alphabet length", len(maxi_alphabet))
	var results []TimingResult
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	for range [15]int{} {
		// Construct suffix tree
		inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:2], inputsize)
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
		dfstime, lookuptime, dfstimebefore := tandemrepeat.DecorateTreeWithVocabulary(st)
		results = append(results,
			TimingResult{
				InputSize:          2,
				Algorithm:          "Linear",
				RunningTime:        time.Duration(dfstime),
				ExpectedComplexity: "n",
				Alphabet:           "Dfs",
			})
		results = append(results,
			TimingResult{
				InputSize:          2,
				Algorithm:          "Linear",
				RunningTime:        time.Duration(lookuptime),
				ExpectedComplexity: "n",
				Alphabet:           "Lookup",
			})
		results = append(results,
			TimingResult{
				InputSize:          2,
				Algorithm:          "Linear",
				RunningTime:        time.Duration(dfstimebefore),
				ExpectedComplexity: "n",
				Alphabet:           "dfs before",
			})
	}
	//iterate all alphabets
	for idx := 1; idx < len(maxi_alphabet); idx++ {
		fmt.Println(idx)
		for range [5]int{} {
			// Construct suffix tree
			inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:idx], inputsize)
			st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
			dfstime, lookuptime, dfstimebefore := tandemrepeat.DecorateTreeWithVocabulary(st)
			results = append(results,
				TimingResult{
					InputSize:          idx,
					Algorithm:          "Linear",
					RunningTime:        time.Duration(dfstime),
					ExpectedComplexity: "n",
					Alphabet:           "Dfs",
				})
			results = append(results,
				TimingResult{
					InputSize:          idx,
					Algorithm:          "Linear",
					RunningTime:        time.Duration(lookuptime),
					ExpectedComplexity: "n",
					Alphabet:           "Lookup",
				})
			results = append(results,
				TimingResult{
					InputSize:          idx,
					Algorithm:          "Linear",
					RunningTime:        time.Duration(dfstimebefore),
					ExpectedComplexity: "n",
					Alphabet:           "dfs before",
				})
		}
		idx += 4

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}

func AverageChildrenIncreasingAlphabetSize(inputsize int) {
	//iterate all alphabettypes
	maxi_alphabet := stringgenerators.AlphabetByte
	fmt.Println("alphabet length", len(maxi_alphabet))
	var results []TimingResult
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("time_csvs/timing_results_%s.csv", currentTime)
	for range [2]int{} {
		// Construct suffix tree
		inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:2], inputsize)
		st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
		childrenCount := DfsCountChildren(st)

		//calculate average number of children
		average := 0
		for i, count := range childrenCount {
			average += i * count
		}
		average = average / st.GetSize()

		results = append(results,
			TimingResult{
				InputSize:          2,
				Algorithm:          "Size of tree",
				RunningTime:        time.Duration(average),
				ExpectedComplexity: "n",
				Alphabet:           "Byte",
			})
	}
	//iterate all alphabets
	for idx := 1; idx < len(maxi_alphabet); idx++ {
		fmt.Println(idx)
		for range [5]int{} {
			// Construct suffix tree
			inputString := stringgenerators.GenerateStringFromGivenAlphabet(maxi_alphabet[:idx], inputsize)
			st := suffixtreeimpl.ConstructMcCreightSuffixTree(inputString)
			childrenCount := DfsCountChildren(st)
			//calculate average number of children
			average := 0.0
			for i, count := range childrenCount {
				if i == 1 {
					continue
				}
				average += float64(i * count)
			}
			average = average / float64(st.GetSize()-inputsize)
			fmt.Println(average)
			//calculate average number of children

			results = append(results,
				TimingResult{
					InputSize:          idx,
					Algorithm:          "Size of tree",
					RunningTime:        time.Duration(average),
					ExpectedComplexity: "n",
					Alphabet:           "Byte",
				})

		}
		idx += 4

	}
	// Save results to a CSV file
	if err := SaveResults(results, filename); err != nil {
		fmt.Println("Error saving results:", err)
	}
}
