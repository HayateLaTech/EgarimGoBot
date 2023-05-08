package command

import "github.com/bwmarrin/discordgo"

func Tempt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID: "prey_select",
			Content:  "Select your Prey!",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID: "prey",
							MenuType: discordgo.UserSelectMenu,
						},
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		panic(err)
	}
}
