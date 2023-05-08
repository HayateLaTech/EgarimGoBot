package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Survive(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userid := i.Member.User.ID

	se := fmt.Sprintf("<@%s> survived! *headpats* **+50 points for Gryffindor!**", userid)

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.Message.ChannelID,
		Components: []discordgo.MessageComponent{},
		Embeds:     []*discordgo.MessageEmbed{},
		Content:    &se,
	})

	if err != nil {
		panic(err)
	}
}
