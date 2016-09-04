package command

import (
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type info struct {
	cp      *CommandPool
	jix     *subCommand
	version *subCommand
	aboot   *subCommand
}

func (T *info) Init() {
	T.jix = &subCommand{
		command:   "!jix",
		numArgs:   0,
		cooldown:  10 * time.Second,
		clearance: channel.VIEWER,
	}
	T.version = &subCommand{
		command:   "!version",
		numArgs:   0,
		cooldown:  15 * time.Second,
		clearance: channel.VIEWER,
	}
	T.aboot = &subCommand{
		command:   "!about",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.VIEWER,
	}
}

func (T *info) ID() string {
	return "info"
}

func (T *info) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.jix.parse(message, clearance)
	if err == nil {
		T.cp.Say("Staff DansGame hide the weed!")
	}

	_, err = T.version.parse(message, clearance)
	if err == nil {
		T.cp.Say("/me is v3.0")
	}

	_, err = T.aboot.parse(message, clearance)
	if err == nil {
		T.cp.Say(`/me was developed by staff member Jix on his free time. ` +
			`Documentation of available commands at https://github.com/jixwanwang/jixbot/blob/master/docs/commands.md . ` +
			`The commands need to be enabled by Jix to work. ` +
			`Source code is available at https://github.com/jixwanwang/jixbot . ` +
			`Any feedback would be appreciated!`)
	}
}
