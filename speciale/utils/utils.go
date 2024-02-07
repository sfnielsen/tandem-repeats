package utils

import "strconv"

//convert a string and index to cigar format
func convertToCigarFormat(x string, read string, index []int) string {
	var result string
	var i int
	for i = 0; i < len(index); i++ {
		result += strconv.Itoa(index[i] + 1)
		result += "="
		result += string(read[index[i]])
	}
	if i == 0 {
		return "*"
	}
	return result
}
