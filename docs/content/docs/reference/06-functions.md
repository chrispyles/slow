---
title: Functions
weight: 6
---

# Functions

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

## Defer Statements

Inside a function, a function call can be deferred so that it runs just before the function exits, instead of wherever in the body the `defer` statement is (like Go's `defer` statement). Statements are accrued but not evaluated as the function's body executes and before the function exits, they are run. **Deferred functions are currently not run if the function throws an error.** This behavior will likely be added in the future.

```
-> func foo() {
..   # without the deferral, this statement would error because the variable does not yet exist
..   defer print(i)
..   var i = 0
..   for _ in :20 {
..     ++i
..   }
.. }
-> foo()
20
```
