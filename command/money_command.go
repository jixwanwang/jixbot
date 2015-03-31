package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type moneyCommand struct {
	baseCommand
	channel      *channel.ViewerList
	currencyName string
}

func (T moneyCommand) ID() string {
	return "money"
}

func (T moneyCommand) Response(username, message string) string {
	viewer, ok := T.channel.InChannel(username)
	if strings.TrimSpace(message) == "!cash" && ok {
		return fmt.Sprintf("@%s You have %d %ss", username, viewer.Money, T.currencyName)
	}

	return ""
}

func (T moneyCommand) GetClearance() channel.Level {
	return T.baseCommand.clearance
}

func (T moneyCommand) String() string {
	return ""
}
