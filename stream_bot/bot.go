package stream_bot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/command"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/pastebin"
)

const (
	creator = "jixwanwang"
)

type Bot struct {
	username string
	oath     string
	client   *irc.Client
	commands *command.CommandPool
	channel  *channel.Channel
	db       *sql.DB
	texter   messaging.Texter
	pasteBin pastebin.Client

	groupclient *irc.Client
	groupchat   string

	shutdown chan int
}

func New(channelName, username, oath, groupchat string, texter messaging.Texter, pb pastebin.Client, db *sql.DB) (*Bot, error) {
	bot := &Bot{
		username:  username,
		oath:      oath,
		shutdown:  make(chan int),
		channel:   channel.New(channelName, db),
		db:        db,
		groupchat: groupchat,
		texter:    texter,
		pasteBin:  pb,
	}

	log.Printf("starting up for %v", channelName)
	bot.startup()

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

func (B *Bot) startup() {
	chatServer := "irc.chat.twitch.tv:80"

	// Retrieve servers for the channel
	resp, err := http.Get("http://tmi.twitch.tv/servers?channel=" + B.channel.GetChannelName())
	if err == nil {
		defer resp.Body.Close()
		var m map[string]interface{}
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&m)
		if err == nil {
			servers, ok := m["servers"].([]interface{})
			if ok {
				chatServer = fmt.Sprintf("%v", servers[rand.Intn(len(servers))])
			}
		}
	}
	log.Printf("chat server for %s: %s", B.channel.GetChannelName(), chatServer)

	B.client, _ = irc.New(chatServer, 10)
	B.groupclient, _ = irc.New("192.16.64.212:443", 10)
	B.reloadClients()
	B.commands = command.NewCommandPool(B.channel, B.client, B.groupclient, B.texter, B.pasteBin, B.db)
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

	B.client.Send(fmt.Sprintf("JOIN #%s", B.username))

	err = B.groupclient.Reload()
	if err != nil {
		log.Printf("%s", err.Error())
	}
	B.groupclient.Send(fmt.Sprintf("PASS %s", B.oath))
	B.groupclient.Send(fmt.Sprintf("NICK %s", B.username))
	B.groupclient.Send("CAP REQ :twitch.tv/commands")
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
					} else {
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
