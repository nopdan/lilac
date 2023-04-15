package double

import (
	"fmt"
	"testing"

	"github.com/flowerime/rose/pkg/rose"
)

func TestMapping(t *testing.T) {
	m := newMapping("../../assets/双拼映射表/星辰双拼.ini", AABC)
	fmt.Println(m)
}

func TestToDoublePinyin(t *testing.T) {
	dict := rose.Parse("../pinyin/test/sogou_bak.bin", "sogou_bin")
	table := ToDoublePinyin(dict.ToPinyinTable(), "test/双拼映射表.ini", AABC)
	fmt.Println(table)
}
