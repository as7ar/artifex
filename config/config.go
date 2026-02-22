package config

import (
	"flag"

	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/voice"
)

var (
	// Default Data

	TOKEN       string
	Prefix      string
	SymbolColor int = 0xffea94

	IsDebug bool

	// System Data

	DISCORD *session.Session
	VOICE   *voice.Session
)

func init() {
	debug := flag.Bool("debug", true, "true/false")
	prefix := flag.String("prefix", "!", "")

	flag.Parse()
	IsDebug = *debug
	Prefix = *prefix
}
