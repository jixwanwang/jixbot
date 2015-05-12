package command

import (
	"fmt"
	"strings"

	"github.com/jixwanwang/jixbot/messaging"
)

type summonCommand struct {
	texter messaging.Texter
	cp     *CommandPool
}

func (T summonCommand) Init() {

}

func (T summonCommand) ID() string {
	return "summon"
}

func (T summonCommand) Response(username, message string) string {
	index := strings.Index(strings.ToLower(message), "jix")
	_, ok := T.cp.channel.InChannel("jixwanwang")
	if index >= 0 && !ok {
		T.texter.SendText(fmt.Sprintf("[%s] %s: %s", T.cp.channel.GetChannelName(), username, message))
		return fmt.Sprintf("Jix has been summoned! PogChamp")
	}
	return ""
}

func (T summonCommand) String() string {
	return ""
}
