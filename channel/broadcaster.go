package channel

import (
	"time"

	"github.com/jixwanwang/jixbot/twitch_api"
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
		tolerance: 2 * time.Minute,
	}

	return b
}

// Tolerance for stream crashes
func (B *Broadcaster) checkOnline() {
	stream := twitch_api.LiveStream(B.username)
	if stream == nil || len(stream.Data) == 0 {
		if time.Since(B.lastOnline) > B.tolerance {
			B.Online = false
		}
		return
	}

	// Don't update the OnlineSince field unless the stream is just coming online.
	if !B.Online {
		B.OnlineSince = stream.Data[0].StartedAt
	}

	B.lastOnline = time.Now()
	B.Online = true
}
