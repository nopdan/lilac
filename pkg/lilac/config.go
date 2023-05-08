package lilac

import (
	"bufio"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flowerime/lilac/pkg/encoder"
	m "github.com/flowerime/lilac/pkg/mapping"
	"github.com/flowerime/pinyin"
	"gopkg.in/ini.v1"
)

type Config struct {
	Rule string `ini:"Rule"`
	Sort bool   `ini:"Sort"`

	Keep     bool   `ini:"保留单字全码"`
	CharRule string `ini:"单字简码规则"`
	WordRule string `ini:"词组简码规则"`

	encoder   *encoder.Encoder
	shortener *shortener
	checker   *checker

	dict  string
	check string

	Dir      string      // 配置文件所在目录
	Result   [][2]string // 最终生成的码表
	MisMatch []*misMatch // 不匹配的编码
	// key: 空码, value: 后续词和编码
	Empty map[string][]string // 空码
}

// 从文件初始化一个 Config，传入 Config 路径和 拼音数据路径
func NewConfig(path string, pydata string) *Config {
	// 跳过解析下列 Section
	cfg, err := ini.LoadSources(ini.LoadOptions{
		UnparseableSections: []string{"Dict", "Char", "Mapping", "Check"},
	}, path)
	if err != nil {
		panic(err)
	}
	c := new(Config)
	config := cfg.Section("Config")
	err = config.MapTo(c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("配置读取成功: %v\n", path)

	c.Dir = filepath.Dir(path)
	c.dict = cfg.Section("Dict").Body()
	c.check = cfg.Section("Check").Body()
	c.initShortener()
	c.initChecker()

	// 初始化编码器
	enc := encoder.NewEncoder(c.Rule)
	var text string
	text = cfg.Section("Char").Body()
	enc.Char = readMap(text, c.Dir)

	text = cfg.Section("Mapping").Body()
	data := readMap(text, c.Dir)
	// 若存在映射表，则加载拼音数据
	py := pinyin.New()
	if len(data) != 0 {
		enc.Mapping = m.NewMapping(data)
		addPinyinData(py, pydata)
	}
	enc.Pinyin = py
	c.encoder = enc

	c.Result = make([][2]string, 0)
	c.MisMatch = make([]*misMatch, 0)
	c.Empty = make(map[string][]string, 0)
	// fmt.Printf("c: %+v\n", c)
	return c
}

// 生成码表
func (c *Config) Build() {
	rd := strings.NewReader(c.dict)
	scan := bufio.NewScanner(rd)
	c.readDict(scan, false)
	// 生成简码
	c.Shorten()
	// 按照编码排序
	if c.Sort {
		sort.SliceStable(c.Result, func(i, j int) bool {
			return c.Result[i][1] < c.Result[j][1]
		})
	}
	// 词库校验
	c.Check()
}

// 这个 flag 表示 ?>>() 格式，
// 里面的词全都为拼音辅助生成编码，
// 即这个词库是一个拼音词库
func (c *Config) readDict(scan *bufio.Scanner, flag bool) {
	for scan.Scan() {
		line := scan.Text()
		// 跳过空行
		if line == "" {
			continue
		}
		// 若符合导入文件的格式则递归执行
		if sc, newFlag, err := include(line, c.Dir); err == nil {
			c.readDict(sc, newFlag)
			continue
		}
		// 读取一行
		tmp := strings.Split(line, "\t")
		// ? 号开头，拼音辅助编码
		word, ok := strings.CutPrefix(tmp[0], "?")
		if flag {
			ok = true
		}
		// 这是自带编码的词条，直接加入
		if !ok && len(tmp) == 2 {
			c.Result = append(c.Result, [2]string{tmp[0], tmp[1]})
			continue
		}
		// 下面处理 ? 开头的词条
		pinyin := []string{}
		if len(tmp) == 2 {
			pinyin = strings.Split(tmp[1], " ")
		}
		// 自动生成的编码，可能有多个
		gen := c.encoder.Encode(word, pinyin)
		for _, code := range gen {
			c.Result = append(c.Result, [2]string{word, code})
		}
	}
}

// 导入拼音数据
func addPinyinData(py *pinyin.Pinyin, dir string) {
	paths := []string{
		"pinyin.txt",
		"duoyin.txt",
		"correct.txt",
	}
	for _, item := range paths {
		path := filepath.Join(dir, item)
		py.AddFile(path)
	}
}

// 追加目录下的所有拼音文件数据
func appendDir(py *pinyin.Pinyin, dir string) {
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
		if !info.IsDir() {
			py.AddFile(path)
		}
		return nil
	})
}
