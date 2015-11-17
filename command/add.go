package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type addCommand struct {
	cp       *CommandPool
	plebComm *subCommand
	modComm  *subCommand
}

func (T *addCommand) Init() {
	T.plebComm = &subCommand{
		command:    "!addcommand",
		numArgs:    1,
		cooldown:   1 * time.Second,
		lastCalled: time.Now().Add(-1 * time.Second),
		clearance:  channel.MOD,
	}
	T.modComm = &subCommand{
		command:    "!addmodcommand",
		numArgs:    1,
		cooldown:   1 * time.Second,
		lastCalled: time.Now().Add(-1 * time.Second),
		clearance:  channel.MOD,
	}
}

func (T *addCommand) ID() string {
	return "add"
}

func (T *addCommand) Response(username, message string, whisper bool) {
	clearance := T.cp.channel.GetLevel(username)
	if T.cp.channel.GetLevel(username) < channel.MOD {
		return
	}

	var comm *textCommand

	args, err := T.plebComm.parse(message, clearance)
	if err == nil && len(args) > 1 && args[0][:1] == "!" {
		comm = &textCommand{
			clearance: channel.VIEWER,
			command:   strings.ToLower(args[0]),
			response:  args[1],
		}
	}

	args, err = T.modComm.parse(message, clearance)
	if err == nil && len(args) > 1 && args[0][:1] == "!" {
		comm = &textCommand{
			clearance: channel.MOD,
			command:   strings.ToLower(args[0]),
			response:  args[1],
		}
	}

	if comm == nil {
		return
	}

	for i, c := range T.cp.commands {
		if c.command == comm.command {
			T.cp.db.Exec("UPDATE textcommands SET message=$1, clearance=$2 WHERE channel=$3 AND command=$4", comm.response, comm.clearance, T.cp.channel.GetChannelName(), comm.command)
			T.cp.commands[i] = comm
			T.cp.Say(fmt.Sprintf("@%s Command %s -> %s updated", username, comm.command, comm.response))
			return
		}
	}

	if T.cp.channel.GetLevel(username) == channel.GOD {
		for i, c := range T.cp.globalcommands {
			if c.command == comm.command {
				T.cp.db.Exec("UPDATE textcommands SET message=$1, clearance=$2 WHERE channel=$3 AND command=$4", comm.response, comm.clearance, "_global", comm.command)
				T.cp.globalcommands[i] = comm
				T.cp.Say(fmt.Sprintf("@%s Global command %s -> %s updated", username, comm.command, comm.response))
				return
			}
		}
	}

	T.cp.db.Exec("INSERT INTO textcommands (channel, command, message, clearance) VALUES ($1,$2,$3,$4)", T.cp.channel.GetChannelName(), comm.command, comm.response, comm.clearance)
	T.cp.commands = append(T.cp.commands, comm)
	T.cp.Say(fmt.Sprintf("@%s Command %s -> %s created", username, comm.command, comm.response))
}
