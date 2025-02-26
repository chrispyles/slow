---
title: Maps
weight: 2
---

# Maps

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

## Map Methods

The `map` type has several built-in methods, each of which is described below.

### `map.get`

The `get` method of a `map` returns the value corresponding to the provided key.

```
var m = {1: 2}
m.get(1)
> 2
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

### `map.set`

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

### `map.to_mutable`

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

### `map.to_immutable`

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