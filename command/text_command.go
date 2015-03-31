package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type textCommand struct {
	baseCommand
	command  string
	response string
}

func (T textCommand) ID() string {
	return "text"
}

func (T textCommand) Response(username, message string) string {
	if strings.ToLower(message) == T.command {
		return T.response
	}
	return ""
}

func (B textCommand) GetClearance() channel.Level {
	return B.baseCommand.clearance
}

func (B textCommand) String() string {
	level := "viewer"
	switch B.baseCommand.clearance {
	case channel.VIEWER:
		level = "viewer"
	default:
		level = "mod"
	}
	return fmt.Sprintf("%s,%s,%s", B.command, level, B.response)
}
