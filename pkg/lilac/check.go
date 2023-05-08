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

type checker struct {
	dict     map[string][]string // 词的所有编码
	allCodes map[string]string   // 所有出现过的编码
}

func (c *Config) initChecker() {
	chk := new(checker)
	chk.dict = make(map[string][]string)
	chk.allCodes = make(map[string]string)
	c.checker = chk
}

// 校验码表，返回错误编码列表和空码列表
func (c *Config) Check() {
	rd := strings.NewReader(c.check)
	scan := bufio.NewScanner(rd)
	chk := c.checker
	enc := c.encoder

	chk.read(scan, c.Dir)
	for word, codes := range chk.dict {
		gen := enc.Encode(word, []string{})
		// 形码方案
		if enc.Mapping == nil {
			if !contains(gen, codes) {
				c.MisMatch = append(c.MisMatch, &misMatch{word, codes, gen, []string{}})
			}
			continue
		}
		yin := enc.Pinyin.Match(word)
		// 编码不匹配，并且拼音长和词长相等
		if !contains(gen, codes) && utf8.RuneCountInString(word) == len(yin) {
			c.MisMatch = append(c.MisMatch, &misMatch{word, codes, gen, yin})
		}
	}
	// 空码
	for code, word := range chk.allCodes {
		if len(code) == 1 {
			continue
		}
		pre := code[:len(code)-1] // 去掉最后一码
		if _, ok := chk.allCodes[pre]; !ok {
			if _, ok := c.Empty[pre]; !ok {
				c.Empty[pre] = make([]string, 0)
			}
			c.Empty[pre] = append(c.Empty[pre], word+"#"+code)
		}
	}
}

// ret 词的所有编码
func (chk *checker) read(scan *bufio.Scanner, dir string) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line, dir); err == nil {
			chk.read(sc, dir)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) < 2 {
			continue
		}
		word, code := tmp[0], tmp[1]
		if _, ok := chk.allCodes[code]; !ok {
			chk.allCodes[code] = word
		}
		// 忽略单字
		if utf8.RuneCountInString(word) == 1 {
			continue
		}
		// 忽略包含非汉字
		if !ku.IsHan(word) {
			continue
		}
		if _, ok := chk.dict[word]; !ok {
			chk.dict[word] = make([]string, 0)
		}
		chk.dict[word] = append(chk.dict[word], code)
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
