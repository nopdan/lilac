package main

import (
	"os"

	"github.com/flowerime/lilac/pkg/lilac"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		return
	}
	input := args[1]
	output := args[2]

	dict := lilac.Generate(input)
	lilac.WriteFile(dict, output)
}
