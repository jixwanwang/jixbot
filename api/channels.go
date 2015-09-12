package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

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
