package lilac

import (
	"bytes"
	"sort"
	"strings"
)

func (c *Config) OutputResult() []byte {
	var buf bytes.Buffer
	for _, entry := range c.Result {
		buf.WriteString(entry[0])
		buf.WriteByte('\t')
		buf.WriteString(entry[1])
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func (c *Config) OutputMisMatch() []byte {
	if len(c.MisMatch) == 0 {
		return []byte{}
	}
	var buf bytes.Buffer
	buf.WriteString("词条\t编码\t生成的编码\t拼音\n")
	buf.WriteString("-----------------------------\n")
	for _, v := range c.MisMatch {
		buf.WriteString(v.Word)
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Codes, " "))
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Gen, " "))
		buf.WriteByte('\t')
		buf.WriteString(strings.Join(v.Pinyin, " "))
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func (c *Config) OutputEmpty() []byte {
	if len(c.Empty) == 0 {
		return []byte{}
	}
	// 排序后输出
	li := make([][2]string, 0, len(c.Empty))
	for pre, entries := range c.Empty {
		li = append(li, [2]string{pre, strings.Join(entries, " ")})
	}
	sort.Slice(li, func(i, j int) bool {
		return li[i][0] < li[j][0]
	})
	var buf bytes.Buffer
	buf.WriteString("空码\t后续\n")
	buf.WriteString("-----------------------------\n")
	for _, v := range li {
		buf.WriteString(v[0])
		buf.WriteByte('\t')
		buf.WriteString(v[1])
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}
