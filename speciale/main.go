package main

import (
	"speciale/suffixtreeimpl"
	"speciale/tandemrepeat"
)

func main() {
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree("ababababbbbbbbbbbbbbbbbbbbbbbbbbbabababababababaaaabababaaaababaaaaaaaaaaaaaaabbbbbbbbbbbaaaaaaaaaaaaaaabbbbbbbbbaaaaaaababababababababbbbbbbbbbababababbacdacdabcbdbcdabcdabcdbadfdsbfdbsfdhsgcdabcdgscvhdsvhcdsbcdsbcdbbcsdcgdscdsgcdgsgcdsggcgggcgdsgcgdsgcgdsgsgcdgsgcgdsgcgsgcgscscscscgsgcdgcdgsgcdgscgdsggcdsgcgdsgcgdsgcgsdgcgdsgcds$")
	print(st.GetRoot())
	print(st.GetInputString())

	st2 := "ababab$"
	// find tandem repeats
	tr := tandemrepeat.FindTandemRepeatsNaive(st2)
	for _, repeat := range tr {
		println(tandemrepeat.GetTandemRepeatSubstring(repeat, st2))
	}
}
