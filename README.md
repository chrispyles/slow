# The Slow Programming Language

[![Run tests](https://github.com/chrispyles/slow/actions/workflows/run-tests.yml/badge.svg)](https://github.com/chrispyles/slow/actions/workflows/run-tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/chrispyles/slow/badge.svg?branch=main)](https://coveralls.io/github/chrispyles/slow?branch=main)

<p align="center">
  <img src="logos/logo.png" width="256" alt="Slow logo">
</p>

The Slow programming language is a dynamically-typed interpreted programming language implemented in Go.

## Installation

Currently, the only way to run the Slow interpreter is to build it from source. This can be done using the `go` CLI:

```console
$ go install github.com/chrispyles/slow/cmd/slow@latest
```

This will add the `slow` binary to your Go `bin` directory, so make sure it is in your `$PATH`. Slow requires Go 1.23 or later. 

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

## Reference

A complete reference of the Slow programming language is available in the [documnetation](https://slowlange.dev).
