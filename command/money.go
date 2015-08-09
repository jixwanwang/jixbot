package command

import (
	"fmt"
	"strings"
)

type money struct {
	cp *CommandPool
}

func (T *money) Init() {

}

func (T *money) ID() string {
	return "money"
}

func (T *money) Response(username, message string) string {
	viewer, ok := T.cp.channel.InChannel(username)
	if strings.TrimSpace(message) == "!cash" && ok {
		return fmt.Sprintf("You have %d %ss", viewer.GetMoney(), T.cp.channel.Currency)
	}

	return ""
}

func (T *money) WhisperOnly() bool {
	return true
}

func (T *money) String() string {
	return ""
}
