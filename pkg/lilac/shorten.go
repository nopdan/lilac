package lilac

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nopdan/ku"
)

type shortener struct {
	charRule jmRule
	wordRule jmRule
	keep     bool
	countMap map[string]int // 每个编码出现的次数
	quanma   [][2]string    // 单字全码
}

func (c *Config) initShortener() {
	s := new(shortener)
	s.keep = c.Keep
	s.charRule = newRule(c.CharRule)
	s.wordRule = newRule(c.WordRule)
	s.countMap = make(map[string]int)
	s.quanma = make([][2]string, 0)
	c.shortener = s
}

// rule 默认1，空无限
// 1:0,2:3,3:2,6: => 1:0,2:3,3:2,4:1,5:1,6:99999999
func (c *Config) Shorten() {
	s := c.shortener
	// fmt.Printf("s: %v\n", s)
	for i := range c.Result {
		word, code := c.Result[i][0], c.Result[i][1]
		// 判断是否是词组
		isWord := utf8.RuneCountInString(word) > 1
		if isWord {
			// 规则为空，保留全码
			if s.wordRule == nil {
				s.countMap[code]++
				continue
			}
			c.Result[i] = s.word(c.Result[i])
			continue
		}
		// 单字
		if s.charRule == nil {
			s.countMap[code]++
			continue
		}
		c.Result[i] = s.char(c.Result[i])
	}
	// 合并单字全码
	if s.keep {
		c.Result = append(c.Result, s.quanma...)
	}
}

// 处理单字词条
func (s *shortener) char(entry [2]string) [2]string {
	code := entry[1]
	for j := range code {
		curr := code[:j+1]        // 截取编码
		count := s.countMap[curr] // 出现次数
		// 若当前编码出现次数小于规则里的
		if count < s.charRule.get(j+1) || j == len(code)-1 {
			s.countMap[curr]++
			if code != curr {
				// 保存简码的全码
				s.quanma = append(s.quanma, entry)
				code = curr
			}
			break
		}
	}
	entry[1] = code
	return entry
}

// 处理词组词条
func (s *shortener) word(entry [2]string) [2]string {
	code := entry[1]
	for j := range code {
		curr := code[:j+1]        // 截取编码
		count := s.countMap[curr] // 出现次数
		// 若当前编码出现次数小于规则里的
		if count < s.wordRule.get(j+1) || j == len(code)-1 {
			s.countMap[curr]++
			code = curr
			break
		}
	}
	entry[1] = code
	return entry
}

type jmRule map[int]int

func (j jmRule) get(length int) int {
	if v, ok := j[length]; ok {
		return v
	}
	return 1
}

func newRule(rule string) jmRule {
	rule = ku.TrimSpace(rule)
	if rule == "" {
		return nil
	}
	j := make(jmRule, 0)
	r := strings.Split(rule, ",")
	for _, v := range r {
		tmp := strings.Split(v, ":")
		if len(tmp) != 2 {
			continue
		}
		length, _ := strconv.Atoi(tmp[0])
		if length < 1 {
			continue
		}
		var val int
		if tmp[1] == "" {
			val = 1e5
		} else {
			val, _ = strconv.Atoi(tmp[1])
		}
		j[length] = val
	}
	return j
}
