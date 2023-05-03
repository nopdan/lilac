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
	ruleSli := strings.Split(rules, ",")
	fmt.Println(ruleSli)

	// 处理特殊规则 ab++
	if len(ruleSli) == 1 {
		if r, ok := strings.CutSuffix(ruleSli[0], "++"); ok {
			e.sRule = parseRule(r)
			return
		}
	}
	for _, r := range ruleSli {
		tmp := strings.Split(r, ":")
		if len(tmp) != 2 {
			fmt.Printf("规则解析错误: %v\n", rules)
			panic("")
		}
		// 冒号前面表示词长
		length := len(tmp[0])
		if length != 0 {
			var err error
			length, err = strconv.Atoi(tmp[0])
			if err != nil {
				fmt.Printf("规则解析错误: %v\n", rules)
				panic(err)
			}
		}
		// A+B+a1+b1
		tmp = strings.Split(tmp[1], "+")
		rl := make([]rule, len(tmp))
		for i, v := range tmp {
			rl[i] = parseRule(v)
		}
		e.Rule[length] = rl
	}
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
