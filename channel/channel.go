package channel

import (
	"database/sql"
	"log"
	"strconv"
	"time"
)

type Channel struct {
	Username    string
	Broadcaster *Broadcaster
	ViewerList  *ViewerList

	// Other properties
	Currency         string
	SubName          string
	ComboTrigger     string
	LineTypedReward  int
	MinuteSpentAward int
	Emotes           []string

	db *sql.DB
}

func New(channel string, db *sql.DB) *Channel {
	c := &Channel{
		Username:    channel,
		Broadcaster: NewBroadcaster(channel),
		ViewerList:  NewViewerList(channel, db),
		db:          db,
	}

	log.Printf("getting channel properties for %s", channel)
	rows, err := db.Query("SELECT k, v FROM channel_properties WHERE channel=$1", channel)
	c.Currency = "Coin"
	c.SubName = "subscribers"
	c.ComboTrigger = "PogChamp"
	if err != nil {
		log.Printf("couldn't get channel_properties %s", err.Error())
	}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		c.SetProperty(k, v)
	}
	rows.Close()

	log.Printf("getting emotes properties for %s", channel)
	emotes := []string{}
	rows, err = db.Query("SELECT emote FROM emotes WHERE channel=$1", channel)
	if err != nil {
		log.Printf("couldn't get emotes %s", err.Error())
	}
	for rows.Next() {
		var emote string
		err := rows.Scan(&emote)
		if err == nil {
			emotes = append(emotes, emote)
		}
	}
	rows.Close()

	c.Emotes = emotes
	log.Printf("Loaded emotes %s for %s", emotes, channel)

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
	} else if k == "line_typed_reward" {
		V.LineTypedReward, _ = strconv.Atoi(v)
	} else if k == "minute_spent_reward" {
		V.MinuteSpentAward, _ = strconv.Atoi(v)
	} else {
		valid = false
	}
	if valid {
		insert := "INSERT INTO channel_properties (channel, k, v) SELECT $1, $2, $3"
		upsert := "UPDATE channel_properties SET v=$3 WHERE k=$2 AND channel=$1"
		V.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", V.GetChannelName(), k, v)
	}
}

func (V *Channel) AddEmote(e string) {
	for _, emote := range V.Emotes {
		if e == emote {
			return
		}
	}
	V.Emotes = append(V.Emotes, e)
	V.db.Exec("INSERT INTO emotes (channel, emote) VALUES ($1, $2)", V.Username, e)
}

func (V *Channel) DeleteEmote(e string) {
	for i, emote := range V.Emotes {
		if e == emote {
			V.Emotes = append(V.Emotes[:i], V.Emotes[i+1:]...)
			V.db.Exec("DELETE FROM emotes WHERE channel=$1 AND emote=$2", V.Username, e)
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
