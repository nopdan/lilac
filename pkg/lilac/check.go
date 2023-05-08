package lilac

import (
	"bufio"
	"strings"
	"unicode/utf8"

	"github.com/nopdan/ku"
)

type misMatch struct {
	Word   string
	Codes  []string // 码表中的编码
	Gen    []string // 生成的编码
	Pinyin []string // 自动标注的读音
}

type Checker struct {
	dict     map[string][]string // 词的所有编码
	allCodes map[string]string   // 所有出现过的编码

	MisMatch []*misMatch // 不匹配的编码
	// key: 空码, value: 后续词和编码
	Empty map[string][]string // 空码
}

func NewChecker() *Checker {
	c := new(Checker)
	c.dict = make(map[string][]string)
	c.allCodes = make(map[string]string)

	c.MisMatch = make([]*misMatch, 0)
	c.Empty = make(map[string][]string, 0)
	return c
}

// 校验码表，返回错误编码列表和空码列表
func (c *Checker) Check(conf *Config) {
	rd := strings.NewReader(conf.check)
	scan := bufio.NewScanner(rd)

	c.read(scan, conf.dir)
	for word, codes := range c.dict {
		gen := conf.encoder.Encode(word, []string{})
		// 形码方案
		if conf.encoder.Mapping == nil {
			if !contains(gen, codes) {
				c.MisMatch = append(c.MisMatch, &misMatch{word, codes, gen, []string{}})
			}
			continue
		}
		yin := conf.encoder.Pinyin.Match(word)
		// 编码不匹配，并且拼音长和词长相等
		if !contains(gen, codes) && utf8.RuneCountInString(word) == len(yin) {
			c.MisMatch = append(c.MisMatch, &misMatch{word, codes, gen, yin})
		}
	}
	// 空码
	for code, word := range c.allCodes {
		if len(code) == 1 {
			continue
		}
		pre := code[:len(code)-1] // 去掉最后一码
		if _, ok := c.allCodes[pre]; !ok {
			if _, ok := c.Empty[pre]; !ok {
				c.Empty[pre] = make([]string, 0)
			}
			c.Empty[pre] = append(c.Empty[pre], word+"#"+code)
		}
	}
}

// ret 词的所有编码
func (c *Checker) read(scan *bufio.Scanner, dir string) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line, dir); err == nil {
			c.read(sc, dir)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) < 2 {
			continue
		}
		word, code := tmp[0], tmp[1]
		if _, ok := c.allCodes[code]; !ok {
			c.allCodes[code] = word
		}
		// 忽略单字
		if utf8.RuneCountInString(word) == 1 {
			continue
		}
		// 忽略包含非汉字
		if !ku.IsHan(word) {
			continue
		}
		if _, ok := c.dict[word]; !ok {
			c.dict[word] = make([]string, 0)
		}
		c.dict[word] = append(c.dict[word], code)
	}
}

// 包含编码，或者简码
func contain(codes []string, code string) bool {
	for _, item := range codes {
		if strings.HasPrefix(code, item) {
			return true
		}
	}
	return false
}

func contains(gen []string, codes []string) bool {
	flag := false
	for _, code := range gen {
		if contain(codes, code) {
			flag = true
		}
	}
	return flag
}
