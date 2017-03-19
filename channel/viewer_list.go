package channel

import "github.com/jixwanwang/jixbot/db"

type Level int

const (
	VIEWER      Level = 0
	FOLLOWER    Level = 2
	SUBSCRIBER  Level = 4
	MOD         Level = 6
	ADMIN       Level = 8
	STAFF       Level = 9
	BROADCASTER Level = 10
	GOD         Level = 12
)

// Represents the list of viewers that are in a channel.
type ViewerList struct {
	channel string

	viewers map[string]*Viewer
	staff   map[string]int
	mods    map[string]int
	subs    map[string]int

	db db.DB
}

func NewViewerList(channel string, db db.DB) *ViewerList {
	viewers := &ViewerList{
		channel: channel,
		db:      db,
		viewers: map[string]*Viewer{},
		staff:   map[string]int{},
		mods:    map[string]int{},
		subs:    map[string]int{},
	}

	return viewers
}

func (V *ViewerList) AddViewer(username string) *Viewer {
	if _, ok := V.viewers[username]; !ok {
		v := V.FindViewer(username)
		if v == nil {
			// Create user if they don't exist
			v = &Viewer{
				id:        -1,
				Username:  username,
				updated:   true,
				brawlsWon: nil,
				manager:   V,
			}
		}
		V.viewers[username] = v
	}
	return V.viewers[username]
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

func (V *ViewerList) SetSubscriber(username string) {
	V.AddViewers([]string{username})
	if _, ok := V.viewers[username]; ok {
		V.subs[username] = 1
	}
}

func (V *ViewerList) InChannel(username string) (*Viewer, bool) {
	v, ok := V.viewers[username]
	return v, ok
}

// FindViewer looks up a user in the database. If not found, nil is returned.
func (V *ViewerList) FindViewer(username string) *Viewer {
	viewer := &Viewer{
		id:        -1,
		updated:   false,
		Username:  username,
		brawlsWon: nil,
		manager:   V,
	}

	id, err := V.db.FindViewer(username, V.channel)
	if err != nil {
		return nil
	}
	viewer.id = id

	return viewer
}

func (V *ViewerList) AllViewers() []*Viewer {
	viewers := []*Viewer{}
	for _, v := range V.viewers {
		viewers = append(viewers, v)
	}
	return viewers
}

func (V *ViewerList) GetLevel(username string) Level {
	if username == JIXBOT_CREATOR {
		return GOD
	} else if username == V.channel {
		return BROADCASTER
	} else if _, ok := V.staff[username]; ok {
		return STAFF
	} else if _, ok := V.mods[username]; ok {
		return MOD
	} else if _, ok := V.subs[username]; ok {
		return SUBSCRIBER
	}
	return VIEWER
}

func (V *ViewerList) Flush() {
	for _, v := range V.viewers {
		if v.updated {
			v.save()
		}
	}
}

func (V *ViewerList) Close() {
	V.Flush()
}
