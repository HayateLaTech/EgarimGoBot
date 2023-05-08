package modal

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func BirthdayResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ModalSubmitData()

	userid := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).CustomID
	val := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	regex, regerr := regexp.Compile("([0-9]{1,3})([hms])")

	if regerr != nil {
		panic(regerr)
	}

	matches := regex.FindAllStringSubmatch(val, 3)

	if len(matches) == 0 {
		return
	}

	var (
		hh = 0
		mm = 0
		ss = 0
	)

	for i := 0; i < len(matches); i++ {
		ci, err := strconv.Atoi(matches[i][1])
		if err != nil {
			continue
		}

		switch matches[i][2] {
		case "h":
			hh += ci
		case "m":
			mm += ci
		case "s":
			ss += ci
		}
	}

	// add birthday role to member
	errole := s.GuildMemberRoleAdd(os.Getenv("DISCORD_GUILD"), userid, os.Getenv("BIRTHDAY_ROLE"))

	if errole != nil {
		panic(errole)
	}

	// create goroutine that runs besides the main thread
	go scheduleBirthdayRoleRemoval(s, userid, hh, mm, ss)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s> receives the Birthday Role for %v Hours : %v Minutes : %v Seconds!", userid, hh, mm, ss),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		panic(err)
	}
}

func scheduleBirthdayRoleRemoval(s *discordgo.Session, uid string, hh int, mm int, ss int) {
	nowTime := time.Now()
	thenTime := nowTime.Add(time.Duration(hh) * time.Hour).Add(time.Duration(mm) * time.Minute).Add(time.Duration(ss) * time.Second)

	time.Sleep(thenTime.Sub(nowTime))

	// remove birthday role again
	err := s.GuildMemberRoleRemove(os.Getenv("DISCORD_GUILD"), uid, os.Getenv("BIRTHDAY_ROLE"))

	if err != nil {
		log.Print(err)
	}
}
