package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type uptimeCommand struct {
	cp        *CommandPool
	lastCheck time.Time
	upComm    *subCommand
}

func (T *uptimeCommand) Init() {
	T.upComm = &subCommand{
		command:    "!uptime",
		numArgs:    0,
		cooldown:   15 * time.Second,
		lastCalled: time.Now().Add(-15 * time.Second),
	}
}

func (T *uptimeCommand) ID() string {
	return "uptime"
}

func (T *uptimeCommand) Response(username, message string) string {
	message = strings.TrimSpace(strings.ToLower(message))

	_, err := T.upComm.parse(message)
	if err == nil {
		if !T.cp.broadcaster.Online {
			return fmt.Sprintf("%s isn't online.", T.cp.channel.GetChannelName())
		}
		uptime := time.Now().UTC().Sub(T.cp.broadcaster.OnlineSince)
		minutes := int(uptime.Minutes())
		return fmt.Sprintf("%s hours, %s minutes", strconv.Itoa(minutes/60), strconv.Itoa(minutes%60))
	}

	return ""
}

func (T *uptimeCommand) String() string {
	return "Uptime Command"
}
