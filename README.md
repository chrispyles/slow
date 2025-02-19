# The Slow Programming Language

[![Run tests](https://github.com/chrispyles/slow/actions/workflows/run-tests.yml/badge.svg)](https://github.com/chrispyles/slow/actions/workflows/run-tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/chrispyles/slow/badge.svg?branch=main)](https://coveralls.io/github/chrispyles/slow?branch=main)

<p align="center">
  <img src="logos/logo.png" width="256" alt="Slow logo">
</p>

The Slow programming language is a dynamically-typed interpreted programming language implemented in Go.

## Installation

Currently, the only way to run the Slow interpreter is to build it from source. This can be done using the `go` CLI:

```console
$ go install github.com/chrispyles/slow@latest
```

This will add the `slow` binary to your Go `bin` directory, so make sure it is in your `$PATH`. Slow uses Go generics, so you must have Go 1.18 or later. 

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
- `str`: an immutable character sequence, `"The quick brown fox jumps over the lazy dog."`
- `bytes`: an immutable byte sequence, `0xDEADBEEF`
- `null`

All numeric data types (except `bool`) are 64 bits (backed by Go's 64-bit number types).

When converting from `bytes` to `float`, `int`, or `uint`, the bytes are decoded using the big-endian binary format. For example, 0xBEEF is converted to a `uint` as `48879u`, not `61374u`.

All values are truthy except `false`, `0` (in all numeric types), `""`, a `bytes` object with all null bytes (e.g. `0x00`), and `null`.

#### Strings

Strings are delimited by double quote `"` characters. They must be on a single line and allow the following escape sequences:

| Sequence | Description               |
|----------|---------------------------|
| `\"`     | double quote character    |
| `\n`     | newline character         |
| `\r`     | carriage return character |
| `\t`     | tab character             |
| `\\`     | backslash character       |

```
-> print("___garoo\rkan\njump")
kangaroo
jump
```

#### Bytes

Bytes are written as case-insensitive hexadecimal values prefixed with `0x` (for example, `0xDEADBEEF`). There must be an even number of characters in a `bytes` literal.

```
-> var b = 0xdeadbeef
0xDEADBEEF
```

### Statements

Statements in Slow are delimited by newlines. Slow ignores indentation. Comments are preceded with the `#` character, after which everything on that line is ignored.

```
# this is a comment
var x       # <- this is a statement
x = 0       # <- this is also a statement
  print(x)  # <- yet another statement
```

Every statement in Slow evaluates to a value. For bocks, this is the return value of the final statement in the block. For example, an `if`/`else` block evaluates to the value of the last statement in whichever branch was executed. A `for` loop evaluates to the value of the last statement in the last iteration.

```
-> var x = 0
0
-> if x {
..   1
.. }
.. else {
..   0
.. }
..
0
-> for x in range(10) {
..   x
.. }
..
9
```

WARNING: Despite the fact that every statement evaluates to a value, not all statements can be bound to a variable and Slow does not support implicit returns.

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

To define a constant variable, use a `const` statement. Constants must be initialized with a value when they are declared and cannot be reassigned (although it is possible to declare a new variable with the same name in a child frame).

```
const x = 1
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

Fields can also be reassigned using dot notation:

```
foo.bar = 1
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

The addition operator `+` can also be used to concatenate strings:

```
-> "foo" + "bar"
"foobar"
```

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

The left operand of a reassignment oeprators may only be an already-declared variable, an object field, or an [index](#indexing).

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

The table below shows the precedence of each binary operator (a lower precedence means the operation is executed sooner). Arithmetic operators follow the [standard order of operations](https://en.wikipedia.org/wiki/Order_of_operations).

| Operator | Precedence |
|----------|------------|
| `**`     | 0          |
| `*`      | 1          |
| `/`      | 1          |
| `//`     | 1          |
| `%`      | 1          |
| `+`      | 2          |
| `-`      | 2          |
| `==`     | 3          |
| `!=`     | 4          |
| `<`      | 4          |
| `<=`     | 4          |
| `>`      | 4          |
| `>=`     | 4          |

All reassignment operators have a higher precedenece than any other operator, and only one reassignment operator may be present in a single statement.

Subexpressions wrapped in parentheses are executed before the rest of the expression, and can be used to override operator precedence:

```
-> 2 + 3 * 4
14
-> (2 + 3) * 4
20
```

#### Unary Operators

Slow supports the following unary operators:

| Operator | Description      |
|----------|------------------|
| `+`      | no-op            |
| `-`      | numeric negation |
| `!`      | logical negation |
| `++`     | increment        |
| `--`     | decrement        |

The unary `+` operator is like multiplying a value by `1u`. The `-` operator flips the sign of its argument, and can only be used with numeric types. The `!` operator returns the logical negation of its operand, can be used with any type of value, and always returns a `bool`.

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

Lists can be either mutable or immutable; all lists are mutable by default, by an immutable copy of any list can be created with the `to_immutable` method described below. (Similarly, a mutable copy of any list can be created with the `to_mutable` method.) Immutable lists do not allow any modification (e.g. index assignment, `list.append`). However, making an immutable list does not make its elements themselves immutable.

#### List Methods

The `list` type has several built-in methods, each of which is described below.

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

##### `list.to_mutable`

The `to_mutable` method of `list` creates a mutable copy of the list. This method can be used on any list (immutable or mutable).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> var l2 = l1.to_mutable()
[1, 2, 3]
-> l2.append(4)
-> l1
[1, 2, 3]
-> l2
[1, 2, 3, 4]
```

##### `list.to_immutable`

The `to_immutable` method of `list` creates an immutable copy of the list. This method can be used on any list (mutable or immutable).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> var l2 = l1.to_immutable()
[1, 2, 3]
-> l1.append(4)
-> l1
[1, 2, 3, 3]
-> l2.append(4)
ValueError: list is immutable
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

Maps can be either mutable or immutable; all maps are mutable by default, by an immutable copy of any map can be created with the `to_immutable` method described below. (Similarly, a mutable copy of any map can be created with the `to_mutable` method.) Immutable maps do not allow any modification (e.g. index assignment, `map.set`). However, making an immutable map does not make its elements themselves immutable.

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

##### `map.to_mutable`

The `to_mutable` method of `map` creates a mutable copy of the map. This method can be used on any map (immutable or mutable).

```
-> var m1 = {1: 2}
{1: 2}
-> var m2 = m1.to_mutable()
{1: 2}
-> m2.set(3, 4)
-> m1
{1: 2}
-> m2
{1: 2, 3: 4}
```

##### `map.to_immutable`

The `to_immutable` method of `map` creates an immutable copy of the map. This method can be used on any map (mutable or immutable).

```
-> var m1 = {1, 2}
{1: 2}
-> var m2 = m1.to_immutable()
{1: 2}
-> m1.set(3, 4)
false
-> m1
{1: 2, 3: 4}
-> m2.set(3, 4)
ValueError: map is immutable
```

### Indexing

Slow supports indexing with square bracket syntax. All indexable types in Slow are zero-indexed.

For all indexable types except `map`s, only a `bool`, `uint`, or `int` may be used as an index. These types also support negative indexing to retrieve elements beginning at the end of the sequence. The index of the last element is `-1`, the second to last is `-2`, etc. `map`s can be indexed with any hashable value.

Note that under the hood, numeric indexes (excluding `map`s) are coerced to Go's `int` type (not `int64`, which is how Slow `int`s are stored). This means it is possible to overflow the range of possible index values if you use a `uint` that's too large. 

Any time an indexable value is indexed with an out-of-bounds or nonexistent index, an `IndexError` is thrown.

#### String Indexing

String characters can be accessed used zero-indexed integers. Indexing returns a single-character `str`.

```
-> var s = "abcdef"
"abcdef"
-> s[0]
"a"
-> s[-1]
"f"
```

Because strings are immutable, index assignments are not supported.

#### Bytes Indexing

Individual bytes can be accessed used zero-indexed integers. Indexing returns a single-byte `bytes`.

```
-> var b = 0xDEADBEEF
0xDEADBEEF
-> b[0]
0xDE
-> b[-1]
0xEF
```

Because `bytes` are immutable, index assignments are not supported.

#### List Indexing

Lists are zero-indexed. To retrieve the element of a `list` at a particular index, use square brackets:

```
-> var l = [1, 2]
[1, 2]
-> l[1]
2
```

List elements can also be updated using indexing:

```
l[0] = 3
l[1] += 1
++l[2]
```


```
-> var l = [1, 2, 3, 4, 5]
[1, 2, 3, 4, 5]
-> l[-1]
5
-> l[-2]
4
```

#### Map Indexing

Map values can be retrieved or updated using indexing. To retrieve the element of a `map` key, use square brackets:

```
-> var m = {1: 2, 2: 4}
{1: 2, 2: 4}
-> m[1]
2
-> m[1] = 3
{1: 3, 2: 4}
-> ++m[2]
4
-> m
{1: 3, 2: 5}
-> m[3] = 6
6
-> m
{1: 3, 2: 5, 3: 6}
```

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

Slow also supports `switch` statements, which match a value to a series of possible `case`s using the logic of the `==` operator. Unlike many other languages, Slow's `switch` cases **do not** fall through by default; the `fallthrough` keyword must be used to trigger fall through. Slow uses curly brackets to wrap `case` bodies. An optional `default` case may be added after all other cases which is executed if no other case is matched.

```
switch x % 3 {
  case 0 {
    print("x is a multiple of 3")
  }
  case 1 { fallthrough }
  case 2 {
    print("x is not a multiple of 3")
  }
  default {
    print("x is not an int")
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

Functions in Slow are declared using the `func` keyword. A function can have zero or more arguments and returns `null` by default unless a `return` statement is included.

The syntax for declaring a function is `func <name>(<arg1>, <arg2>, <etc>)` followed by its body enclosed in curly brackets.

```
func isFactor(x, y) {
  return x % y == 0
}
```

Like everything in Slow, functions are also values and can be treated as such.

```
-> func mysteryFunc(x) { return x % 2 == 0 }
<func mysteryFunc>
-> var isEven = mysteryFunc
<func mysteryFunc>
-> isEven(2)
true
```

#### Defer Statements

Inside a function, a function call can be deferred so that it runs just before the function exits, instead of wherever in the body the `defer` statement is (like Go's `defer` statement). Statements are accrued but not evaluated as the function's body executes and before the function exits, they are run. **Deferred functions are currently not run if the function throws an error.** This behavior will likely be added in the future.

```
-> func foo() {
..   # without the deferral, this statement would error because the variable does not yet exist
..   defer print(i)
..   var i = 0
..   for _ in range(20) {
..     ++i
..   }
.. }
-> foo()
20
```

#### Built-in Functions

Slow has a few functions built into the language. They are declared in a frozen frame that is the parent of the frame that the global environment is declared in.

##### `exit`

The `exit` function exits the Slow interpreter. It takes 1 optional argument, an integer indicating the exit code, which defaults to 0.

```
-> exit()
Exiting with code 0
```

```
-> exit(1)
Exiting with code 1
```

If the argument passed to `exit` is not an integer, it is cast to an integer. If this cast succeeds, the integer value is used as the exit code. If the cast fails (i.e. if the argument is non-numeric), the argument is cast to a `bool` first, which is then cast to an `int` (since all values can be cast to `bool`).

```
-> exit(2.1)
Exiting with code 2
```

```
-> exit([])  # lists are truthy, so bool([]) == true
Exiting with code 1
```

##### `import`

The `import` function imports a module. It takes 1 argument, either a path to another Slow file (with extesnion `.slo`) or the name a of a built-in module (e.g. `fs`) and returns a [`module`](#modules). If the argument is a path, the file is read and executed, and the resulting global environment is converted to a `module`.

```
const fs = import("fs")
const utils = import("utils.slo")
```

##### `len`

The `len` function returns the length of the provided value if it is supported. Currently, the only types in Slow that have lengths are strings, `bytes`, `list`s, and `map`s. `len` always returns a `uint`.

```
-> var l = [1, 2, 3]
[1, 2, 3]
-> len(l)
3u
-> len({1: 2, 3: 4})
2u
```

##### `print`

`print` prints its arguments to stdout followed by a newline character. It can accept any number of arguments, converts them to their string representation, and concatenates those strings.

```
-> print(1u, 2, 3., "foo")
1u23.0foo
```

##### `range`

The `range` function creates a generator for a range of numbers. It can take either 1, 2, or 3 arguments.

- If it receives 1 argument `n`, the generator it returns yields the sequence `0, 1, 2, ..., n-1`.
- If it receives two arguments `m` and `n`, the generator yields the sequence `m, m+1, m+2, ..., n-1`. If `m > n`, the generator is infinite.
- If it receives three arguments `m`, `n`, and `k`, the generator yields the sequence `m, m+k, m+2*k, ..., m+jk` where `m+jk` is the largest possible value less than `n`.

```
range(10)        # yields 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
range(2, 10)     # yields 2, 3, 4, 5, 6, 7, 8, 9
range(2, 10, 3)  # yields 2, 5, 8
```

##### `type`

The `type` function takes a single argument and returns a string representing the type of its argument.

```
-> type(true)
"bool"
-> type(1.)
"float"
-> type(range)
"func"
-> type(range(1))
"generator"
-> type(1)
"int"
-> type([])
"list"
-> type({})
"map"
-> type(null)
"null"
-> type("")
"str"
-> type(1u)
"uint"
```

### Modules

Slow supports importing objects and functions from other Slow files or from built-in modules via the `module` type. `module`s are frozen; that is, new variables can't be declared in them nor can variables be rebound. A `module` can be created using the [`import` function](#import) and variables in its environment can be accessed use dot syntax.

For example, to import a `module` from another Slow file:

```
# logging.slo
func info(s) {
  print(s)
}
```

```
const logging = import("logging.slo")

logging.info("successfully imported logging.slo")
```

Or, to import a built-in module:

```
const fs = import("fs")

print(fs.read("foo.txt"))
```

#### Built-in Modules

This section describes Slow's built-in modules. Any built-in module can be imported by passing a string with its name to `import`:

```
const fs = import("fs")
```

##### `fs`

The `fs` module provides functions for interacting with the file system.

###### `fs.read`

`fs.read` takes a path to a file and returns a `str` with its contents. The file must be UTF-8 encoded.

```
const fs = import("fs")
fs.read("foo.txt")
```

###### `fs.readBytes`

`fs.readBytes` takes a path to a file and returns a `bytes` with its contents.

```
const fs = import("fs")
fs.readBytes("foo.txt")
```

## Planned Features

### Generators

### List Slicing

### Classes

### Format Strings

### Imports

### Error Handling

### Type Casting

### Variadic and Keyword Function Arguments

## Planned APIS

### `fs` Module

#### `fs.append`

#### `fs.write`

### `json` Module

### `path` Module
