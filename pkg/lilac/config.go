package lilac

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/flowerime/lilac/pkg/encoder"
	m "github.com/flowerime/lilac/pkg/mapping"
	"gopkg.in/ini.v1"
)

type config struct {
	Type int    `ini:"Type"`
	Rule string `ini:"Rule"`
	Sort int    `ini:"Sort"`

	dict    string
	encoder *encoder.Encoder
}

func newConfig(path string) *config {
	// 手动解析下列 Section
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"Dict", "Correct", "Char", "Mapping"},
	}, path)
	if err != nil {
		panic(err)
	}

	c := new(config)
	config := cfg.Section("Config")
	err = config.MapTo(c)
	if err != nil {
		fmt.Println(err)
	}

	var text string

	enc := encoder.NewEncoder(c.Rule)
	text = cfg.Section("Correct").Body()
	enc.Correct = HandleText(text)
	text = cfg.Section("Char").Body()
	enc.Char = HandleText(text)
	text = cfg.Section("Mapping").Body()
	data := HandleText(text)
	if len(data) != 0 {
		enc.Mapping = m.NewMapping(data)
	}
	c.encoder = enc

	c.dict = cfg.Section("Dict").Body()
	fmt.Printf("c: %+v\n", c)
	return c
}

// 生成码表
func Generate(path string) [][]string {
	c := newConfig(path)

	rd := strings.NewReader(c.dict)
	scan := bufio.NewScanner(rd)
	ret := make([][]string, 0)
	ret = c.run(scan, ret, false)

	if c.Sort == 1 {
		sort.SliceStable(ret, func(i, j int) bool {
			return ret[i][1] < ret[j][1]
		})
	}
	return ret
}

// 递归
func (c *config) run(scan *bufio.Scanner, ret [][]string, flag bool) [][]string {
	for scan.Scan() {
		line := scan.Text()
		if sc, newFlag, err := include(line); err == nil {
			ret = c.run(sc, ret, newFlag)
			continue
		}

		tmp := strings.Split(line, "\t")
		word, ok := strings.CutPrefix(tmp[0], "?")
		if flag {
			ok = true
		}
		entry := []string{word}
		if len(tmp) == 1 {
			gen := c.encoder.Encode(word, []string{})
			// fmt.Printf("? 词组: %v, 生成: %v\n", word, gen)
			entry = append(entry, gen...)
			ret = append(ret, flat(entry)...)
			continue
		}

		if len(tmp) == 2 {
			if !ok {
				ret = append(ret, tmp)
				continue
			}
			// ? 号开头为自动造词
			py := strings.Split(tmp[1], " ")
			gen := c.encoder.Encode(word, py)
			// fmt.Printf("? 词组: %v, 拼音: %v, 生成: %v\n", word, py, gen)
			entry = append(entry, gen...)
			ret = append(ret, flat(entry)...)
		} else {
			fmt.Println("错误:", line)
		}
	}
	return ret
}

// 展开一词多编码
func flat(entry []string) [][]string {
	if len(entry) <= 2 {
		return [][]string{entry}
	}
	ret := make([][]string, 0, len(entry)-1)
	for i := 1; i < len(entry); i++ {
		ret = append(ret, []string{entry[0], entry[i]})
	}
	// fmt.Println(entry, ret)
	return ret
}
