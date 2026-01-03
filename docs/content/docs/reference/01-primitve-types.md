---
title: 'Primitive Types'
weight: 1
---

# Primitive Types

Slow has the following primitive data types:

- `bool`: `true` or `false`
- `float`: `1.`, `2.1`, `3.14`, `.02` etc.
- `int`: `1`, `2`, `-1`, etc.
- `uint`: `0u`, `1u`, `2u`, etc.
- `str`: an immutable character sequence, `"The quick brown fox jumps over the lazy dog."`
- `bytes`: an immutable byte sequence, `0xDEADBEEF`
- `null`

All numeric data types (except `bool`) are 64 bits (backed by Go's 64-bit number types).

When converting from `bytes` to `float`, `int`, or `uint`, the bytes are decoded using the big-endian binary format. For example, `0xBEEF` is converted to a `uint` as `48879u`, not `61374u`.

All values are truthy except `false`, `0` (in all numeric types), `""`, a `bytes` object with all null bytes (e.g. `0x00`), and `null`.

## Strings

Strings are delimited by double quote `"` characters. They must be on a single line and allow the following escape sequences:

| Sequence | Description               |
|----------|---------------------------|
| `\"`     | double quote character    |
| `\n`     | newline character         |
| `\r`     | carriage return character |
| `\t`     | tab character             |
| `\\`     | backslash character       |

```
-> print("___garoo\rkan\njump")
kangaroo
jump
```

## Bytes

Bytes are written as case-insensitive hexadecimal values prefixed with `0x` (for example, `0xDEADBEEF`). There must be an even number of characters in a `bytes` literal.

```
-> var b = 0xdeadbeef
0xDEADBEEF
```

## Type Casting

Values of one primitive type can be cast to another using the `as` keyword.

```
-> true as int
1
-> 1u as float
1.0
```

The list of types that can be cast to is:

- `bool`
- `bytes`
- `float`
- `int`
- `str`
- `uint`

While there are other types (`func`, `list`, `module`, etc.), these types do not support type casting.

When numeric types are cast to `str`, the returned string contains the number printed in decimal format.

```
-> 1u as str
"1"
-> 1 as str
"1"
-> 1.0 as str
"1.000000"
```

When a `str` is cast to a numeric type, Slow attempts to parse the string as a number of that type and fails if the string is not a valid decimal number. Note that even though `uint`s are represented with the `u` suffix in Slow, `"1u"` will error if you try to cast it to a `uint`.

```
-> "1" as uint
1u
-> "1" as int
1
-> "1" as float
1.0
```
