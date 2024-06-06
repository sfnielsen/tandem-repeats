package main

import (
	"fmt"
	"os/exec"
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
	ExpectedComplexity: "nlogn"}
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

var linearTrOutput = utils.AlgorithmBase{
	Name: "TandemRepeat Linear",
	Algorithm: func(args ...interface{}) interface{} {
		tandemrepeat.GetAllTandemRepeatsFromDecoratedTree(args[0].(suffixtree.SuffixTreeInterface))
		return []tandemrepeat.TandemRepeat{}
	},
	ExpectedComplexity: "n^2"}
var linearTrOutputStruct utils.AlgorithmInterface = &utils.AlgorithmTandemrepeatVOCOutput{linearTrOutput}

var branchingTrOutput = utils.AlgorithmBase{
	Name: "TandemRepeat Logarithmic",
	Algorithm: func(args ...interface{}) interface{} {
		tandemrepeat.GetAllTandemRepeats(args[0].([]tandemrepeat.TandemRepeat), args[1].(suffixtree.SuffixTreeInterface))

		return []tandemrepeat.TandemRepeat{}
	},
	ExpectedComplexity: "n^2"}
var braTrOutputStruct utils.AlgorithmInterface = &utils.AlgorithmTandemrepeatBRAOutput{branchingTrOutput}

func main() {
	functionSlice := []utils.AlgorithmInterface{linearTR, logTR}

	utils.TakeTimeAllAlphabets(functionSlice, 1000000, 30)

	pythonScript := "../visualization.py"
	//scriptArgs := []string{}

	// Build the command to execute the Python script
	cmd := exec.Command("python", pythonScript)

	// Capture the output of the Python script
	_, err := cmd.Output()
	if err != nil {
		// Handle error
		fmt.Println("Error:", err)
		return
	}

}
