package command

import (
	"fmt"
	"strings"
)

type moneyCommand struct {
	cp *CommandPool
}

func (T *moneyCommand) Init() {

}

func (T *moneyCommand) ID() string {
	return "money"
}

func (T *moneyCommand) Response(username, message string) string {
	viewer, ok := T.cp.channel.InChannel(username)
	if strings.TrimSpace(message) == "!cash" && ok {
		return fmt.Sprintf("You have %d %ss", viewer.GetMoney(), T.cp.channel.Currency)
	}

	return ""
}

func (T *moneyCommand) WhisperOnly() bool {
	return true
}

func (T *moneyCommand) String() string {
	return ""
}
