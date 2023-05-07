package main

import (
	"testing"
)

func TestBuild(t *testing.T) {
	// path := "sample/86五笔.ini"
	path := "sample/星辰双拼.ini"
	// path := "./sample/星空键道6.ini"
	// path := "test/哲豆音形.ini"
	conf := newConfig(path)
	conf.Build()
}
