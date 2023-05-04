package lilac

import (
	"bufio"
	"fmt"
	"strings"
)

func (c *config) runCheck(scan *bufio.Scanner) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line); err == nil {
			c.runCheck(sc)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) < 2 {
			continue
		}
		word, code := tmp[0], tmp[1]
		pinyin := []string{}
		if len(tmp) == 3 {
			py := tmp[2]
			pinyin = strings.Split(py, " ")
		}
		codes := c.encoder.Encode(word, pinyin)
		if !contain(code, codes) {
			fmt.Printf("编码错误！词组: %v 编码: %v 可能的正确编码: %v\n", word, code, codes)
		}
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
