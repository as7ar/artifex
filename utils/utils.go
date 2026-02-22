package utils

import (
	"github.com/as7ar/noori/logger"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

func GetUserVoiceChannelID(s *state.State, guildID discord.GuildID, userID discord.UserID) (discord.ChannelID, bool) {
	/*voiceStates, err := s.Cabinet.VoiceStates(guildID)
	for _, voiceState := range voiceStates {
		if !voiceState.ChannelID.IsValid() {
			return 0, false
		}

		voiceState.
	}
	if err != nil {
		logger.Err(err)
		return 0, false
	}

	return 0, false*/

	voiceState, err := s.Cabinet.VoiceState(guildID, userID)
	if err != nil {
		logger.Err(err) // item not found in store?
		return 0, false
	}

	//fmt.Println(voiceState.ChannelID)

	if !voiceState.ChannelID.IsValid() {
		return 0, false
	}

	return voiceState.ChannelID, true
}
