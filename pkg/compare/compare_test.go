// Package compare
package compare

import (
	"fmt"
	"testing"
)

type A struct {
	a string
	A string `json:"a" mapstructure:"a" yaml:"a"`
	B B      `json:"b" mapstructure:"b" yaml:"b"`
}
type B struct {
	b string
	B string `json:"b" mapstructure:"b" yaml:"b"`
	c C
	C C
}
type C struct {
	c string
	C string
}

func TestCompareStruct(t *testing.T) {
	a1 := A{
		a: "a",
		A: "A",
		B: B{
			b: "b1",
			B: "B",
			c: C{
				c: "c1",
				C: "C",
			},
			C: C{
				c: "c2",
				C: "C",
			},
		},
	}
	a2 := A{
		a: "a",
		A: "A",
		B: B{
			b: "b1",
			B: "B",
			c: C{
				c: "c2",
				C: "C",
			},
			C: C{
				c: "c2",
				C: "C",
			},
		},
	}
	// s := cmp.Diff(a1, a2, cmpopts.IgnoreUnexported())
	// fmt.Println(s)
	diff, diffCount := CompareStruct(a1, a2)
	fmt.Println(diff, diffCount)
}
