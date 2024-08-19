# The Slow Programming Language

The Slow programming language is a dynamically-typed interpreted programming language implemented in Go.

## Installation

Currently, the only way to run the Slow interpreter is to build it from source. This can be done by cloning this repository and using the `go` CLI to build it.

```console
$ git clone https://github.com/chrispyles/slow
$ cd slow
$ go install .
```

This will add the `slow` binary to your Go `bin` directory, so make sure it is in your `$PATH`.

## Usage

The `slow` interpreter has two modes: script execution and live interpretation. To run a Slow script (idiomatically a `.slo` file), pass the path to the file to the `slow` command:

```console
$ slow main.slo
```

To launch the Slow interpreter, run `slow -i`. Like the Python CLI, you can run a file and then start an interpreter with that file's environment using the `-i` flag with the script path:

```console
$ slow -i main.slo
```

## Reference

This section contains a comprehensive reference of the entire Slow programming language.

Some (not all) of the examples below use the interpreter format. In this format, statements are prefixed with `->`, the interpreter prompt. If a statement takes up multiple lines, all lines after the first will be prefixed with `..`. This format mimicks the Slow interpreter.

```
# This is the non-interpreter format
var x = 2
print(x)
```

```
-> # This is the interpreter format
-> var x = 2
2
-> print(x)
2
```

### Data Types

Slow has the following primitive data types:

- `bool`: `true` or `false`
- `float`: `1.`, `2.1`, `3.14`, `.02` etc.
- `int`: `1`, `2`, `-1`, etc.
- `uint`: `0u`, `1u`, `2u`, etc.
- `str`: `"The quick brown fox jumps over the lazy dog."`
- `null`

All numeric data types (except `bool`) are 64 bits (backed by Go's 64-bit number types).

All values are truthy except `false`, `0` (in all numeric types), `""`, and `null`.

#### Strings

<!-- TODO: escape sequences -->

### Statements

Statements in Slow are delimited by newlines. Slow ignores indentation. Comments are preceded with the `#` character, after which everything on that line is ignored.

```
# this is a comment
var x       # <- this is a statement
x = 0       # <- this is also a statement
  print(x)  # <- yet another statement
```

<!-- TODO: all statements evaluate to a value -->
<!-- TODO: subexpressions wrapped in parens are evaluated first -->

### Variables

Variables must be declared with the `var` statement.

```
var x
```

They can also be assigned in the same statement.

```
-> var x = 1
1
```

Variables are scoped to the frame they're declared in. Setting a variable in a child frame will update the value of the variable in the frame in which it was declared.

```
-> var x = 1
1
-> func f() {
..   ++x
.. }
..
-> f()
-> x
2
```

Fields or methods of values are accessed using dot notation:

```
# Get the "bar" attribute of the variable "foo"
foo.bar

# Call the "baz" method of the variable "foo"
foo.baz()
```

### Operators

#### Binary Operators

Slow supports the following arithmetic operators:

| Operator | Description    |
|----------|----------------|
| `+`      | addition       |
| `-`      | subtraction    |
| `*`      | multiplication |
| `/`      | division       |
| `%`      | modulus        |
| `//`     | floor division |
| `**`     | exponentiation |

When using arithmetic operators, the precedence of types is `float`, then `int`, then `uint`. That is, if you add (or subtract, multiply, etc.) a `float` and an `int` (or `uint`), the result is a `float`. If you combine an `int` and a `uint`, the result is an `int`. `bool` values can also be used in arithmetic expressions, where `true` becomes `1` and `false` becomes `0` (both treated like `uint`s). There are some exceptions to this:

- Divison (`/`) always returns a `float`.
- Modulus (`%`) always returns an `int`.
- Floor division (`//`) always returns an `int`.

Also note that the expotentiation operator, `**`, is backed by Go's [`math.Pow` function](https://pkg.go.dev/math#Pow), meaning its operands are converted to `float64`s and the result is converted from `float64` to the correct result type.

Each of the arithmetic operators has a reassignment variant that reassigns its left operand to the value of the expression; the syntax for these variants is the arithmetic operator suffixed with an `=` (`+=`, `-=`, `*=`, `/=`, `%=`, `//=`, `**=`).

```
-> var x = 2
2
-> x += 1
3
-> x //= 2
1
-> print(x)
1
```

The left operand of a reassignment oeprators may only be an already-declared variable, an object field, or an [index](#list-indexing).

Slow supports the following logical operators:

| Operator | Description |
|----------|-------------|
| `&&`     | logical and |
| `\|\|`   | logical or  |
| `^^`     | logical xor |

Both the `&&` and `||` operators short-circuit; that is, the second operand of `&&` and `||` are only evaluated if the first is falsey and truthy, respectively. `&&` and `||` also do not change the types of their operands (e.g. `1 && 2` returns `1`, not `true`), but `^^` always returns a `bool`. All logical operators also have a reassignment variant (`&&=`, `||=`, `^^=`).

Slow supports the following comparison operators:

| Operator | Description               |
|----------|---------------------------|
| `==`     | equal                     |
| `!=`     | not equal                 |
| `<`      | less then                 |
| `<=`     | less than or equal to     |
| `>`      | greater than              |
| `>=`     | greather than or equal to |

The `==` and `!=` support all types. (Note that all non-primitive types, like [lists](#lists), are compared by reference and not by value.) The other comparison operators only support numeric types and strings.

#### Unary Operators

<!-- TODO -->

The unary reassignment operators (`++` and `--`) return the value of the variable **before** the operation but set the value of the variable/field to the value after applying the operation.

```
-> var x = 1
1
-> ++x
1
-> x
2
```

Note that all unary operators must precede their operand; that is `++x` is valid, but `x++` is not.

#### Ternary Conditional Operator

Slow supports a ternary conditional operator of the form `<condition> ? <value if true> : <value if false>`. Only one branch of the operator is ever evaluated.

```
# In the example below, if i is even, only f is called, otherwise
# only g is called.
var x = i % 2 == 0 ? f() : g()
```

### Lists

Slow has a built-in `list` type backed by a Go slice. Lists literals are declared using square brackets:

```
var l = []
```

You can also specify elements when writing a list literal:

```
l = [1, 2, 3]
```

#### List Indexing

Lists are zero-indexed. To retrieve the element of a list at a particular index, use square brackets:

```
-> var l = [1, 2]
[1, 2]
-> l[1]
2
```

#### List Methods

The `list` type has several built-in methods, each of which is described below.

<!-- TODO: more methods -->

##### `list.append`

The `append` method of `list` adds an element to the end of that `list` in-place.

```
-> var l = [1, 2]
[1, 2]
-> l.append(3)
-> print(l)
[1, 2, 3]
```

##### `list.equals`

The `equals` method of `list` compares it against another value. If the other value is also a list and each element of the two lists are equal (either by the `==` operator or `list.equals` if the corresponding elements are both themselves `list`s).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> l1.equals(1)
false
-> l1.equals([1, 2])
false
-> l1.equals([1, 2, 3])
true
-> [[1, 2], [2, 3]].equals([[1, 2], [2, 3]])
true
```

### Maps

Slow has a built-in `map` type that is implemented using a hash table constructed from Go's `map` type. Map literals are declared using curly brackets:

```
var m = {}
```

Like lists, you can also specify elements in a map literal:

```
m = {1: 2, 3: 4, true: 1, "foo": "bar"}
```

Only hashable types may be used as `map` keys; the only types that are currently hashable are primitives. Any type may be used as a value in a `map`.

#### Map Methods

The `map` type has several built-in methods, each of which is described below.

##### `map.get`

The `get` method of a `map` returns the value corresponding to the provided key.

```
var m = {1: 2}
m.get(1)
#> 2
```

If the specified key is not in the map, a `KeyError` is thrown. You can specify a default value for a key by passing in a second value:

```
-> var m = {1: 2}
{1: 2}
-> m.get(1)
2
-> m.get(3, 4)
4
```

##### `map.set`

The `set` method of a `map` creates a new key-value pair in the map. It returns `true` if the key was already present in the map (i.e. if it was overwritten) and `false` if it was not.

```
-> var m = {1: 2}
{1: 2}
-> m.set(3, 4)
false
-> m.set(1, 3)
true
```

The provided key must be of a hashable type.

### Conditionals

Slow supports `if` statements to conditionally execute code blocks. To run a block should the condition be falsey, use an `else` statement.

```
if x % 3 == 0 {
  print("x is a multiple of 3")
}
else if x % 3  == 1 {
  print("x mod 3 is 1")
}
else {
  print("x mod 3 is 2")
}
```

Conditions do not need to evaluate to `bool`s; the truthiness of the condition's value will be used to determine whether to run the `if` body. 

```
if x % 2 {
  print("x is odd")
}
else {
  print("x is even")
}
```

Because every statement in Slow evaluates to a value, the return value of the last statement in an `if`/`else if`/`else` block's body is the value of the statement.

```
-> var x = 1
1
-> if x % 2 == 0 {
..   "even"
.. }
.. else {
..   "odd"
.. }
..
"odd"
```

Slow also supports `switch` statements, which match a value to a series of possible `case`s using the logic of the `==` operator. Unlike many other languages, Slow's `switch` cases **do not** fall through by default; the `fallthrough` keyword must be used to trigger fall through. Slow uses curly brackets to wrap `case` bodies.

```
switch x % 3 {
  case 0 {
    print("x is a multiple of 3")
  }
  case 1 { fallthrough }
  case 2 {
    print("x is not a multiple of 3")
  }
}
```

### Control Flow

Slow supports control flow with `for` and `while` loops.

`for` loops iterate over a pre-defined set of values returned by an [iterator](#iterators). The `for` loop has the syntax `for <loop variable> in <iterator>` followed by a body enclosed in curly brackets.

```
-> for i in range(5) {
..   print(i)
.. }
..
0
1
2
3
4
```

`while` loops iterate while a condition evaluates to a truthy value and have the syntax `while <condition>` followed by a body enclosed in curly brackets.

```
var x = 10
while x > 0 {
  print(x)
  --x
}
```

A loop can run indefinitely by setting its condition to a value that is always truthy:

```
while true {
  print("still running...")
}
```

In either loop type, the rest of the current iteration can be skipped using the `continue` keyword.

```
for x in range(20) {
  if x % 2 == 1 { continue }
  print(x)
}
```

Loops can be broken early using the `break` keyword.

```
var x = 0
while true {
  ++x
  if x > 100 { break }
}
```

#### Iterators

The iterator in a `for` loop is a built-in type in Slow. `list`s and `str`s come with iterators, and there is also a [generator type](#generators) that backs built-in functions like [`range`](#range).

```
# To iterate over each character in a string:
var s = "foobar"
for c in s {
  print(c)
}

# To iterate over each element of a list:
for e in [1, 2, 3] {
  print(e)
}
```

### Functions

<!-- TODO -->

```
func isFactor(x, y) {
  return x % y == 0
}
```

#### Built-in Functions

##### `exit`

##### `len`

##### `print`

##### `range`

##### `type`

## Planned Features

### Generators

### List Slicing

### Map Indexing

### Classes

### Format Strings

### Imports

### Error Handling

### Type Casting
