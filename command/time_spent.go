package command

import (
	"fmt"
	"log"
	"strings"
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
		cooldown:  1 * time.Second,
		clearance: channel.VIEWER,
	}

	T.stats = &subCommand{
		command:   "!longestwatchers",
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

	viewer, ok := T.cp.channel.InChannel(username)
	if !ok {
		return
	}

	clearance := T.cp.channel.GetLevel(username)

	_, err := T.time.parse(message, clearance)
	if err == nil {
		T.cp.Whisper(username, fmt.Sprintf("You have spent %s watching %s", timeSpentString(viewer.GetTimeSpent()), T.cp.channel.Username))
	}

	_, err = T.stats.parse(message, clearance)
	if err == nil {
		T.cp.Say(T.calculateLongest())
	}
}

func timeSpentString(minutes int) string {
	days := minutes / (60 * 24)
	hours := (minutes - days*60*24) / 60
	mins := minutes - days*60*24 - hours*60
	parts := []string{}
	if days == 1 {
		parts = append(parts, fmt.Sprintf("%d day", days))
	}
	if days > 1 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}

	if hours == 1 {
		parts = append(parts, fmt.Sprintf("%d hour", hours))
	}
	if hours > 1 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}

	parts = append(parts, fmt.Sprintf("%d minutes", mins))
	return strings.Join(parts, ", ")
}

func (T *timeSpent) calculateLongest() string {
	rows, err := T.cp.db.Query(`SELECT sum(c.count) as time, v.username FROM counts AS c `+
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
