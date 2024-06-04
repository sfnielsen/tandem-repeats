package suffixtree

//function takes a string with $ as input and returns a string
func InputStringToInternalString(str string) (newString string, alphabetSize int) {
	const ASCIISIZE = 128

	//make a bytesize array of bools
	isInString := make([]bool, ASCIISIZE) // default value is false

	//iterate input string and set the corresponding index in the bool array to true
	for _, c := range str {
		isInString[c] = true
	}

	//iterate the bool array and count the number of true values
	mappingArray := make([]byte, ASCIISIZE)
	number := 0
	for i, b := range isInString {
		if b {
			mappingArray[i] = byte(number)
			number++
		}
	}

	//make a new array of ints with the same length as the input string
	newStrByte := make([]byte, len(str))

	//iterate the input string and set the corresponding index in the new array to the value of the mapping array
	for i, c := range str {
		newStrByte[i] = mappingArray[byte(c)]
	}

	return string(newStrByte), number

}
