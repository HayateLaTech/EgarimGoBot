package command

import (
	"fmt"
	"log"
	"noeru/egarim/internal/model"
	"noeru/egarim/internal/util"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Contract(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var user discordgo.User

	if (i.Member != nil) {
		user = *i.Member.User
	} else {
		user = *i.User
	}

	subject := model.FirstOrCreateProfile(user)

	options := i.ApplicationCommandData().Options

	switch options[0].Name {

	case "sign":
		sign(s, i, subject)

	case "redeem":
		key := fmt.Sprintf("%v", options[0].Options[0].Value)
		redeem(s, i, subject, key)
	}
}

func redeem(s *discordgo.Session, i *discordgo.InteractionCreate, subject model.Subject, key string) {
	target, err := model.FindSubjectByKey(key)

	if err != nil {

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf(":warning: The Secret Key `%s` is **invalid**, baka!!", key),
			},
		})

		if err != nil {
			log.Fatal(err)
		}

	} else {

		if target.ID == subject.ID {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: ":warning: You can't redeem yourself, Baka!!",
				},
			})
	
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		// +1 Point for the Kill
		oldPoints := subject.Points
		subject.Points = subject.Points + 1
		util.DB.Save(&subject)

		// create new takedown notice
		takedown := model.RegNewKill(*target, subject)

		// reset secret of prey
		target.Secret = model.GenerateSecret(model.Size)
		util.DB.Save(&target)

		// send message in Notification Channel
		msg, err := s.ChannelMessageSendComplex(
			os.Getenv("NOTIFICATION_CHANNEL"),
			&discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:     "Takedown",
						Color:     13836063,
						Timestamp: takedown.CreatedAt.Format(time.RFC3339),
						Footer: &discordgo.MessageEmbedFooter{
							Text: "Mirage is proud!",
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "😈 Succubus",
								Value:  fmt.Sprintf("<@%s>", takedown.SuccId),
								Inline: true,
							},
							{
								Name:   "🤤 Prey",
								Value:  fmt.Sprintf("<@%s>", takedown.PreyId),
								Inline: true,
							},
							{
								Name:   "Point Update",
								Value:  fmt.Sprintf("%v >> **%v**", oldPoints, subject.Points),
								Inline: false,
							},
						},
					},
				},
			},
		)

		if err != nil {
			log.Fatal(err)
		}

		// save msgId in Takedown
		takedown.NotifMsgId = msg.ID
		util.DB.Save(&takedown)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Redeemed!",
			},
		})

		if err != nil {
			log.Fatal(err)
		}

	}
}

func sign(s *discordgo.Session, i *discordgo.InteractionCreate, subject model.Subject) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Your Contract",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Your Secret Key",
							Value: subject.Secret,
						},
						{
							Name: "Your Points",
							Value: fmt.Sprintf("**%v**", subject.Points),
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Pledge your Soul to Mirage~",
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		log.Fatalf("sign %s > error!", subject.Userid)
	}
}
