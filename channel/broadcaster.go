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
	Online      bool
	OnlineSince time.Time

	username   string
	lastOnline time.Time
	tolerance  time.Duration
}

func NewBroadcaster(channel string) *Broadcaster {
	b := &Broadcaster{
		username:  channel,
		tolerance: 1 * time.Minute,
	}

	return b
}

// Tolerance for stream crashes
func (B *Broadcaster) checkOnline() {
	resp, err := http.Get(baseURL + "/" + B.username)
	// Don't change state if the request fails
	if err != nil {
		return
	}

	b, _ := ioutil.ReadAll(resp.Body)

	var v map[string]interface{}

	json.Unmarshal(b, &v)

	if v["stream"] == nil {
		if time.Since(B.lastOnline) > B.tolerance {
			B.Online = false
		}
		return
	}

	stream_info, ok := v["stream"].(map[string]interface{})

	// Don't update the OnlineSince field unless the stream is just coming online.
	if !B.Online && ok && stream_info["created_at"] != nil {
		t, _ := time.Parse("2006-01-02T15:04:05Z", stream_info["created_at"].(string))
		B.OnlineSince = t
	}

	B.lastOnline = time.Now()
	B.Online = true
}
