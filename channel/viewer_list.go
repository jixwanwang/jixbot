package channel

import (
	"database/sql"
	"log"
	"time"
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
	db                  *sql.DB
	viewers             map[string]*Viewer
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
		channel:             channel,
		db:                  db,
		viewers:             map[string]*Viewer{},
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
			viewers.Flush()
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
	V.db.Exec("INSERT INTO emotes (channel, emote) VALUES ($1, $2)", V.channel, e)
}

func (V *ViewerList) GetChannelName() string {
	return V.channel
}

func (V *ViewerList) AddViewer(username string) {
	if _, ok := V.viewers[username]; !ok {
		v := V.FindViewer(username)
		if v == nil {
			// Create user if they don't exist
			v = &Viewer{
				id:         -1,
				Username:   username,
				updated:    true,
				linesTyped: -1,
				money:      -1,
				brawlsWon:  nil,
				manager:    V,
			}
		}
		V.viewers[username] = v
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

func (V *ViewerList) InChannel(username string) (*Viewer, bool) {
	v, ok := V.viewers[username]
	return v, ok
}

func (V *ViewerList) FindViewer(username string) *Viewer {
	row := V.db.QueryRow(`SELECT v.id, c.count FROM (SELECT id, username FROM viewers WHERE channel=$1 AND username=$2) as v `+
		`JOIN counts as c on c.viewer_id=v.id and type='money'`, V.channel, username)

	var id, money int
	err := row.Scan(&id, &money)
	if err == nil {
		return &Viewer{
			id:         id,
			updated:    false,
			Username:   username,
			linesTyped: -1,
			money:      money,
			brawlsWon:  nil,
			manager:    V,
		}
	} else {
		return nil
	}
}

func (V *ViewerList) AllViewers() []*Viewer {
	viewers := []*Viewer{}
	for _, v := range V.viewers {
		viewers = append(viewers, v)
	}
	return viewers
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
	for _, v := range V.viewers {
		if v.updated {
			v.save()
		}
	}
}

func (V *ViewerList) Close() {
	V.Flush()
}
