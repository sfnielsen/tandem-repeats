package main

import (
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"speciale/utils"
)

func main() {
	utils.TakeTimeAndSave(utils.AlgorithmSetup{
		SuffixTreeConstructor: suffixtreeimpl.ConstructNaiveSuffixTree,
		TandemRepeatFinder:    tandemrepeat.FindTandemRepeatsLogarithmic,
	}, 20000, 40)
}
