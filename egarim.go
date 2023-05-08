package main

import (
	"log"
	"noeru/egarim/internal/discord/command"
	"noeru/egarim/internal/discord/component"
	"noeru/egarim/internal/discord/modal"
	"noeru/egarim/internal/model"
	"noeru/egarim/internal/util"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Bot parameters
var (
	RemoveCommands        = false
	InitRemoveAppCommands = false
	dbFile                = "egarim.db"
)

var s *discordgo.Session

func init() {
	log.Print("Opening DB Connection..")
	DB, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	util.DB = DB

	// Migrate the schemas
	util.DB.AutoMigrate(
		&model.Subject{},
		&model.Takedown{},
	)
}

func init() {
	var err error

	// loading .env
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal(".env could not be loaded!")
	}

	s, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	birthdayCommandPermissions int64 = discordgo.PermissionManageRoles
	nsfw                       bool  = true

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "Birthday",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Type:                     discordgo.UserApplicationCommand,
			DefaultMemberPermissions: &birthdayCommandPermissions,
			GuildID:                  os.Getenv("DISCORD_GUILD"),
		},
		{
			Name:        "tempt",
			Description: ".. use Egarim to tempt a Member!",
			Type:        discordgo.ChatApplicationCommand,
			NSFW:        &nsfw,
		},
		{
			Name:        "leaderboard",
			Description: "Who got the most Points~? I wonder..~",
			Type:        discordgo.ChatApplicationCommand,
			NSFW:        &nsfw,
		},
		{
			Name:        "takedowns",
			Description: "Display the last 20 Takedowns",
			Type:        discordgo.ChatApplicationCommand,
			NSFW:        &nsfw,
		},
		{
			Name:        "contract",
			Description: "Contracts are like Profiles and totally not like.. evil.. or stuff!",
			Type:        discordgo.ChatApplicationCommand,
			NSFW:        &nsfw,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "sign",
					Description: "Do you want to sign a Contract with me~?",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "redeem",
					Description: "Redeem a Secret Key to gain Points~!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "key",
							Description: "The Secret Key you wanna redeem",
							Required:    true,
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"Birthday": command.Birthday,

		"tempt": command.Tempt,

		"contract": command.Contract,

		"leaderboard": command.Leaderboard,

		"takedowns": command.Takedowns,
	}

	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"nut": component.Nut,

		"survive": component.Survive,

		"prey": component.Prey,
	}

	modalHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"birthday_response": modal.BirthdayResponse,
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		switch i.Type {

		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case discordgo.InteractionModalSubmit:
			if h, ok := modalHandlers[i.ModalSubmitData().CustomID]; ok {
				h(s, i)
			}

		case discordgo.InteractionMessageComponent:
			customId := i.MessageComponentData().CustomID

			split := strings.Split(customId, "_")

			if len(split) > 1 && split[1] == "id" {
				trigU := i.Member.User.ID
				if trigU != split[2] {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You're not allowed to use this Command!",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					return
				} else {
					if h, ok := componentsHandlers[split[0]]; ok {
						h(s, i)
					}
					break
				}
			}

			if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}

		}

	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	if InitRemoveAppCommands {
		globalCommands, err := s.ApplicationCommands(s.State.User.ID, "")
		if err != nil {
			panic(err)
		}
		// remove global commands
		for i := 0; i < len(globalCommands); i++ {
			log.Printf("Removing %s from Global Commands", globalCommands[i].Name)
			s.ApplicationCommandDelete(s.State.User.ID, "", globalCommands[i].ID)
		}

		guildCommands, err := s.ApplicationCommands(s.State.User.ID, os.Getenv("DISCORD_GUILD"))
		if err != nil {
			panic(err)
		}
		// remove guild commands
		for i := 0; i < len(guildCommands); i++ {
			log.Printf("Removing %s from Guild Commands", guildCommands[i].Name)
			s.ApplicationCommandDelete(s.State.User.ID, os.Getenv("DISCORD_GUILD"), guildCommands[i].ID)
		}

		os.Exit(0)
	} else {

		log.Println("Adding commands..")
		registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
		for i, v := range commands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			log.Printf("Created Command with Name \"%s\"", v.Name)
			registeredCommands[i] = cmd
		}
		log.Printf("Added %v Commands!", len(registeredCommands))

		defer s.Close()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		log.Println("Press Ctrl+C to exit")
		<-stop

		if RemoveCommands {
			log.Println("Removing commands...")
			// // We need to fetch the commands, since deleting requires the command ID.
			// // We are doing this from the returned commands on line 375, because using
			// // this will delete all the commands, which might not be desirable, so we
			// // are deleting only the commands that we added.
			// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
			// if err != nil {
			// 	log.Fatalf("Could not fetch registered commands: %v", err)
			// }

			for _, v := range registeredCommands {
				err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
				if err != nil {
					log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
				}
			}
		}

		log.Println("Gracefully shutting down.")
	}
}
