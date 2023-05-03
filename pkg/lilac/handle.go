package lilac

import (
	"bufio"
	"strings"
)

func handleText(text string) map[string][]string {
	ret := make(map[string][]string)
	rd := strings.NewReader(text)
	scan := bufio.NewScanner(rd)
	handle(scan, ret)
	return ret
}

// 递归处理
func handle(scan *bufio.Scanner, ret map[string][]string) {
	for scan.Scan() {
		line := scan.Text()
		if sc, _, err := include(line); err == nil {
			handle(sc, ret)
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
