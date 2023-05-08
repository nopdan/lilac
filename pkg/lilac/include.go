package lilac

import (
	"bufio"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nopdan/ku"
)

// 处理特殊行，引入其他文件，flag 表示是否为拼音词库
func include(line, dir string) (*bufio.Scanner, bool, error) {
	var flag, ok bool
	if line, ok = strings.CutPrefix(line, "?>>("); ok {
		flag = true
	} else if line, ok = strings.CutPrefix(line, ">>("); !ok {
		return nil, flag, errors.New(line + " doesn't have prefix")
	}
	if line, ok = strings.CutSuffix(line, ")"); !ok {
		return nil, flag, errors.New(line + " doesn't have suffix")
	}
	// 相对路径
	if !filepath.IsAbs(line) {
		line = filepath.Join(dir, line)
	}
	rd, err := ku.Read(line)
	if err != nil {
		fmt.Printf("导入文件失败: %v\n", line)
		return nil, flag, err
	}
	fmt.Printf("导入文件: %v\n", line)
	scan := bufio.NewScanner(rd)
	return scan, flag, err
}
