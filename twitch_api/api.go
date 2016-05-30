package twitch_api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type emoticonsAPIResponse struct {
	Emoticons []Emote `"json:emoticons"`
}

type Emote struct {
	Regex          string `json:"regex`
	State          string `json:"state"`
	SubscriberOnly bool   `json:"subscriber_only"`
}

func GetEmotes(channel string) []string {
	resp, err := http.Get("http://api.twitch.tv/kraken/chat/hotform/emoticons?on_site=1")
	if err != nil {
		log.Printf("failed to do GET %v", err)
		return []string{}
	}

	defer resp.Body.Close()
	var emotes emoticonsAPIResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&emotes)
	if err != nil {
		log.Printf("failed to parse emotes %v", err)
		return []string{}
	}

	subEmotes := []string{}
	for _, e := range emotes.Emoticons {
		if e.SubscriberOnly && strings.ToLower(e.State) == "active" {
			subEmotes = append(subEmotes, e.Regex)
		}
	}

	return subEmotes
}

type ircServers struct {
	Servers []string `json:"servers"`
}

func GetIRCServer(channel, def string) string {
	resp, err := http.Get("http://tmi.twitch.tv/servers?channel=" + channel)
	if err == nil {
		defer resp.Body.Close()
		var m ircServers
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&m)
		if err == nil {
			return m.Servers[rand.Intn(len(m.Servers))]
		}
	}
	return def
}

func GetIRCCluster(def string) string {
	resp, err := http.Get("http://tmi.twitch.tv/servers?cluster=group")
	if err == nil {
		defer resp.Body.Close()
		var m ircServers
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&m)
		if err == nil {
			return m.Servers[rand.Intn(len(m.Servers))]
		}
	}
	return def
}

type KrakenStream struct {
	Stream *Stream `json:"stream"`
}

type Stream struct {
	CreatedAt time.Time `json:"created_at"`
}

func LiveStream(channel string) *KrakenStream {
	resp, err := http.Get("https://api.twitch.tv/kraken/streams/" + B.username)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var s KrakenStream
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&s)

	return &s
}
