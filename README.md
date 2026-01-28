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

And lint that the newer [`cmp.Equal`][cmp-equal] should be used.
For more use cases and examples, check [equality-comparison](analyzer/testdata/src/equality_comparison).

### [Got before Want](https://go.dev/wiki/TestComments#got-before-want)

Test outputs should output the actual value that the function returned before printing the value that was expected.
A usual format for printing test outputs is `YourFunc(%v) = %v, want %v`.

### [Identify the Input](https://go.dev/wiki/TestComments#identify-the-input)

In most tests, your test failure messages should include the function inputs if they are short.
If the relevant properties of the inputs are not obvious (for example, because the inputs are large or opaque),
you should name your test cases with a description of what’s being tested and print the description as part of your error message.

Do not use the index of the test in the test table as a substitute for naming your tests or printing the inputs.
Nobody wants to go through your test table and count the entries to figure out which test case is failing.

### [Print Diffs](https://go.dev/wiki/TestComments#print-diffs)

If your function returns a large output, then it can be hard for someone reading the failure message to
find the differences when your test fails.
Instead of printing both the returned value and the wanted value, make a diff.

[cmp-equal]: https://pkg.go.dev/github.com/google/go-cmp/cmp#Equal
