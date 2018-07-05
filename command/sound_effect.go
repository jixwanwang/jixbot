package command

import (
	"fmt"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/twitch_api"
)

type soundEffect struct {
	cp *CommandPool

	queue *subCommand
}

func (T *soundEffect) Init() {
	T.queue = &subCommand{
		command:   "!jukebox",
		numArgs:   1,
		cooldown:  5 * time.Second,
		clearance: channel.VIEWER,
	}
}

func (T *soundEffect) ID() string {
	return "sound_effect"
}

func (T *soundEffect) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	viewer, ok := T.cp.channel.InChannel(username)
	if !ok {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	args, err := T.queue.parse(message, clearance)
	if err == nil {
		if viewer.GetMoney() < 5000 {
			T.cp.Say(fmt.Sprintf("You need 5000 %ss to buy a sound effect", T.cp.channel.Currency))
			return
		}
		viewer.AddMoney(-5000)
		twitch_api.QueueSoundEffect(args[0])
	}
}
