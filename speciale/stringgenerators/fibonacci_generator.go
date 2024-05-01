package stringgenerators

type FibonacciStringGenerator struct {
	First  string // not used
	Second string // not used
}

func (g *FibonacciStringGenerator) SetSeed(providedSeed int) {}

func (g *FibonacciStringGenerator) GenerateString(n int) string {
	if n <= 0 {
		return "$"
	}

	// The first two characters of the string are always "1"
	fn1 := "a"
	fn := "ab"
	temp := ""

	// Iterate and concatenate strings according to the Fibonacci sequence
	for len(fn) < n {
		// Append the current Fibonacci string to the builder
		temp = fn
		fn += fn1
		fn1 = temp
	}

	//trim so it has length n
	fn = fn[:n]

	// Return the built string with the sentinel character
	return fn + "$"
}
