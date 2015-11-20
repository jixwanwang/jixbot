package command

import (
	"fmt"
	"log"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type textCommand struct {
	cp        *CommandPool
	clearance channel.Level

	command  string
	response string

	comm *subCommand
}

func (T *textCommand) Init() {
	T.comm = &subCommand{
		command:   T.command,
		numArgs:   0,
		cooldown:  100 * time.Millisecond,
		clearance: T.clearance,
	}
	log.Printf("%v", T.comm)
}

func (T *textCommand) ID() string {
	return "text"
}

func (T *textCommand) Response(username, message string, whisper bool) {
	if whisper {
		return
	}
	clearance := T.cp.channel.GetLevel(username)
	_, err := T.comm.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.response)
	}
}

func (T *textCommand) String() string {
	level := "viewer"
	switch T.clearance {
	case channel.VIEWER:
		level = "viewer"
	default:
		level = "mod"
	}
	return fmt.Sprintf("%s,%s,%s", T.command, level, T.response)
}
