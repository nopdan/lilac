package main

import (
	"bufio"
	"strings"

	util "github.com/flowerime/goutil"
)

type Config struct {
	Zcgz     string `ini:"组词规则"`
	Jmgz     string `ini:"简码规则"`
	Scjm     int    `ini:"生成简码"`
	Bldzqm   int    `ini:"保留单字全码"`
	Blczqm   int    `ini:"保留词组全码"`
	Sort     int    `ini:"重新按字母排序"`
	External string `ini:"使用单独的码表造词"`

	Correct string
	Data    string
	Special string

	data    [][]string
	goucima map[rune][]string
}

func NewConfig() *Config {
	p := new(Config)
	p.goucima = make(map[rune][]string)
	return p
}

func (c *Config) handleData() {
	c.data = make([][]string, 0, len(c.Data)/8)
	c.goucima = make(map[rune][]string)

	rd := strings.NewReader(c.Data)
	scan := bufio.NewScanner(rd)

	// 单独的单字码表
	if c.External != "" {
		extRd, err := util.Read(c.External)
		if err != nil {
			panic(err)
		}
		extScan := bufio.NewScanner(extRd)
		for extScan.Scan() {
			tmp := strings.Split(extScan.Text(), "\t")
			if len(tmp) != 2 {
				continue
			}
			chars := []rune(tmp[0])
			if len(chars) != 1 {
				continue
			}
			c.appendGoucima(chars[0], tmp[1])
		}
		for scan.Scan() {
			tmp := strings.Split(scan.Text(), "\t")
			if len(tmp) == 0 {
				continue
			}
			c.data = append(c.data, tmp)
		}
	}

	for scan.Scan() {
		tmp := strings.Split(scan.Text(), "\t")
		if len(tmp) == 0 {
			continue
		}
		var char rune
		if len(tmp[0]) <= 5 {
			chars := []rune(tmp[0])
			if len(chars) == 1 {
				char = chars[0]
			}
		}
		if char != 0 {
			if len(tmp) == 2 {
				c.appendGoucima(char, tmp[1])
			}
			if len(tmp) == 3 {
				c.appendGoucima(char, tmp[2])
			}
		}
		c.data = append(c.data, tmp)
	}
}

func (c *Config) appendGoucima(char rune, code string) {
	if _, ok := c.goucima[char]; ok {
		c.goucima[char] = append(c.goucima[char], code)
	} else {
		c.goucima[char] = []string{code}
	}
}

func (c *Config) getDzmb() map[rune][]string {

	dzmb := make(map[rune][]string)
	b := strings.NewReader(c.Data)
	scan := bufio.NewScanner(b)
	for scan.Scan() {
		entry := strings.Split(scan.Text(), "\t")
		if len(entry) != 2 {
			continue
		}
		tmp := []rune(entry[0])
		if len(tmp) > 1 {
			continue
		}
		dz := tmp[0]
		if _, ok := dzmb[dz]; !ok {
			dzmb[dz] = []string{entry[1]}
			continue
		}
		dzmb[dz] = append(dzmb[dz], entry[1])
	}
	return dzmb
}
