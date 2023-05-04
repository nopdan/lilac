package encoder

import (
	"errors"
	"fmt"
	"strings"

	m "github.com/flowerime/lilac/pkg/mapping"
	"github.com/flowerime/lilac/pkg/util"
	"github.com/flowerime/pinyin"
)

type Encoder struct {
	Correct map[string][]string

	Rule    map[int][]rule
	sRule   []rule // 特殊规则 a1+a2..
	Char    map[string][]string
	Mapping *m.Mapping
	py      *pinyin.Pinyin
}

func NewEncoder(rules string) *Encoder {
	e := new(Encoder)
	e.Rule = make(map[int][]rule)
	e.initRule(rules)
	e.py = pinyin.New()
	return e
}

// pinyin is not required
func (e *Encoder) Encode(word string, pinyin []string) []string {
	if word == "" {
		return []string{}
	}
	if codes, ok := e.Correct[word]; ok {
		return codes
	}
	chars := []rune(word)
	length := len(chars)
	rl := e.getRule(length)

	// 形码用不到音
	if e.Mapping == nil {
		one := e.encodeOne(chars, pinyin, rl)
		// fmt.Printf("? 词组: %v, 生成: %v\n", word, one)
		return one
	}

	// 处理音码
	ret := make([]string, 0)
	var pycodes [][]string
	// 词库中没有拼音，需要自动注音
	if len(pinyin) == 0 {
		pinyin = e.py.Match(word)
		// fmt.Println(word, pinyin)
	}
	// zhi shi => [[ai ui], [ai vi], [ei ui], [ei vi]]
	pycodes = e.Mapping.FromPinyin(pinyin)
	// fmt.Printf("注音: %v\t%v\t%v\n", word, pinyin, pycodes)

	for _, pycode := range pycodes {
		// [ai, ui]
		one := e.encodeOne(chars, pycode, rl)
		ret = append(ret, one...)
	}
	ret = util.RmRepeat(ret)
	// fmt.Printf("? 词组: %v, 拼音: %v, 转换后: %v, 生成: %v\n", word, pinyin, pycodes, ret)
	return ret
}

func (e *Encoder) getRule(length int) []rule {
	// [A, B, a1, b1]
	if rl, ok := e.Rule[length]; ok {
		return rl
	}
	// 单字默认规则
	if length == 1 {
		return []rule{
			{isYin: true, idxChar: 0, idxCode: 0},  // 音部分
			{isYin: false, idxChar: 0, idxCode: 0}, // 形部分
		}
	}
	// 特殊规则
	if len(e.sRule) != 0 {
		return e.sRule
	}
	return e.Rule[0]
}

// 一组拼音生成的编码
func (e *Encoder) encodeOne(chars []rune, pycode []string, rl []rule) []string {
	// fmt.Printf("rl: %v\n", rl)
	tmp := make([][]string, 0)
	for _, r := range rl {
		idx := r.idxChar
		var codes []string
		if r.isYin {
			if idx == -1 {
				idx = len(pycode) - 1
			}
			if idx < len(pycode) {
				codes = []string{pycode[idx]}
			}
		} else {
			if idx == -1 {
				idx = len(chars) - 1
			}
			if idx < len(chars) {
				codes = e.Char[string(chars[idx])]
			}
		}
		// fmt.Printf("codes: %v\n", codes)
		// 等于 0 时取全部编码
		if r.idxCode != 0 {
			var err error
			codes, err = cut(codes, r.idxCode)
			if err != nil {
				fmt.Println(err, "编码错误", string(chars), pycode)
				continue
			}
		}
		tmp = append(tmp, codes)
	}
	tmp = util.Product(tmp)
	ret := make([]string, len(tmp))
	for i := range tmp {
		ret[i] = merge(tmp[i])
	}
	return ret
}

// [ai, ui], 1 => [a, u]，idx是索引+1
func cut(codes []string, idx int) ([]string, error) {
	ret := make([]string, 0, len(codes))
	for i := range codes {
		if idx > len(codes[i]) {
			return ret, errors.New("index out of range")
		}
		ret = append(ret, string(codes[i][idx-1]))
	}
	return ret, nil
}

func merge(sli []string) string {
	var s strings.Builder
	for i := range sli {
		s.WriteString(sli[i])
	}
	return s.String()
}
