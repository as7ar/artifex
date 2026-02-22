package noori

import (
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
)

type App struct {
	SESSION *session.Session
	STATE   *state.State
	VOICE   *voice.Session
}
