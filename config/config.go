package config

import (
	"flag"

	noori2 "github.com/as7ar/noori/noori"
)

var (
	// Default Data
	TOKEN       string
	IsDebug     bool
	Prefix      string
	SymbolColor int = 0xffea94

	NOORI *noori2.App
)

func init() {
	debug := flag.Bool("debug", true, "true/false")
	prefix := flag.String("prefix", "!", "")

	flag.Parse()
	IsDebug = *debug
	Prefix = *prefix
}
