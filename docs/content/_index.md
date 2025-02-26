+++
title = 'Home'
type = 'docs'
+++

# The Slow Programming Language

<div style="display: flex; justify-content: center; margin: 2rem;">
  {{< logo size="200" >}}
</div>

The Slow programming language is a dynamically-typed interpreted programming language implemented in Go.

## Installation

Currently, the only way to run the Slow interpreter is to build it from source. This can be done using the `go` CLI:

```console
$ go install github.com/chrispyles/slow@latest
```

This will add the `slow` binary to your Go `bin` directory, so make sure it is in your `$PATH`. Slow uses Go generics, so you must have Go 1.18 or later. 

## Usage

The `slow` interpreter has two modes: script execution and live interpretation. To launch the Slow interpreter, just run the `slow` command:

```console
$ slow
```

To run a Slow script (idiomatically a `.slo` file), pass the path to the file to the `slow` command:

```console
$ slow main.slo
```

To launch the Slow interpreter after executing a script, use `slow -i`, like the Python CLI:

```console
$ slow -i main.slo
```
