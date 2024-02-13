package main

import (
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"speciale/utils"
)

func main() {
	utils.TakeTimeAndSave(utils.AlgorithmSetup{
		SuffixArrayConstructor: suffixtreeimpl.ConstructNaiveSuffixTree,
		TandemRepeatFinder:     tandemrepeat.FindTandemRepeatsLogarithmic,
	}, 150000, 150)
}
