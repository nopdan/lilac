package util

import (
	"fmt"
	"testing"
)

func TestProduct(t *testing.T) {
	sli := [][]byte{{'a', 'b'}, {'c'}, {'d', 'e', 'f'}}
	fn := func(sli [][]byte) {
		new := Product(sli)
		for _, v := range new {
			fmt.Println(string(v))
		}
	}
	fn(sli)

	sli = [][]byte{{'a', 'b'}}
	fn(sli)

	sli = [][]byte{{'a', 'b'}, {}}
	fn(sli)
}
