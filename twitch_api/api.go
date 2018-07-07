package twitch_api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var clientID = ""

func SetClientID(id string) {
	clientID = id
}

type emoticonsAPIResponse struct {
	Emoticons []Emote `"json:emoticons"`
}

type Emote struct {
	Regex          string `json:"regex`
	State          string `json:"state"`
	SubscriberOnly bool   `json:"subscriber_only"`
}

func makeRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", clientID)

	return http.DefaultClient.Do(req)
}

func GetEmotes(channel string) []string {
	resp, err := makeRequest("GET", "http://api.twitch.tv/kraken/chat/"+channel+"/emoticons?on_site=1")
	if err != nil {
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
	resp, err := makeRequest("GET", "http://tmi.twitch.tv/servers?channel="+channel)
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
	resp, err := makeRequest("GET", "http://tmi.twitch.tv/servers?cluster=group")
	if err != nil {
		return def
	}
	defer resp.Body.Close()

	var m ircServers
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&m)
	if err != nil {
		return def
	}

	return m.Servers[rand.Intn(len(m.Servers))]
}

type KrakenStream struct {
	Stream *Stream `json:"stream"`
}

type Stream struct {
	CreatedAt time.Time `json:"created_at"`
}

func LiveStream(channel string) *KrakenStream {
	resp, err := makeRequest("GET", "https://api.twitch.tv/kraken/streams/"+channel)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var s KrakenStream
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&s)

	return &s
}

func QueueSoundEffect(name string) {
	log.Printf("%s/enqueue/%s/%s", os.Getenv("SOUND_EFFECT_URL"), name, os.Getenv("SOUND_EFFECT_TOKEN"))
	makeRequest("GET", fmt.Sprintf("%s/enqueue/%s/%s", os.Getenv("SOUND_EFFECT_URL"), name, os.Getenv("SOUND_EFFECT_TOKEN")))
}

func ListSoundEffects() []string {
	resp, err := makeRequest("GET", fmt.Sprintf("%s/tracks/%s", os.Getenv("SOUND_EFFECT_URL"), os.Getenv("SOUND_EFFECT_TOKEN")))
	if err != nil {
		return []string{}
	}

	var sounds []string
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&sounds)

	return sounds
}
