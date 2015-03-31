package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/messaging"
)

type summonCommand struct {
	baseCommand
	texter  messaging.Texter
	channel *channel.ViewerList
}

func (T summonCommand) ID() string {
	return "summon"
}

func (T summonCommand) Response(username, message string) string {
	index := strings.Index(strings.ToLower(message), "jixwanwang")
	_, ok := T.channel.InChannel("jixwanwang")
	if index >= 0 && !ok {
		T.texter.SendText(fmt.Sprintf("[%s] %s: %s", T.channel, username, message))
		return fmt.Sprintf("Jix has been summoned! PogChamp")
	}
	return ""
}

func (T summonCommand) GetClearance() channel.Level {
	return T.baseCommand.clearance
}

func (T summonCommand) String() string {
	return ""
}
