package lilac

import (
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/nopdan/ku"
)

func (c *config) runCheck(scan *bufio.Scanner) {
	dict := make(map[string][]string)
	readCheck(scan, dict)
	for word, codes := range dict {
		gen := c.encoder.Encode(word, []string{})
		yin := c.encoder.Pinyin.Match(word)
		// 编码不匹配，并且拼音长和词长相等
		if !contains(gen, codes) && utf8.RuneCountInString(word) == len(yin) {
			fmt.Printf("Check Error! 词组: %v %v\t生成: %v 读音: %v\n", word, codes, gen, c.encoder.Pinyin.Match(word))
		}
	}
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
