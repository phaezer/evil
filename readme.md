# EVIL

Evil provides go code generation that creates or sets top level scoped variables with
the results of function calls.

## Motivation

Originally inspired by zig's `comptime` this package aims to fill a gap in the go
programming language: build-time meta programming.

## Usage

Set a go generate declaration on a target file with an optional output filename using the `-o` option.

```go
//go:generate evil -o evilfile.go
```

If no output file is specified the default is `<target filename>_evil.go`

When `go generate` is called, the command scans through the file's function declarations, looking for
functions with signatures equal to:

```go
func (g evil.Generator)
```

For each function with the given signature, a temporary main.go file and package
(if the generator was not called on a main package) is created that:

- inits a new `evil.Generator`
- calls each evil generator function found in the file
- generates a new file that creates or sets the value based on the result of the function calls

### Creating a var or const

To create a new var with an evil generator function, use the `evil.Generator.CreateVar` method:

```go
func createEvilVar(g evil.Generator) {
  // create a new package-scoped variable with the type and value of the output
  //  from the fibonacci function
  g.CreateVar("evilVar", fibonacci(1000))
}
```

The same thing can be done for constants with `evil.Generator.CreateConst`

```go
func createEvilConst(g evil.Generator) {
  // create a new package exported constant with the type and value of the output
  //  from the constantGenerator function
  g.CreateVar("EvilConst", constantGenerator())
}
```

The above will result in the generated file setting a const named `EvilConst` at the package level.

### Setting var values

To create a var initializer for an existing package scoped variable use the `evil.Generator.Init` method.

```go
var fibValues []int

func createEvilConst(g evil.Generator) {
  // create an initializer function that includes setting
  //  fibValues to the result of calling fibonacci(1000)
  g.Init(fibValues, fibonacci(1000))
}
```

The above will result in the generated file having an `init()` function that sets the `fibValues` var to the result of the `fibonacci(1000)` result.

As long as the `createEvilConst` function isn't referenced anywhere else, the compiler will discard the function.

### Advanced usage

#### Build tags

To ensure that files used by evil's genertor are not compiled for release, add a build tag:

```go
// +build evil
```

## Limitations

Generation only works for package / global scoped vars and consts, it cannot set values in other scopes.