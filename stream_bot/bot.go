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

	groupclient *irc.Client
	groupchat   string

	shutdown chan int
}

func New(channelName, username, oath, groupchat string, texter messaging.Texter, db *sql.DB) (*Bot, error) {
	bot := &Bot{
		channel:     channelName,
		username:    username,
		oath:        oath,
		texter:      texter,
		shutdown:    make(chan int),
		broadcaster: channel.NewBroadcaster(channelName),
		db:          db,
		groupchat:   groupchat,
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

func (B *Bot) GetActiveCommands() []string {
	return B.commands.GetActiveCommands()
}

func (B *Bot) AddActiveCommand(c string) {
	B.commands.ActivateCommand(c)
}

func (B *Bot) DeleteCommand(c string) {
	B.commands.DeleteCommand(c)
}

func (B *Bot) startup() {
	B.viewerlist = channel.NewViewerList(B.channel, B.db)
	B.client, _ = irc.New("irc.twitch.tv:6667", 10)
	B.groupclient, _ = irc.New("192.16.64.212:443", 10)
	B.reloadClients()
	B.commands = command.NewCommandPool(B.viewerlist, B.broadcaster, B.client, B.groupclient, B.texter, B.db)
}

func (B *Bot) reloadClients() {
	err := B.client.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}

	B.client.Send(fmt.Sprintf("PASS %s", B.oath))
	B.client.Send(fmt.Sprintf("NICK %s", B.username))
	B.client.Send(fmt.Sprintf("JOIN #%s", B.channel))
	B.client.Send("CAP REQ :twitch.tv/membership")
	B.client.Send("CAP REQ :twitch.tv/tags")

	err = B.groupclient.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	B.groupclient.Send(fmt.Sprintf("PASS %s", B.oath))
	B.groupclient.Send(fmt.Sprintf("NICK %s", B.username))
	B.groupclient.Send("CAP REQ :twitch.tv/commands")
}

func (B *Bot) Start() {
	reads := B.client.ReadLoop()
	groupreads := B.groupclient.ReadLoop()

	for {
		select {
		case <-B.shutdown:
			log.Printf("shut down!")
			return
		case e := <-reads:
			if e.Err != nil {
				// TODO: flush commands, reload everything
				log.Printf("Error %s, reloading irc client", e.Err.Error())
				B.reloadClients()
				B.viewerlist.Flush()
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
				isMod, ok := e.Tags["user-type"]
				if ok && isMod == "mod" {
					B.viewerlist.AddMod(username)
					log.Printf("%s did + as a mod", username)
				}
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
				} else if msg == e.Message {
					// Not of the channel, must be group chat
					msg = strings.TrimPrefix(e.Message, "#"+B.groupchat+" :")
					log.Printf("%s said in group chat: %s", username, msg)
				} else {
					B.processMessage(username, msg)
					log.Printf("%s said: %s", username, msg)
				}

			default: //ignore
				log.Printf("Unknown: %v", e)
			}
		case e := <-groupreads:
			if e.Err != nil {
				// TODO: flush commands, reload everything
				log.Printf("Error %s, reloading irc client", e.Err.Error())
				B.reloadClients()
				continue
			}

			switch e.Kind {
			case "PRIVMSG": // Message
				username := fromToUsername(e.From)
				msg := strings.TrimPrefix(e.Message, "#"+B.groupchat+" :")
				B.processMessage(username, msg)
				log.Printf("%s said in group chat: %s", username, msg)
			case "WHISPER":
				from := fromToUsername(e.From)
				space := strings.Index(e.Message, " :")
				to := e.Message[:space]
				msg := strings.TrimPrefix(e.Message, to+" :")
				log.Printf("%s whispered to %s: %s", from, to, msg)
				B.processWhisper(from, msg)
			default: //ignore
				log.Printf("Don't care about this group chat message: %v", e)
			}
		}
	}
}

func (B *Bot) Shutdown() {
	B.shutdown <- 1
	log.Printf("shutting down for %s", B.channel)
	B.viewerlist.Close()
}

func fromToUsername(from string) string {
	exclam := strings.Index(from, "!")
	if exclam < 0 {
		exclam = len(from)
	}
	return strings.ToLower(from[1:exclam])
}

func (B *Bot) processWhisper(username, msg string) {
	B.viewerlist.RecordMessage(username, msg)
	response := B.commands.GetWhisperResponse(username, msg)
	if len(response) > 0 {
		B.groupclient.Whisper(B.channel, username, response)
	}
}

func (B *Bot) processMessage(username, msg string) {
	B.viewerlist.RecordMessage(username, msg)
	response := B.commands.GetResponse(username, msg)
	if len(response) > 0 {
		B.client.Say("#"+B.channel, response)
	}
}
