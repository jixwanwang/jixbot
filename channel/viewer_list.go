package channel

import (
	"database/sql"
	"log"
	"time"

	"github.com/jixwanwang/jixbot/stats"
)

type Level int

const (
	VIEWER      Level = 0
	TURBO       Level = 1
	FOLLOWER    Level = 2
	SUBSCRIBER  Level = 4
	MOD         Level = 6
	ADMIN       Level = 7
	STAFF       Level = 9
	BROADCASTER Level = 10
	GOD         Level = 12
)

func init() {
	// TODO: Load known staff
}

type ViewerList struct {
	channel             string
	db                  *stats.ViewerManager
	realDB              *sql.DB
	viewers             map[string]*stats.Viewer
	staff               map[string]int
	mods                map[string]int
	lotteryContributers map[string]int

	// Other properties
	Currency     string
	SubName      string
	ComboTrigger string
	Emotes       []string
}

func NewViewerList(channel string, db *sql.DB) *ViewerList {
	viewers := &ViewerList{
		channel: channel,
		// TODO: this is bad, this class should be the viewer_manager
		db:                  stats.Init(channel, db),
		realDB:              db,
		viewers:             map[string]*stats.Viewer{},
		staff:               map[string]int{},
		mods:                map[string]int{},
		lotteryContributers: map[string]int{},
	}

	// TODO: should be in channel class, not in viewlist class
	rows, err := db.Query("SELECT k, v FROM channel_properties WHERE channel=$1", channel)
	viewers.Currency = "Coin"
	viewers.SubName = "subscribers"
	viewers.ComboTrigger = "PogChamp"
	if err != nil {
		log.Printf("couldn't get channel_properties %s", err.Error())
	}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		viewers.SetProperty(k, v)
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
		if err == nil {
			emotes = append(emotes, emote)
		}
	}
	rows.Close()

	viewers.Emotes = emotes
	log.Printf("Loaded emotes %s for %s", emotes, channel)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			<-ticker.C
			log.Printf("Saving EVERYTHING")
			viewers.db.Flush()
		}
	}()

	return viewers
}

func (V *ViewerList) SetProperty(k, v string) {
	if k == "currency" {
		V.Currency = v
	}
	if k == "subname" {
		V.SubName = v
	}
	if k == "combo_trigger" {
		V.ComboTrigger = v
	}
}

func (V *ViewerList) AddEmote(e string) {
	for _, emote := range V.Emotes {
		if e == emote {
			return
		}
	}
	V.Emotes = append(V.Emotes, e)
	V.realDB.Exec("INSERT INTO emotes (channel, emote) VALUES ($1, $2)", V.channel, e)
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
	for _, v := range V.viewers {
		v.AddMoney(1)
	}
}

func (V *ViewerList) Flush() {
	log.Printf("Flushing viewerlist")
	V.db.Flush()
}

func (V *ViewerList) Close() {
	V.db.Flush()
}
