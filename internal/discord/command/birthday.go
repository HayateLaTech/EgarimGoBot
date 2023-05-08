package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Birthday(s *discordgo.Session, i *discordgo.InteractionCreate) {

	users := i.ApplicationCommandData().Resolved.Users
	var user *discordgo.User
	var name string

	for _, u := range users {
		user = u
	}

	name = user.Username

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "birthday_response",
			Title:    fmt.Sprintf("Give %s the Birthday Role", name),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    user.ID,
							Label:       "Duration",
							Placeholder: "24h 30m 5s",
							Style:       discordgo.TextInputShort,
							Required:    true,
						},
					},
				},
			},
		},
	})

	if err != nil {
		panic(err)
	}

}
