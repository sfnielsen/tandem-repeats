package main

import (
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"speciale/utils"
	"speciale/stringgenerators"
	"fmt"
	"os/exec"
)

func main() {
	utils.TakeTimeAndSave(utils.AlgorithmSetup{
		SuffixTreeConstructor: suffixtreeimpl.ConstructNaiveSuffixTree,
		TandemRepeatFinder:    tandemrepeat.FindAllBranchingTandemRepeatsLogarithmic,
	}, 40000, 20, stringgenerators.AlphabetAB)

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
