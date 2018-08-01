package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/twitch_api"
)

type soundEffect struct {
	cp *CommandPool

	queue *subCommand
	list  *subCommand
}

func (T *soundEffect) Init() {
	T.queue = &subCommand{
		command:   "!jukebox",
		numArgs:   1,
		cooldown:  5 * time.Second,
		clearance: channel.VIEWER,
	}

	T.list = &subCommand{
		command:   "!jukeboxsounds",
		numArgs:   0,
		cooldown:  15 * time.Second,
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
		if viewer.GetMoney() < 2500 {
			T.cp.Say(fmt.Sprintf("You need 2500 %ss to buy a sound effect", T.cp.channel.Currency))
			return
		}
		viewer.AddMoney(-2500)
		twitch_api.QueueSoundEffect(args[0])
	}

	_, err = T.list.parse(message, clearance)
	if err == nil {
		sounds := twitch_api.ListSoundEffects()
		message := "Hey - are those Hot Coins burning up your wallet? For unlimited fun (at 2500 Hot Coins each), type !jukebox and then any of the following phrases to hear your clip! CHOOSE WISELY Jebaited "
		userSounds := []string{}
		for _, sound := range sounds {
			userSounds = append(userSounds, strings.TrimSuffix(sound, ".mp3"))
		}
		T.cp.Say(fmt.Sprintf("%s %s", message, strings.Join(userSounds, ", ")))
		return
	}
}
