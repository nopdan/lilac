package main

import (
	"testing"

	"github.com/flowerime/lilac/pkg/lilac"
)

func TestGenerate(t *testing.T) {
	// path := "sample/86五笔.ini"
	// path := "sample/星辰双拼.ini"
	path := "test/哲豆音形.ini"
	lilac.Run(path)
}
