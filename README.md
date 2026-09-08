# Test Comments Lint

[![CI](https://github.com/manuelarte/testcomments/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/testcomments/actions/workflows/ci.yml)
![version](https://img.shields.io/github/v/release/manuelarte/testcomments)

Go Lint that follows standards described in [TestComments](https://go.dev/wiki/TestComments).

## ⬇️  Getting Started

### Run it as a standalone linter

To install it, run:

```bash
go install github.com/manuelarte/testcomments@latest
```

And then use it with

```bash
testcomments [-equality-comparison.reflect=true|false] [-equality-comparison.equal=true|false]
[-got-before-want=true|false] [-identify-function=true|false]
[-table-driven-format.type=map|slice] [-table-driven-format.inlined=true|false] ./...
```

Parameters:

- `equality-comparison.reflect`: `true|false` (default `true`) Checks `reflect.DeepEqual` can be replaced by newer `cmp.Equal`.
- `equality-comparison.equal`: `true|false` (default `true`) Checks helper test functions to compare two structs that
can be replaced by either `cmp.Equal` or `cmp.Diff`.
- `got-before-want`: `true|false` (default `true`) Check that the failure message outputs the actual value that the
function returned before printing the value that was expected.
- `identify-function`: `true|false` (default `true`) Check that the failure messages in `t.Errorf` contains the function name.
- `table-driven-format.type`: `map|slice` (default ``) Check that the table-driven tests are either Map or Slice, empty to leave it as it is.
- `table-driven-format.inlined`: `true|false` (default `false`) Check that the table-driven tests are inlined in the `for` loop.

### Run it as a module plugin in golangci-lint

You can integrate this linter with [golangci-lint](https://golangci-lint.run/)
by using the [module plugin](https://golangci-lint.run/docs/plugins/module-plugins/).

Example of a `custom-gcl.yml` file that includes this linter:

```yaml
version: v2.13.2
plugins:
  - module: "github.com/manuelarte/testcomments"
    import: "github.com/manuelarte/testcomments/plugin"
    version: latest
```

## 🚀 Features

### Equality Comparison

#### [Reflect](https://go.dev/wiki/TestComments#equality-comparison-and-diffs)

This linter detects the expression:

<!-- markdownlint-disable -->
```go
if !reflect.DeepEqual(got, want) {
	t.Errorf("MyFunction got %v, want %v", got, want)
}
```
<!-- markdownlint-enable -->

And lint that the newer [`cmp.Equal`][cmp-equal] or [`cmp.Diff`][cmp-diff] should be used.
For more use cases and examples, check [equality-comparison](analyzer/testdata/src/reflect-deepequal).

> [!NOTE]
> Suggested Fix can't be supported since it could potentially imply adding go-cmp dependency
> and `reflect.DeepEqual` can't be directly replaced by `cmp.Equal` or `cmp.Diff`.

#### Equal

This linter detects helper functions like:

```go
func areEqual(a, b MyStruct) bool {
 return a.Name && b.Name && a.Surname == b.Surname
}
```

And propose to use `cmp.Equal` or `cmp.Diff`.

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
