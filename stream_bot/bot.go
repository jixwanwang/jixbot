package stream_bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/command"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

const (
	creator = "jixwanwang"
)

type Bot struct {
	channel    string
	username   string
	oath       string
	texter     messaging.Texter
	client     *irc.Client
	commands   *command.CommandPool
	viewerlist *channel.ViewerList

	shutdown chan int
}

func New(channelName, username, oath string, texter messaging.Texter) (*Bot, error) {
	bot := &Bot{
		channel:  channelName,
		username: username,
		oath:     oath,
		texter:   texter,
		shutdown: make(chan int),
	}

	bot.reload()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C

			bot.viewerlist.Tick()
		}
	}()

	return bot, nil
}

func (B *Bot) reload() {
	B.viewerlist = channel.NewViewerList(B.channel)

	B.commands = command.NewCommandPool(B.viewerlist, B.texter)

	B.reloadClient()
}

func (B *Bot) reloadClient() {
	client, err := irc.New("irc.twitch.tv:6667", 5)
	if err != nil {
		log.Fatalf("Couldn't connect to client")
	}

	client.Send(fmt.Sprintf("PASS %s", B.oath))
	client.Send(fmt.Sprintf("NICK %s", B.username))
	client.Send(fmt.Sprintf("JOIN #%s", B.channel))

	B.client = client
}

func (B *Bot) Start() {
	for {
		select {
		case <-B.shutdown:
			log.Printf("shut down!")
			return
		default:
			e, err := B.client.ReadEvent()
			if err != nil {
				// TODO: flush commands, reload everything
				log.Printf("Error %s, reloading irc client", err.Error())
				B.commands.FlushTextCommands()
				B.reloadClient()
			}

			switch e.Kind {
			case "353": // Add viewers
				colon := strings.Index(e.Message, ":")
				usernames := strings.Split(e.Message[colon+1:], " ")
				B.viewerlist.AddViewers(usernames)
			case "MODE": // Mods
				lastSpace := strings.LastIndex(e.Message, " ")
				username := e.Message[lastSpace+1:]
				plus := e.Message[lastSpace-2 : lastSpace-1]
				log.Printf("%s did %s as a mod", username, plus)
				if plus == "+" {
					B.viewerlist.AddMod(username)
				} else {
					B.viewerlist.RemoveMod(username)
				}
			case "JOIN": // Viewers
				B.viewerlist.AddViewers([]string{fromToUsername(e.From)})
			case "PART": // Leaving
				B.viewerlist.RemoveViewer(fromToUsername(e.From))
			case "PRIVMSG": // Message
				username := fromToUsername(e.From)
				msg := strings.TrimPrefix(e.Message, "#"+B.channel+" :")
				log.Printf("%s said: %s", username, msg)
				B.processMessage(username, msg)
			default: //ignore
			}
		}
	}
}

func (B *Bot) Shutdown() {
	B.shutdown <- 1
	log.Printf("shutting down for %s", B.channel)
	B.commands.FlushTextCommands()
}

func fromToUsername(from string) string {
	exclam := strings.Index(from, "!")
	return from[1:exclam]
}

func (B *Bot) processMessage(username, msg string) {
	B.viewerlist.RecordMessage(username, msg)
	response := B.commands.GetResponse(username, msg)
	if len(response) > 0 {
		B.client.Say("#"+B.channel, response)
	}
}
