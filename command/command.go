package command

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

const (
	commandFilePath = "data/textcommands/"
	configFilePath  = "data/config/"
	globalChannel   = "_global"
)

type Command interface {
	ID() string
	Init()
	Response(username, message string) string
	String() string
}

func NewCommandPool(channel *channel.ViewerList, broadcaster *channel.Broadcaster, irc *irc.Client, texter messaging.Texter) *CommandPool {
	cp := &CommandPool{
		channel:      channel,
		broadcaster:  broadcaster,
		irc:          irc,
		texter:       texter,
		currencyName: "HotCoin",
	}

	globals := loadTextCommands(globalChannel)
	channels := loadTextCommands(channel.GetChannelName())

	specials := cp.specialCommands()
	filtered := filterCommands(channel.GetChannelName(), specials)

	cp.specials = filtered
	for _, c := range cp.specials {
		c.Init()
	}
	cp.commands = channels
	cp.globalcommands = globals

	return cp
}

func loadTextCommands(channel string) []*textCommand {
	commands := []*textCommand{}

	commandsRaw, _ := ioutil.ReadFile(commandFilePath + channel)
	commandList := strings.Split(string(commandsRaw), "\n")

	for _, c := range commandList {
		comm, err := parseTextCommand(c)
		if err != nil {
			continue
		}
		commands = append(commands, &comm)
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
		clearance: perm,
		command:   parts[0],
		response:  parts[2],
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

var BadCommandError = fmt.Errorf("Bad command")

type subCommand struct {
	command    string
	numArgs    int
	cooldown   time.Duration
	lastCalled time.Time
}

func (C *subCommand) parse(message string) ([]string, error) {
	args := []string{}

	// Rate limit
	if C.cooldown.Nanoseconds() > 0 && time.Since(C.lastCalled).Nanoseconds() < C.cooldown.Nanoseconds() {
		return args, BadCommandError
	}

	parts := strings.Split(message, " ")
	if parts[0] != C.command {
		return args, BadCommandError
	}

	if len(parts)-1 < C.numArgs {
		return args, BadCommandError
	}

	prefix := strings.Join(parts[:C.numArgs+1], " ")
	// There will be a trailing space on the prefix if the message has more than enough parts
	if C.numArgs < len(parts)-1 {
		prefix = prefix + " "
	}

	if strings.HasPrefix(message, prefix) {
		args = parts[1 : C.numArgs+1]
		remaining := strings.TrimPrefix(message, prefix)
		if len(remaining) > 0 {
			args = append(args, strings.TrimPrefix(message, prefix))
		}
		C.lastCalled = time.Now()
		return args, nil
	}

	return args, BadCommandError
}
