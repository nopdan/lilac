package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowerime/lilac/pkg/lilac"
	"github.com/flowerime/pinyin"
)

func main() {
	var input string
	output := "output.txt"

	args := os.Args
	if len(args) >= 2 {
		input = args[1]
	}
	if len(args) >= 3 {
		output = args[2]
	}

	conf := newConfig(input)
	dict := conf.Build()
	lilac.WriteFile(dict, output)
	fmt.Printf("output: %v\n", output)

	check := conf.Check()
	var buf bytes.Buffer
	buf.WriteString("词条\t编码\t生成的编码\t拼音")
	for _, v := range check {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Codes, " "))
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Gen, " "))
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Pinyin, " "))
		buf.WriteByte('\n')
	}
	os.WriteFile("check.txt", buf.Bytes(), 0666)
}

func newConfig(path string) *lilac.Config {
	py := pinyin.New()
	py.AddFile("./pinyin-data/pinyin.txt")
	py.AddFile("./pinyin-data/duoyin.txt")
	py.AddFile("./pinyin-data/correct.txt")
	py.AddFile("correct.txt")
	// appendDir(py, "data")

	conf := lilac.NewConfig(path, py)
	return conf
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
