package command

import (
	"fmt"
	"strconv"
	"strings"
)

type giveMoneyCommand struct {
	cp           *CommandPool
	currencyName string
}

func (T giveMoneyCommand) Init() {

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

	to_viewer, ok := T.cp.channel.InChannel(to_user)
	if !ok {
		return fmt.Sprintf("@%s that user isn't in the chat.", username)
	}

	viewer, _ := T.cp.channel.InChannel(username)
	if viewer.Money < amount {
		return fmt.Sprintf("@%s you don't have enough %ss", username, T.cp.currencyName)
	}

	viewer.Money = viewer.Money - amount
	to_viewer.Money = to_viewer.Money + amount

	// if to_user == "jixbot" {
	// 	total := T.channel.AddToLottery(username, amount)
	// 	return fmt.Sprintf("@%s you have purchased %d lottery tickets! You have a total of %d.", username, amount, total)
	// }

	return fmt.Sprintf("@%s gave @%s %d %ss!", username, to_user, amount, T.cp.currencyName)
}

func (T giveMoneyCommand) String() string {
	return ""
}
