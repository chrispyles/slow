---
title: 'Planned Features'
---

# Planned Features

The features listed below are not yet implemented but are planned and may include proposed syntax.

## Ternary Conditional Operator

Slow will support a ternary conditional operator of the form `<condition> ? <value if true> : <value if false>`. Only one branch of the operator is ever evaluated.

```
# In the example below, if i is even, only f is called, otherwise
# only g is called.
var x = i % 2 == 0 ? f() : g()
```

## Generators

```
func myRange(stop) {
  for i in range(stop) {
    yield i
  }
}

for i in myRange(20) {
  print(i)
}
```

## List Slicing

```
l[0:2]
l[1:]
l[:2]
l[:-1]
l[-2:]
```

## Comprehensions

```
var l2 = [2 * i for i in l if i % 2 == 0]
var m2 = {i: 2 * i for i in l}

for i in (2 * i for i in l) {
  print(i)
}
```

## Classes

```
class MyClass {
  var publicInts = []

  private var someInts = []

  private readonly var aValue = "foo"

  func :init(firstInt) {
    this.someInts.append(firstInt)
  }

  func :str() {
    return "<MyClass>"
  }

  func add(anotherInt) {
    this.someInts.append(anotherInt)
  }
}

var c = MyClass(1)
```

## String Interpolation and Formatting

```
var name = "John"
print("Hello, {{ name }}")
var aFloat = 1.0000000003
print("{{ aFloat:.3f }}")
```

## Error Throwing and Handling

```
try {
  throw Error()
} catch {
  handleError()
}
```

## Variadic and Keyword Function Arguments

```
func foo(x, y = 2, *args, **kwargs) {
  ...
}
```

## List and Map Destructuring

```
var [a, b, c] = [1, 2, 3]
var {d, e, f} = {"d": 4, "e": 5, "f": 6}
```

## Sets

Both mutable and immutable variants.

## Planned APIs

### `fs` Module

#### `fs.append`

#### `fs.write`

### `json` Module

### `path` Module
