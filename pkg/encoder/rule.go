package encoder

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
2:A+B+a1+b1,
3:A1+B1+C1+a1+b1+c1,
:A1+B1+C1+Z1+a1+b1

special   ab++
*/

// 一条规则
type rule struct {
	dc  []unit // 定长规则
	zj  []unit // 不定长（整句）规则
	sep string // 整句编码分隔符
}

// 规则的最小单位，如 a1 B
type unit struct {
	isYin   bool // 大写字母表示音
	idxChar int  // 字的索引，-1 表示最后一个字
	idxCode int
}

// 解析规则
func (e *Encoder) initRule(rules string) {
	rules = TrimSpace(rules)
	lines := strings.Split(rules, ",")
	fmt.Println(lines)

	for _, line := range lines {
		tmp := strings.Split(line, ":")
		if len(tmp) != 2 {
			fmt.Printf("规则解析错误: %v, line: %v\n", rules, line)
			panic("")
		}
		// 冒号前面表示词长
		length, err := strconv.Atoi(tmp[0])
		if err != nil && tmp[0] != "" {
			fmt.Printf("规则解析错误: %v, line: %v\n", rules, line)
			panic(err)
		}
		rl := rule{}
		tmp = strings.Split(tmp[1], "..")
		// 整句规则
		if len(tmp) == 2 {
			rl.zj = parseUnits(tmp[0])
			if tmp[1] == "_" {
				tmp[1] = " "
			}
			rl.sep = tmp[1]
		} else {
			// 定长规则
			rl.dc = parseUnits(tmp[0])
		}
		e.Rule[length] = rl
	}
}

// 解析规则 A+a+A2+a11
func parseUnits(r string) []unit {
	tmp := strings.Split(r, "+")
	rl := make([]unit, len(tmp))
	for i, v := range tmp {
		rl[i] = parseUnit(v)
	}
	return rl
}

// 解析规则 A a A2 a11
func parseUnit(r string) unit {
	unit := unit{}

	parse := func(a, z byte) {
		unit.idxChar = int(r[0] - a)
		if r[0] == z {
			unit.idxChar = -1 // 最后一个字
		}
		if len(r) >= 2 {
			var err error
			unit.idxCode, err = strconv.Atoi(r[1:]) // 编码的索引
			if err != nil {
				fmt.Printf("err rule part: %v\n", r)
			}
		}
	}
	if 'A' <= r[0] && r[0] <= 'Z' {
		unit.isYin = true
		parse('A', 'Z')
	} else if 'a' <= r[0] && r[0] <= 'z' {
		// rl.isYin = false
		parse('a', 'z')
	} else {
		fmt.Printf("err rule part: %v\n", r)
	}
	return unit
}

func TrimSpace(s string) string {
	new := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			new = append(new, r)
		}
	}
	return string(new)
}
