package lilac

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/flowerime/lilac/pkg/encoder"
	m "github.com/flowerime/lilac/pkg/mapping"
	"github.com/flowerime/pinyin"
	"gopkg.in/ini.v1"
)

type Config struct {
	Rule string `ini:"Rule"`
	Sort bool   `ini:"Sort"`

	Keep     bool   `ini:"保留单字全码"`
	CharRule string `ini:"单字简码规则"`
	WordRule string `ini:"词组简码规则"`

	dict    string
	check   string
	encoder *encoder.Encoder
}

func NewConfig(path string, py *pinyin.Pinyin) *Config {
	// 手动解析下列 Section
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"Dict", "Char", "Mapping", "Check"},
	}, path)
	if err != nil {
		panic(err)
	}

	c := new(Config)
	config := cfg.Section("Config")
	err = config.MapTo(c)
	if err != nil {
		fmt.Println(err)
	}

	var text string

	enc := encoder.NewEncoder(c.Rule)
	enc.Pinyin = py
	text = cfg.Section("Char").Body()
	enc.Char = HandleText(text)
	text = cfg.Section("Mapping").Body()
	data := HandleText(text)
	if len(data) != 0 {
		enc.Mapping = m.NewMapping(data)
	}
	c.encoder = enc

	c.dict = cfg.Section("Dict").Body()
	c.check = cfg.Section("Check").Body()
	// fmt.Printf("c: %+v\n", c)
	return c
}

// 生成码表
func (c *Config) Build() [][2]string {
	rd := strings.NewReader(c.dict)
	scan := bufio.NewScanner(rd)
	ret := make([][2]string, 0)
	ret = c.run(scan, ret, false)

	s := newShortener(c)
	ret = s.Shorten(ret)

	if c.Sort {
		sort.SliceStable(ret, func(i, j int) bool {
			return ret[i][1] < ret[j][1]
		})
	}
	return ret
}

// 递归
func (c *Config) run(scan *bufio.Scanner, ret [][2]string, flag bool) [][2]string {
	for scan.Scan() {
		line := scan.Text()
		if line == "" {
			continue
		}

		if sc, newFlag, err := include(line); err == nil {
			ret = c.run(sc, ret, newFlag)
			continue
		}

		tmp := strings.Split(line, "\t")
		word, ok := strings.CutPrefix(tmp[0], "?")
		if flag {
			ok = true
		}
		if !ok && len(tmp) == 2 {
			ret = append(ret, [2]string{tmp[0], tmp[1]})
			continue
		}

		entry := []string{word}
		pinyin := []string{}
		if len(tmp) == 2 {
			pinyin = strings.Split(tmp[1], " ")
		}
		gen := c.encoder.Encode(word, pinyin)
		entry = append(entry, gen...)
		ret = append(ret, flat(entry)...)
	}
	return ret
}

// 展开一词多编码
func flat(entry []string) [][2]string {
	if len(entry) < 2 {
		return [][2]string{}
	}
	ret := make([][2]string, 0, len(entry)-1)
	for i := 1; i < len(entry); i++ {
		ret = append(ret, [2]string{entry[0], entry[i]})
	}
	// fmt.Println(entry, ret)
	return ret
}
