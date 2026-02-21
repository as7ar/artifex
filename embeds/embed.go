package embeds

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	LimitTitle       = 256
	LimitDescription = 2048
	LimitFieldName   = 256
	LimitFieldValue  = 1024
	LimitFields      = 25
	LimitFooter      = 2048
)

type Builder struct {
	embed *discordgo.MessageEmbed
}

func New() *Builder {
	return &Builder{
		embed: &discordgo.MessageEmbed{},
	}
}

func (b *Builder) Title(v string) *Builder {
	b.embed.Title = truncate(v, LimitTitle)
	return b
}

func (b *Builder) Description(v string) *Builder {
	b.embed.Description = truncate(v, LimitDescription)
	return b
}

func (b *Builder) Color(v int) *Builder {
	b.embed.Color = v
	return b
}

func (b *Builder) URL(v string) *Builder {
	b.embed.URL = v
	return b
}

func (b *Builder) Author(name, iconURL, url string) *Builder {
	b.embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
		URL:     url,
	}
	return b
}

func (b *Builder) Footer(text, iconURL string) *Builder {
	b.embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    truncate(text, LimitFooter),
		IconURL: iconURL,
	}
	return b
}

func (b *Builder) Thumbnail(url string) *Builder {
	b.embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
	return b
}

func (b *Builder) Image(url string) *Builder {
	b.embed.Image = &discordgo.MessageEmbedImage{URL: url}
	return b
}

func (b *Builder) Field(name, value string, inline bool) *Builder {
	if len(b.embed.Fields) >= LimitFields {
		return b
	}

	name = truncate(name, LimitFieldName)

	chunks := splitValue(value, LimitFieldValue)

	for i, chunk := range chunks {
		fieldName := name
		if i > 0 {
			fieldName = name + " (cont.)"
		}

		b.embed.Fields = append(b.embed.Fields, &discordgo.MessageEmbedField{
			Name:   fieldName,
			Value:  chunk,
			Inline: inline,
		})
	}

	return b
}

func (b *Builder) Build() *discordgo.MessageEmbed {
	return b.embed
}

func truncate(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit]
}

func splitValue(s string, limit int) []string {
	if len(s) <= limit {
		return []string{s}
	}

	var result []string
	for len(s) > limit {
		cut := strings.LastIndexAny(s[:limit], " \n-")
		if cut <= 0 {
			cut = limit
		}
		result = append(result, s[:cut])
		s = s[cut:]
	}
	if len(s) > 0 {
		result = append(result, s)
	}
	return result
}

func Generic(title, desc string) *discordgo.MessageEmbed {
	return New().
		Title(title).
		Description(desc).
		Color(0x1c1c1c).
		Build()
}

func Error(title, desc string) *discordgo.MessageEmbed {
	return New().
		Title(title).
		Description(desc).
		Color(0xb40000).
		Build()
}
