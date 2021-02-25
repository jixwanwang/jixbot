package command

import (
	"time"
	"fmt"

	"github.com/jixwanwang/jixbot/channel"
)

type emoteRate struct {
	cp      *CommandPool
	countEmote   *subCommand
	emotesPerMinute *subCommand
	emoteRates map[int64]int
}

func (T *emoteRate) Init() {
	T.countEmote = &subCommand{
		command:   "pogchamppogchamp",
		numArgs:   0,
		cooldown:  1 * time.Millisecond,
		clearance: channel.VIEWER,
	}
	T.emotesPerMinute = &subCommand{
		command: "!epm",
		numArgs: 0,
		cooldown: 10 * time.Second,
		clearance: channel.SUBSCRIBER,
	}
	T.emoteRates = map[int64]int{}
}

func (T *emoteRate) ID() string {
	return "emoteRate"
}

func (T *emoteRate) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.countEmote.parse(message, clearance)
	if err == nil {
		now := time.Now().Unix()
		if count, ok := T.emoteRates[now]; ok {
			T.emoteRates[now] = count + 1
		} else {
			T.emoteRates[now] = 1
		}
		if len(T.emoteRates) > 120 {
			newEmoteRates := map[int64]int{}
			for timestamp, count := range T.emoteRates{
				if now - timestamp <= 60 {
					newEmoteRates[timestamp] = count
				}
			}
			T.emoteRates = newEmoteRates
		}
	}
	
	_, err = T.emotesPerMinute.parse(message, clearance)
	if err == nil {
		now := time.Now().Unix()
		emoteRate := 0
		for i := 0; i < 60; i++ {
			if count, ok := T.emoteRates[now - int64(i)]; ok {
				emoteRate += count
			}
		}
		if emoteRate < 5 {
			T.cp.Say(fmt.Sprintf("Only %d emotes have been sent in the last minute. Stop slacking! SwiftRage", emoteRate)) 
		} else {
			T.cp.Say(fmt.Sprintf("%d emotes have been sent in the last minute!", emoteRate)) 
		}	
	}
}
