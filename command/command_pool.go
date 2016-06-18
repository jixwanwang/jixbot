package command

import (
	"log"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/pastebin"
)

type CommandPool struct {
	channel  *channel.Channel
	irc      *irc.Client
	ircW     *irc.Client
	texter   messaging.Texter
	pasteBin pastebin.Client
	db       db.DB

	specials       []Command
	enabled        map[string]bool
	commands       []*textCommand
	globalcommands []*textCommand

	modOnly bool
	subOnly bool
}

func NewCommandPool(channel *channel.Channel, irc, ircW *irc.Client, texter messaging.Texter, pasteBin pastebin.Client, db db.DB) *CommandPool {
	cp := &CommandPool{
		channel:  channel,
		irc:      irc,
		ircW:     ircW,
		db:       db,
		texter:   texter,
		pasteBin: pasteBin,
	}

	globals := cp.loadTextCommands(globalChannel)
	channels := cp.loadTextCommands(channel.GetChannelName())

	specials := cp.specialCommands()
	enabled := cp.enabledCommands()

	cp.enabled = enabled
	cp.specials = specials
	for _, c := range cp.specials {
		if _, ok := cp.enabled[c.ID()]; ok {
			c.Init()
		}
	}
	cp.commands = channels
	cp.globalcommands = globals

	return cp
}

func (C *CommandPool) loadTextCommands(channelName string) []*textCommand {
	comms, err := C.db.GetTextCommands(channelName)
	if err != nil {
		comms = []db.TextCommand{}
	}

	commands := []*textCommand{}
	for _, c := range comms {
		command := &textCommand{
			cp:        C,
			clearance: channel.Level(c.Clearance),
			command:   c.Command,
			response:  c.Response,
			cooldown:  c.Cooldown,
		}
		command.Init()
		commands = append(commands, command)
	}

	return commands
}

func (C *CommandPool) enabledCommands() map[string]bool {
	allowed, err := C.db.GetCommands(C.channel.GetChannelName())
	if err != nil {
		log.Printf("Couldn't read commands: %s", err.Error())
		allowed = map[string]bool{}
	}

	// Always enable info command
	allowed["info"] = true

	return allowed
}

func (C *CommandPool) specialCommands() []Command {
	return []Command{
		&info{
			cp: C,
		},
		&addCommand{
			cp: C,
		},
		&deleteCommand{
			cp: C,
		},
		&summon{
			cp: C,
		},
		&money{
			cp: C,
		},
		&brawl{
			cp: C,
		},
		&uptime{
			cp: C,
		},
		&subMessage{
			cp: C,
		},
		&fail{
			cp: C,
		},
		&combo{
			cp: C,
		},
		&emotes{
			cp: C,
		},
		&timeSpent{
			cp: C,
		},
		&conversation{
			cp: C,
		},
		&modonly{
			cp: C,
		},
		&linesTyped{
			cp: C,
		},
		&commandList{
			cp: C,
		},
		&questions{
			cp: C,
		},
	}
}

func (C *CommandPool) GetActiveCommands() []string {
	comms := []string{}
	for _, c := range C.specials {
		if _, ok := C.enabled[c.ID()]; ok {
			comms = append(comms, c.ID())
		}
	}
	return comms
}

func (C *CommandPool) ActivateCommand(command string) {
	exists := false
	for _, c := range C.specials {
		if c.ID() == command {
			exists = true
			c.Init()
			C.enabled[command] = true
			break
		}
	}

	if !exists {
		return
	}

	C.db.AddCommand(C.channel.GetChannelName(), command)
}

func (C *CommandPool) DeleteCommand(command string) {
	exists := false
	for _, c := range C.specials {
		if c.ID() == command {
			exists = true
			C.enabled[command] = false
			break
		}
	}

	if !exists {
		return
	}

	C.db.DeleteCommand(C.channel.GetChannelName(), command)
}

func (C *CommandPool) Say(message string) {
	C.irc.Say("#"+C.channel.GetChannelName(), message)
}

func (C *CommandPool) Whisper(username, message string) {
	log.Printf("Attempting to whisper %s to %s", message, username)
	C.ircW.Whisper(username, message)
}

func (C *CommandPool) GetResponse(username, message string, whisper bool) {
	if C.modOnly && !whisper && C.channel.GetLevel(username) < channel.MOD {
		return
	}

	if C.subOnly && !whisper && C.channel.GetLevel(username) < channel.SUBSCRIBER {
		return
	}

	for _, c := range C.specials {
		if _, ok := C.enabled[c.ID()]; ok {
			c.Response(username, message, whisper)
		}
	}
	for _, c := range C.globalcommands {
		c.Response(username, message, whisper)
	}
	for _, c := range C.commands {
		c.Response(username, message, whisper)
	}
}
