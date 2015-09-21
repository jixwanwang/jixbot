package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	c "github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
)

const statsFilePath = "data/stats/"
const textFilePath = "data/textcommands/"

type oldStruct struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
	Money      int    `json:"money"`
	BrawlsWon  int    `json:"brawls_won"`
}

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	db, err := db.New(host, port, dbname, user, password)
	if err != nil {
		log.Printf("couldn't connect to db %s", err.Error())
	}

	channels := []string{}
	rows, err := db.Query("SELECT DISTINCT(channel) FROM commands")
	if err != nil {
		log.Fatalf("Failed to get channel list. %s", err.Error())
	}
	for rows.Next() {
		var channel string
		err := rows.Scan(&channel)
		if err == nil {
			channels = append(channels, channel)
		}
	}
	rows.Close()

	for _, channel := range channels {
		db.Query("SELECT c.count, c.viewer_id FROM counts AS c JOIN viewers AS v ON c.viewer_id=v.id AND c.type='money' AND v.channel=$1", channel)
	}

	for _, channel := range channels {
		db.Exec("INSERT INTO channels (username) VALUES ($1)", channel)

		statsRaw, _ := ioutil.ReadFile(statsFilePath + channel + "_stats")
		statLines := strings.Split(string(statsRaw), "\n")

		for _, line := range statLines {
			var old oldStruct

			err := json.Unmarshal([]byte(line), &old)

			if err != nil {
				continue
			}

			db.Exec("INSERT INTO viewers (username, channel) VALUES ($1, $2)", old.Username, channel)
			row := db.QueryRow("SELECT id FROM viewers WHERE username=$1 AND channel=$2", old.Username, channel)
			var id int
			row.Scan(&id)
			if old.BrawlsWon > 0 {
				db.Exec("INSERT INTO brawlwins (season, viewer_id, wins) VALUES ($1, $2, $3)", 1, id, old.BrawlsWon)
			}
			if old.LinesTyped > 0 {
				db.Exec("INSERT INTO counts (type, viewer_id, count) VALUES ($1, $2, $3)", "lines_typed", id, old.LinesTyped)
			}
			if old.Money > 0 {
				db.Exec("INSERT INTO counts (type, viewer_id, count) VALUES ($1, $2, $3)", "money", id, old.Money)
			}
		}
	}

	channels = append(channels, "_global")
	for _, channel := range channels {
		// Migrate text commands
		textCommandsRaw, _ := ioutil.ReadFile(textFilePath + channel)
		lines := strings.Split(string(textCommandsRaw), "\n")
		for _, line := range lines {
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				continue
			}

			var perm c.Level
			switch parts[1] {
			case "viewer":
				perm = c.VIEWER
			case "mod":
				perm = c.MOD
			}

			db.Exec("INSERT INTO textcommands (channel, command, message, clearance) VALUES ($1, $2, $3, $4)", channel, parts[0], parts[2], perm)
		}
	}
}
