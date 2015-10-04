package stream_bot

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/command"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

const (
	creator = "jixwanwang"
)

type Bot struct {
	username string
	oath     string
	texter   messaging.Texter
	client   *irc.Client
	commands *command.CommandPool
	channel  *channel.Channel
	db       *sql.DB

	groupclient *irc.Client
	groupchat   string

	shutdown chan int
}

func New(channelName, username, oath, groupchat string, texter messaging.Texter, db *sql.DB) (*Bot, error) {
	bot := &Bot{
		username:  username,
		oath:      oath,
		texter:    texter,
		shutdown:  make(chan int),
		channel:   channel.New(channelName, db),
		db:        db,
		groupchat: groupchat,
	}

	bot.startup()

	// ticker := time.NewTicker(1 * time.Minute)
	// go func() {
	// 	for {
	// 		<-ticker.C

	// 		if bot.broadcaster.Online {
	// 			bot.viewerlist.Tick()
	// 		}
	// 	}
	// }()

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

func (B *Bot) GetEmotes() []string {
	return B.channel.Emotes
}

func (B *Bot) AddEmote(e string) {
	B.channel.AddEmote(e)
}

func (B *Bot) SetProperty(k, v string) {
	B.channel.SetProperty(k, v)
}

func (B *Bot) startup() {
	B.client, _ = irc.New("irc.twitch.tv:6667", 10)
	B.groupclient, _ = irc.New("192.16.64.212:443", 10)
	B.reloadClients()
	B.commands = command.NewCommandPool(B.channel, B.client, B.groupclient, B.texter, B.db)
}

func (B *Bot) reloadClients() {
	err := B.client.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}

	B.client.Send(fmt.Sprintf("PASS %s", B.oath))
	B.client.Send(fmt.Sprintf("NICK %s", B.username))
	B.client.Send(fmt.Sprintf("JOIN #%s", B.channel.GetChannelName()))
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
				log.Printf("Error %s, reloading irc client", e.Err.Error())
				B.reloadClients()
				B.channel.ViewerList.Flush()
				continue
			}

			switch e.Kind {
			case "353": // Add viewers
				colon := strings.Index(e.Message, ":")
				usernames := strings.Split(e.Message[colon+1:], " ")
				B.channel.ViewerList.AddViewers(usernames)
			case "MODE": // Mods
				lastSpace := strings.LastIndex(e.Message, " ")
				username := e.Message[lastSpace+1:]
				plus := e.Message[lastSpace-2 : lastSpace-1]
				log.Printf("%s did %s as a mod", username, plus)
				if plus == "+" {
					B.channel.ViewerList.AddMod(username)
				} else {
					B.channel.ViewerList.RemoveMod(username)
				}
			case "JOIN": // Viewers
				B.channel.ViewerList.AddViewers([]string{fromToUsername(e.From)})
			case "PART": // Leaving
				B.channel.ViewerList.RemoveViewer(fromToUsername(e.From))
			case "PRIVMSG": // Message
				username := fromToUsername(e.From)
				isMod, ok := e.Tags["user-type"]
				if ok && isMod == "mod" {
					B.channel.ViewerList.AddMod(username)
					log.Printf("%s did + as a mod", username)
				}
				msg := strings.TrimPrefix(e.Message, "#"+B.channel.GetChannelName()+" :")
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
	log.Printf("shutting down for %s", B.channel.GetChannelName())
	B.channel.ViewerList.Close()
}

func fromToUsername(from string) string {
	exclam := strings.Index(from, "!")
	if exclam < 0 {
		exclam = len(from)
	}
	return strings.ToLower(from[1:exclam])
}

func (B *Bot) processWhisper(username, msg string) {
	B.commands.GetResponse(username, msg)
}

func (B *Bot) processMessage(username, msg string) {
	B.channel.RecordMessage(username, msg)
	B.commands.GetResponse(username, msg)
}
