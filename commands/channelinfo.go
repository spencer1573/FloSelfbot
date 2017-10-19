package commands

import (
	"fmt"
	"strconv"
	"time"

	"strings"

	"github.com/bwmarrin/discordgo"
)

// Cinfo struct handles Cinfo Command
type Cinfo struct{}

func (c *Cinfo) message(ctx *Context) {
	em := createEmbed(ctx)
	em.Fields = make([]*discordgo.MessageEmbedField, 0)
	if ctx.Channel.Type != discordgo.ChannelTypeDM && ctx.Channel.Type != discordgo.ChannelTypeGroupDM {
		em.Author = &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Channel Info: %s", ctx.Channel.Name),
			IconURL: "https://twemoji.maxcdn.com/36x36/2139.png",
		}
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "ID",
			Value:  ctx.Channel.ID,
			Inline: true,
		})
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Type",
			Value:  strconv.FormatInt(int64(ctx.Channel.Type), 10),
			Inline: true,
		})
		var msg *discordgo.Message
		msgs, _ := ctx.Sess.ChannelMessages(ctx.Channel.ID, 1, ctx.Mess.ID, "", "")
		msg = msgs[0]
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Last Message",
			Value:  fmt.Sprintf("<@%s>: %s", msg.Author.ID, msg.Content),
			Inline: true,
		})
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Position",
			Value:  strconv.Itoa(ctx.Channel.Position),
			Inline: true,
		})
		if ctx.Channel.Topic != "" {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   "Topic",
				Value:  ctx.Channel.Topic,
				Inline: true,
			})
		}
		var hidden int
		for _, m := range ctx.Guild.Members {
			perms, err := ctx.Sess.State.UserChannelPermissions(m.User.ID, ctx.Channel.ID)
			if err == nil && perms&discordgo.PermissionReadMessages == 0 {
				hidden++
			}
		}
		em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
			Name:   "Members",
			Value:  fmt.Sprintf("%d/%d members can see this channel", len(ctx.Guild.Members)-hidden, len(ctx.Guild.Members)),
			Inline: true,
		})
		if t, err := ctx.GetCreationTime(ctx.Channel.ID); err == nil {
			em.Fields = append(em.Fields, &discordgo.MessageEmbedField{
				Name:   "Creation",
				Value:  fmt.Sprintf("%s (%.2f days ago)", t.Format(time.ANSIC), time.Now().Sub(t).Hours()/24),
				Inline: true,
			})
		}
	} else {
		recipients := make([]string, 0, len(ctx.Channel.Recipients))
		for _, recipient := range ctx.Channel.Recipients {
			recipients = append(recipients, recipient.Username)
		}
		em.Author = &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("DM Info: %s", strings.Join(recipients, ", ")),
			IconURL: "https://twemoji.maxcdn.com/36x36/2139.png",
		}
	}
	ctx.SendEm(em)
}

func (c *Cinfo) description() string             { return "Channel info" }
func (c *Cinfo) usage() string                   { return "" }
func (c *Cinfo) detailed() string                { return "Returns channel info for the channel you are currently in." }
func (c *Cinfo) subcommands() map[string]Command { return make(map[string]Command) }
