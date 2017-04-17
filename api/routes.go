package api

import (
	"encoding/json"
	"net/http"

	"github.com/jixwanwang/jixbot/stream_bot"
	"github.com/zenazn/goji/web"
)

type API struct {
	bots *stream_bot.BotPool
}

func NewAPI(bots *stream_bot.BotPool) (http.Handler, *API, error) {
	api := &API{
		bots: bots,
	}

	mux := web.New()

	// Channel modification
	mux.Post("/channels/:channel", api.newChannelBot)
	mux.Get("/channels", api.getChannels)
	mux.Get("/channels/:channel", api.getChannelInfo)

	// Channel properties
	mux.Put("/channels/properties/:channel", api.setProperty)

	// Command modification
	mux.Get("/commands/text/:channel", api.getTextCommands)
	mux.Get("/commands/:channel", api.getCommands)
	mux.Put("/commands/:channel", api.addCommands)
	mux.Delete("/commands/:channel", api.deleteCommands)

	return mux, api, nil
}

func (T *API) Close() {
	T.bots.Shutdown()
}

func serveError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}

func serveJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		serveError(w, err)
		return
	}

	w.Write(b)
}
