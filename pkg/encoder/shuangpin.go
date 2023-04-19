package encoder

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/flowerime/lilac/pkg/util"
)

type Shuangpin struct {
	Key    map[string][]byte   // 声母或韵母可以对应多个键
	Yinjie map[string][]string // 音节可以对应多个按键组合
}

// 从文件初始化双拼方案
func NewShuangpin(path string) *Shuangpin {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	brd := bytes.NewReader(data)
	scan := bufio.NewScanner(brd)

	s := new(Shuangpin)
	s.Key = make(map[string][]byte)
	s.Yinjie = make(map[string][]string)

	// 读取基本映射
	for scan.Scan() {
		text := scan.Text()
		if strings.HasPrefix(text, "//") {
			continue
		}
		if strings.HasPrefix(text, "###") {
			break
		}
		tmp := strings.Split(text, " ")
		key := strings.ToLower(tmp[0])
		for i := 1; i < len(tmp); i++ {
			if _, ok := s.Key[tmp[i]]; !ok {
				s.Key[tmp[i]] = make([]byte, 0, 1)
			}
			s.Key[tmp[i]] = append(s.Key[tmp[i]], key[0])
		}
	}

	// 读取自定义音节
	for scan.Scan() {
		text := scan.Text()
		if strings.HasPrefix(text, "//") {
			continue
		}
		tmp := strings.Split(text, " ")
		if _, ok := s.Yinjie[tmp[0]]; !ok {
			s.Yinjie[tmp[0]] = make([]string, 0, 1)
		}
		for i := 1; i < len(tmp); i++ {
			if len(tmp[i]) != 2 {
				continue
			}
			s.Yinjie[tmp[0]] = append(s.Yinjie[tmp[0]], tmp[i])
		}
	}
	return s
}

// zhi shi => [ [ai ei], [ui vi] ] => [[ai ui], [ai vi], [ei ui], [ei vi]]
func (s *Shuangpin) FromPinyin(py []string) [][]string {
	ret := make([][]string, len(py))
	for i, yinjie := range py {
		ret[i] = s.FromYinjie(yinjie)
	}
	ret = util.Product(ret)
	return ret
}

// 一个音节转换为按键组合 shi => ["ui", "vi"]
func (s *Shuangpin) FromYinjie(yinjie string) []string {
	// 自定义音节
	if v, ok := s.Yinjie[yinjie]; ok {
		return v
	}
	//
	var sm, ym string
	switch yinjie[0] {
	case 'a', 'o', 'e':
		sm = "#"
		ym = yinjie
	default:
		if yinjie[1] == 'h' {
			sm = yinjie[:2]
			ym = yinjie[2:]
		} else {
			sm = yinjie[:1]
			ym = yinjie[1:]
		}
	}
	smKeys := s.Key[sm]
	ymKeys := s.Key[ym]
	keys := util.Product([][]byte{smKeys, ymKeys})
	ret := make([]string, len(keys))
	for i, v := range keys {
		ret[i] = string(v)
	}
	return ret
}
