package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/stream_bot"
)

const ()

func main() {
	log.SetFlags(0)

	nickname := os.Getenv("NICKNAME")
	oath := os.Getenv("OATH_TOKEN")
	twilioAccount := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioSecret := os.Getenv("TWILIO_SECRET")
	twilioNumber := os.Getenv("TWILIO_NUMBER")
	myNumber := os.Getenv("JIX_NUMBER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")

	channels := strings.Split(os.Getenv("CHANNELS"), ",")

	db, err := db.New(host, port, dbname, user)
	if err != nil {
		log.Fatalf("Failed to create db: %s", err.Error())
	}

	texter := messaging.NewTexter(twilioAccount, twilioSecret, twilioNumber, myNumber)

	bots := []*stream_bot.Bot{}
	for _, channel := range channels {
		b, err := stream_bot.New(channel, nickname, oath, texter, db)

		if err != nil {
			log.Fatalf("Failed to create client for %s: %s", channel, err.Error())
		}
		bots = append(bots, b)
		go b.Start()
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, os.Kill)

	go func() {
		<-quit
		log.Printf("quitting")
		for _, b := range bots {
			b.Shutdown()
		}
		os.Exit(0)
	}()

	// Never finish
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
