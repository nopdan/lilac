package lilac

import (
	"bufio"
	"errors"
	"strings"

	util "github.com/flowerime/goutil"
)

// 处理特殊行，引入其他文件，flag 表示是否为辅助生成的码表
func include(line string) (*bufio.Scanner, bool, error) {
	var flag, ok bool
	if line, ok = strings.CutPrefix(line, "?>>("); ok {
		flag = true
	} else if line, ok = strings.CutPrefix(line, ">>("); !ok {
		return nil, flag, errors.New(line + " doesn't have prefix")
	}
	if line, ok = strings.CutSuffix(line, ")"); !ok {
		return nil, flag, errors.New(line + " doesn't have suffix")
	}
	rd, err := util.Read(line)
	if err != nil {
		return nil, flag, err
	}
	scan := bufio.NewScanner(rd)
	return scan, flag, err
}
