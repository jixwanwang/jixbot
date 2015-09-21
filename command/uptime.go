package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type uptime struct {
	cp        *CommandPool
	lastCheck time.Time
	upComm    *subCommand
}

func (T *uptime) Init() {
	T.upComm = &subCommand{
		command:    "!uptime",
		numArgs:    0,
		cooldown:   15 * time.Second,
		lastCalled: time.Now().Add(-15 * time.Second),
		clearance:  channel.VIEWER,
	}
}

func (T *uptime) ID() string {
	return "uptime"
}

func (T *uptime) Response(username, message string) string {
	message = strings.TrimSpace(strings.ToLower(message))
	clearance := T.cp.channel.GetLevel(username)

	_, err := T.upComm.parse(message, clearance)
	if err == nil {
		if !T.cp.channel.Broadcaster.Online {
			return fmt.Sprintf("%s isn't online.", T.cp.channel.GetChannelName())
		}
		uptime := time.Now().UTC().Sub(T.cp.channel.Broadcaster.OnlineSince)
		minutes := int(uptime.Minutes())
		return fmt.Sprintf("%s hours, %s minutes", strconv.Itoa(minutes/60), strconv.Itoa(minutes%60))
	}

	return ""
}

func (T *uptime) WhisperOnly() bool {
	return false
}

func (T *uptime) String() string {
	return "Uptime Command"
}
