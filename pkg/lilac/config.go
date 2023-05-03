package lilac

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/flowerime/lilac/pkg/encoder"
	m "github.com/flowerime/lilac/pkg/mapping"
	"gopkg.in/ini.v1"
)

type Config struct {
	// Zcgz     string `ini:"组词规则"`
	// Jmgz     string `ini:"简码规则"`
	// Scjm     int    `ini:"生成简码"`
	// Bldzqm   int    `ini:"保留单字全码"`
	// Blczqm   int    `ini:"保留词组全码"`
	// Sort     int    `ini:"重新按字母排序"`
	// External string `ini:"使用单独的码表造词"`

	Type int    `ini:"Type"`
	Rule string `ini:"Rule"`

	Correct map[string][]string
	Char    map[string][]string
	Mapping *m.Mapping

	Dict    string
	encoder *encoder.Encoder
}

func NewConfig(path string) *Config {
	// 手动解析下列 Section
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"Dict", "Correct", "Char", "Mapping"},
	}, path)
	if err != nil {
		panic(err)
	}

	c := new(Config)
	config := cfg.Section("Config")
	err = config.MapTo(c)
	if err != nil {
		fmt.Println(err)
	}

	var text string
	text = cfg.Section("Correct").Body()
	c.Correct = handleText(text)

	c.encoder = encoder.NewEncoder(c.Rule)
	text = cfg.Section("Char").Body()
	c.encoder.Char = handleText(text)
	text = cfg.Section("Mapping").Body()
	data := handleText(text)
	if len(data) != 0 {
		c.encoder.Mapping = m.NewMapping(data)
	}

	c.Dict = cfg.Section("Dict").Body()
	return c
}

// 生成码表
func (c *Config) Do() [][]string {
	rd := strings.NewReader(c.Dict)
	scan := bufio.NewScanner(rd)
	ret := make([][]string, 0)
	for scan.Scan() {
		line := scan.Text()
		tmp := strings.Split(line, "\t")

		word, ok := strings.CutPrefix(tmp[0], "?")
		entry := []string{word}
		if len(tmp) == 1 {
			gen := c.encoder.Encode(word, []string{})
			fmt.Printf("? 词组: %v, 生成: %v\n", word, gen)
			entry = append(entry, gen...)
			ret = append(ret, entry)
			continue
		}

		if len(tmp) == 2 {
			if !ok {
				ret = append(ret, tmp)
				continue
			}
			// ? 号开头为自动造词
			py := strings.Split(tmp[1], " ")
			gen := c.encoder.Encode(word, py)
			fmt.Printf("? 词组: %v, 拼音: %v, 生成: %v\n", word, py, gen)
			entry = append(entry, gen...)
			ret = append(ret, entry)
		} else {
			fmt.Println("错误:", line)
		}
	}
	return ret
}

// func (c *Config) handleData() {
// 	c.data = make([][]string, 0, len(c.Dict)/8)
// 	c.goucima = make(map[rune][]string)

// 	rd := strings.NewReader(c.Dict)
// 	scan := bufio.NewScanner(rd)

// 	// 单独的单字码表
// 	if c.External != "" {
// 		extRd, err := util.Read(c.External)
// 		if err != nil {
// 			panic(err)
// 		}
// 		extScan := bufio.NewScanner(extRd)
// 		for extScan.Scan() {
// 			tmp := strings.Split(extScan.Text(), "\t")
// 			if len(tmp) != 2 {
// 				continue
// 			}
// 			chars := []rune(tmp[0])
// 			if len(chars) != 1 {
// 				continue
// 			}
// 			c.appendGoucima(chars[0], tmp[1])
// 		}
// 		for scan.Scan() {
// 			tmp := strings.Split(scan.Text(), "\t")
// 			if len(tmp) == 0 {
// 				continue
// 			}
// 			c.data = append(c.data, tmp)
// 		}
// 	}

// 	for scan.Scan() {
// 		tmp := strings.Split(scan.Text(), "\t")
// 		if len(tmp) == 0 {
// 			continue
// 		}
// 		var char rune
// 		if len(tmp[0]) <= 5 {
// 			chars := []rune(tmp[0])
// 			if len(chars) == 1 {
// 				char = chars[0]
// 			}
// 		}
// 		if char != 0 {
// 			if len(tmp) == 2 {
// 				c.appendGoucima(char, tmp[1])
// 			}
// 			if len(tmp) == 3 {
// 				c.appendGoucima(char, tmp[2])
// 			}
// 		}
// 		c.data = append(c.data, tmp)
// 	}
// }

// func (c *Config) appendGoucima(char rune, code string) {
// 	if _, ok := c.goucima[char]; ok {
// 		c.goucima[char] = append(c.goucima[char], code)
// 	} else {
// 		c.goucima[char] = []string{code}
// 	}
// }

// func (c *Config) getDzmb() map[rune][]string {

// 	dzmb := make(map[rune][]string)
// 	b := strings.NewReader(c.Dict)
// 	scan := bufio.NewScanner(b)
// 	for scan.Scan() {
// 		entry := strings.Split(scan.Text(), "\t")
// 		if len(entry) != 2 {
// 			continue
// 		}
// 		tmp := []rune(entry[0])
// 		if len(tmp) > 1 {
// 			continue
// 		}
// 		dz := tmp[0]
// 		if _, ok := dzmb[dz]; !ok {
// 			dzmb[dz] = []string{entry[1]}
// 			continue
// 		}
// 		dzmb[dz] = append(dzmb[dz], entry[1])
// 	}
// 	return dzmb
// }
