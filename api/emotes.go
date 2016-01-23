package api

import (
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

func (T *API) getEmotes(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	emotes := bot.GetEmotes()

	serveJSON(w, emotes)
}

func (T *API) addEmotes(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	r.ParseForm()
	emotes := strings.Split(r.FormValue("emotes"), ",")
	for _, e := range emotes {
		bot.AddEmote(e)
	}

	e := bot.GetEmotes()

	serveJSON(w, e)
}

func (T *API) deleteEmotes(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	bot, ok := T.bots[channel]
	if !ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	r.ParseForm()
	emotes := strings.Split(r.FormValue("emotes"), ",")
	for _, e := range emotes {
		bot.DeleteEmote(e)
	}

	e := bot.GetEmotes()

	serveJSON(w, e)
}
