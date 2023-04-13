package main

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func main() {
	fp := "test.ini"
	config := LoadFile(fp)
	config.handleData()
	fmt.Printf("%+v", config)
}

// 读取配置文件
func LoadFile(fp string) (p *Config) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"Correct", "Data", "Special"},
	}, fp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg.SectionStrings())

	p = NewConfig()
	config := cfg.Section("Config")
	err = config.MapTo(p)
	if err != nil {
		fmt.Println(err)
	}

	p.Correct = cfg.Section("Correct").Body()
	p.Data = cfg.Section("Data").Body()
	p.Special = cfg.Section("Special").Body()
	return
}
