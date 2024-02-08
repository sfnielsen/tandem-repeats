package tandemrepeat

//simple function to find tandem repeats in O(n^3)
func FindTandemRepeatsNaive(s string) []TandemRepeat {
	var repeats []TandemRepeat
	var i, j, k int

	//keep two pointers, i and j, and a third pointer, k, to keep track of the length of the tandem repeat
	for i = 0; i < len(s); i++ {
		for j = i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				for k = 0; k < j-i; k++ {
					//if the characters at i+k and j+k are not the same, we break the loop
					//as this is not a tandem repeat
					if s[i+k] != s[j+k] {
						break
					}
				}
				//if we ran k to the end, we have a tandem repeat
				if k == j-i {
					repeats = append(repeats, TandemRepeat{i, k, 2})
				}
			}
		}
	}
	return repeats
}
