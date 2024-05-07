package main

import (
	"fmt"
	"os/exec"
	"speciale/stringgenerators"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"speciale/utils"
)

var logTandemAlgo = utils.AlgorithmBase{
	Name: "TandemRepeat Logarithmic",
	Algorithm: func(args ...interface{}) interface{} {
		return tandemrepeat.FindAllBranchingTandemRepeatsLogarithmic(args[0].(suffixtree.SuffixTreeInterface))
	},
	ExpectedComplexity: "nlogn"}
var logTR utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{logTandemAlgo}

var lineraTandemAlgobase = utils.AlgorithmBase{
	Name: "TandemRepeat Linear",
	Algorithm: func(args ...interface{}) interface{} {
		tandemrepeat.DecorateTreeWithVocabulary(args[0].(suffixtree.SuffixTreeInterface))
		return []tandemrepeat.TandemRepeat{}
	},
	ExpectedComplexity: "n"}
var linearTR utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{lineraTandemAlgobase}

var naiveSuffixTreeAlgo utils.AlgorithmInterface = &utils.AlgorithmBase{
	Name: "SuffixTree Naive",
	Algorithm: func(args ...interface{}) interface{} {
		return suffixtreeimpl.ConstructNaiveSuffixTree(args[0].(string))
	},
	ExpectedComplexity: "n^2"}

var mcCregightSuffixTreeAlgo utils.AlgorithmInterface = &utils.AlgorithmBase{
	Name: "SuffixTree McCreight",
	Algorithm: func(args ...interface{}) interface{} {
		return suffixtreeimpl.ConstructMcCreightSuffixTree(args[0].(string))
	},
	ExpectedComplexity: "n"}

func main() {
	functionSlice := []utils.AlgorithmInterface{logTR}

	utils.TakeTimeAndSave(functionSlice, 200000, 50, stringgenerators.AlphabetA)

	pythonScript := "../visualization_all_alphabets.py"
	//scriptArgs := []string{}

	// Build the command to execute the Python script
	cmd := exec.Command("python3", pythonScript)

	// Capture the output of the Python script
	_, err := cmd.Output()
	if err != nil {
		// Handle error
		fmt.Println("Error:", err)
		return
	}

}
