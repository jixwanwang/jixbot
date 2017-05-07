package command

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

var quoteKinds = map[string]string{
	"quote": "Quote",
	"clip":  "Clip",
}

type quotes struct {
	cp *CommandPool

	getQuote    *subCommand
	addQuote    *subCommand
	deleteQuote *subCommand

	getClip    *subCommand
	addClip    *subCommand
	deleteClip *subCommand
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
	T.deleteQuote = &subCommand{
		command:   "!deletequote",
		numArgs:   1,
		cooldown:  5 * time.Second,
		clearance: channel.MOD,
	}

	T.getClip = &subCommand{
		command:   "!clip",
		numArgs:   0,
		cooldown:  5 * time.Second,
		clearance: channel.VIEWER,
	}
	T.addClip = &subCommand{
		command:   "!addclip",
		numArgs:   0,
		cooldown:  10 * time.Second,
		clearance: channel.MOD,
	}
	T.deleteClip = &subCommand{
		command:   "!deleteclip",
		numArgs:   1,
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
		T.getQuoteHelper(args, "quote")
		return
	}

	args, err = T.addQuote.parse(message, clearance)
	if err == nil {
		T.addQuoteHelper(args, "quote")
		return
	}

	args, err = T.deleteQuote.parse(message, clearance)
	if err == nil {
		T.deleteQuoteHelper(args[0], "quote")
		return
	}

	args, err = T.getClip.parse(message, clearance)
	if err == nil {
		T.getQuoteHelper(args, "clip")
		return
	}

	args, err = T.addClip.parse(message, clearance)
	if err == nil {
		T.addQuoteHelper(args, "clip")
		return
	}

	args, err = T.deleteClip.parse(message, clearance)
	if err == nil {
		T.deleteQuoteHelper(args[0], "clip")
		return
	}
}

func (T *quotes) getQuoteHelper(args []string, kind string) {
	if len(args) == 0 {
		args = []string{""}
	}

	if args[0] == "list" {
		commands := []string{fmt.Sprintf("#\t| %s", quoteKinds[kind])}

		quotes, err := T.cp.db.AllQuotes(T.cp.channel.GetChannelName(), kind)
		if err != nil {
			return
		}

		for _, q := range quotes {
			commands = append(commands, fmt.Sprintf("%v.\t| %s", q.Rank, q.Quote))
		}

		paste := T.cp.pasteBin.Paste(fmt.Sprintf("%ss for %s", quoteKinds[kind], T.cp.channel.GetChannelName()), strings.Join(commands, "\n"), "1", "10M")
		if len(paste) > 0 {
			T.cp.Say(fmt.Sprintf("Saved %ss: %s", quoteKinds[kind], paste))
		}
		return
	}

	rank, _ := strconv.Atoi(args[0])

	var (
		quote     string
		quoteRank int
		err       error
	)

	// Do a search if the argument exists but can't be parsed into an int
	if rank == 0 && len(args[0]) > 0 {
		quote, quoteRank, err = T.cp.db.SearchQuote(T.cp.channel.GetChannelName(), kind, args[0])
	} else {
		quote, quoteRank, err = T.cp.db.GetQuote(T.cp.channel.GetChannelName(), kind, rank)
	}

	if err == nil {
		T.cp.Say(fmt.Sprintf("%v #%v: %v", quoteKinds[kind], quoteRank, quote))
	} else if err == sql.ErrNoRows {
		T.cp.Say("None available")
	}
}

func (T *quotes) addQuoteHelper(args []string, kind string) {
	if len(args) == 0 || len(args[0]) == 0 {
		return
	}
	quote := args[0]

	rank, err := T.cp.db.AddQuote(T.cp.channel.GetChannelName(), kind, quote)
	if err == nil {
		T.cp.Say(fmt.Sprintf("%v #%v added!", quoteKinds[kind], rank))
	} else {
		log.Printf("%v", err)
		T.cp.Say("Failed to add quote, please try again!")
	}
}

func (T *quotes) deleteQuoteHelper(rankString string, kind string) {
	rank, _ := strconv.Atoi(rankString)

	rank, err := T.cp.db.DeleteQuote(T.cp.channel.GetChannelName(), kind, rank)
	if err == nil {
		T.cp.Say(fmt.Sprintf("%v #%v deleted", quoteKinds[kind], rank))
	} else {
		T.cp.Say("Failed to add quote, please try again!")
	}
}
