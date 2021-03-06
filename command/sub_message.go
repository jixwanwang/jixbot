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

	if index := strings.Index(msg, "subscribed for "); index > 0 {
		name := msg[:strings.Index(msg, " ")]
		sub := msg[index:]
		sub = strings.TrimPrefix(sub, "subscribed for ")
		months, _ := strconv.Atoi(sub[:strings.Index(sub, " ")])
		if months == 0 {
			emotes := strings.Join(T.cp.channel.Emotes, " ")
			if T.cp.channel.BotIsSubbed {
				T.cp.Say(fmt.Sprintf("@%s, Thank you for re-subscribing! %s", name, emotes))
			} else {
				T.cp.FancySay(fmt.Sprintf("@%s, Thank you for re-subscribing! %s", name, emotes))
			}
			return
		}

		emote := "PogChamp"
		if len(T.cp.channel.Emotes) > 0 {
			emote = T.cp.channel.Emotes[rand.Intn(len(T.cp.channel.Emotes))]
		}
		emotes := ""
		for i := 0; i < months; i++ {
			emotes = emotes + emote + " "
		}

		if T.cp.channel.BotIsSubbed {
			T.cp.Say(fmt.Sprintf("@%s, Thank you for re-subscribing, for %v months! %s", name, months, emotes))
		} else {
			T.cp.FancySay(fmt.Sprintf("@%s, Thank you for re-subscribing, for %v months! %s", name, months, emotes))
		}
	} else if strings.Index(msg, "gifted a") > 0 {
		name := msg[:strings.Index(msg, " ")]
		targetName := msg[strings.LastIndex(msg, " ")+1:]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		if T.cp.channel.BotIsSubbed {
			T.cp.Say(fmt.Sprintf("@%s, Thank you for gifting a sub, welcome to the %s, %s! %s", name, T.cp.channel.SubName, targetName, emotes))
		} else {
			T.cp.FancySay(fmt.Sprintf("@%s, Thank you for gifting a sub, welcome to the %s, %s! %s", name, T.cp.channel.SubName, targetName, emotes))
		}
	} else if strings.Index(msg, "an anonymous gifter gave ") > 0 {
		leftover := strings.TrimPrefix(msg, "an anonymous gifter gave ")
		targetName := leftover[:strings.Index(leftover, " ")]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		if T.cp.channel.BotIsSubbed {
			T.cp.Say(fmt.Sprintf("Welcome to the %s, %s! %s", T.cp.channel.SubName, targetName, emotes))
		} else {
			T.cp.FancySay(fmt.Sprintf("Welcome to the %s, %s! %s", T.cp.channel.SubName, targetName, emotes))
		}
	} else if strings.Index(msg, "just subscribed") > 0 || strings.Index(msg, "twitch prime") > 0 {
		name := msg[:strings.Index(msg, " ")]
		emotes := strings.Join(T.cp.channel.Emotes, " ")
		if T.cp.channel.BotIsSubbed {
			T.cp.Say(fmt.Sprintf("@%s, Thank you for subscribing, welcome to the %s! %s", name, T.cp.channel.SubName, emotes))
		} else {
			T.cp.FancySay(fmt.Sprintf("@%s, Thank you for subscribing, welcome to the %s! %s", name, T.cp.channel.SubName, emotes))
		}
	}
}
