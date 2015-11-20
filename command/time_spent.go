package command

import (
	"fmt"
	"log"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type timeSpent struct {
	cp *CommandPool

	time  *subCommand
	stats *subCommand
}

func (T *timeSpent) Init() {
	T.time = &subCommand{
		command:   "!timespent",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.VIEWER,
	}

	T.stats = &subCommand{
		command:   "!longestviewers",
		numArgs:   0,
		cooldown:  1 * time.Minute,
		clearance: channel.MOD,
	}
}

func (T *timeSpent) ID() string {
	return "timespent"
}

func (T *timeSpent) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)
	_, err := T.stats.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.calculateLongest())
	}
}

func (T *timeSpent) calculateLongest() string {
	rows, err := T.cp.db.Query(`SELECT sum(c.count) time, v.username FROM counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE c.type='time' AND v.channel=$1 `+
		`GROUP BY v.username ORDER BY time DESC LIMIT 10`, T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		return ""
	}

	var viewer string
	var time int
	output := "Longest watchers: "
	for rows.Next() {
		rows.Scan(&time, &viewer)
		output = fmt.Sprintf("%s%s - %d minutes, ", output, viewer, time)
	}
	rows.Close()
	return output[:len(output)-2]
}
