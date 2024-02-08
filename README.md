# go-boolexpr

A basic boolean expression evaluator for a string and a map of boolean variables.

```golang
  vars := map[string][]bool{
    "a": []bool{true},
    "b": []bool{false},
    "c": []bool{true, false},
  }

  v, err := boolexpr.Parse("a or b", vars)  // v = true

  v, err = boolexpr.Parse("a and any c", vars)  // v = true

  v, err = boolexpr.Parse("b xor (all c and b)", vars)  // v = false
```
