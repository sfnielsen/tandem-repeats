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

func main() {
	var suffixTreeAlgo utils.AlgorithmInterface = &utils.AlgorithmBase{
		Name: "SuffixTree",
		Algorithm: func(args ...interface{}) interface{} {
			return suffixtreeimpl.ConstructNaiveSuffixTree(args[0].(string))
		},
		ExpectedComplexity: "nlogn"}
	tandemAlgo := utils.AlgorithmBase{
		Name: "TandemRepeat",
		Algorithm: func(args ...interface{}) interface{} {
			return tandemrepeat.FindAllBranchingTandemRepeatsLogarithmic(args[0].(suffixtree.SuffixTreeInterface))
		},
		ExpectedComplexity: "nlogn"}

	var tdalg utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{tandemAlgo}

	functionSlice := []utils.AlgorithmInterface{suffixTreeAlgo, tdalg}

	utils.TakeTimeAndSave(functionSlice, 70000, 70, stringgenerators.AlphabetAB)

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
