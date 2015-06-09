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
	broadcaster *channel.Broadcaster
	texter      messaging.Texter
	db          *sql.DB

	specials       []Command
	commands       []*textCommand
	globalcommands []*textCommand
}

func (C *CommandPool) specialCommands() []Command {
	return []Command{
		&addCommandCommand{
			cp: C,
		},
		&deleteCommandCommand{
			cp: C,
		},
		summonCommand{
			cp: C,
		},
		moneyCommand{
			cp: C,
		},
		giveMoneyCommand{
			cp: C,
		},
		&brawlCommand{
			cp: C,
		},
		&uptimeCommand{
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

func (C *CommandPool) GetResponse(username, message string) string {
	for _, c := range C.specials {
		res := c.Response(username, message)
		if len(res) > 0 {
			return res
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
