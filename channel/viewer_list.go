package channel

import "database/sql"

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

	return viewers
}

func (V *ViewerList) AddViewer(username string) *Viewer {
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
	v, ok := V.viewers[username]
	if ok {
		v.subscriber = true
	}
}

func (V *ViewerList) InChannel(username string) (*Viewer, bool) {
	v, ok := V.viewers[username]
	return v, ok
}

// FindViewer looks up a user in the database. If not found, nil is returned.
func (V *ViewerList) FindViewer(username string) *Viewer {
	row := V.db.QueryRow(`SELECT id FROM viewers WHERE channel=$1 AND username=$2`, V.channel, username)

	viewer := &Viewer{
		id:         -1,
		updated:    false,
		Username:   username,
		linesTyped: -1,
		timeSpent:  -1,
		money:      -1,
		brawlsWon:  nil,
		manager:    V,
	}

	var id int
	err := row.Scan(&id)
	if err == nil {
		viewer.id = id
	} else {
		return nil
	}

	var count int
	var kind string
	rows, err := V.db.Query(`SELECT count, type FROM counts WHERE viewer_id=$1 AND (type='money' OR type='time' OR type='lines_typed')`, viewer.id)
	if err != nil {
		return viewer
	}
	for rows.Next() {
		err := rows.Scan(&count, &kind)
		if err == nil {
			switch {
			case kind == "money":
				viewer.money = count
			case kind == "time":
				viewer.timeSpent = count
			case kind == "lines_typed":
				viewer.linesTyped = count
			}
		}
	}
	rows.Close()

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
	if username == "jixwanwang" {
		return GOD
	} else if username == V.channel {
		return BROADCASTER
	} else if _, ok := V.staff[username]; ok {
		return STAFF
	} else if _, ok := V.mods[username]; ok {
		return MOD
	} else {
		u, ok := V.InChannel(username)
		if ok && u.subscriber {
			return SUBSCRIBER
		}
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
