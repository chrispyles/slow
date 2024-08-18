# The Slow Programming Language

<!-- TODO -->

## Reference

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

### Variables

Variables must be declared with the `var` statement.

```
var x
```

They can also be assigned in the same statement.

```
var x = 1
```

Variables are scoped to the frame they're declared in. Setting a variable in a child frame that will update the value of the variable in the frame in which it was declared.

### Operators

#### Binary Operators

Slow supports the following binary operators:

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
var l = [1, 2]
l.append(3)
print(l)
#> [1, 2, 3]
```

#### `equals`

<!-- TODO -->

### Control Flow

<!-- for, while, switch -->
Slow supports control flow with `for` and `while` loops.

`for` loops iterate over a pre-defined set of values returned by an [iterator](#iterators). The `for` loop has the syntax `for <loop variable> in <iterator>` followed by a body enclosed in curly brackets.

```
for i in range(20) {
  print(i)
}
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
