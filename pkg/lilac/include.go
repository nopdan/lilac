package lilac

import (
	"bufio"
	"errors"
	"strings"

	util "github.com/flowerime/goutil"
)

// 处理特殊行，引入其他文件
func include(line string) (*bufio.Scanner, error) {
	var ok bool
	if line, ok = strings.CutPrefix(line, ">>("); !ok {
		return nil, errors.New(line + " doesn't have prefix")
	}
	if line, ok = strings.CutSuffix(line, ")"); !ok {
		return nil, errors.New(line + " doesn't have suffix")
	}
	rd, err := util.Read(line)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(rd)
	return scan, err
}
