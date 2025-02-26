---
title: fs
---

# `fs`

The `fs` module provides functions for interacting with the file system.

## `fs.read`

`fs.read` takes a path to a file and returns a `str` with its contents. The file must be UTF-8 encoded.

```
const fs = import("fs")
fs.read("foo.txt")
```

## `fs.readBytes`

`fs.readBytes` takes a path to a file and returns a `bytes` with its contents.

```
const fs = import("fs")
fs.readBytes("foo.txt")
```
