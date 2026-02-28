package utils

import (
	"github.com/as7ar/noori/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

func GetUserVoiceChannelID(s *state.State, guildID discord.GuildID, userID discord.UserID) (discord.ChannelID, bool) {

	voiceState, err := s.Cabinet.VoiceState(guildID, userID)
	if err != nil {
		logger.Err(err) // item not found in store?
		return 0, false
	}

	logger.Debug("vcs found:", voiceState.ChannelID)

	return voiceState.ChannelID, true
}
