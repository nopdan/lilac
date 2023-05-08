package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flowerime/lilac/pkg/lilac"
	"github.com/flowerime/pinyin"
)

func main() {
	var input, output string
	args := os.Args
	if len(args) >= 2 {
		input = args[1]
	}
	if len(args) >= 3 {
		output = args[2]
	}
	run(input, output)
}

func run(input, output string) {
	if output == "" {
		output = "输出/结果.txt"
	}
	// 创建文件夹
	dir := filepath.Dir(output)
	os.MkdirAll(dir, os.ModePerm)

	conf := newConfig(input)
	dict := conf.Build()
	lilac.WriteFile(dict, output)
	fmt.Printf("output: %v\n", output)

	// 错码校验
	chk := lilac.NewChecker()
	chk.Check(conf)
	var buf bytes.Buffer
	if len(chk.MisMatch) != 0 {
		buf.WriteString("词条\t编码\t生成的编码\t拼音\n")
		buf.WriteString("-----------------------------\n")
		for _, v := range chk.MisMatch {
			buf.WriteString(v.Word)
			buf.WriteByte('\t')
			buf.WriteString(strings.Join(v.Codes, " "))
			buf.WriteByte('\t')
			buf.WriteString(strings.Join(v.Gen, " "))
			buf.WriteByte('\t')
			buf.WriteString(strings.Join(v.Pinyin, " "))
			buf.WriteByte('\n')
		}
		os.WriteFile(filepath.Join(dir, "错码校验.txt"), buf.Bytes(), 0666)
		buf.Reset()
	}

	// 空码校验
	if len(chk.Empty) != 0 {
		os.MkdirAll(filepath.Dir(output), os.ModePerm)
		buf.WriteString("空码\t后续\n")
		buf.WriteString("-----------------------------\n")
		li := make([][2]string, 0, len(chk.Empty))
		for pre, entries := range chk.Empty {
			li = append(li, [2]string{pre, strings.Join(entries, " ")})
		}
		sort.Slice(li, func(i, j int) bool {
			return li[i][0] < li[j][0]
		})
		for i := range li {
			buf.WriteString(li[i][0])
			buf.WriteByte('\t')
			buf.WriteString(li[i][1])
			buf.WriteByte('\n')
		}
		os.WriteFile(filepath.Join(dir, "空码校验.txt"), buf.Bytes(), 0666)
	}
}

func newConfig(path string) *lilac.Config {
	py := pinyin.New()
	py.AddFile("./pinyin-data/pinyin.txt")
	py.AddFile("./pinyin-data/duoyin.txt")
	py.AddFile("./pinyin-data/correct.txt")
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
