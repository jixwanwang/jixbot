package stream_bot

import (
	"log"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/command"
	"github.com/jixwanwang/jixbot/config"
	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/pastebin"
)

type Bot struct {
	chat     *irc.TwitchChatClient
	commands *command.CommandPool
	channel  *channel.Channel
	db       db.DB
	texter   messaging.Texter
	pasteBin pastebin.Client

	shutdown chan int
}

func NewBot(channelName string, chat *irc.TwitchChatClient, texter messaging.Texter, pb pastebin.Client, db db.DB) *Bot {
	bot := &Bot{
		chat:     chat,
		shutdown: make(chan int),
		channel:  channel.New(channelName, db),
		db:       db,
		texter:   texter,
		pasteBin: pb,
	}

	bot.reloadClients()
	bot.commands = command.NewCommandPool(bot.channel, chat, texter, pb, db)

	return bot
}

func (B *Bot) GetTextCommands() []db.TextCommand {
	return B.commands.GetTextCommands()
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
func (B *Bot) SetProperty(k, v string) {
	B.channel.SetProperty(k, v)
}
func (B *Bot) GetProperties() map[string]interface{} {
	return B.channel.GetProperties()
}

func (B *Bot) reloadClients() {
	err := B.chat.Reload()
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
	reads, groupreads := B.chat.ReadLoops()

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
				} else if username == "twitchnotify" {
					B.processMessage(username, msg)
				} else {
					B.processMessage(username, msg)
				}
			// Sub notification
			case "USERNOTICE":
				username := fromToUsername(e.From)
				log.Printf("%v", username, e.Tags)
				if username == "tmi.twitch.tv" {
					B.processMessage(username, strings.Replace(e.Tags["system-msg"], `\s`, " ", -1))
					log.Printf("%v", strings.Replace(e.Tags["system-msg"], `\s`, " ", -1))
				}
			default: //ignore
			}
		// Whispers
		case e := <-groupreads:
			if e.Err != nil {
				log.Printf("Error %s, reloading group irc client", e.Err.Error())
				B.reloadClients()
				continue
			}

			switch e.Kind {
			case "PRIVMSG": // Message
				username := fromToUsername(e.From)
				msg := strings.TrimPrefix(e.Message, "#"+config.GroupChat+" :")
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
