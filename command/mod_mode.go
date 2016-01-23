package command

import (
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type modonly struct {
	cp *CommandPool

	modOnly *subCommand
	subOnly *subCommand
	mode    *subCommand
}

func (T *modonly) Init() {
	T.modOnly = &subCommand{
		command:   "!modonly",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.MOD,
	}
	T.subOnly = &subCommand{
		command:   "!subonly",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.MOD,
	}
	T.mode = &subCommand{
		command:   "!commandmode",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.MOD,
	}
}

func (T *modonly) ID() string {
	return "modonly"
}

func (T *modonly) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)
	_, err := T.modOnly.parse(message, clearance)
	if err == nil {
		if T.cp.modOnly {
			T.cp.modOnly = false
			T.cp.Say("/me is no longer in mod only mode. Freedom!")
		} else {
			T.cp.modOnly = true
			T.cp.Say("/me is now in mod only mode. Only mods can use my commands.")
		}
	}

	_, err = T.subOnly.parse(message, clearance)
	if err == nil {
		if T.cp.subOnly {
			T.cp.subOnly = false
			T.cp.Say("/me is no longer in sub only mode. Plebs rejoice!")
		} else {
			T.cp.subOnly = true
			T.cp.Say("/me is now in sub only mode. Only subs can use my commands.")
		}
	}

	_, err = T.mode.parse(message, clearance)
	if err == nil {
		if T.cp.subOnly {
			T.cp.Say("/me is currently in sub only mode. Only subscribers can use my commands")
		} else if T.cp.modOnly {
			T.cp.Say("/me is currently in mod only mode. Only mods can use my commands")
		} else {
			T.cp.Say("/me is currently in normal mode, everyone can use commands!")
		}
	}
}
