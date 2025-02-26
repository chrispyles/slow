---
title: Built-in Functions
weight: 8
---

# Built-in Functions

Slow has a few functions built into the language. They are declared in a frozen frame that is the parent of the frame that the global environment is declared in.

## `exit`

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

## `import`

The `import` function imports a module. It takes 1 argument, either a path to another Slow file (with extesnion `.slo`) or the name a of a built-in module (e.g. `fs`) and returns a [`module`]({{< relref "09-modules" >}}). If the argument is a path, the file is read and executed, and the resulting global environment is converted to a `module`.

```
const fs = import("fs")
const utils = import("utils.slo")
```

## `len`

The `len` function returns the length of the provided value if it is supported. Currently, the only types in Slow that have lengths are strings, `bytes`, `list`s, and `map`s. `len` always returns a `uint`.

```
-> var l = [1, 2, 3]
[1, 2, 3]
-> len(l)
3u
-> len({1: 2, 3: 4})
2u
```

## `print`

`print` prints its arguments to stdout followed by a newline character. It can accept any number of arguments, converts them to their string representation, and concatenates those strings.

```
-> print(1u, 2, 3., "foo")
1u23.0foo
```

## `range`

The `range` function creates a generator for a range of numbers. It can take either 1, 2, or 3 arguments.

- If it receives 1 argument `n`, the generator it returns yields the sequence `0, 1, 2, ..., n-1`.
- If it receives two arguments `m` and `n`, the generator yields the sequence `m, m+1, m+2, ..., n-1`. If `m > n`, the generator is infinite.
- If it receives three arguments `m`, `n`, and `k`, the generator yields the sequence `m, m+k, m+2*k, ..., m+jk` where `m+jk` is the largest possible value less than `n`.

```
range(10)        # yields 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
range(2, 10)     # yields 2, 3, 4, 5, 6, 7, 8, 9
range(2, 10, 3)  # yields 2, 5, 8
```

## `type`

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