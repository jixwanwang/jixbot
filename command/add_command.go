package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type addCommandCommand struct {
	baseCommand
	cp   *CommandPool
	key  string
	perm channel.Level
}

func (T addCommandCommand) ID() string {
	return "add"
}

func (T addCommandCommand) Response(username, message string) string {
	remaining := strings.TrimPrefix(message, T.key)
	if remaining == message {
		return ""
	}

	if remaining[:1] != "!" {
		return ""
	}

	space := strings.Index(remaining, " ")
	if space < 0 {
		return ""
	}

	command := strings.ToLower(remaining[:space])
	text := remaining[space+1:]

	defer T.cp.FlushTextCommands()

	existing := T.cp.hasTextCommand(command)

	if existing < 0 {
		T.cp.commands = append(T.cp.commands, textCommand{
			baseCommand: baseCommand{
				clearance: T.perm,
			},
			command:  command,
			response: text,
		})

		return fmt.Sprintf("@%s Command %s -> %s created!", username, command, text)
	} else {
		T.cp.commands[existing] = textCommand{
			baseCommand: baseCommand{
				clearance: T.perm,
			},
			command:  command,
			response: text,
		}

		return fmt.Sprintf("@%s Command %s -> %s updated", username, command, text)
	}
}

func (T addCommandCommand) GetClearance() channel.Level {
	return T.baseCommand.clearance
}

func (T addCommandCommand) String() string {
	return ""
}
