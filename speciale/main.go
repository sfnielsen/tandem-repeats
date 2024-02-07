package main

import (
	"speciale/suffixtreeimpl"
)

func main() {
	// Create a NaiveSuffixTree instance
	st := suffixtreeimpl.ConstructNaiveSuffixTree("ababababbbbbbbbbbbbbbbbbbbbbbbbbbabababababababaaaabababaaaababaaaaaaaaaaaaaaabbbbbbbbbbbaaaaaaaaaaaaaaabbbbbbbbbaaaaaaababababababababbbbbbbbbbababababbacdacdabcbdbcdabcdabcdbadfdsbfdbsfdhsgcdabcdgscvhdsvhcdsbcdsbcdbbcsdcgdscdsgcdgsgcdsggcgggcgdsgcgdsgcgdsgsgcdgsgcgdsgcgsgcgscscscscgsgcdgcdgsgcdgscgdsggcdsgcgdsgcgdsgcgsdgcgdsgcds$")
	print(st.GetRoot())
	print(st.GetInputString())
}
