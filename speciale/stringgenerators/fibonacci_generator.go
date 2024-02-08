package stringgenerators

type FibonacciStringGenerator struct {
	First  string
	Second string
}

func (g *FibonacciStringGenerator) GenerateString(n int) string {
	if n <= 0 {
		return ""
	}

	// The first two characters of the string are always "1"
	prev1 := g.First
	prev2 := g.Second

	// Iterate and concatenate strings according to the Fibonacci sequence
	for i := 2; i < n+1; i++ {
		current := prev2 + prev1

		// Update previous strings for the next iteration
		prev1 = prev2
		prev2 = current
	}

	// Add the sentinel character ('$') at the end
	prev2 += "$"

	return prev2
}
