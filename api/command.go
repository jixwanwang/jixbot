package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

func (T *API) getCommands(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]

	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()

	comms := bot.GetActiveCommands()

	serveJSON(w, comms)
}

func (T *API) addCommands(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]

	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	commands := strings.Split(r.FormValue("commands"), ",")
	log.Printf("adding %s to commands for %s", commands, channel)
	for _, c := range commands {
		bot.AddActiveCommand(c)
	}

	comms := bot.GetActiveCommands()

	serveJSON(w, comms)
}

func (T *API) deleteCommands(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]

	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.ParseForm()
	commands := strings.Split(r.FormValue("commands"), ",")
	log.Printf("adding %s to commands for %s", commands, channel)
	for _, c := range commands {
		bot.DeleteCommand(c)
	}

	comms := bot.GetActiveCommands()

	serveJSON(w, comms)
}
