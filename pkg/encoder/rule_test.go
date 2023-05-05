package encoder

import (
	"fmt"
	"testing"
)

func TestRule(t *testing.T) {

	test := func(rules string) {
		e := NewEncoder(rules)
		fmt.Println(e.Rule)
	}

	rules := `
		2:A+a1+B+b1+a2+b2,
		3:A+B+C+a2+b2+c2,
		:A+B+C+Z+a2+b2`
	test(rules)

	rules = ":a1+a2.._"
	test(rules)
	rules = ":a1+a2..'"
	test(rules)

	// rules = `
	// 	2:A+a1+B+b1+a2+b2
	// 	3:A+B+C+Z,
	// 	:A+B+C+Z+a2+b2`
	// test(rules)
}
