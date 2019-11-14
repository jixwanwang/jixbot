package twitch_api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var clientID = ""

func SetClientID(id string) {
	clientID = id
}

type emoticonsAPIResponse struct {
	Plans  EmotePlans `json:"plans"`
	Emotes []Emote    `json:"emotes"`
}

type EmotePlans struct {
	Plan5  string `json:"$4.99"`
	Plan10 string `json:"$9.99"`
	Plan25 string `json:"$24.99"`
}

type Emote struct {
	Code        string `json:"code"`
	EmoticonSet string `json:"emoticon_set,int"`
}

func makeRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", clientID)

	return http.DefaultClient.Do(req)
}

type userAPIResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func getUserID(channel string) string {
	resp, err := makeRequest("GET", "http://api.twitch.tv/helix/users?login="+channel)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var user userAPIResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&user)
	if err != nil {
		log.Printf("failed to parse user %v", err)
		return ""
	}

	if len(user.Data) == 0 {
		return ""
	}

	return user.Data[0].ID
}

func GetEmotes(channel string) []string {
	channelID := getUserID(channel)

	resp, err := makeRequest("GET", "https://api.twitchemotes.com/api/v4/channels/"+channelID)
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
	for _, e := range emotes.Emotes {
		if e.EmoticonSet == emotes.Plans.Plan5 {
			subEmotes = append(subEmotes, e.Code)
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

type HelixStream struct {
	Data []Stream `json:"data"`
}

type Stream struct {
	StartedAt time.Time `json:"started_at"`
}

func LiveStream(channel string) *HelixStream {
	resp, err := makeRequest("GET", "https://api.twitch.tv/helix/streams?user_login="+channel)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var s HelixStream
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
