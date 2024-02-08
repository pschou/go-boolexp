package boolexpr_test

import (
	"fmt"

	boolexpr "github.com/pschou/go-boolexp"
)

func ExampleParse() {
	vars := map[string][]bool{
		"a": []bool{true},
		"b": []bool{false},
		"c": []bool{true, false},
		"d": []bool{true, true, true},
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
		"(a | b) & b",
		"(a | b) or b",
		"(a and b) & b",
		"(a and b) or b",
		"a or b or a or b",
		"a or (a and b) or b",
		"a and (a and b) or b",
		"a and (a and b) ! or b",
		"a and not b",
		"a and !(b)",
		"a xor b",
		"a^a",
		"any c",
		"all c",
		"none c",
		"not c",
		"none a",
		"none b",
		"all d and any c",
		"b xor (all c and b)",
	}

	fmt.Printf("vars = %#v\n", vars)

	for _, test := range tests {
		v, err := boolexpr.Parse(test, vars)
		fmt.Println("testing:", test, "->", v, err)
	}

	used := make(map[string]bool)
	boolexpr.ParseWithUsed("a | a", vars, used)
	fmt.Printf("testing for used %#v\n", used)
	// Output:
	// vars = map[string][]bool{"a":[]bool{true}, "b":[]bool{false}, "c":[]bool{true, false}, "d":[]bool{true, true, true}}
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
	// testing: (a | b) & b -> false <nil>
	// testing: (a | b) or b -> true <nil>
	// testing: (a and b) & b -> false <nil>
	// testing: (a and b) or b -> false <nil>
	// testing: a or b or a or b -> true <nil>
	// testing: a or (a and b) or b -> true <nil>
	// testing: a and (a and b) or b -> false <nil>
	// testing: a and (a and b) ! or b -> false boolexp: missing logical expression "! or b"
	// testing: a and not b -> true <nil>
	// testing: a and !(b) -> true <nil>
	// testing: a xor b -> true <nil>
	// testing: a^a -> false <nil>
	// testing: any c -> true <nil>
	// testing: all c -> false <nil>
	// testing: none c -> false <nil>
	// testing: not c -> false boolexp: multiple values for "c" with no aggregate operator
	// testing: none a -> false <nil>
	// testing: none b -> true <nil>
	// testing: all d and any c -> true <nil>
	// testing: b xor (all c and b) -> false <nil>
	// testing for used map[string]bool{"a":true}
}
