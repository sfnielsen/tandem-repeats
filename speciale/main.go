package main

import (
	"fmt"
	"speciale/naive"
	"speciale/suffixtree"
)

func main() {
	// Create an instance of NaiveTree
	naiveTree := &naive.NaiveTree{}
	random_node := new(suffixtree.SuffixTreeNode)

	// Call functions on the NaiveTree instance
	result := naiveTree.EdgeLength(random_node) // Adjust the function name based on your implementation

	// Print or use the result as needed
	fmt.Println("Result:", result)
}
