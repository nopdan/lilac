package lilac

import (
	"bufio"
	"strings"
)

// 读取tsv码表
func HandleText(text, dir string) map[string][]string {
	ret := make(map[string][]string)
	rd := strings.NewReader(text)
	scan := bufio.NewScanner(rd)
	handle(scan, ret, dir)
	return ret
}

// 递归处理
func handle(scan *bufio.Scanner, ret map[string][]string, dir string) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line, dir); err == nil {
			handle(sc, ret, dir)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) != 2 {
			continue
		}
		key := tmp[0]
		vals := strings.Split(tmp[1], " ")
		if _, ok := ret[key]; !ok {
			ret[key] = make([]string, 0)
		}
		ret[key] = append(ret[key], vals...)
	}
}
