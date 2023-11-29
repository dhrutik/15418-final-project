package main

import (
	"fmt"
	"main/seq_tree"
)

func main() {
	tree := seq_tree.NewTree()
	values := []byte("test")
	for i := 0; i < 1000000; i++ {
		tree.Insert(i, values)
	}
	fmt.Println("Done")
}
