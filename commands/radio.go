package commands

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	"github.com/as7ar/noori/logger"
	"github.com/as7ar/noori/utils"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/diamondburned/arikawa/v3/voice/voicegateway"
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

func RadioCommand(st *state.State, c *gateway.MessageCreateEvent, url string) {
	guildID := c.GuildID

	vscID, ok := utils.GetUserVoiceChannelID(st, guildID, c.Author.ID)
	if !ok || vscID == 0 {
		return
	}

	guildVoiceManager.Lock()
	if vc, exists := guildVoiceManager.sessions[guildID]; exists {
		vc.cancel()
		_ = vc.vs.Leave(context.Background())
		delete(guildVoiceManager.sessions, guildID)
	}
	guildVoiceManager.Unlock()

	v := voice.NewSessionCustom(st, c.Author.ID)

	ctx, cancel := context.WithCancel(context.Background())

	ready := make(chan struct{})
	v.AddHandler(func(*voicegateway.ReadyEvent) {
		close(ready)
	})

	v.AddHandler(func(e *voice.ReconnectError) {
		logger.Err("Voice reconnect error:", e.Err)
	})

	wait := make(chan struct{})
	go func() {
		time.Sleep(500 * time.Millisecond)
		close(wait)
	}()
	<-wait
	if err := v.JoinChannelAndSpeak(ctx, vscID, false, false); err != nil {
		cancel()
		logger.Err("join failed:", err)
		return
	}

	<-ready

	guildVoiceManager.Lock()
	guildVoiceManager.sessions[guildID] = &VoiceController{
		vs:     v,
		cancel: cancel,
	}
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
			logger.Err("play error:", err)
		}
	}()
}

func StopCommand(c *gateway.MessageCreateEvent) {
	guildVoiceManager.Lock()
	defer guildVoiceManager.Unlock()

	if vc, ok := guildVoiceManager.sessions[c.GuildID]; ok {
		vc.cancel()
		_ = vc.vs.Leave(context.Background())
		delete(guildVoiceManager.sessions, c.GuildID)
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
