package command

import (
	"fmt"
	"strings"
	"time"
)

type summon struct {
	cp       *CommandPool
	lastSent time.Time
}

func (T *summon) Init() {
	T.lastSent = time.Now().Add(-10 * time.Minute)
}

func (T *summon) ID() string {
	return "summon"
}

func (T *summon) Response(username, message string) string {
	index := strings.Index(strings.ToLower(message), "jix")
	indexbot := strings.Index(strings.ToLower(message), "jixbot")
	_, ok := T.cp.channel.InChannel("jixwanwang")
	if index >= 0 && indexbot < 0 && !ok && time.Since(T.lastSent).Seconds() > 600 {
		T.lastSent = time.Now()
		T.cp.texter.SendText(fmt.Sprintf("[%s] %s: %s", T.cp.channel.GetChannelName(), username, message))
		return fmt.Sprintf("Jix has been summoned! PogChamp")
	}
	return ""
}

func (T *summon) WhisperOnly() bool {
	return false
}

func (T *summon) String() string {
	return ""
}
