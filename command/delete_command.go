package command

import (
	"fmt"
	"strings"
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
	remaining := strings.TrimPrefix(message, "!deletecommand ")
	if remaining == message {
		return ""
	}

	remaining = strings.TrimSpace(remaining)

	if remaining[:1] != "!" {
		return ""
	}

	existing := T.cp.hasTextCommand(remaining)

	if existing < 0 {
		return ""
	}

	T.cp.commands = append(T.cp.commands[:existing], T.cp.commands[existing+1:]...)

	T.cp.FlushTextCommands()

	return fmt.Sprintf("@%s Command %s deleted", username, remaining)
}

func (T *deleteCommandCommand) String() string {
	return ""
}
