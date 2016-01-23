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

	if username != "twitchnotify" {
		return
	}

	msg := strings.ToLower(message)

	if strings.Index(msg, "just subscribed!") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		T.cp.Say(fmt.Sprintf("Thank you for subscribing %s, welcome to the %s! %s", sub, T.cp.channel.SubName, emotes))
		viewer, ok := T.cp.channel.InChannel(username)
		if !ok {
			return
		}

		T.cp.Say(fmt.Sprintf("%s has spent %s watching and has typed %d lines of chat.", sub, timeSpentString(viewer.GetTimeSpent()), viewer.GetLinesTyped()))
	} else if strings.Index(msg, "subscribed for ") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		msg = msg[strings.Index(msg, " ")+1:]
		msg = strings.TrimPrefix(msg, "subscribed for ")
		emote := T.cp.channel.Emotes[rand.Intn(len(T.cp.channel.Emotes))]
		monthIndex := strings.Index(msg, " months in a row!")
		if monthIndex == -1 {
			T.cp.Say(fmt.Sprintf("Thank you for re-subscribing %s for 1 month! %s", sub, emote))
			return
		}

		months, err := strconv.Atoi(msg[:monthIndex])
		if err != nil {
			emotes := strings.Join(T.cp.channel.Emotes, " ")
			T.cp.Say(fmt.Sprintf("Thank you %s for re-subscribing, for %d months! %s", sub, months, emotes))
			return
		}

		emotes := ""
		for i := 0; i < months; i++ {
			emotes = emotes + emote + " "
		}

		T.cp.Say(fmt.Sprintf("Thank you %s for re-subscribing, for %v months! %s", sub, months, emotes))

		viewer, ok := T.cp.channel.InChannel(username)
		if !ok {
			return
		}
		T.cp.Say(fmt.Sprintf("%s has spent %s watching and has typed %d lines of chat.", sub, timeSpentString(viewer.GetTimeSpent()), viewer.GetLinesTyped()))
	}
}
