package encoder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nopdan/ku"
	m "github.com/nopdan/lilac/pkg/mapping"
	"github.com/nopdan/pinyin"
)

type Encoder struct {
	Rule    map[int]rule
	Char    map[string][]string
	Mapping *m.Mapping
	Pinyin  *pinyin.Pinyin
}

func NewEncoder(rules string) *Encoder {
	e := new(Encoder)
	e.Rule = make(map[int]rule)
	e.initRule(rules)
	return e
}

// 自动编码
// pinyin is not required
func (e *Encoder) Encode(word string, pinyin []string) []string {
	if word == "" {
		return []string{}
	}
	chars := []rune(word)
	length := len(chars)
	rl := e.getRule(length)

	// 形码用不到音
	if e.Mapping == nil || e.Pinyin == nil {
		one := e.encodeOne(chars, pinyin, rl)
		// fmt.Printf("? 词组: %v, 生成: %v\n", word, one)
		return one
	}

	// 处理音码
	ret := make([]string, 0)
	var pycodes [][]string
	// 词库中没有拼音，需要自动注音
	if len(pinyin) == 0 {
		pinyin = e.Pinyin.Match(word)
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
	ret = ku.Unique(ret)
	// fmt.Printf("? 词组: %v, 拼音: %v, 转换后: %v, 生成: %v\n", word, pinyin, pycodes, ret)
	return ret
}

func (e *Encoder) getRule(length int) rule {
	// [A, B, a1, b1]
	if rl, ok := e.Rule[length]; ok {
		return rl
	}
	// 单字默认规则
	if length == 1 {
		return rule{dc: parseUnits("A+a")}
	}
	return e.Rule[0]
}

// 一组拼音生成的编码
func (e *Encoder) encodeOne(chars []rune, pycode []string, rl rule) []string {
	// fmt.Printf("rl: %v\n", rl)
	tmp := make([][]string, 0)

	// 取一码
	encode := func(idxChar, idxCode int, isYin bool) {
		var codes []string
		if isYin {
			if idxChar == -1 {
				idxChar = len(pycode) - 1
			}
			if idxChar < len(pycode) {
				codes = []string{pycode[idxChar]}
			}
		} else {
			if idxChar == -1 {
				idxChar = len(chars) - 1
			}
			if idxChar < len(chars) {
				codes = e.Char[string(chars[idxChar])]
			}
		}
		// fmt.Printf("codes: %v\n", codes)
		// 等于 0 时取全部编码
		if idxCode != 0 {
			var err error
			codes, err = cut(codes, idxCode)
			if err != nil {
				fmt.Println(err, "编码错误", string(chars), pycode)
				return
			}
		}
		if len(codes) != 0 {
			tmp = append(tmp, codes)
		}
	}

	// 定长规则为空，则为整句规则
	if len(rl.dc) != 0 {
		for _, unit := range rl.dc {
			encode(unit.idxChar, unit.idxCode, unit.isYin)
		}
	} else {
		for i := range chars {
			for _, unit := range rl.zj {
				encode(i, unit.idxCode, unit.isYin)
			}
		}
	}

	tmp = ku.Product(tmp)
	ret := make([]string, len(tmp))
	for i := range tmp {
		ret[i] = strings.Join(tmp[i], rl.sep)
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
