package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jixwanwang/jixbot/db"
)

const file = "output.txt"
const file2 = "output2.txt"

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
		quote := line[strings.Index(line, ".")+2:]

		_, err := database.AddQuote("hotform", "quote", quote)
		if err != nil {
			log.Printf("%v", err)
		}
	}

	clips, _ := ioutil.ReadFile(file2)
	lines = strings.Split(string(clips), "\n")
	for _, line := range lines {
		quote := line[strings.Index(line, ".")+2:]

		_, err := database.AddQuote("hotform", "clip", quote)
		if err != nil {
			log.Printf("%v", err)
		}
	}
}
