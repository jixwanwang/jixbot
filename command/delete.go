package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type deleteCommand struct {
	cp *CommandPool
}

func (T *deleteCommand) Init() {

}

func (T *deleteCommand) ID() string {
	return "delete"
}

func (T *deleteCommand) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	if T.cp.channel.GetLevel(username) < channel.MOD {
		return
	}

	remaining := strings.TrimPrefix(message, "!deletecommand ")
	if remaining == message {
		return
	}

	remaining = strings.ToLower(strings.TrimSpace(remaining))

	for i, c := range T.cp.commands {
		if c.command == remaining {
			T.cp.db.Exec("DELETE FROM textcommands WHERE channel=$1 AND command=$2", T.cp.channel.GetChannelName(), remaining)
			T.cp.commands = append(T.cp.commands[:i], T.cp.commands[i+1:]...)
			T.cp.Say(fmt.Sprintf("@%s Command %s deleted", username, remaining))
			return
		}
	}
}
