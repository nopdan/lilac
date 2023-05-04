package main

import (
	"testing"

	"github.com/flowerime/lilac/pkg/lilac"
)

func TestDo(t *testing.T) {

	do := func(path string) {
		c := lilac.NewConfig(path)
		c.Do()
		// di := c.Do()
		// fmt.Println(di)
	}
	// path := "sample/86五笔.ini"
	// path := "sample/星辰双拼.ini"
	path := "sample/哲豆音形.ini"
	do(path)
}
