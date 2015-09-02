package command

import (
	"fmt"
	"time"
)

type combo struct {
	cp *CommandPool

	lastCombo time.Time
	active    bool
	comboers  map[string]bool
}

func (T *combo) Init() {
	T.lastCombo = time.Now().Add(-10 * time.Minute)
	T.active = false
}

func (T *combo) ID() string {
	return "combo"
}

func (T *combo) Response(username, message string) string {
	// TODO: make sub only
	if message == T.cp.channel.ComboTrigger {
		if !T.active {
			T.comboers = map[string]bool{}
			T.lastCombo = time.Now()
			T.comboers[username] = true
			T.active = true
		} else if time.Since(T.lastCombo).Seconds() < 10 {
			if _, ok := T.comboers[username]; !ok {
				T.comboers[username] = true
				T.lastCombo = time.Now()
			}
		} else {
			T.active = false
			return fmt.Sprintf("%s %s %d COMBO %s %s",
				T.cp.channel.ComboTrigger,
				T.cp.channel.ComboTrigger,
				len(T.comboers),
				T.cp.channel.ComboTrigger,
				T.cp.channel.ComboTrigger)
		}
	}
	return ""
}

func (T *combo) WhisperOnly() bool {
	return false
}

func (T *combo) String() string {
	return ""
}
