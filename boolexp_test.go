package boolexpr_test

import (
	"fmt"

	boolexpr "github.com/pschou/go-boolexp"
)

func ExampleParse() {
	vars := map[string]bool{
		"a": true,
		"b": false,
	}

	tests := []string{
		"a & b",
		"a&b",
		"a &b",
		"a& b",
		"a and b",
		"a AND b",
		"a | b",
		"a|b",
		"a |b",
		"a| b",
		"a or b",
		"a OR b",
		"a and (a | b)",
		"a & (a | b)",
		"a or (a | b)",
		"b and (a | b)",
		"b & (a | b)",
		"b or (a | b)",
	}

	fmt.Printf("vars = %#v\n", vars)

	for _, test := range tests {
		v, err := boolexpr.Parse(test, vars)
		fmt.Println("testing:", test, "->", v, err)
	}

	_, used, _ := boolexpr.ParseWithUsed("a | a", vars)
	fmt.Printf("testing for used %#v\n", used)
	// Output:
	// vars = map[string]bool{"a":true, "b":false}
	// testing: a & b -> false <nil>
	// testing: a&b -> false <nil>
	// testing: a &b -> false <nil>
	// testing: a& b -> false <nil>
	// testing: a and b -> false <nil>
	// testing: a AND b -> false <nil>
	// testing: a | b -> true <nil>
	// testing: a|b -> true <nil>
	// testing: a |b -> true <nil>
	// testing: a| b -> true <nil>
	// testing: a or b -> true <nil>
	// testing: a OR b -> true <nil>
	// testing: a and (a | b) -> true <nil>
	// testing: a & (a | b) -> true <nil>
	// testing: a or (a | b) -> true <nil>
	// testing: b and (a | b) -> false <nil>
	// testing: b & (a | b) -> false <nil>
	// testing: b or (a | b) -> true <nil>
	// testing for used map[string]bool{"a":true}
}
