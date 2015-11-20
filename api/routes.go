package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jixwanwang/jixbot/messaging"
	"github.com/jixwanwang/jixbot/stream_bot"
	"github.com/zenazn/goji/web"
)

type API struct {
	nickname  string
	oath      string
	groupchat string
	texter    messaging.Texter
	db        *sql.DB

	bots map[string]*stream_bot.Bot
}

func NewAPI(channels []string, nickname, oath, groupchat string, texter messaging.Texter, db *sql.DB) (http.Handler, *API, error) {
	api := &API{
		nickname:  nickname,
		oath:      oath,
		groupchat: groupchat,
		texter:    texter,
		db:        db,
		bots:      map[string]*stream_bot.Bot{},
	}
	for _, channel := range channels {
		log.Printf("loading bot for %s", channel)
		b, err := stream_bot.New(channel, nickname, oath, groupchat, texter, db)

		if err != nil {
			return nil, nil, err
		}
		api.bots[channel] = b
		go b.Start()
	}

	mux := web.New()

	// Channel modification
	mux.Post("/channels/:channel", api.newChannelBot)
	// mux.Get("/channels", api.getChannels)
	// mux.Get("/channels/:channel", api.getChannelInfo)

	// Channel properties
	mux.Put("/channels/properties/:channel", api.setProperty)

	// Command modification
	mux.Get("/commands/:channel", api.getCommands)
	mux.Put("/commands/:channel", api.addCommands)
	mux.Delete("/commands/:channel", api.deleteCommands)

	// Emote modification
	mux.Get("/emotes/:channel", api.getEmotes)
	mux.Put("/emotes/:channel", api.addEmotes)
	mux.Delete("/emotes/:channel", api.deleteEmotes)

	return mux, api, nil
}

func (T *API) Close() {
	for _, b := range T.bots {
		b.Shutdown()
	}
}

func serveError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}

func (T *API) newChannelBot(C web.C, w http.ResponseWriter, r *http.Request) {
	channel := C.URLParams["channel"]
	if _, ok := T.bots[channel]; ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	b, err := stream_bot.New(channel, T.nickname, T.oath, T.groupchat, T.texter, T.db)

	if err != nil {
		serveError(w, err)
		return
	}

	T.bots[channel] = b
	go b.Start()

	w.WriteHeader(http.StatusOK)
}
