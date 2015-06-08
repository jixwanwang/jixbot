package stream_bot

import (
	"database/sql"
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
	channel     string
	username    string
	oath        string
	texter      messaging.Texter
	client      *irc.Client
	commands    *command.CommandPool
	viewerlist  *channel.ViewerList
	broadcaster *channel.Broadcaster
	db          *sql.DB

	shutdown chan int
}

func New(channelName, username, oath string, texter messaging.Texter, db *sql.DB) (*Bot, error) {
	bot := &Bot{
		channel:     channelName,
		username:    username,
		oath:        oath,
		texter:      texter,
		shutdown:    make(chan int),
		broadcaster: channel.NewBroadcaster(channelName),
		db:          db,
	}

	bot.startup()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C

			if bot.broadcaster.Online {
				bot.viewerlist.Tick()
			}
		}
	}()

	return bot, nil
}

func (B *Bot) startup() {
	B.viewerlist = channel.NewViewerList(B.channel, B.db)
	B.client, _ = irc.New("irc.twitch.tv:6667", 10)
	B.reloadClient()
	B.commands = command.NewCommandPool(B.viewerlist, B.broadcaster, B.client, B.texter, B.db)
}

func (B *Bot) reloadClient() {
	err := B.client.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}

	B.client.Send(fmt.Sprintf("PASS %s", B.oath))
	B.client.Send(fmt.Sprintf("NICK %s", B.username))
	B.client.Send(fmt.Sprintf("JOIN #%s", B.channel))
	B.client.Send("TWITCHCLIENT 2")
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
				continue
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
				if username == "jtv" {
					special := strings.TrimPrefix(e.Message, "jixbot :")
					if strings.HasPrefix(special, "USERCOLOR") {

					} else if strings.HasPrefix(special, "EMOTESET") {

					} else if strings.HasPrefix(special, "SPECIALUSER") {
						// parts := strings.Split(special, " ")
						// log.Printf("NOTICE: %s is a %s", parts[1], parts[2])
					} else {
						// log.Printf("jtv said: %s", special)
					}
				} else if username == "twitchnotify" {
					log.Printf("TWITCHNOTIFY SAYS: %s", msg)
					B.processMessage(username, msg)
				} else {
					B.processMessage(username, msg)
					log.Printf("%s said: %s", username, msg)
				}
			default: //ignore
				// log.Printf("Unknown: %v", e)
			}
		}
	}
}

func (B *Bot) Shutdown() {
	B.shutdown <- 1
	log.Printf("shutting down for %s", B.channel)
	B.commands.FlushTextCommands()
	B.viewerlist.Close()
}

func fromToUsername(from string) string {
	exclam := strings.Index(from, "!")
	if exclam < 0 {
		exclam = len(from)
	}
	return strings.ToLower(from[1:exclam])
}

func (B *Bot) processMessage(username, msg string) {
	B.viewerlist.RecordMessage(username, msg)
	response := B.commands.GetResponse(username, msg)
	if len(response) > 0 {
		B.client.Say("#"+B.channel, response)
	}
}
