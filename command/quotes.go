package command

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type quotes struct {
	cp *CommandPool

	getQuote *subCommand
	addQuote *subCommand
}

func (T *quotes) Init() {
	T.getQuote = &subCommand{
		command:   "!quote",
		numArgs:   0,
		cooldown:  5 * time.Second,
		clearance: channel.VIEWER,
	}
	T.addQuote = &subCommand{
		command:   "!addquote",
		numArgs:   0,
		cooldown:  10 * time.Second,
		clearance: channel.MOD,
	}
}

func (T *quotes) ID() string {
	return "quotes"
}

func (T *quotes) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	args, err := T.getQuote.parse(message, clearance)
	if err == nil {
		if len(args) == 0 {
			args = []string{""}
		}

		rank, _ := strconv.Atoi(args[0])
		quote, r, err := T.cp.db.GetQuote(T.cp.channel.GetChannelName(), rank)
		if err == nil {
			T.cp.Say(fmt.Sprintf("Quote #%v: %v", r, quote))
		} else if err == sql.ErrNoRows {
			T.cp.Say("No quote available")
		}
		return
	}

	args, err = T.addQuote.parse(message, clearance)
	if err == nil {
		if len(args) == 0 || len(args[0]) == 0 {
			return
		}
		quote := args[0]

		rank, err := T.cp.db.AddQuote(T.cp.channel.GetChannelName(), quote)
		if err == nil {
			T.cp.Say(fmt.Sprintf("Quote #%v added!", rank))
		} else {
			T.cp.Say("Failed to add quote, please try again!")
		}
		return
	}
}
