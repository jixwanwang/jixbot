package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
)

var inputFile string
var channelName string

func main() {
	flag.StringVar(&inputFile, "input-file", "", "path to exported revlo file")
	flag.StringVar(&channelName, "name", "", "name of the channel to import for")
	flag.Parse()

	if inputFile == "" || channelName == "" {
		log.Fatalf("needs arguments")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	database, err := db.New(host, port, dbname, user, password)
	if err != nil {
		log.Fatalf("couldn't connect to db %s", err.Error())
	}

	// format is
	///tacorocco_	100518531	8133
	raw, _ := ioutil.ReadFile(inputFile)
	lines := strings.Split(string(raw), "\n")[1:]
	viewerList := channel.NewViewerList(channelName, database)
	for _, line := range lines {
		pieces := strings.Split(line, ",")
		if len(pieces) != 4 {
			continue
		}

		username := strings.ToLower(pieces[0])
		points, _ := strconv.Atoi(pieces[2])

		user := viewerList.AddViewer(username)
		user.AddMoney(points)
	}

	viewerList.Flush()
}
