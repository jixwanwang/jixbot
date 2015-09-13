package command

import (
	"fmt"
	"strings"
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
	// TODO: make sub only maybe?
	if index := strings.Index(message, T.cp.channel.ComboTrigger); index >= 0 {
		if !T.active {
			if time.Since(T.lastCombo).Minutes() > 3 {
				T.comboers = map[string]bool{}
				T.lastCombo = time.Now()
				T.comboers[username] = true
				T.active = true
			}
		} else if time.Since(T.lastCombo).Seconds() < 15 {
			if _, ok := T.comboers[username]; !ok {
				T.comboers[username] = true
				T.lastCombo = time.Now()
				if len(T.comboers) == 5 {
					return fmt.Sprintf("%s The combo begins! Type %s to keep the combo going and get 100 %ss!",
						T.cp.channel.ComboTrigger,
						T.cp.channel.ComboTrigger,
						T.cp.channel.Currency)
				}
				if len(T.comboers)%10 == 0 {
					numCombo := len(T.comboers) / 10
					if numCombo > 5 {
						numCombo = 5
					}
					comboSpam := ""
					for i := 0; i < numCombo; i++ {
						comboSpam = comboSpam + T.cp.channel.ComboTrigger + " "
					}
					return fmt.Sprintf("%s %d COMBO %s",
						comboSpam,
						len(T.comboers),
						comboSpam)
				}
			}
		} else {
			T.lastCombo = time.Now()
			if len(T.comboers) < 5 {
				T.comboers = map[string]bool{}
				T.comboers[username] = true
				viewer, in := T.cp.channel.InChannel(username)
				if in {
					viewer.AddMoney(100)
				}
				return ""
			} else {
				T.active = false
				return fmt.Sprintf("%s C-C-C-C-COMBO BREAKER (%d combo achieved!)", T.cp.channel.ComboTrigger, len(T.comboers))
			}
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
