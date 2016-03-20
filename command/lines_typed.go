package command

import (
	"fmt"
	"log"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type linesTyped struct {
	cp *CommandPool

	time  *subCommand
	stats *subCommand
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
}

func (T *linesTyped) calculateChattiest() string {
	rows, err := T.cp.db.Query(`SELECT sum(c.count) as lines, v.username FROM counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE c.type='lines_typed' AND v.channel=$1 `+
		`GROUP BY v.username ORDER BY lines DESC LIMIT 10`, T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return ""
	}

	var viewer string
	var lines int
	output := "Chattiest users: "
	for rows.Next() {
		rows.Scan(&lines, &viewer)
		output = fmt.Sprintf("%s%s - %d lines, ", output, viewer, lines)
	}
	rows.Close()
	return output[:len(output)-2]
}
