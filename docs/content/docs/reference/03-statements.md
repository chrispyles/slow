---
title: 'Statements and Variables'
weight: 3
---

# Statements

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

## Variables

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
