package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
)

type giveMoneyCommand struct {
	baseCommand
	channel      *channel.ViewerList
	currencyName string
}

func (T giveMoneyCommand) ID() string {
	return "money"
}

func (T giveMoneyCommand) Response(username, message string) string {
	remaining := strings.TrimPrefix(message, "!givecash ")
	if remaining == message {
		return ""
	}

	space := strings.Index(remaining, " ")
	if space < 0 {
		return ""
	}

	to_user := strings.ToLower(remaining[:space])
	amount, _ := strconv.Atoi(remaining[space+1:])

	if amount <= 0 {
		return ""
	}

	to_viewer, ok := T.channel.InChannel(to_user)
	if !ok {
		return fmt.Sprintf("@%s that user isn't in the chat.", username)
	}

	viewer, _ := T.channel.InChannel(username)
	if viewer.Money < amount {
		return fmt.Sprintf("@%s you don't have enough %ss", username, T.currencyName)
	}

	viewer.Money = viewer.Money - amount
	to_viewer.Money = to_viewer.Money + amount

	return fmt.Sprintf("@%s gave @%s %d %ss!", username, to_user, amount, T.currencyName)
}

func (T giveMoneyCommand) GetClearance() channel.Level {
	return T.baseCommand.clearance
}

func (T giveMoneyCommand) String() string {
	return ""
}
