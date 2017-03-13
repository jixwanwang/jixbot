package api

import (
	"net/http"

	"github.com/jixwanwang/jixbot/stream_bot"
	"github.com/zenazn/goji/web"
)

func (T *API) newChannelBot(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	if _, ok := T.bots[channel]; ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	b, err := stream_bot.New(channel, T.nickname, T.oath, T.groupchat, T.texter, T.pasteBin, T.db)

	if err != nil {
		serveError(w, err)
		return
	}

	T.bots[channel] = b
	go b.Start()

	w.WriteHeader(http.StatusOK)
}

func (T *API) getChannelInfo(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot, ok := T.bots[channel]
	if !ok {
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
	bot, ok := T.bots[channel]
	if !ok {
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
	channels := []string{}
	for channel := range T.bots {
		channels = append(channels, channel)
	}

	serveJSON(w, channels)
}
