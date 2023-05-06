package mapping

import (
	"strings"

	"github.com/nopdan/ku"
)

type Mapping struct {
	Sm_ym  map[string][]byte   // 声母或韵母可以对应多个键
	Yinjie map[string][]string // 音节可以对应多个按键组合
}

func NewMapping(data map[string][]string) *Mapping {
	s := new(Mapping)
	s.Sm_ym = make(map[string][]byte)
	s.Yinjie = make(map[string][]string)

	for k, v := range data {
		if len(k) == 1 {
			// 小写表示音节，大写表示按键
			switch k {
			case "a", "e", "o":
				goto Yinjie
			}
			key := strings.ToLower(k)
			for _, val := range v {
				if _, ok := s.Sm_ym[val]; !ok {
					s.Sm_ym[val] = make([]byte, 0)
				}
				s.Sm_ym[val] = append(s.Sm_ym[val], key[0])
			}
			continue
		}
	Yinjie:
		yj := strings.ToLower(k)
		if _, ok := s.Yinjie[yj]; !ok {
			s.Yinjie[yj] = make([]string, 0)
		}
		s.Yinjie[yj] = append(s.Yinjie[yj], v...)
	}
	return s
}

// zhi shi => [ [ai ei], [ui vi] ] => [[ai ui], [ai vi], [ei ui], [ei vi]]
func (s *Mapping) FromPinyin(py []string) [][]string {
	ret := make([][]string, len(py))
	for i, yinjie := range py {
		ret[i] = s.FromYinjie(yinjie)
	}
	ret = ku.Product(ret)
	return ret
}

// 一个音节转换为按键组合 shi => ["ui", "vi"]
func (s *Mapping) FromYinjie(yinjie string) []string {
	// 自定义音节
	if v, ok := s.Yinjie[yinjie]; ok {
		return v
	}
	// 分离声母和韵母
	// shi => sh, i
	// an => a, an
	var sm, ym string
	switch yinjie[0] {
	case 'a', 'o', 'e':
		sm = string(yinjie[0])
		ym = yinjie
	default:
		if len(yinjie) == 1 {
			sm = yinjie
			break
		}
		if yinjie[1] == 'h' {
			sm = yinjie[:2]
			ym = yinjie[2:]
		} else {
			sm = yinjie[:1]
			ym = yinjie[1:]
		}
	}

	smKeys := s.Sm_ym[sm]
	ymKeys := s.Sm_ym[ym]
	keys := ku.Product([][]byte{smKeys, ymKeys})
	// fmt.Println(sm, smKeys, ym, ymKeys, keys)

	ret := make([]string, len(keys))
	for i, v := range keys {
		ret[i] = string(v)
	}
	return ret
}
