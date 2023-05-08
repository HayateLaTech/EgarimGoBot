package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Prey(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()

	userid := data.Values[0]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Sent!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	_, err := s.ChannelMessageSendComplex(
		i.ChannelID,
		&discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s>", userid),
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Mirage attempts to tempt you..",
					Description: "You're trapped, paralyzed beneath a low level Imp..\nThrusting into this obvious Trap wasn't the brightest idea you had today, yet it feels so good..\n\nYour eyes roll back as her pussy begins to convulse and milk your sensitive cock..\nIt feels so good.. so heavenly..\nLet her milk it all out..\n**Submit**",
					Footer: &discordgo.MessageEmbedFooter{
						Text: ".. will you survive .. or will you nut~?",
					},
					Image: &discordgo.MessageEmbedImage{
						URL: "https://cdn.discordapp.com/attachments/617082313888628763/1091084666443939971/image0.gif",
					},
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("survive_id_%s", userid),
							Style:    discordgo.SuccessButton,
							Label:    "Survive",
						},
						discordgo.Button{
							CustomID: fmt.Sprintf("nut_id_%s", userid),
							Style:    discordgo.DangerButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🥜",
							},
							Label: "Nut",
						},
					},
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}
}
