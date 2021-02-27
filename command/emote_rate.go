package command

import (
	"time"
	"fmt"
	"strconv"

	"github.com/jixwanwang/jixbot/channel"
)

type emoteRate struct {
	cp      *CommandPool
	countEmote   *subCommand
	emotesPerMinute *subCommand
	emoteRates map[int64]int
}

const MaxDuration = 60

func (T *emoteRate) Init() {
	T.countEmote = &subCommand{
		command:   "pogchamppogchamp",
		numArgs:   0,
		cooldown:  1 * time.Millisecond,
		clearance: channel.VIEWER,
	}
	T.emotesPerMinute = &subCommand{
		command: "!emotecheck",
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
		now := time.Now().Unix() / int64(60)
		if count, ok := T.emoteRates[now]; ok {
			T.emoteRates[now] = count + 1
		} else {
			T.emoteRates[now] = 1
		}
		if len(T.emoteRates) > MaxDuration * 2 {
			newEmoteRates := map[int64]int{}
			for timestamp, count := range T.emoteRates{
				if now - timestamp <= MaxDuration {
					newEmoteRates[timestamp] = count
				}
			}
			T.emoteRates = newEmoteRates
		}
	}
	
	args, err := T.emotesPerMinute.parse(message, clearance)
	if err == nil {
		offset := 30
		if len(args) > 0 {
			minutes, ok := strconv.Atoi(args[0])
			if ok == nil && (minutes <= MaxDuration) && (minutes > 1) {
				offset = minutes
			}
		}
		now := time.Now().Unix() / int64(60)
		emoteRate := 0
		for i := 0; i < offset; i++ {
			if count, ok := T.emoteRates[now - int64(i)]; ok {
				emoteRate += count
			}
		}
		durationString := fmt.Sprintf("%d minutes", offset)
		if emoteRate < offset {
			T.cp.Say(fmt.Sprintf("Only %d emotes have been sent in the last %s. Stop slacking! SwiftRage", emoteRate, durationString)) 
		} else {
			T.cp.Say(fmt.Sprintf("%d emotes have been sent in the last %s!", emoteRate, durationString)) 
		}	
	}
}
