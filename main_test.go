package main

import (
	"testing"
)

func TestBuild(t *testing.T) {

	run("sample/86五笔.ini")
	run("sample/星辰双拼.ini")
	run("sample/哲豆音形.ini")
	run("sample/星空键道6.ini")
}
