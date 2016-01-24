package command

import (
	"database/sql"
	"log"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

type CommandPool struct {
	channel *channel.Channel
	irc     *irc.Client
	ircW    *irc.Client
	texter  messaging.Texter
	db      *sql.DB

	specials       []Command
	enabled        map[string]bool
	commands       []*textCommand
	globalcommands []*textCommand

	modOnly bool
	subOnly bool
}

func NewCommandPool(channel *channel.Channel, irc, ircW *irc.Client, texter messaging.Texter, db *sql.DB) *CommandPool {
	cp := &CommandPool{
		channel: channel,
		irc:     irc,
		ircW:    ircW,
		db:      db,
		texter:  texter,
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
	commands := []*textCommand{}

	rows, err := C.db.Query("SELECT command, message, clearance FROM textcommands WHERE channel=$1", channelName)
	if err != nil {
		log.Printf("Couldn't read text commands")
		return commands
	}
	for rows.Next() {
		var comm, message string
		var clearance int
		rows.Scan(&comm, &message, &clearance)

		command := &textCommand{
			cp:        C,
			clearance: channel.Level(clearance),
			command:   comm,
			response:  message,
		}
		command.Init()
		commands = append(commands, command)
	}
	rows.Close()

	return commands
}

func (C *CommandPool) enabledCommands() map[string]bool {
	allowed := map[string]bool{}

	// Always enable info command
	allowed["info"] = true

	rows, err := C.db.Query("SELECT command FROM commands WHERE channel=$1", C.channel.GetChannelName())
	if err != nil {
		log.Printf("Couldn't read commands")
		return allowed
	}

	for rows.Next() {
		var comm string
		if err := rows.Scan(&comm); err == nil {
			allowed[comm] = true
		}
	}

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

	C.db.Exec("INSERT INTO commands (channel, command) VALUES ($1, $2)", C.channel.GetChannelName(), command)
}

func (C *CommandPool) DeleteCommand(command string) {
	C.db.Exec("DELETE FROM commands WHERE channel=$1 AND command=$2", C.channel.GetChannelName(), command)
}

func (C *CommandPool) Say(message string) {
	C.irc.Say("#"+C.channel.GetChannelName(), message)
}

func (C *CommandPool) Whisper(username, message string) {
	C.ircW.Whisper(C.channel.GetChannelName(), username, message)
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
