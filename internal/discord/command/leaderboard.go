package command

import (
	"fmt"
	"log"
	"noeru/egarim/internal/model"
	"noeru/egarim/internal/util"

	"github.com/bwmarrin/discordgo"
)

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var subjects []model.Subject
	var fields []*discordgo.MessageEmbedField

	util.DB.Where("Points > 0").Order("Points DESC").Limit(10).Find(&subjects)

	for i, subject := range subjects {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%v.", i+1),
			Value: fmt.Sprintf("<@%s>\n\n**%v** <:MirageClose:712522550592143391>", subject.Userid, subject.Points),
		})
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  "Leaderboard",
					Fields: fields,
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
