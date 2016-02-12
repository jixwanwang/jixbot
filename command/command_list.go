package command

import (
	"fmt"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type commandList struct {
	cp *CommandPool

	commands *subCommand
}

func (T *commandList) Init() {
	T.commands = &subCommand{
		command:   "!commands",
		numArgs:   0,
		cooldown:  5 * time.Minute,
		clearance: channel.VIEWER,
	}
}

func (T *commandList) ID() string {
	return "command_list"
}

func (T *commandList) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)
	_, err = T.commands.parse(message, clearance)
	if err == nil {
		commands := "command \t| response"
		for _, c := range T.cp.commands {
			commands = fmt.Sprintf("%s\n%s\t| %s", commands, c.command, c.response)
		}

		paste := T.cp.pasteBin.Paste("Jixbot commands", commands, "1", "10M")
		if len(paste) > 0 {
			T.cp.Say(fmt.Sprintf("Current jixbot commands: %s", paste))
		}
	}
}
