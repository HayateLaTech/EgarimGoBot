package command

import (
	"fmt"
	"log"
	"noeru/egarim/internal/model"
	"noeru/egarim/internal/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Takedowns(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var takedowns []model.Takedown
	var fields []*discordgo.MessageEmbedField

	util.DB.Order("created_at DESC").Limit(20).Find(&takedowns)

	for _, takedown := range takedowns {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%v", takedown.CreatedAt.Format(time.ANSIC)),
			Value: fmt.Sprintf("<@%v>\n ⚔️ 🔽 \n<@%v>", takedown.SuccId, takedown.PreyId),
			Inline: true,
		})
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  "20 Latest Takedowns",
					Fields: fields,
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

}
