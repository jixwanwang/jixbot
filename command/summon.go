package command

import (
	"fmt"
	"strings"
	"time"
)

type summon struct {
	cp *CommandPool

	mentions []time.Time
	lastSent time.Time
}

func (T *summon) Init() {
	T.lastSent = time.Now().Add(-10 * time.Minute)
}

func (T *summon) ID() string {
	return "summon"
}

func (T *summon) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	index := strings.Index(strings.ToLower(message), "jix")
	indexbot := strings.Index(strings.ToLower(message), "jixbot")
	_, ok := T.cp.channel.InChannel("jix_bot")
	if index >= 0 && indexbot < 0 && !ok {
		mentions := []time.Time{}
		for _, t := range T.mentions {
			if time.Since(t) < 2*time.Minute {
				mentions = append(mentions, t)
			}
		}

		// Send text
		if len(mentions) >= 2 && time.Since(T.lastSent) > 10*time.Minute {
			T.lastSent = time.Now()
			T.cp.texter.SendText(fmt.Sprintf("[%s] %s: %s", T.cp.channel.GetChannelName(), username, message))
			T.cp.Say(fmt.Sprintf("Jix has been summoned! PogChamp"))
		}
	}
}
