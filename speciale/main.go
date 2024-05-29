package main

import (
	"fmt"
	"os/exec"
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"speciale/utils"
)

var mccreight utils.AlgorithmInterface = &utils.AlgorithmBase{
	Name: "McCreight",
	Algorithm: func(args ...interface{}) interface{} {
		return suffixtreeimpl.ConstructMcCreightSuffixTree(args[0].(string))
	},
	ExpectedComplexity: "n"}

var naive utils.AlgorithmInterface = &utils.AlgorithmBase{
	Name: "Naive",
	Algorithm: func(args ...interface{}) interface{} {
		return suffixtreeimpl.ConstructNaiveSuffixTree(args[0].(string))
	},
	ExpectedComplexity: "nlogn"}

var tandemrepeat_nlogn = utils.AlgorithmBase{
	Name: "nlogn stoye gusfield",
	Algorithm: func(args ...interface{}) interface{} {
		return tandemrepeat.FindAllBranchingTandemRepeatsLogarithmic(args[0].(suffixtree.SuffixTreeInterface))
	},
	ExpectedComplexity: "nlogn"}
var tandemrepeat_n = utils.AlgorithmBase{
	Name: "n stoye gusfield",
	Algorithm: func(args ...interface{}) interface{} {
		tandemrepeat.DecorateTreeWithVocabulary(args[0].(suffixtree.SuffixTreeInterface))
		return nil
	},
	ExpectedComplexity: "n"}

var nstoygus utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{tandemrepeat_n}
var nlognstoygus utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{tandemrepeat_nlogn}

func main() {
	functionSlice := []utils.AlgorithmInterface{nstoygus}

	utils.MeasureSizeOfTrees(functionSlice, 200000)

	pythonScript := "../visualization.py"
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
