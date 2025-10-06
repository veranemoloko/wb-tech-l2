package main

import (
	"fmt"
	"my_sort/internal/sorting"
)

func main() {
	opts := sorting.ParseFlags()

	err := sorting.SortLines(opts)
	if err != nil {
		fmt.Println(err)
	}

}
