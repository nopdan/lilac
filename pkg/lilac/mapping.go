package lilac

import (
	"bufio"
	"strings"
)

// 读取 tsv 码表，Char 和 Mapping
func readMap(text, dir string) map[string][]string {
	ret := make(map[string][]string)
	rd := strings.NewReader(text)
	scan := bufio.NewScanner(rd)
	_readMap(scan, dir, ret)
	return ret
}

// 递归处理
func _readMap(scan *bufio.Scanner, dir string, ret map[string][]string) {
	for scan.Scan() {
		line := scan.Text()
		// 引入其他文件
		if sc, _, err := include(line, dir); err == nil {
			_readMap(sc, dir, ret)
			continue
		}
		tmp := strings.Split(line, "\t")
		if len(tmp) != 2 {
			continue
		}
		key := tmp[0]
		values := strings.Split(tmp[1], " ")
		if _, ok := ret[key]; !ok {
			ret[key] = make([]string, 0)
		}
		ret[key] = append(ret[key], values...)
	}
}
