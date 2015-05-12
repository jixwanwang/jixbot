package channel

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.twitch.tv/kraken/streams"
)

type Broadcaster struct {
	Username    string
	OnlineSince time.Time
}

func NewBroadcaster(channel string) *Broadcaster {
	return &Broadcaster{Username: channel}
}

func (B *Broadcaster) Online() bool {
	resp, _ := http.Get(baseURL + "/" + B.Username)
	b, _ := ioutil.ReadAll(resp.Body)

	var v map[string]interface{}

	json.Unmarshal(b, &v)

	if v["stream"] == nil {
		return false
	}

	stream_info, ok := v["stream"].(map[string]interface{})

	if ok && stream_info["created_at"] != nil {
		t, _ := time.Parse("2006-01-02T15:04:05Z", stream_info["created_at"].(string))
		B.OnlineSince = t
	}

	return true
}
