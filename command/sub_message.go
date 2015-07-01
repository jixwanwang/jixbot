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

func (T *subMessage) Response(username, message string) string {
	if username != "twitchnotify" {
		return ""
	}

	msg := strings.ToLower(message)

	if strings.Index(msg, "just subscribed!") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		return fmt.Sprintf("Thank you for subscribing %s, welcome to the %s! %s", sub, T.cp.channel.SubName, emotes)
	} else if strings.Index(msg, "subscribed for ") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emote := T.cp.channel.Emotes[rand.Intn(len(T.cp.channel.Emotes))]
		monthIndex := strings.Index(msg, " months in a row!")
		if monthIndex == -1 {
			return fmt.Sprintf("Thank you for re-subscribing %s for 1 month! %s", sub, emote)
		}

		months, err := strconv.Atoi(msg[monthIndex-1 : monthIndex])
		if err != nil {
			emotes := strings.Join(T.cp.channel.Emotes, " ")
			return fmt.Sprintf("Thank you %s for re-subscribing, for %s months! %s", sub, msg[monthIndex-1:monthIndex], emotes)
		}

		emotes := ""
		for i := 0; i < months; i++ {
			emotes = emotes + emote + " "
		}

		return fmt.Sprintf("Thank you %s for re-subscribing, for %v months! %s", sub, months, emotes)
	}

	return ""
}

func (T *subMessage) WhisperOnly() bool {
	return false
}

func (T *subMessage) String() string {
	return ""
}
