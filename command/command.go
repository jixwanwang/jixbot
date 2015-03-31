package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/messaging"
)

const (
	commandFilePath = "data/textcommands/"
	configFilePath  = "data/config/"
	globalChannel   = "_global"
)

type Command interface {
	ID() string
	GetClearance() channel.Level
	Response(username, message string) string
	String() string
}

type baseCommand struct {
	clearance channel.Level
}

func NewCommandPool(channel *channel.ViewerList, texter messaging.Texter) *CommandPool {
	cp := &CommandPool{channel: channel}

	globals := loadTextCommands(globalChannel)
	channels := loadTextCommands(channel.GetChannelName())

	specials := cp.specialCommands(texter)
	filtered := filterCommands(channel.GetChannelName(), specials)

	cp.specials = filtered
	cp.commands = channels
	cp.globalcommands = globals

	return cp
}

func loadTextCommands(channel string) []textCommand {
	commands := []textCommand{}

	commandsRaw, _ := ioutil.ReadFile(commandFilePath + channel)
	commandList := strings.Split(string(commandsRaw), "\n")

	for _, c := range commandList {
		comm, err := parseTextCommand(c)
		if err != nil {
			continue
		}
		commands = append(commands, comm)
	}
	return commands
}

func parseTextCommand(line string) (textCommand, error) {
	parts := strings.Split(line, ",")

	if len(parts) < 3 {
		return textCommand{}, fmt.Errorf("Not a valid command")
	}

	var perm channel.Level
	switch parts[1] {
	case "viewer":
		perm = channel.VIEWER
	case "mod":
		perm = channel.MOD
	}

	return textCommand{
		baseCommand: baseCommand{
			clearance: perm,
		},
		command:  parts[0],
		response: parts[2],
	}, nil
}

func filterCommands(channel string, commands []Command) []Command {
	newcommands := []Command{}

	configRaw, _ := ioutil.ReadFile(configFilePath + channel)

	ids := map[string]int{}
	for _, id := range strings.Split(string(configRaw), "\n") {
		ids[id] = 1
	}

	for _, c := range commands {
		if _, ok := ids[c.ID()]; ok {
			newcommands = append(newcommands, c)
		}
	}

	return newcommands
}
