package command

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const emotesFilePath = "data/emotes/"

type subMessage struct {
	cp *CommandPool
}

func (T *subMessage) Init() {
}

func (T *subMessage) ID() string {
	return "submessage"
}

func (T *subMessage) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	if username != "tmi.twitch.tv" {
		return
	}

	msg := strings.ToLower(message)

	if strings.Index(msg, "just subscribed!") > 0 || strings.Index(msg, "twitch prime") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		T.cp.FancySay(fmt.Sprintf("@%s, Thank you for subscribing, welcome to the %s! %s", sub, T.cp.channel.SubName, emotes))
	} else if strings.Index(msg, "subscribed for ") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		msg = msg[strings.Index(msg, " ")+1:]
		msg = strings.TrimPrefix(msg, "subscribed for ")
		emote := T.cp.channel.Emotes[rand.Intn(len(T.cp.channel.Emotes))]
		monthIndex := strings.Index(msg, " months in a row!")
		if monthIndex == -1 {
			T.cp.FancySay(fmt.Sprintf("@%s, Thank you for re-subscribing for 1 month! %s", sub, emote))
			return
		}

		months, err := strconv.Atoi(msg[:monthIndex])
		if err != nil {
			emotes := strings.Join(T.cp.channel.Emotes, " ")
			T.cp.FancySay(fmt.Sprintf("@%s, Thank you for re-subscribing! %s", sub, emotes))
			return
		}

		emotes := ""
		for i := 0; i < months; i++ {
			emotes = emotes + emote + " "
		}

		T.cp.FancySay(fmt.Sprintf("@%s, Thank you for re-subscribing, for %v months! %s", sub, months, emotes))
	}
}
