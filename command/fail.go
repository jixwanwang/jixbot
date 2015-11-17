package command

import (
	"fmt"
	"strings"
	"time"
)

type fail struct {
	cp       *CommandPool
	lastUsed time.Time
}

func (T fail) Init() {
	T.lastUsed = time.Now()
}

func (T fail) ID() string {
	return "failfish"
}

func (T fail) Response(username, message string, whisper bool) {
	if time.Since(T.lastUsed).Seconds() < 2 {
		return
	}

	index := strings.Index(strings.ToLower(message), "failfish")
	if index == -1 {
		return
	}

	emote := message[index : index+8]
	if emote != "FailFish" {
		T.lastUsed = time.Now()
		T.cp.Say(fmt.Sprintf("@%s FailFish", username))
	}
}
