---
title: Operators
weight: 4
---

# Operators

## Binary Operators

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

The left operand of a reassignment oeprators may only be an already-declared variable, an object field, or an [index]({{< relref "07-indexing.md" >}}).

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

The `==` and `!=` support all types. (Note that all [non-primitive types]({{< relref "02-non-primitive-types" >}}) are compared by reference and not by value.) The other comparison operators only support numeric types and strings.

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

## Unary Operators

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
