package encoder

import (
	"errors"
	"fmt"
	"strings"

	m "github.com/flowerime/lilac/pkg/mapping"
	"github.com/flowerime/lilac/pkg/util"
	"github.com/flowerime/rose/pkg/zhuyin"
)

type Encoder struct {
	Correct map[string][]string

	Rule    map[int][]rule
	sRule   rule // 特殊规则 ab++
	Char    map[string][]string
	Mapping *m.Mapping
}

func NewEncoder(rules string) *Encoder {
	e := new(Encoder)
	e.Rule = make(map[int][]rule)
	e.initRule(rules)
	return e
}

// pinyin is not required
func (e *Encoder) Encode(word string, pinyin []string) []string {
	if codes, ok := e.Correct[word]; ok {
		return codes
	}

	chars := []rune(word)
	length := len(chars)
	if length == 0 {
		return []string{}
	}

	// [A, B, a1, b1]
	rl, ok := e.Rule[length]
	if !ok {
		rl = e.Rule[0]
		// 单字规则
		if length == 1 {
			rl = []rule{
				{isYin: true, idxChar: 0, idxCode: 0},  // 音部分
				{isYin: false, idxChar: 0, idxCode: 0}, // 形部分
			}
		}
	}

	// 形码用不到音
	if e.Mapping == nil {
		one := e.encodeOne(chars, []string{}, rl)
		// fmt.Printf("? 词组: %v, 生成: %v\n", word, one)
		return one
	}

	// 处理音码
	ret := make([]string, 0)
	var pycodes [][]string
	// 词库中没有拼音，需要自动注音
	if len(pinyin) == 0 {
		pinyin = zhuyin.Get(word)
	}
	// zhi shi => [[ai ui], [ai vi], [ei ui], [ei vi]]
	pycodes = e.Mapping.FromPinyin(pinyin)

	for _, pycode := range pycodes {
		// [ai, ui]
		one := e.encodeOne(chars, pycode, rl)
		ret = append(ret, one...)
	}
	ret = util.RmRepeat(ret)
	// fmt.Printf("? 词组: %v, 拼音: %v, 转换后: %v, 生成: %v\n", word, pinyin, pycodes, ret)
	return ret
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
			codes = []string{pycode[idx]}
		} else {
			if idx == -1 {
				idx = len(chars) - 1
			}
			codes = e.Char[string(chars[idx])]
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
