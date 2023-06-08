# go-boolexp

A basic boolean expression evaluator for a string and a map of boolean variables.

```golang
  vars := map[string]bool{
    "a": true,
    "b": false,
  }

  v, err := boolexpr.Parse("a or b", vars)  // v = true
```
