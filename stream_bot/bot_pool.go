package stream_bot

import (
	"fmt"
	"log"

	"github.com/jixwanwang/jixbot/config"
	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/pastebin"
	"github.com/jixwanwang/jixbot/twitch_api"
)

type BotPool struct {
	bots map[string]*Bot

	texter        messaging.Texter
	pasteBin      pastebin.Client
	db            db.DB
	whisperClient *irc.Client
}

func NewPool() (*BotPool, error) {
	database, err := db.New(config.DBHost, config.DBPort, config.DBName, config.DBUser, config.DBPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %v", err)
	}

	channels, err := database.GetAllChannels()
	if err != nil {
		return nil, fmt.Errorf("failed to get all channels: %v", err)
	}
	log.Printf("Channels being loaded: %v", channels)
	channels = []string{"hotform"}

	texter := messaging.NewTexter(config.TwilioAccount, config.TwilioSecret, config.TwilioNumber, config.JixNumber)
	pasteBin := pastebin.NewClient(config.PastebinKey)

	twitch_api.SetClientID(config.ClientID)

	groupServer := twitch_api.GetIRCCluster("irc.chat.twitch.tv:80")
	whisperClient := irc.New(groupServer, config.GroupChat, config.OauthToken, config.Nickname, 10)
	log.Printf("connected to group irc at: %s", groupServer)

	pool := &BotPool{
		bots:          map[string]*Bot{},
		texter:        texter,
		pasteBin:      pasteBin,
		db:            database,
		whisperClient: whisperClient,
	}

	for _, channel := range channels {
		pool.AddBot(channel)
	}

	return pool, nil
}

func (B *BotPool) AddBot(channel string) {
	log.Printf("loading bot for %s", channel)

	chatServer := twitch_api.GetIRCServer(channel, "irc.chat.twitch.tv:80")

	client := irc.New(chatServer, channel, config.OauthToken, config.Nickname, 10)

	props, err := B.db.GetChannelProperties(channel)

	var fancyClient *irc.Client
	if err == nil && props["fancy_name"] != "" && props["fancy_oauth"] != "" {
		fancyClient = irc.New(chatServer, channel, props["fancy_oauth"], props["fancy_name"], 3)
	} else {
		fancyClient = irc.New(chatServer, channel, config.FancyOauthToken, config.FancyNickname, 3)
	}

	chat := irc.NewTwitchChat(channel, client, fancyClient, B.whisperClient)

	b := NewBot(channel, chat, B.texter, B.pasteBin, B.db)
	go b.Start()

	B.bots[channel] = b
}

func (B *BotPool) GetBot(channel string) *Bot {
	if b, ok := B.bots[channel]; ok {
		return b
	}

	return nil
}

func (B *BotPool) GetChannels() []string {
	channels := []string{}
	for k := range B.bots {
		channels = append(channels, k)
	}

	return channels
}

func (B *BotPool) Shutdown() {
	for _, bot := range B.bots {
		bot.Shutdown()
	}
}
