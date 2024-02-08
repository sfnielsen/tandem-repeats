package stringgenerators

type StringGenerator interface {
	GenerateString(n int) string
}

const (
	AlphabetA       string = "A"
	AlphabetAB      string = "AB"
	AlphabetDNA     string = "ACGT"
	AlphabetProtein string = "ACDEFGHIKLMNPQRSTVWY"
	AlphabetASCII   string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)
