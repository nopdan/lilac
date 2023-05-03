package encoder

import (
	"fmt"
	"testing"
)

func TestRule(t *testing.T) {
	rules := `
		2:A+a1+B+b1+a2+b2,
		3:A+B+C+a2+b2+c2,
		:A+B+C+Z+a2+b2
	`
	e := NewEncoder(rules)
	fmt.Println(e.Rule, e.sRule)
}
