package command

import (
	"database/sql"
	"log"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

type CommandPool struct {
	channel     *channel.ViewerList
	irc         *irc.Client
	ircW        *irc.Client
	broadcaster *channel.Broadcaster
	texter      messaging.Texter
	db          *sql.DB

	specials       []Command
	enabled        map[string]bool
	commands       []*textCommand
	globalcommands []*textCommand
}

func NewCommandPool(channel *channel.ViewerList, broadcaster *channel.Broadcaster, irc, ircW *irc.Client, texter messaging.Texter, db *sql.DB) *CommandPool {
	cp := &CommandPool{
		channel:     channel,
		broadcaster: broadcaster,
		irc:         irc,
		ircW:        ircW,
		db:          db,
		texter:      texter,
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

		commands = append(commands, &textCommand{
			clearance: channel.Level(clearance),
			command:   comm,
			response:  message,
		})
	}
	rows.Close()

	return commands
}

func (C *CommandPool) enabledCommands() map[string]bool {
	allowed := map[string]bool{}

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
		&slots{
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

func (C *CommandPool) GetResponse(username, message string) string {
	for _, c := range C.specials {
		if _, ok := C.enabled[c.ID()]; ok {
			if !c.WhisperOnly() {
				res := c.Response(username, message)
				if len(res) > 0 {
					return res
				}
			}
		}
	}
	for _, c := range C.globalcommands {
		res := c.Response(username, message)
		if len(res) > 0 {
			return res
		}
	}
	for _, c := range C.commands {
		res := c.Response(username, message)
		if len(res) > 0 {
			return res
		}
	}

	return ""
}

func (C *CommandPool) GetWhisperResponse(username, message string) string {
	for _, c := range C.specials {
		if _, ok := C.enabled[c.ID()]; ok {
			if c.WhisperOnly() {
				res := c.Response(username, message)
				if len(res) > 0 {
					return res
				}
			}
		}
	}
	return ""
}
