# Test Comments Lint

[![CI](https://github.com/manuelarte/testcommentslint/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/testcommentslint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/testcommentslint)](https://goreportcard.com/report/github.com/manuelarte/testcommentslint)
![version](https://img.shields.io/github/v/release/manuelarte/testcommentslint)

Go Lint that follows standards described in [TestComments](https://go.dev/wiki/TestComments).

## ⬇️  Getting Started

To install it, run:

```bash
go install github.com/manuelarte/testcommentslint@latest
```

And then use it with

```bash
testcommentslint [-equality-comparison=true|false] [-failure-message=true|false] [-table-driven-format.type=map|slice] [-table-driven-format.inlined=true|false] ./...
```

Parameters:

- `equality-comparison`: `true|false` (default `true`) Checks `reflect.DeepEqual` can be replaced by newer `cmp.Equal`.
- `failure-message`: `true|false` (default `true`) Check that the failure messages in `t.Errorf` follow the format expected.
- `table-driven-format.type`: `map|slice` (default ``) Check that the table-driven tests are either Map or Slice, empty to leave it as it is.
- `table-driven-format.inlined`: `true|false` (default `false`) Check that the table-driven tests are inlined in the `for` loop.

## 🚀 Features

### [Equality Comparison and Diffs](https://go.dev/wiki/TestComments#equality-comparison-and-diffs)

This linter detects expressions like:

```go
if !reflect.DeepEqual(got, want) {
    t.Errorf("MyFunction got %v, want %v", got, want)
}
```

And lint that the newer [`cmp.Equal`][cmp-equal] or [`cmp.Diff`][cmp-diff] should be used.
For more use cases and examples, check [equality-comparison](analyzer/testdata/src/equality_comparison).

> [!NOTE]
> Suggested Fix can't be supported since it could potentially imply adding go-cmp dependency
> and `reflect.DeepEqual` can't be directly replaced by `cmp.Equal` or `cmp.Diff`.

### [Got before Want](https://go.dev/wiki/TestComments#got-before-want)

Test outputs should output the actual value that the function returned before printing the value that was expected.
A usual format for printing test outputs is `YourFunc(%v) = %v, want %v`.

For diffs, directionality is less apparent, and thus it is important to include a key to aid in interpreting the failure.
See [Print Diffs](https://go.dev/wiki/TestComments#print-diffs).

Whichever order you use in your failure messages, you should explicitly indicate the ordering as a part of the failure message,
because existing code is inconsistent about the ordering.

### Table-Driven Test Format

Checks whether the table-driven test follow the specified format.

#### Map non-inlined

```go
tests := map[string]struct {
    in int
    out int
} {
    "test1": {
        in: 1,
        out: 1,
    },
}
for name, test := range tests {
    t.Run(name, func(t *testing.T) {
		got := abs(test.in)
		if got != test.out {
			t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
		}
    })
}
```

#### Map inlined

```go
for name, test := range map[string]struct {
	in int
    out int
} {
    "test1": {
		in: 1, 
		out: 1,
	},
} {
	t.Run(name, func(t *testing.T) {
		got := abs(test.in)
		if got != test.out {
			t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
		}
	})
}
```

#### Slice non-inlined

```go
tests := []struct {
	name string
	in int
    out int
} {
	{
		name: "test1",
        in: 1,
        out: 1,
    },
}
for _, test := range tests {
    t.Run(test.name, func(t *testing.T) {
        got := abs(test.in)
        if got != test.out {
            t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
        }
 })
}
```

#### Slice inlined

```go
for _, test := range []struct {
	name string
    in int
    out int
} {
	{
		name: "test1", 
		in: 1, 
		out: 1,
	},
} {
	t.Run(test.name, func(t *testing.T) {
		got := abs(test.in)
		if got != test.out {
			t.Errorf("abs(%d) = %d, want %d", test.in, got, test.out)
		}
	})
}
```

[cmp-equal]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal
[cmp-diff]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff
