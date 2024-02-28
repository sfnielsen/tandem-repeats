package tandemrepeat

//simple function to find tandem repeats in O(n^3)
func FindTandemRepeatsNaive(s string) []TandemRepeat {
	var repeats []TandemRepeat
	var i, j, k int

	//keep two pointers, i and j, and a third pointer, k, to keep track of the length of the tandem repeat
	for i = 0; i < len(s); i++ {
		for j = i + 1; j < len(s); j++ {
			k = j - i //length of the tandem repeat
			if j+k > len(s)-1 {
				break
			}
			if s[i:j] == s[j:j+k] {
				//add tandem repeat to list
				repeats = append(repeats, TandemRepeat{i, k, 2})

			}
		}
	}
	return repeats
}
