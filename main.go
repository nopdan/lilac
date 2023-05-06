package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flowerime/lilac/pkg/lilac"
	"github.com/flowerime/pinyin"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		return
	}
	input := args[1]
	output := args[2]

	conf := lilac.NewConfig(input, NewPinyin())
	dict := conf.Build()
	lilac.WriteFile(dict, output)
}

func NewPinyin() *pinyin.Pinyin {
	py := pinyin.New()
	filepath.Walk("data", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		if !info.IsDir() {
			py.AddFile(path)
		}
		return nil
	})
	return py
}
