package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

func (T *API) newChannelBot(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot := T.bots.GetBot(channel)
	if bot != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	T.bots.AddBot(channel)

	w.WriteHeader(http.StatusOK)
}

func (T *API) getChannelInfo(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot := T.bots.GetBot(channel)
	if bot == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	commands := bot.GetActiveCommands()
	props := bot.GetProperties()

	serveJSON(w, map[string]interface{}{
		"commands":   commands,
		"properties": props,
	})
}

func (T *API) setProperty(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot := T.bots.GetBot(channel)
	if bot == nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	r.ParseForm()
	key := r.FormValue("key")
	value := r.FormValue("value")
	bot.SetProperty(key, value)

	w.WriteHeader(http.StatusOK)
}

func (T *API) getChannels(c web.C, w http.ResponseWriter, r *http.Request) {
	serveJSON(w, T.bots.GetChannels())
}
