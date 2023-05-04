package lilac

import (
	"bytes"
	"os"
)

func Output(dict [][]string) []byte {
	var buf bytes.Buffer
	for _, entry := range dict {
		buf.WriteString(entry[0])
		buf.WriteByte('\t')
		buf.WriteString(entry[1])
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

func WriteFile(dict [][]string, path string) {
	data := Output(dict)
	err := os.WriteFile(path, data, 0666)
	if err != nil {
		panic(err)
	}
}
