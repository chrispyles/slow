---
title: Indexing
weight: 7
---

# Indexing

Slow supports indexing with square bracket syntax. All indexable types in Slow are zero-indexed.

For all indexable types except `map`s, only a `bool`, `uint`, or `int` may be used as an index. These types also support negative indexing to retrieve elements beginning at the end of the sequence. The index of the last element is `-1`, the second to last is `-2`, etc. `map`s can be indexed with any hashable value.

Note that under the hood, numeric indexes (excluding `map`s) are coerced to Go's `int` type (not `int64`, which is how Slow `int`s are stored). This means it is possible to overflow the range of possible index values if you use a `uint` that's too large. 

Any time an indexable value is indexed with an out-of-bounds or nonexistent index, an `IndexError` is thrown.

## String Indexing

String characters can be accessed used zero-indexed integers. Indexing returns a single-character `str`.

```
-> var s = "abcdef"
"abcdef"
-> s[0]
"a"
-> s[-1]
"f"
```

Because strings are immutable, index assignments are not supported.

## Bytes Indexing

Individual bytes can be accessed used zero-indexed integers. Indexing returns a single-byte `bytes`.

```
-> var b = 0xDEADBEEF
0xDEADBEEF
-> b[0]
0xDE
-> b[-1]
0xEF
```

Because `bytes` are immutable, index assignments are not supported.

## List Indexing

Lists are zero-indexed. To retrieve the element of a `list` at a particular index, use square brackets:

```
-> var l = [1, 2]
[1, 2]
-> l[1]
2
```

List elements can also be updated using indexing:

```
l[0] = 3
l[1] += 1
++l[2]
```


```
-> var l = [1, 2, 3, 4, 5]
[1, 2, 3, 4, 5]
-> l[-1]
5
-> l[-2]
4
```

### List Slicing

[Ranges]({{< relref "05-control-flow.md#ranges" >}}), or any [generator]({{< relref "05-control-flow.md#generators" >}}) that yields numeric values, can be used for slicing lists. The start, stop, and step values of a range are optional when slicing under certain circumstances, described by the rules below.

- `step` is always optional and defaults to `1u`
- `start` defaults to `0u` if `step` is non-negative or the list length if it is negative
- `stop` defaults to the list length if `step` is non-negative or `0u` if it is negative

Here are some example slices using ranges, with annotations indicating whether the ranges are valid in other contexts:

{{< inputOutput >}}

{{< codeWithCaption cmd="cat" file="slicing.slo" >}}
var l = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
print(l[:])         # only valid for slicing
print(l[2:4])
print(l[4:2:-1])
print(l[4:2])
print(l[:5])
print(l[5:])        # only valid for slicing
print(l[::-1])      # only valid for slicing
print(l[:5:-1])     # only valid for slicing
print(l[0::2])      # only valid for slicing
print(l[:8:2])
print(l[range(3)])
{{< /codeWithCaption >}}

{{< codeWithCaption cmd="slow" file="slicing.slo" >}}
[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
[2, 3]
[4, 3]
[]
[0, 1, 2, 3, 4]
[5, 6, 7, 8, 9]
[9, 8, 7, 6, 5, 4, 3, 2, 1]
[9, 8, 7, 6]
[0, 2, 4, 6, 8]
[0, 2, 4, 6]
[0, 1, 2]
{{< /codeWithCaption >}}

{{< /inputOutput >}}

## Map Indexing

Map values can be retrieved or updated using indexing. To retrieve the element of a `map` key, use square brackets:

```
-> var m = {1: 2, 2: 4}
{1: 2, 2: 4}
-> m[1]
2
-> m[1] = 3
{1: 3, 2: 4}
-> ++m[2]
4
-> m
{1: 3, 2: 5}
-> m[3] = 6
6
-> m
{1: 3, 2: 5, 3: 6}
```
