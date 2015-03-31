package command

import (
	"fmt"
	"io/ioutil"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/messaging"
)

type CommandPool struct {
	channel        *channel.ViewerList
	specials       []Command
	commands       []textCommand
	globalcommands []textCommand
}

func (C *CommandPool) specialCommands(texter messaging.Texter) []Command {
	return []Command{
		addCommandCommand{
			baseCommand: baseCommand{
				clearance: channel.MOD,
			},
			cp:   C,
			key:  "!addcommand ",
			perm: channel.VIEWER,
		},
		addCommandCommand{
			baseCommand: baseCommand{
				clearance: channel.MOD,
			},
			cp:   C,
			key:  "!addmodcommand ",
			perm: channel.MOD,
		},
		deleteCommandCommand{
			baseCommand: baseCommand{
				clearance: channel.MOD,
			},
			cp: C,
		},
		summonCommand{
			baseCommand: baseCommand{
				clearance: channel.VIEWER,
			},
			texter:  texter,
			channel: C.channel,
		},
		moneyCommand{
			baseCommand: baseCommand{
				clearance: channel.VIEWER,
			},
			channel:      C.channel,
			currencyName: "HotCoin",
		},
		giveMoneyCommand{
			baseCommand: baseCommand{
				clearance: channel.VIEWER,
			},
			channel:      C.channel,
			currencyName: "HotCoin",
		},
	}
}

func (C *CommandPool) FlushTextCommands() {
	data := ""

	for _, c := range C.commands {
		if c.ID() == "text" {
			data = fmt.Sprintf("%s%v\n", data, c)
		}
	}
	ioutil.WriteFile(commandFilePath+C.channel.GetChannelName(), []byte(data), 0666)
}

func (C *CommandPool) GetResponse(username, message string) string {
	clearance := C.channel.GetLevel(username)
	for _, c := range C.specials {
		if clearance >= c.GetClearance() {
			res := c.Response(username, message)
			if len(res) > 0 {
				return res
			}
		}
	}
	for _, c := range C.globalcommands {
		if clearance >= c.GetClearance() {
			res := c.Response(username, message)
			if len(res) > 0 {
				return res
			}
		}
	}
	for _, c := range C.commands {
		if clearance >= c.GetClearance() {
			res := c.Response(username, message)
			if len(res) > 0 {
				return res
			}
		}
	}

	return ""
}

func (C *CommandPool) hasTextCommand(comm string) int {
	for i, c := range C.commands {
		if c.command == comm {
			return i
		}
	}
	return -1
}
