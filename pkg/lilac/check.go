package lilac

import (
	"bufio"
	"strings"
	"unicode/utf8"

	"github.com/nopdan/ku"
)

type CheckResult struct {
	Word   string
	Codes  []string // 码表中的编码
	Gen    []string // 生成的编码
	Pinyin []string // 自动标注的读音
}

// 校验码表
func (c *config) Check() []*CheckResult {
	rd := strings.NewReader(c.check)
	scan := bufio.NewScanner(rd)
	dict := make(map[string][]string)
	readCheck(scan, dict)

	wrong := make([]*CheckResult, 0)
	for word, codes := range dict {
		gen := c.encoder.Encode(word, []string{})
		// 形码方案
		if c.encoder.Mapping == nil {
			if !contains(gen, codes) {
				wrong = append(wrong, &CheckResult{word, codes, gen, []string{}})
			}
			continue
		}
		yin := c.encoder.Pinyin.Match(word)
		// 编码不匹配，并且拼音长和词长相等
		if !contains(gen, codes) && utf8.RuneCountInString(word) == len(yin) {
			wrong = append(wrong, &CheckResult{word, codes, gen, yin})
		}
	}
	return wrong
}

func readCheck(scan *bufio.Scanner, ret map[string][]string) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line); err == nil {
			readCheck(sc, ret)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) < 2 {
			continue
		}
		word, code := tmp[0], tmp[1]
		// 忽略单字
		if utf8.RuneCountInString(word) == 1 {
			continue
		}
		// 忽略包含非汉字
		if !ku.IsHan(word) {
			continue
		}
		if _, ok := ret[word]; !ok {
			ret[word] = make([]string, 0)
		}
		ret[word] = append(ret[word], code)
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
