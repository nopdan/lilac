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

map[int][]rule
*/
type rule struct {
	isYin   bool // 大写字母表示音
	idxChar int  // 字的索引，-1 表示最后一个字
	idxCode int
}

// 解析规则
func (e *Encoder) initRule(rules string) {
	rules = TrimSpace(rules)
	lines := strings.Split(rules, ",")
	fmt.Println(lines)

	// 处理特殊规则 a1+a2..
	special := false
	if len(lines) == 1 {
		lines[0], special = strings.CutSuffix(lines[0], "..")
	}
	for _, line := range lines {
		tmp := strings.Split(line, ":")
		if len(tmp) == 1 {
			// 没有 : 特殊规则
			if special {
				e.sRule = parse(line)
				return
			}
			fmt.Printf("规则解析错误: %v, line: %v\n", rules, line)
			panic("")
		}
		// 冒号前面表示词长
		length, err := strconv.Atoi(tmp[0])
		if len(tmp) != 2 || err != nil && tmp[0] != "" {
			fmt.Printf("规则解析错误: %v, line: %v\n", rules, line)
			panic(err)
		}
		e.Rule[length] = parse(tmp[1])
	}
}

// 解析规则 A+a+A2+a11
func parse(r string) []rule {
	tmp := strings.Split(r, "+")
	rl := make([]rule, len(tmp))
	for i, v := range tmp {
		rl[i] = parseRule(v)
	}
	return rl
}

// 解析规则 A a A2 a11
func parseRule(r string) rule {
	rl := rule{}

	parse := func(a, z byte) {
		rl.idxChar = int(r[0] - a)
		if r[0] == z {
			rl.idxChar = -1 // 最后一个字
		}
		if len(r) >= 2 {
			var err error
			rl.idxCode, err = strconv.Atoi(r[1:]) // 编码的索引
			if err != nil {
				fmt.Printf("err rule part: %v\n", r)
			}
		}
	}
	if 'A' <= r[0] && r[0] <= 'Z' {
		rl.isYin = true
		parse('A', 'Z')
	} else if 'a' <= r[0] && r[0] <= 'z' {
		// rl.isYin = false
		parse('a', 'z')
	} else {
		fmt.Printf("err rule part: %v\n", r)
	}
	return rl
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
