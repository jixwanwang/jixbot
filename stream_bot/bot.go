package stream_bot

import (
	"database/sql"
	"log"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/command"
	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/pastebin"
	"github.com/jixwanwang/jixbot/twitch_api"
)

const (
	creator = "jixwanwang"
)

type Bot struct {
	username    string
	groupchat   string
	client      *irc.Client
	groupclient *irc.Client
	commands    *command.CommandPool
	channel     *channel.Channel
	db          db.DB
	texter      messaging.Texter
	pasteBin    pastebin.Client

	shutdown chan int
}

func New(channelName, username, oauth, groupchat string, texter messaging.Texter, pb pastebin.Client, sqlDB *sql.DB) (*Bot, error) {
	dbInterface := db.NewDB(sqlDB)

	bot := &Bot{
		username:  username,
		groupchat: groupchat,
		shutdown:  make(chan int),
		channel:   channel.New(channelName, dbInterface),
		db:        dbInterface,
		texter:    texter,
		pasteBin:  pb,
	}

	log.Printf("starting up for %v", channelName)

	chatServer := twitch_api.GetIRCServer(channelName, "irc.chat.twitch.tv:80")
	log.Printf("chat server for %s: %s", channelName, chatServer)

	groupServer := twitch_api.GetIRCCluster("irc.chat.twitch.tv:80")
	log.Printf("group chat server for %s: %s", channelName, groupServer)

	bot.client, _ = irc.New(chatServer, channelName, oauth, username, 10)
	log.Printf("Connected to chat irc")
	bot.groupclient, _ = irc.New(groupServer, groupchat, oauth, username, 10)
	log.Printf("connected to group irc")

	bot.reloadClients()
	bot.commands = command.NewCommandPool(bot.channel, bot.client, bot.groupclient, texter, pb, dbInterface)

	return bot, nil
}

func (B *Bot) AddActiveCommand(c string) {
	B.commands.ActivateCommand(c)
}
func (B *Bot) GetActiveCommands() []string {
	return B.commands.GetActiveCommands()
}
func (B *Bot) DeleteCommand(c string) {
	B.commands.DeleteCommand(c)
}

func (B *Bot) AddEmote(e string) {
	B.channel.AddEmote(e)
}
func (B *Bot) GetEmotes() []string {
	return B.channel.Emotes
}
func (B *Bot) DeleteEmote(e string) {
	B.channel.DeleteEmote(e)
}

func (B *Bot) SetProperty(k, v string) {
	B.channel.SetProperty(k, v)
}
func (B *Bot) GetProperties() map[string]interface{} {
	return B.channel.GetProperties()
}

func (B *Bot) reloadClients() {
	err := B.client.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}

	err = B.groupclient.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}
}

func (B *Bot) Shutdown() {
	B.shutdown <- 1
	log.Printf("shutting down for %s", B.channel.GetChannelName())
	B.channel.ViewerList.Close()
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
				}
				isSub, ok := e.Tags["subscriber"]
				if ok && isSub == "1" {
					B.channel.ViewerList.SetSubscriber(username)
				}
				msg := strings.TrimPrefix(e.Message, "#"+B.channel.GetChannelName()+" :")
				if username == "jtv" {
					special := strings.TrimPrefix(e.Message, "jixbot :")
					if strings.HasPrefix(special, "USERCOLOR") {

					} else if strings.HasPrefix(special, "EMOTESET") {

					} else if strings.HasPrefix(special, "SPECIALUSER") {
						// parts := strings.Split(special, " ")
						// log.Printf("NOTICE: %s is a %s", parts[1], parts[2])
					}
				} else if username == "tmi.twitch.tv" {
					if e.Kind == "USERNOTICE" {
						B.processMessage(username, strings.Replace(e.Tags["system-msg"], `\s`, " "))
					}
				} else if username == "twitchnotify" {
					B.processMessage(username, msg)
				} else if msg == e.Message {
					// Not of the channel, must be group chat
					msg = strings.TrimPrefix(e.Message, "#"+B.groupchat+" :")
				} else {
					B.processMessage(username, msg)
				}
			default: //ignore
			}
		// Whispers
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
			case "WHISPER":
				from := fromToUsername(e.From)
				space := strings.Index(e.Message, " :")
				to := e.Message[:space]
				msg := strings.TrimPrefix(e.Message, to+" :")
				B.processWhisper(from, msg)
			default: //ignore
			}
		}
	}
}

func fromToUsername(from string) string {
	exclam := strings.Index(from, "!")
	if exclam < 0 {
		exclam = len(from)
	}
	return strings.ToLower(from[1:exclam])
}

func (B *Bot) processWhisper(username, msg string) {
	B.commands.GetResponse(username, msg, true)
}

func (B *Bot) processMessage(username, msg string) {
	B.channel.RecordMessage(username, msg)
	B.commands.GetResponse(username, msg, false)
}
