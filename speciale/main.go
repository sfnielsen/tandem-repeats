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
	println(mccreight, naive)

	tandemrepeat_nlogn := utils.AlgorithmBase{
		Name: "nlogn stoye gusfield",
		Algorithm: func(args ...interface{}) interface{} {
			return tandemrepeat.FindAllTandemRepeatsLogarithmic(args[0].(suffixtree.SuffixTreeInterface))
		},
		ExpectedComplexity: "nlogn"}

	tandemrepeat_n := utils.AlgorithmBase{
		Name: "n stoye gusfield",
		Algorithm: func(args ...interface{}) interface{} {
			return tandemrepeat.FindAllTandemRepeatsLogarithmic(args[0].(suffixtree.SuffixTreeInterface))
		},
		ExpectedComplexity: "n^2"}

	var nstoygus utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{tandemrepeat_n}
	var nlognstoygus utils.AlgorithmInterface = &utils.AlgorithmTandemrepeat{tandemrepeat_nlogn}

	functionSlice := []utils.AlgorithmInterface{nlognstoygus, nstoygus}

	utils.TakeTimeAndSave(functionSlice, 30000, 50, stringgenerators.AlphabetDNA)

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
