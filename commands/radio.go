package commands

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	"github.com/as7ar/noori/config"
	"github.com/as7ar/noori/embeds"
	"github.com/as7ar/noori/logger"
	"github.com/as7ar/noori/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/oggreader"
)

var guildVoiceManager = struct {
	sync.Mutex
	sessions map[discord.GuildID]*VoiceController
}{
	sessions: make(map[discord.GuildID]*VoiceController),
}

type VoiceController struct {
	vs     *voice.Session
	cancel context.CancelFunc
}

func RadioCommand(s *session.Session, c *gateway.MessageCreateEvent, url string) {
	guildID := c.GuildID
	chatID := c.ChannelID

	vscID, err := utils.GetUserVoiceChannelID(config.NOORI.STATE, guildID, c.Author.ID)
	if err != nil || vscID == 0 {
		_, _ = s.SendEmbedReply(chatID, c.Message.ID,
			embeds.New().
				Color(config.SymbolColor).
				Title("🚨 Warning").
				Description("You must be entered in `voice chat`").Build())
		return
	}

	guildVoiceManager.Lock()
	if vc, ok := guildVoiceManager.sessions[guildID]; ok {
		vc.cancel()
		_ = vc.vs.Leave(context.Background())
		delete(guildVoiceManager.sessions, guildID)
	}
	guildVoiceManager.Unlock()

	v, err := voice.NewSession(s)
	if err != nil {
		logger.Err("Failed to create voice session", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := v.JoinChannelAndSpeak(ctx, vscID, false, true); err != nil {
		logger.Err("Failed to join voice channel", err)
		cancel()
		return
	}

	vc := &VoiceController{vs: v, cancel: cancel}
	guildVoiceManager.Lock()
	guildVoiceManager.sessions[guildID] = vc
	guildVoiceManager.Unlock()

	go func() {
		defer func() {
			cancel()
			_ = v.Leave(context.Background())
			guildVoiceManager.Lock()
			delete(guildVoiceManager.sessions, guildID)
			guildVoiceManager.Unlock()
		}()

		if err := playYT(ctx, v, url); err != nil && !errors.Is(err, context.Canceled) {
			logger.Err("Failed to play YouTube audio", err)
		}
	}()
}

func StopCommand(c *gateway.MessageCreateEvent) {
	guildVoiceManager.Lock()
	defer guildVoiceManager.Unlock()

	if vc, ok := guildVoiceManager.sessions[c.GuildID]; ok {
		vc.cancel()
	}
}

func playYT(ctx context.Context, vs *voice.Session, url string) error {
	ytdl := exec.CommandContext(ctx, "yt-dlp", "-f", "bestaudio", "-g", url)
	videoURLBytes, err := ytdl.Output()
	if err != nil {
		return err
	}
	videoURL := string(videoURLBytes)

	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-i", videoURL,
		"-f", "opus",
		"-ar", "48000",
		"-ac", "2",
		"-loglevel", "quiet",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := oggreader.DecodeBuffered(vs, stdout); err != nil {
		return err
	}

	return cmd.Wait()
}
