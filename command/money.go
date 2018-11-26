package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type money struct {
	cp *CommandPool

	cash    *subCommand
	stats   *subCommand
	give    *subCommand
	giveAll *subCommand
}

func (T *money) Init() {
	T.cash = &subCommand{
		command:   "!cash",
		numArgs:   0,
		cooldown:  200 * time.Millisecond,
		clearance: channel.VIEWER,
	}

	T.stats = &subCommand{
		command:   "!toponepercent",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.MOD,
	}

	T.give = &subCommand{
		command:   "!give",
		numArgs:   2,
		cooldown:  200 * time.Millisecond,
		clearance: channel.VIEWER,
	}

	T.giveAll = &subCommand{
		command:   "!giveall",
		numArgs:   1,
		cooldown:  1 * time.Minute,
		clearance: channel.BROADCASTER,
	}
}

func (T *money) ID() string {
	return "money"
}

func (T *money) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	viewer, ok := T.cp.channel.InChannel(username)
	if !ok {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.stats.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.calculateRichest())
		return
	}

	args, err := T.give.parse(message, clearance)
	if err == nil {
		amount, _ := strconv.Atoi(args[1])
		if amount <= 0 {
			T.cp.Whisper(username, "Please enter a valid amount.")
			return
		}

		to_viewer, ok := T.cp.channel.InChannel(strings.ToLower(args[0]))
		if !ok {
			T.cp.Whisper(username, "That user isn't in the chat.")
			return
		}

		viewer, _ := T.cp.channel.InChannel(username)
		if clearance != channel.GOD && clearance != channel.BROADCASTER {
			if viewer.GetMoney() < amount {
				T.cp.Whisper(username, fmt.Sprintf("You don't have enough %ss", T.cp.channel.Currency))
				return
			}
			viewer.AddMoney(-amount)
		}
		to_viewer.AddMoney(amount)

		T.cp.Whisper(username, fmt.Sprintf("You gave %s %d %ss!", to_viewer.Username, amount, T.cp.channel.Currency))
		T.cp.Whisper(to_viewer.Username, fmt.Sprintf("You received %d %ss from %s!", amount, T.cp.channel.Currency, username))
		return
	}

	_, err = T.cash.parse(message, clearance)
	if err == nil {
		T.cp.Whisper(username, fmt.Sprintf("You have %d %ss in %s's channel", viewer.GetMoney(), T.cp.channel.Currency, T.cp.channel.Username))
		return
	}

	args, err = T.giveAll.parse(message, clearance)
	if err == nil {
		amount, _ := strconv.Atoi(args[0])
		T.cp.channel.AddMoney(amount)
		T.cp.channel.Flush()
		T.cp.Say(fmt.Sprintf("Everyone gets %d %ss, go nuts!", amount, T.cp.channel.Currency))
	}

	return
}

func (T *money) calculateRichest() string {
	counts, err := T.cp.db.HighestCount(T.cp.channel.GetChannelName(), "money")
	if err != nil {
		return ""
	}

	output := "Richest people: "
	for _, c := range counts {
		output = fmt.Sprintf("%s%s - %d %ss, ", output, c.Username, c.Count, T.cp.channel.Currency)
	}
	return output[:len(output)-2]
}
