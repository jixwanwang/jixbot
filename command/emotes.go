package command

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

var adjectives = []string{
	"beautiful",
	"wonderful",
	"amazing",
	"unique",
	"awesome",
	"exciting",
	"fabulous",
	"lit",
	"smashing",
}

type emotes struct {
	cp *CommandPool

	emotes *subCommand
}

func (T *emotes) Init() {
	T.emotes = &subCommand{
		command:   "!emotes",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.VIEWER,
	}
}

func (T *emotes) ID() string {
	return "emotes"
}

func (T *emotes) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.emotes.parse(message, clearance)
	if err != nil || len(T.cp.channel.Emotes) == 0 {
		return
	}

	adjective := adjectives[rand.Intn(len(adjectives))]
	T.cp.Say(fmt.Sprintf("Subscribe to get access to these %s emotes: %s", adjective, strings.Join(T.cp.channel.Emotes, " ")))
}
