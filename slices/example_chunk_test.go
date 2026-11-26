package slices_test

import (
	"fmt"

	"github.com/tomcz/gotools/slices"
)

func ExampleChunk() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	output := slices.Chunk(input, 3)

	fmt.Println(input)
	fmt.Println(output)

	// Output:
	// [1 2 3 4 5 6 7 8 9 10]
	// [[1 2 3] [4 5 6] [7 8 9 10]]
}
