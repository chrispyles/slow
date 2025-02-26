---
title: Control Flow
weight: 5
---

# Control Flow

## Conditionals

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

## Iteration

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

### Iterators

The iterator in a `for` loop is a built-in type in Slow. `list`s and `str`s come with iterators, and there is also a [generator type]({{< relref "planned.md#generators" >}}) that backs built-in functions like [`range`]({{< relref "08-builtins.md#range" >}}).

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