package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type deleteCommandCommand struct {
	cp *CommandPool
}

func (T *deleteCommandCommand) Init() {

}

func (T *deleteCommandCommand) ID() string {
	return "delete"
}

func (T *deleteCommandCommand) Response(username, message string) string {
	if T.cp.channel.GetLevel(username) < channel.MOD {
		return ""
	}

	remaining := strings.TrimPrefix(message, "!deletecommand ")
	if remaining == message {
		return ""
	}

	remaining = strings.ToLower(strings.TrimSpace(remaining))

	if remaining[:1] != "!" {
		return ""
	}

	for i, c := range T.cp.commands {
		if c.command == remaining {
			T.cp.db.Exec("DELETE FROM textcommands WHERE channel=$1 AND command=$2", T.cp.channel.GetChannelName(), remaining)
			T.cp.commands = append(T.cp.commands[:i], T.cp.commands[i+1:]...)
			return fmt.Sprintf("@%s Command %s deleted", username, remaining)
		}
	}

	return ""
	// existing := T.cp.hasTextCommand(remaining)

	// if existing < 0 {
	// 	return ""
	// }

	// T.cp.commands = append(T.cp.commands[:existing], T.cp.commands[existing+1:]...)

	// T.cp.FlushTextCommands()

	// return fmt.Sprintf("@%s Command %s deleted", username, remaining)
}

func (T *deleteCommandCommand) String() string {
	return ""
}
