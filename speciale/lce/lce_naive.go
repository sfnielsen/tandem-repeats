package lce

// find the longest common extension of two suffixes that starts at i and j
func FindLCEForwardSlow(s string, i, j, alphabetsize int) int {
	if alphabetsize == 2 {
		if i > j {
			return len(s) - i - 1
		} else {
			return len(s) - j - 1
		}
	}
	lce := 0

	//match letters until we have a mismatch
	for i < len(s) && j < len(s) {
		if s[i] != s[j] {
			return lce
		} else {
			i++
			j++
			lce++
		}
	}
	return lce

}

// find the longest common extension of two suffixes that ends at i and j
func FindLCEBackwardSlow(s string, i, j, alphabetsize int) int {
	if alphabetsize == 2 {
		if i > j {
			return j + 1
		} else {
			return i + 1
		}
	}

	lce := 0

	//match letters until we have a mismatch
	for i >= 0 && j >= 0 {
		if s[i] != s[j] {
			return lce
		} else {
			i--
			j--
			lce++
		}
	}
	return lce
}
