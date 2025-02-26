---
title: Modules
weight: 9
---

# Modules

Slow supports importing objects and functions from other Slow files or from built-in modules via the `module` type. `module`s are frozen; that is, new variables can't be declared in them nor can variables be rebound. A `module` can be created using the [`import` function]({{< relref "../08-builtins.md#import" >}}) and variables in its environment can be accessed use dot syntax.

For example, to import a `module` from another Slow file:

```
# logging.slo
func info(s) {
  print(s)
}
```

```
const logging = import("logging.slo")

logging.info("successfully imported logging.slo")
```

Or, to import a built-in module:

```
const fs = import("fs")

print(fs.read("foo.txt"))
```

