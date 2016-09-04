package channel

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/db"
	"github.com/jixwanwang/jixbot/twitch_api"
)

type Channel struct {
	Username    string
	Broadcaster *Broadcaster
	ViewerList  *ViewerList

	// Other properties
	Currency         string
	SubName          string
	ComboTrigger     string
	ComboTriggers    []string
	LineTypedReward  int
	MinuteSpentAward int
	Emotes           []string

	db db.DB
}

func New(channel string, db db.DB) *Channel {
	c := &Channel{
		Username:    channel,
		Broadcaster: NewBroadcaster(channel),
		ViewerList:  NewViewerList(channel, db),
		db:          db,
	}

	properties, err := db.GetChannelProperties(channel)
	c.Currency = "Coin"
	c.SubName = "subscribers"
	c.ComboTrigger = "PogChamp"
	c.ComboTriggers = []string{"PogChamp"}
	if err == nil {
		for k, v := range properties {
			c.SetProperty(k, v)
		}
	} else {
		log.Printf("Failed to get channel properties for hotform!")
	}

	c.Emotes = twitch_api.GetEmotes(channel)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			<-ticker.C
			c.Broadcaster.checkOnline()
		}
	}()
	c.Broadcaster.checkOnline()

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			<-ticker.C
			if c.Broadcaster.Online {
				c.AddTime(1)
			}
			c.ViewerList.Flush()
		}
	}()

	return c
}

func (V *Channel) GetChannelName() string {
	return V.Username
}

func (V *Channel) GetLevel(username string) Level {
	return V.ViewerList.GetLevel(username)
}

func (V *Channel) InChannel(username string) (*Viewer, bool) {
	return V.ViewerList.InChannel(username)
}

func (V *Channel) SetProperty(k, v string) {
	valid := true
	if k == "currency" {
		V.Currency = v
	} else if k == "subname" {
		V.SubName = v
	} else if k == "combo_trigger" {
		V.ComboTrigger = v
	} else if k == "combo_triggers" {
		V.ComboTriggers = strings.Split(v, ",")
	} else if k == "line_typed_reward" {
		V.LineTypedReward, _ = strconv.Atoi(v)
	} else if k == "minute_spent_reward" {
		V.MinuteSpentAward, _ = strconv.Atoi(v)
	} else {
		valid = false
	}
	if valid {
		V.db.SetChannelProperty(V.GetChannelName(), k, v)
	}
}

func (V *Channel) GetProperties() map[string]interface{} {
	return map[string]interface{}{
		"currency":            V.Currency,
		"subname":             V.SubName,
		"combo_trigger":       V.ComboTrigger,
		"line_typed_reward":   V.LineTypedReward,
		"minute_spent_reward": V.MinuteSpentAward,
	}
}

func (V *Channel) AddEmote(e string) {
	for _, emote := range V.Emotes {
		if e == emote {
			return
		}
	}
	V.Emotes = append(V.Emotes, e)
	V.db.AddChannelEmote(V.GetChannelName(), e)
}

func (V *Channel) DeleteEmote(e string) {
	for i, emote := range V.Emotes {
		if e == emote {
			V.Emotes = append(V.Emotes[:i], V.Emotes[i+1:]...)
			V.db.DeleteChannelEmote(V.GetChannelName(), e)
		}
	}
}

func (V *Channel) RecordMessage(username, msg string) {
	v, ok := V.ViewerList.InChannel(username)
	if !ok {
		v = V.ViewerList.AddViewer(username)
	}

	if V.Broadcaster.Online {
		v.AddLineTyped()
		v.AddMoney(V.LineTypedReward)
	}
}

func (V *Channel) AddTime(minutes int) {
	for _, v := range V.ViewerList.viewers {
		v.AddTimeSpent(minutes)
		v.AddMoney(V.MinuteSpentAward * minutes)
	}
}

func (V *Channel) Flush() {
	V.ViewerList.Flush()
}
