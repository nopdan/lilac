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
	var input string
	output := "output.txt"

	args := os.Args
	if len(args) >= 2 {
		input = args[1]
	} else if len(args) >= 3 {
		output = args[2]
	}

	dict := build(input)
	lilac.WriteFile(dict, output)
}

func build(path string) [][]string {
	py := pinyin.New()
	py.AddFile("./pinyin-data/pinyin.txt")
	py.AddFile("./pinyin-data/duoyin.txt")
	py.AddFile("./pinyin-data/correct.txt")
	py.AddFile("correct.txt")
	// appendDir(py, "data")

	conf := lilac.NewConfig(path, py)
	return conf.Build()
}

// 追加目录下的所有文件数据
func AppendDir(py *pinyin.Pinyin, dir string) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		if !info.IsDir() {
			py.AddFile(path)
		}
		return nil
	})
}
