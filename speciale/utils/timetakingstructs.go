package utils

import (
	"speciale/suffixtree"
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
	"time"
)

type Algorithm func(inputSize int)

type TimingResult struct {
	InputSize          int
	Algorithm          string
	RunningTime        time.Duration
	ExpectedComplexity string
	Alphabet           string
}

type SuffixTreeConstructionType func(string) suffixtree.SuffixTreeInterface
type TandemRepeatFinderType func(suffixtree.SuffixTreeInterface) []tandemrepeat.TandemRepeat

type AlgorithmInterface interface {
	GetTime(args ...interface{}) time.Duration
	GetName() string
	GetExpectedComplexity() string
}
type AlgorithmBase struct {
	Name               string
	Algorithm          func(...interface{}) interface{}
	ExpectedComplexity string // can either be "nlogn", "n^2" or "n" at the moment.
}

// basic get time method
func (a *AlgorithmBase) GetTime(args ...interface{}) time.Duration {
	start := time.Now()
	a.Algorithm(args[0])
	return time.Since(start)
}

// basic get name method
func (a *AlgorithmBase) GetName() string {
	return a.Name
}

// basic get name method
func (a *AlgorithmBase) GetExpectedComplexity() string {
	return a.ExpectedComplexity
}

// tandem repeat type algorithm that needs to create a suffix tree before taking time
type AlgorithmTandemrepeat struct {
	AlgorithmBase
}

// Altered GetTime algorithm for tandem repeats that first creates a suffix tree and then takes time
func (a *AlgorithmTandemrepeat) GetTime(args ...interface{}) time.Duration {
	var st suffixtree.SuffixTreeInterface = suffixtreeimpl.ConstructMcCreightSuffixTree(args[0].(string))
	start := time.Now()
	a.Algorithm(st)
	return time.Since(start)
}

type AlgorithmSetup struct {
	SuffixTreeConstructor SuffixTreeConstructionType
	TandemRepeatFinder    TandemRepeatFinderType
}
