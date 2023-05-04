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

	c := lilac.NewConfig(input)
	dict := c.Do()
	data := c.Write(dict)

	os.WriteFile(output, data, 0666)
}
