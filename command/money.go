package command

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/stats"
)

type money struct {
	cp    *CommandPool
	stats *subCommand
	give  *subCommand
}

func (T *money) Init() {
	T.stats = &subCommand{
		command:   "!toponepercent",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.MOD,
	}

	T.give = &subCommand{
		command:   "!give",
		numArgs:   2,
		cooldown:  1 * time.Millisecond,
		clearance: channel.VIEWER,
	}
}

func (T *money) ID() string {
	return "money"
}

func (T *money) Response(username, message string) string {
	viewer, ok := T.cp.channel.InChannel(username)
	clearance := T.cp.channel.GetLevel(username)
	_, err := T.stats.parse(message, clearance)
	if err == nil {
		T.cp.irc.Say("#"+T.cp.channel.GetChannelName(), T.calculateRichest())
		return ""
	}

	args, err := T.give.parse(message, clearance)
	if err == nil {
		amount, _ := strconv.Atoi(args[1])
		if amount <= 0 {
			T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), username, "Please enter a valid amount.")
			return ""
		}

		to_viewer, ok := T.cp.channel.InChannel(strings.ToLower(args[0]))
		if !ok {
			T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), username, "That user isn't in the chat.")
			return ""
		}

		viewer, _ := T.cp.channel.InChannel(username)
		if viewer.GetMoney() < amount {
			T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), username, fmt.Sprintf("You don't have enough %ss", T.cp.channel.Currency))
			return ""
		}

		viewer.AddMoney(-amount)
		to_viewer.AddMoney(amount)

		T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), username, fmt.Sprintf("You gave %s %d %ss!", to_viewer.Username, amount, T.cp.channel.Currency))
		T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), to_viewer.Username, fmt.Sprintf("You received %d %ss from %s!", amount, T.cp.channel.Currency, username))
		return ""
	}

	if strings.TrimSpace(message) == "!cash" && ok {
		T.cp.ircW.Whisper(T.cp.channel.GetChannelName(), username, fmt.Sprintf("You have %d %ss", viewer.GetMoney(), T.cp.channel.Currency))
		return ""
	}

	return ""
}

type moneySortInterface struct {
	viewers []stats.Viewer
}

func (V *moneySortInterface) Len() int {
	return len(V.viewers)
}

func (V *moneySortInterface) Less(i, j int) bool {
	return V.viewers[i].GetMoney() > V.viewers[j].GetMoney()
}

func (V *moneySortInterface) Swap(i, j int) {
	oldi := V.viewers[i]
	V.viewers[i] = V.viewers[j]
	V.viewers[j] = oldi
}

func (T *money) calculateRichest() string {
	rows, err := T.cp.db.Query(`SELECT sum(c.count) money, v.username FROM counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE c.type='money' AND v.channel=$1 `+
		`GROUP BY v.username ORDER BY money DESC LIMIT 10`, T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return ""
	}

	var viewer string
	var money int
	output := "Richest people: "
	for rows.Next() {
		rows.Scan(&money, &viewer)
		output = fmt.Sprintf("%s%s - %d %ss, ", output, viewer, money, T.cp.channel.Currency)
	}
	return output[:len(output)-2]
}

func (T *money) WhisperOnly() bool {
	return true
}

func (T *money) String() string {
	return ""
}
