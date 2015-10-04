package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type textCommand struct {
	cp        *CommandPool
	clearance channel.Level
	command   string
	response  string
}

func (T textCommand) Init() {

}

func (T textCommand) ID() string {
	return "text"
}

func (T textCommand) Response(username, message string) {
	if strings.ToLower(message) == T.command {
		T.cp.Say(T.response)
	}
}

func (B textCommand) String() string {
	level := "viewer"
	switch B.clearance {
	case channel.VIEWER:
		level = "viewer"
	default:
		level = "mod"
	}
	return fmt.Sprintf("%s,%s,%s", B.command, level, B.response)
}
