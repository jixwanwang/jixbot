package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/jixwanwang/jixbot/api"
	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/messaging"
	"github.com/zenazn/goji/graceful"
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
	password := os.Getenv("DB_PASS")
	groupchat := os.Getenv("GROUPCHAT")

	channels := strings.Split(os.Getenv("CHANNELS"), ",")

	db, err := db.New(host, port, dbname, user, password)
	if err != nil {
		log.Fatalf("Failed to create db: %s", err.Error())
	}

	texter := messaging.NewTexter(twilioAccount, twilioSecret, twilioNumber, myNumber)

	mux, api, err := api.NewAPI(channels, nickname, oath, groupchat, texter, db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	http.Handle("/", mux)

	graceful.HandleSignals()
	graceful.PreHook(api.Close)

	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[handler] Failed to bind to %s: %s", addr, err)
	}
	log.Printf("[handler] Listening on %s", addr)

	err = graceful.Serve(l, http.DefaultServeMux)
	if err != nil {
		log.Fatalf("[handler] %s", err)
	}
	log.Printf("[handler] Draining server")
	graceful.Wait()
	log.Printf("[handler] Draining complete")

	// bots := []*stream_bot.Bot{}
	// for _, channel := range channels {
	// 	b, err := stream_bot.New(channel, nickname, oath, groupchat, texter, db)

	// 	if err != nil {
	// 		log.Fatalf("Failed to create client for %s: %s", channel, err.Error())
	// 	}
	// 	bots = append(bots, b)
	// 	go b.Start()
	// }

	// quit := make(chan os.Signal, 1)

	// signal.Notify(quit, os.Interrupt, os.Kill)

	// go func() {
	// 	<-quit
	// 	log.Printf("quitting")
	// 	for _, b := range bots {
	// 		b.Shutdown()
	// 	}
	// 	os.Exit(0)
	// }()

	// // Never finish
	// var wg sync.WaitGroup
	// wg.Add(1)
	// wg.Wait()
}
