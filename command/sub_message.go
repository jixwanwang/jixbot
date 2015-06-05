package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const emotesFilePath = "data/emotes/"

type subMessage struct {
	cp     *CommandPool
	emotes []string
	active bool
}

func (T *subMessage) Init() {
	emotesRaw, _ := ioutil.ReadFile(emotesFilePath + T.cp.channel.GetChannelName())
	T.emotes = strings.Split(string(emotesRaw), "\n")

	// Safeguard against no emotes
	if len(T.emotes) == 0 {
		T.active = false
	} else {
		T.active = true
	}

	log.Printf("%v", T.emotes)
}

func (T *subMessage) ID() string {
	return "submessage"
}

func (T *subMessage) Response(username, message string) string {
	if !T.active {
		return ""
	}

	if username != "twitchnotify" {
		return ""
	}

	msg := strings.ToLower(message)

	if strings.Index(msg, "just subscribed!") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emotes := strings.Join(T.emotes, " ")
		return fmt.Sprintf("Thank you for subscribing %s, welcome to the HotShots! %s", sub, emotes)
	} else if strings.Index(msg, "subscribed for ") > 0 {
		sub := msg[:strings.Index(msg, " ")]
		emote := T.emotes[rand.Intn(len(T.emotes))]
		monthIndex := strings.Index(msg, " months in a row!")
		if monthIndex == -1 {
			return fmt.Sprintf("Thank you for re-subscribing %s for 1 month! %s", sub, emote)
		}

		months, err := strconv.Atoi(msg[monthIndex-1 : monthIndex])
		if err != nil {
			emotes := strings.Join(T.emotes, " ")
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

func (T *subMessage) String() string {
	return ""
}
