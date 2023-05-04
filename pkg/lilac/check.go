package lilac

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func (c *config) runCheck(scan *bufio.Scanner) {
	dict := make(map[string][]string)
	readCheck(scan, dict)
	for word, codes := range dict {
		gen := c.encoder.Encode(word, []string{})
		if !contains(gen, codes) {
			fmt.Printf("Check Error! 词组: %v %v\t生成: %v\n", word, codes, gen)
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
		if !isHan(word) {
			continue
		}
		if _, ok := ret[word]; !ok {
			ret[word] = make([]string, 0)
		}
		ret[word] = append(ret[word], code)
	}
}

func contain(code string, codes []string) bool {
	for i := range codes {
		if strings.HasPrefix(codes[i], code) {
			return true
		}
	}
	return false
}

func contains(gen []string, codes []string) bool {
	flag := false
	for _, code := range gen {
		if contain(code, codes) {
			flag = true
		}
	}
	return flag
}

func isHan(word string) bool {
	flag := true
	for _, r := range word {
		if !unicode.Is(unicode.Han, r) {
			flag = false
		}
	}
	return flag
}
