package command

import (
	"fmt"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type linesTyped struct {
	cp *CommandPool

	time    *subCommand
	stats   *subCommand
	actives *subCommand
}

func (T *linesTyped) Init() {
	T.time = &subCommand{
		command:   "!linestyped",
		numArgs:   0,
		cooldown:  1 * time.Second,
		clearance: channel.VIEWER,
	}

	T.stats = &subCommand{
		command:   "!chattiest",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.MOD,
	}

	T.actives = &subCommand{
		command:   "!mostactive",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.MOD,
	}
}

func (T *linesTyped) ID() string {
	return "lines_typed"
}

func (T *linesTyped) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	viewer, ok := T.cp.channel.InChannel(username)
	if !ok {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.time.parse(message, clearance)
	if err == nil {
		T.cp.Whisper(username, fmt.Sprintf("You have typed %v lines in %s's chat", viewer.GetLinesTyped(), T.cp.channel.Username))
	}

	_, err = T.stats.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.calculateChattiest())
	}

	_, err = T.actives.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.calculateMostActive())
	}
}

func (T *linesTyped) calculateChattiest() string {
	counts, err := T.cp.db.HighestCount(T.cp.channel.GetChannelName(), "lines_typed")
	if err != nil {
		return ""
	}

	output := "Chattiest users: "
	for _, c := range counts {
		output = fmt.Sprintf("%s%s - %d lines, ", output, c.Username, c.Count)
	}
	return output[:len(output)-2]
}

func (T *linesTyped) calculateMostActive() string {
	ratios, err := T.cp.db.HighestRatio(T.cp.channel.GetChannelName(), "lines_typed", "time_spent")
	if err != nil {
		return ""
	}

	output := "Most active users: "
	for _, c := range ratios {
		output = fmt.Sprintf("%s%s - %.1f lines/s, ", output, c.Username, c.Ratio)
	}
	return output[:len(output)-2]
}
