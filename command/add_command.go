package command

import (
	"fmt"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type addCommandCommand struct {
	cp       *CommandPool
	plebComm *subCommand
	modComm  *subCommand
}

func (T *addCommandCommand) Init() {
	T.plebComm = &subCommand{
		command:    "!addcommand",
		numArgs:    1,
		cooldown:   1 * time.Second,
		lastCalled: time.Now().Add(-1 * time.Second),
	}
	T.modComm = &subCommand{
		command:    "!addmodcommand",
		numArgs:    1,
		cooldown:   1 * time.Second,
		lastCalled: time.Now().Add(-1 * time.Second),
	}
}

func (T *addCommandCommand) ID() string {
	return "add"
}

func (T *addCommandCommand) Response(username, message string) string {
	if T.cp.channel.GetLevel(username) < channel.MOD {
		return ""
	}

	var comm *textCommand

	args, err := T.plebComm.parse(message)
	if err == nil && args[0][:1] == "!" {
		comm = &textCommand{
			clearance: channel.VIEWER,
			command:   args[0],
			response:  args[1],
		}
	}

	args, err = T.modComm.parse(message)
	if err == nil && args[0][:1] == "!" {
		comm = &textCommand{
			clearance: channel.MOD,
			command:   args[0],
			response:  args[1],
		}
	}

	if comm == nil {
		return ""
	}

	existing := T.cp.hasTextCommand(comm.command)

	if existing < 0 {
		T.cp.commands = append(T.cp.commands, comm)
		return fmt.Sprintf("@%s Command %s -> %s created", username, comm.command, comm.response)
	} else {
		T.cp.commands[existing] = comm
		return fmt.Sprintf("@%s Command %s -> %s updated", username, comm.command, comm.response)
	}

	return ""
}

func (T *addCommandCommand) String() string {
	return ""
}
