package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
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
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)
	if clearance < channel.MOD {
		return
	}

	var comm *textCommand
	comm = &textCommand{
		cp:       T.cp,
		cooldown: defaultCooldown,
	}

	args, err := T.plebComm.parse(message, clearance)
	if err == nil && len(args) > 1 {
		comm.clearance = channel.VIEWER
		comm.command = strings.ToLower(args[0])
	} else {
		args, err = T.modComm.parse(message, clearance)
		if err == nil && len(args) > 1 {
			comm.clearance = channel.MOD
			comm.command = strings.ToLower(args[0])
		}
	}

	if args == nil || len(args) <= 1 {
		return
	}

	response := args[1]
	firstArg := strings.Split(response, " ")[0]
	if strings.Index(firstArg, "-cd=") == 0 {
		cooldown := strings.TrimPrefix(firstArg, "-cd=")
		response = strings.TrimPrefix(response, firstArg+" ")
		cd, err := strconv.Atoi(cooldown)
		if err == nil {
			comm.cooldown = time.Duration(cd) * time.Second
		}
	}

	comm.response = response
	if !comm.ValidateArguments() {
		T.cp.Say(fmt.Sprintf("@%s arguments malformed, must start at $0$ and be consecutive", username))
		return
	}
	comm.Init()

	for i, c := range T.cp.commands {
		if c.command == comm.command {
			tc := db.TextCommand{
				Clearance: int(comm.clearance),
				Command:   comm.command,
				Response:  comm.response,
				Cooldown:  comm.cooldown,
			}
			T.cp.db.UpdateTextCommand(T.cp.channel.GetChannelName(), tc)
			T.cp.commands[i] = comm
			T.cp.commands[i].Init()

			T.cp.Say(fmt.Sprintf("@%s Command %s updated", username, comm.command))
			return
		}
	}

	if T.cp.channel.GetLevel(username) == channel.GOD {
		for i, c := range T.cp.globalcommands {
			if c.command == comm.command {
				tc := db.TextCommand{
					Clearance: int(comm.clearance),
					Command:   comm.command,
					Response:  comm.response,
					Cooldown:  comm.cooldown,
				}
				T.cp.db.UpdateTextCommand("_global", tc)
				T.cp.globalcommands[i] = comm
				T.cp.Say(fmt.Sprintf("@%s Global command %s updated", username, comm.command))
				return
			}
		}
	}

	tc := db.TextCommand{
		Clearance: int(comm.clearance),
		Command:   comm.command,
		Response:  comm.response,
		Cooldown:  comm.cooldown,
	}
	T.cp.db.AddTextCommand(T.cp.channel.GetChannelName(), tc)
	T.cp.commands = append(T.cp.commands, comm)
	T.cp.Say(fmt.Sprintf("@%s Command %s created", username, comm.command))
}
