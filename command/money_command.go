package command

import (
	"fmt"
	"strings"
)

type moneyCommand struct {
	cp           *CommandPool
	currencyName string
}

func (T moneyCommand) Init() {

}

func (T moneyCommand) ID() string {
	return "money"
}

func (T moneyCommand) Response(username, message string) string {
	viewer, ok := T.cp.channel.InChannel(username)
	if strings.TrimSpace(message) == "!cash" && ok {
		return fmt.Sprintf("@%s You have %d %ss", username, viewer.Money, T.cp.currencyName)
	}

	return ""
}

func (T moneyCommand) String() string {
	return ""
}
