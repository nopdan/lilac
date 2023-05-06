package shorten

import (
	"strconv"
	"strings"
)

// rule 默认1，空无限
// 1:0,2:3,3:2,6: => 1:0,2:3,3:2,4:1,5:1,6:99999999
func Shorten(table [][]string, rule string) {
	rl := parseRule(rule)
	countMap := make(map[string]int)
	for i := 0; i < len(table); i++ {
		wc := table[i]
		if len(wc) != 2 {
			continue
		}
		for j := 1; j <= len(wc[1]); j++ {
			curr := wc[1][:j]
			count := countMap[curr]
			if count < rl[j] {
				wc[1] = curr
				table[i] = wc
				countMap[curr]++
				break
			}
		}
	}
	// fmt.Println(countMap)
}

// [0,3,2,1,1,1e5]
func parseRule(rule string) []int {
	ret := make([]int, 0)
	rule = strings.ReplaceAll(rule, " ", "")
	r := strings.Split(rule, ",")
	for _, v := range r {
		tmp := strings.Split(v, ":")
		if len(tmp) != 2 {
			continue
		}
		pos, _ := strconv.Atoi(tmp[0])
		if pos < 1 {
			continue
		}
		var val int
		if tmp[1] == "" {
			val = 1e5
		} else {
			val, _ = strconv.Atoi(tmp[1])
		}
		setVal(&ret, pos, val)
	}
	// fmt.Println(ret)
	return ret
}

func setVal(sli *[]int, pos int, val int) {
	for pos > len(*sli)-1 {
		*sli = append(*sli, 1)
	}
	(*sli)[pos] = val
}
