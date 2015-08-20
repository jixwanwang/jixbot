package command

import (
	"database/sql"

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
	commands       []*textCommand
	globalcommands []*textCommand
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
	}
}

func (C *CommandPool) GetActiveCommands() []string {
	comms := []string{}
	for _, c := range C.specials {
		comms = append(comms, c.ID())
	}
	return comms
}

func (C *CommandPool) GetResponse(username, message string) string {
	for _, c := range C.specials {
		if !c.WhisperOnly() {
			res := c.Response(username, message)
			if len(res) > 0 {
				return res
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
		if c.WhisperOnly() {
			res := c.Response(username, message)
			if len(res) > 0 {
				return res
			}
		}
	}
	return ""
}
