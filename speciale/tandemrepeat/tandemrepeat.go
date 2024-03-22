package tandemrepeat

type TandemRepeat struct {
	// The start index of the tandem repeat
	Start int
	// The length of the tandem repeat
	Length int
	// The number of times the tandem repeat is repeated
	Repeats int
}

// return the actual tandem repeat as a string
func GetTandemRepeatSubstring(tr TandemRepeat, inputString string) string {
	return inputString[tr.Start : tr.Start+tr.Repeats*tr.Length]
}
