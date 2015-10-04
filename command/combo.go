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

func (T *combo) Response(username, message string) {
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
					T.cp.Say(fmt.Sprintf("%s The combo begins! Type %s to keep the combo going and get 100 %ss!",
						T.cp.channel.ComboTrigger,
						T.cp.channel.ComboTrigger,
						T.cp.channel.Currency))
					return
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
					T.cp.Say(fmt.Sprintf("%s %d COMBO %s", comboSpam, len(T.comboers), comboSpam))
					return
				}
			}
		} else {
			T.lastCombo = time.Now()
			if len(T.comboers) < 5 {
				T.comboers = map[string]bool{}
				T.comboers[username] = true
			} else {
				for c := range T.comboers {
					viewer, in := T.cp.channel.InChannel(c)
					if in {
						viewer.AddMoney(100)
					}
				}
				T.active = false
				T.cp.Say(fmt.Sprintf("%s C-C-C-C-COMBO BREAKER (%d combo achieved!)", T.cp.channel.ComboTrigger, len(T.comboers)))
			}
		}
	}
}
