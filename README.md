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
testcommentslint [-equality-comparison=true|false] ./...
```

Parameters:

- `equality-comparison`: `true|false` (default `true`) Checks `reflect.DeepEqual` can be replaced by newer `cmp.Equal`.
- `failure-message`: `true|false` (default `true`) Check that the failure messages in `t.Errorf` follow the format expected.
- `table-driven-format.type`: `map|slice` (default ``) Check that the table-driven tests are either Map or Slice, empty to leave it as it is.
- `table-driven-format.inline`: `true|false` (default `false`) Check that the table-driven tests are inlined in the `for` loop.

## 🚀 Features

### [Compare Full Structures](https://go.dev/wiki/TestComments#compare-full-structures)

TODO(manuelarte): Think about this

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

TODO: do that if a failure message is something like:
check that one of the param is the output of the function above the if.

- expected/want: %s, actual/got: %s,
then it should be
- `YourFunc(%v) = %v, want %v`.

Possible extra:
check if there are many if with t.Errorf that used the same variable, and the recommend using cmp.Diff.

Test outputs should output the actual value that the function returned before printing the value that was expected.
A usual format for printing test outputs is `YourFunc(%v) = %v, want %v`.

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

### [Identify the Input](https://go.dev/wiki/TestComments#identify-the-input)

In most tests, your test failure messages should include the function inputs if they are short.
If the relevant properties of the inputs are not obvious (for example, because the inputs are large or opaque),
you should name your test cases with a description of what’s being tested and print the description as part of your
error message.

Do not use the index of the test in the test table as a substitute for naming your tests or printing the inputs.
Nobody wants to go through your test table and count the entries to figure out which test case is failing.

### [Print Diffs](https://go.dev/wiki/TestComments#print-diffs)

If your function returns a large output, then it can be hard for someone reading the failure message to
find the differences when your test fails.
Instead of printing both the returned value and the wanted value, make a diff.

[cmp-equal]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal
[cmp-diff]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Diff
