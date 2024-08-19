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

<!-- TODO -->

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

<!-- TODO: more methods -->

#### `append`

The `append` method of a list adds an element to the end of that list in-place.

```
_> var l = [1, 2]
#> [1, 2]
_> l.append(3)
_> print(l)
#> [1, 2, 3]
```

#### `equals`

<!-- TODO -->

### Control Flow

<!-- for, while, switch -->
Slow supports control flow with `for` and `while` loops.

`for` loops iterate over a pre-defined set of values returned by an [iterator](#iterators). The `for` loop has the syntax `for <loop variable> in <iterator>` followed by a body enclosed in curly brackets.

```
_> for i in range(5) {
..   print(i)
.. }
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

<!-- TODO -->

### Conditionals

<!-- TODO: if, switch -->

### Functions

<!-- TODO -->

```
func isFactor(x, y) {
  return x % y == 0
}
```

#### Built-in Functions

##### `exit`

##### `help`

##### `len`

##### `print`

##### `range`

##### `type`
