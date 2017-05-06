package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jixwanwang/jixbot/db"
)

const file = "output.txt"

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")

	database, err := db.New(host, port, dbname, user, password)
	if err != nil {
		log.Printf("couldn't connect to db %s", err.Error())
	}

	quotes, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(quotes), "\n")
	for _, line := range lines {
		// rank, _ := strconv.Atoi(line[:strings.Index(line, ".")])
		quote := line[strings.Index(line, ".")+2:]

		database.AddQuote("hotform", quote)
	}
}
