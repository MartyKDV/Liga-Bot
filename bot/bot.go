package bot

import (
	"fmt"
	"liga-bot/config"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var MainInteraction *discordgo.Interaction

type PollData struct {
	ParticipantsYes map[string]discordgo.User
	ParticipantsNo  map[string]discordgo.User
	ResultYes       string
	ResultNo        string
}

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "liga",
		Description: "Ya wanna play league nerds?",
	},
	{
		Name:        "ping",
		Description: "Play a game?",
	},
	{
		Name:        "no",
		Description: "Play a game?",
	},
}
var ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, pollData *PollData){
	"button_accept": handleButtonAccept,
	"button_deny":   handleButtonDeny,
}
var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, pollData *PollData){

	"liga": handleLiga,
}

func handleLiga(s *discordgo.Session, i *discordgo.InteractionCreate, pollData *PollData) {

	// Role - 	&593811313512153090
	// Marty - 	260102898610864129
	// Waleri - 260099615674728451
	// role test2 - <@&1015660476367110234>

	MainInteraction = i.Interaction

	pollData.fillPollResult()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "<@&593811313512153090>",
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeArticle,
					Title: "Liga Now?",
					Description: "Yes: " + pollData.ResultYes + "\n" +
						"No: " + pollData.ResultNo,
					Color: 66773,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Roles: []string{
					"1014561722859786282",
					"593811313512153090",
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Ye",
							Style:    discordgo.PrimaryButton,
							CustomID: "button_accept",
						},
						discordgo.Button{
							Label:    "Nah",
							Style:    discordgo.DangerButton,
							CustomID: "button_deny",
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
func handleButtonAccept(s *discordgo.Session, i *discordgo.InteractionCreate, pollData *PollData) {

	if _, ok := pollData.ParticipantsYes[i.Member.User.Username]; !ok {

		if _, ok = pollData.ParticipantsNo[i.Member.User.Username]; ok {
			delete(pollData.ParticipantsNo, i.Member.User.Username)
		}

		pollData.ParticipantsYes[i.Member.User.Username] = *i.Member.User

		pollData.fillPollResult()

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = s.InteractionResponseEdit(MainInteraction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeArticle,
					Title: "Liga Now?",
					Description: "Yes: " + pollData.ResultYes + "\n" +
						"No: " + pollData.ResultNo,
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Error",
						Description: "You have already voted for this option! Dumbass",
						Color:       16711680,
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func handleButtonDeny(s *discordgo.Session, i *discordgo.InteractionCreate, pollData *PollData) {

	if _, ok := pollData.ParticipantsNo[i.Member.User.Username]; !ok {

		if _, ok = pollData.ParticipantsYes[i.Member.User.Username]; ok {
			delete(pollData.ParticipantsYes, i.Member.User.Username)
		}

		pollData.ParticipantsNo[i.Member.User.Username] = *i.Member.User

		pollData.fillPollResult()

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = s.InteractionResponseEdit(MainInteraction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeArticle,
					Title: "Liga Now?",
					Description: "Yes: " + pollData.ResultYes + "\n" +
						"No: " + pollData.ResultNo,
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					{
						Type:        discordgo.EmbedTypeArticle,
						Title:       "Error",
						Description: "You have already voted for this option! Dumbass",
						Color:       16711680,
					},
				},
			},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func Start() {

	goBot, err := discordgo.New("Bot " + config.Token)
	pollData := PollData{
		ParticipantsYes: make(map[string]discordgo.User),
		ParticipantsNo:  make(map[string]discordgo.User),
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", goBot.State.User.Username, goBot.State.User.Discriminator)
	})
	goBot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i, &pollData)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := ComponentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i, &pollData)
			}
		}
	})

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Adding commands...")

	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := goBot.ApplicationCommandCreate(goBot.State.User.ID, "", v)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		registeredCommands[i] = cmd
		fmt.Println("Added ", cmd.Name)
	}

	fmt.Println("Bot is running...")

	defer goBot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop

	fmt.Println("Removing Commands...")
	for _, v := range registeredCommands {
		err := goBot.ApplicationCommandDelete(goBot.State.User.ID, "", v.ID)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("Shutting Down...")
}

func (p *PollData) fillPollResult() {

	p.ResultYes = ""
	p.ResultNo = ""

	for name := range p.ParticipantsYes {
		if p.ResultYes == "" {
			p.ResultYes += "**" + name + "**"
		} else {
			p.ResultYes += ", **" + name + "**"
		}
	}

	for name := range p.ParticipantsNo {
		if p.ResultNo == "" {
			p.ResultNo += "**" + name + "**"
		} else {
			p.ResultNo += ", **" + name + "**"
		}
	}
}
