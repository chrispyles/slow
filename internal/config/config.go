package config

import "flag"

var (
	Debug = flag.Bool("debug", false, "print asts and values")
)
