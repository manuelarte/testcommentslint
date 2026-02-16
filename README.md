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
testcommentslint [-equality-comparison=true|false] [-got-before-want=true|false] [-identify-function=true|false]
[-table-driven-format.type=map|slice] [-table-driven-format.inlined=true|false] ./...
```

Parameters:

- `equality-comparison`: `true|false` (default `true`) Checks `reflect.DeepEqual` can be replaced by newer `cmp.Equal`.
- `got-before-want`: `true|false` (default `true`) Check that output the actual value that the function returned before
printing the value that was expected.
- `identify-function`: `true|false` (default `true`) Check that the failure messages in `t.Errorf` contains the function name.
- `table-driven-format.type`: `map|slice` (default ``) Check that the table-driven tests are either Map or Slice, empty to leave it as it is.
- `table-driven-format.inlined`: `true|false` (default `false`) Check that the table-driven tests are inlined in the `for` loop.

## 🚀 Features

### [Equality Comparison and Diffs](https://go.dev/wiki/TestComments#equality-comparison-and-diffs)

This linter detects the expression:

<!-- markdownlint-disable -->
```go
if !reflect.DeepEqual(got, want) {
    t.Errorf("MyFunction got %v, want %v", got, want)
}
```
<!-- markdownlint-enable -->

And lint that the newer [`cmp.Equal`][cmp-equal] or [`cmp.Diff`][cmp-diff] should be used.
For more use cases and examples, check [equality-comparison](analyzer/testdata/src/equality_comparison).

> [!NOTE]
> Suggested Fix can't be supported since it could potentially imply adding go-cmp dependency
> and `reflect.DeepEqual` can't be directly replaced by `cmp.Equal` or `cmp.Diff`.

### [Got before Want](https://go.dev/wiki/TestComments#got-before-want)

Test outputs should output the actual value that the function returned before printing the value that was expected.
So prefer failure messages like `YourFunc(%v) = %v, want %v` over `want: %v, got: %v`.

> [!NOTE]
> Suggested Fix can't be supported since it would imply changing the original failure message.

### [Identify The Function](https://go.dev/wiki/TestComments#identify-the-function)

In most tests, failure messages should include the name of the function that failed, even though it seems obvious
from the name of the test function.

Prefer:

`t.Errorf("YourFunc(%v) = %v, want %v", in, got, want)`

and not:

`t.Errorf("got %v, want %v", got, want)`

> [!NOTE]
> Suggested Fix may be supported.

### Table-Driven Test Format

Feature that checks consistency when declaring your table-driven tests.
The options are:

#### Map non-inlined

<!-- markdownlint-disable -->
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
<!-- markdownlint-enable -->

#### Map inlined

<!-- markdownlint-disable -->
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
<!-- markdownlint-enable -->

#### Slice non-inlined

<!-- markdownlint-disable -->
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
<!-- markdownlint-enable -->

#### Slice inlined

<!-- markdownlint-disable -->
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
<!-- markdownlint-enable -->

[cmp-equal]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal
[cmp-diff]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff
