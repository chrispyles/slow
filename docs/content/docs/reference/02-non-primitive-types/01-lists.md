---
title: Lists
weight: 1
---

# Lists

Slow has a built-in `list` type backed by a Go slice. Lists literals are declared using square brackets:

```
var l = []
```

You can also specify elements when writing a list literal:

```
l = [1, 2, 3]
```

Lists can be either mutable or immutable; all lists are mutable by default, by an immutable copy of any list can be created with the `to_immutable` method described below. (Similarly, a mutable copy of any list can be created with the `to_mutable` method.) Immutable lists do not allow any modification (e.g. index assignment, `list.append`). However, making an immutable list does not make its elements themselves immutable.

## List Methods

The `list` type has several built-in methods, each of which is described below.

### `list.append`

The `append` method of `list` adds an element to the end of that `list` in-place.

```
-> var l = [1, 2]
[1, 2]
-> l.append(3)
-> print(l)
[1, 2, 3]
```

### `list.equals`

The `equals` method of `list` compares it against another value. If the other value is also a list and each element of the two lists are equal (either by the `==` operator or `list.equals` if the corresponding elements are both themselves `list`s).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> l1.equals(1)
false
-> l1.equals([1, 2])
false
-> l1.equals([1, 2, 3])
true
-> [[1, 2], [2, 3]].equals([[1, 2], [2, 3]])
true
```

### `list.to_mutable`

The `to_mutable` method of `list` creates a mutable copy of the list. This method can be used on any list (immutable or mutable).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> var l2 = l1.to_mutable()
[1, 2, 3]
-> l2.append(4)
-> l1
[1, 2, 3]
-> l2
[1, 2, 3, 4]
```

### `list.to_immutable`

The `to_immutable` method of `list` creates an immutable copy of the list. This method can be used on any list (mutable or immutable).

```
-> var l1 = [1, 2, 3]
[1, 2, 3]
-> var l2 = l1.to_immutable()
[1, 2, 3]
-> l1.append(4)
-> l1
[1, 2, 3, 3]
-> l2.append(4)
ValueError: list is immutable
```