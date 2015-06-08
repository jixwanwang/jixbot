package channel

import (
	"database/sql"
	"log"

	"github.com/jixwanwang/jixbot/stats"
)

type Level int

const (
	VIEWER      Level = 0
	MOD         Level = 1
	STAFF       Level = 2
	BROADCASTER Level = 3
	GOD         Level = 4
)

func init() {
	// TODO: Load known staff
}

type ViewerList struct {
	channel             string
	db                  *stats.ViewerManager
	viewers             map[string]*stats.Viewer
	staff               map[string]int
	mods                map[string]int
	lotteryContributers map[string]int

	// Other properties
	Currency string
	Emotes   []string
}

func NewViewerList(channel string, db *sql.DB) *ViewerList {
	viewers := &ViewerList{
		channel:             channel,
		db:                  stats.Init(channel, db),
		viewers:             map[string]*stats.Viewer{},
		staff:               map[string]int{},
		mods:                map[string]int{},
		lotteryContributers: map[string]int{},
	}

	rows, err := db.Query("SELECT k, v FROM channel_properties WHERE channel=$1", channel)
	if err != nil {
		log.Printf("couldn't get channel_properties %s", err.Error())
	}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		if k == "currency" {
			viewers.Currency = v
		}
	}
	rows.Close()

	emotes := []string{}
	rows, err = db.Query("SELECT emote FROM emotes WHERE channel=$1", channel)
	if err != nil {
		log.Printf("couldn't get emotes %s", err.Error())
	}
	for rows.Next() {
		var emote string
		err := rows.Scan(&emote)
		if err != nil {
			emotes = append(emotes, emote)
		}
	}
	rows.Close()

	viewers.Emotes = emotes

	return viewers
}

func (V *ViewerList) GetChannelName() string {
	return V.channel
}

func (V *ViewerList) AddViewer(username string) {
	if _, ok := V.viewers[username]; !ok {
		V.viewers[username] = V.db.FindViewer(username)
	}
}

func (V *ViewerList) AddViewers(usernames []string) {
	for _, u := range usernames {
		V.AddViewer(u)
	}
}

func (V *ViewerList) RemoveViewer(username string) {
	delete(V.viewers, username)
	delete(V.mods, username)
}

func (V *ViewerList) AddMod(username string) {
	V.AddViewers([]string{username})
	if _, ok := V.viewers[username]; ok {
		V.mods[username] = 1
	}
}

func (V *ViewerList) RemoveMod(username string) {
	delete(V.mods, username)
}

func (V *ViewerList) InChannel(username string) (*stats.Viewer, bool) {
	v, ok := V.viewers[username]
	return v, ok
}

func (V *ViewerList) AllViewers() []*stats.Viewer {
	return V.db.AllViewers()
}

func (V *ViewerList) GetLevel(username string) Level {
	if username == "jixwanwang" {
		return GOD
	} else if username == V.channel {
		return BROADCASTER
	} else if _, ok := V.staff[username]; ok {
		return STAFF
	} else if _, ok := V.mods[username]; ok {
		return MOD
	}
	return VIEWER
}

func (V *ViewerList) RecordMessage(username, msg string) {
	v, ok := V.viewers[username]
	if !ok {
		V.AddViewer(username)
		v = V.viewers[username]
	}

	v.AddLineTyped()
}

func (V *ViewerList) Tick() {
	// TODO: make money work
	for _, v := range V.viewers {
		v.AddMoney(1)
	}
	V.db.Flush()
}

func (V *ViewerList) Close() {
	V.db.Flush()
}
