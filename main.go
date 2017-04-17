package main

import (
	"log"
	"net"
	"net/http"

	"github.com/jixwanwang/jixbot/api"
	"github.com/jixwanwang/jixbot/stream_bot"
	"github.com/zenazn/goji/graceful"

	_ "net/http/pprof"
)

const ()

func main() {
	log.SetFlags(0)

	botPool, err := stream_bot.NewPool()
	if err != nil {
		log.Fatalf("failed to create bots: %v", err)
	}

	mux, api, err := api.NewAPI(botPool)
	if err != nil {
		log.Fatalf(err.Error())
	}

	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("./public"))))
	http.Handle("/", mux)

	graceful.HandleSignals()
	graceful.PreHook(api.Close)

	addr := ":8080"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[handler] Failed to bind to %s: %s", addr, err)
	}

	err = graceful.Serve(l, http.DefaultServeMux)
	if err != nil {
		log.Fatalf("[handler] %s", err)
	}
	graceful.Wait()
}
