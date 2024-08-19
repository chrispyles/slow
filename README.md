# The Slow Programming Language

<!-- TODO -->

## Reference

<!-- TODO: intro paragraph -->

Some (not all) of the examples below use the interpreter format. In this format, statements are prefixed with `_>`, the interpreter prompt, and outputs are prefixed with `#>`. If a statement takes up multiple lines, all lines after the first will be prefixed with `..`. This format mimicks the Slow interpreter.

```
# This is the non-interpreter format
var x = 2
print(x)
```

```
_> # This is the interpreter format
_> var x = 2
#> 2
_> print(x)
#> 2
```

### Data Types

Slow has the following primitive data types:

- `bool`: `true` or `false`
- `float`: `1.`, `2.1`, `3.14`, `.02` etc.
- `int`: `1`, `2`, `-1`, etc.
- `null`
- `str`: `"The quick brown fox jumps over the lazy dog."`
- `uint`: `0u`, `1u`, `2u`, etc.

All numeric data types (except `bool`) are 64 bits.

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

### Variables

Variables must be declared with the `var` statement.

```
var x
```

They can also be assigned in the same statement.

```
_> var x = 1
#> 1
```

Variables are scoped to the frame they're declared in. Setting a variable in a child frame that will update the value of the variable in the frame in which it was declared.

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

Each of the arithmetic operators has a reassignment variant that reassigns its left operand to the value of the expression; the syntax for these variants is the arithmetic operator suffixed with an `=`.

```
_> var x = 2
#> 2
_> x += 1
#> 3
_> x //= 2
#> 1
_> print(x)
#> 1
```

Slow supports the following logical operators:

| Operator | Description |
|----------|-------------|
| `&&`     | logical and |
| `\|\|`   | logical or  |
| `^^`     | logical xor |

Both the `&&` and `||` operators short-circuit; that is, the second operand of `&&` and `||` are only evaluated if the first is falsely and truthy, respectively. `&&` and `||` also do not change the types of their operands (e.g. `1 && 2` returns `1`, not `true`), but `^^` always returns a `bool`. All logical operators also have a reassignment variant (`&&=`, `||=`, `^^=`).

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

#### Ternary Operator

Slow supports a ternary operator of the form `<condition> ? <value if true> : <value if false>`. Only one branch of the operator is ever evaluated.

```
# In the example below, if i is even, only f is called, otherwise
# only g is called.
var x = i % 2 == 0 ? f() : g()
```

### Lists

Slow has a built-in `list` type backed by a Go slice. Lists are declared using square brackets

```
var l = []
```

You can also specify elements when declaring a list:

```
l = [1, 2, 3]
```

#### List Methods

The `list` type has several built-in methods, each of which is described below.

<!-- TODO: more methods -->

#### `list.append`

The `append` method of `list` adds an element to the end of that `list` in-place.

```
_> var l = [1, 2]
#> [1, 2]
_> l.append(3)
_> print(l)
#> [1, 2, 3]
```

#### `list.equals`

The `equals` method of `list` compares it against another value. If the other value is also a list and each element of the two lists are equal (either by the `==` operator or `list.equals` if the corresponding elements are both themselves `list`s).

```
_> var l1 = [1, 2, 3]
#> [1, 2, 3]
_> l1.equals(1)
#> false
_> l1.equals([1, 2])
#> false
_> l1.equals([1, 2, 3])
#> true
_> [[1, 2], [2, 3]].equals([[1, 2], [2, 3]])
#> true
```

### Conditionals

Slow supports `if` statements to conditionally execute code blocks. To run a block should the condition evaluate to `false`, use an `else` statement.

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

Because every statement in Slow evaluates to a value, the return value of the last statement in an `if`/`else if`/`else` block's body is the value of the statement.

```
_> var x = 1
#> 1
_> if x % 2 == 0 {
..   "even"
.. }
.. else {
..   "odd"
.. }
..
#> "odd"
```

Slow also supports `switch` statements, which match a value to a series of possible `case`s using the `==` operator. Unlike many other languages, Slow's `switch` cases **do not** fall through by default; the `fallthrough` keyword must be used to trigger fall through. Slow also uses curly braces to wrap `case` bodies.

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
_> for i in range(5) {
..   print(i)
.. }
..
#> 0
#> 1
#> 2
#> 3
#> 4
```

`while` loops iterate while a condition evaluates to a truthy value and have the syntax `while <condition>` followed by a body enclosed in curly brackets.

```
var x = 10
while x > 0 {
  print(x)
  x--
}
```

A loop can run indefinitely by setting its condition to a value that is always truthy:

```
while true {
  print("still running...")
}
```

#### Iterators

The iterator in a `for` loop is a built-in type in Slow. `list`s and `str`s come with iterators, and there is also a [generator type](#generators) that backs built-in functions like [`range`](#range).

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

### Maps

### Classes

### Format Strings

### Imports

### Error Handling
