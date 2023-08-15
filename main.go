package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nopdan/ku"
	"github.com/nopdan/lilac/pkg/lilac"
)

func main() {
	var input string
	args := os.Args
	if len(args) >= 2 {
		input = args[1]
	}
	run(input)
}

func run(input string) {
	conf := lilac.NewConfig(input, "pinyin-data")
	conf.Build()
	name := ku.GetFileName(input)

	write := func(data []byte, suffix string) {
		if len(data) == 0 {
			return
		}
		path := filepath.Join(conf.Dir, name+"."+suffix+".txt")
		err := os.WriteFile(path, data, 0666)
		if err == nil {
			fmt.Printf("输出%s: %v\n", suffix, path)
		} else {
			fmt.Printf("err: %v\n", err)
		}
	}

	write(conf.OutputResult(), "结果")
	write(conf.OutputMisMatch(), "错码")
	write(conf.OutputEmpty(), "空码")
	fmt.Println()
}
