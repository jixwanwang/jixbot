// ViewerManager is responsible for handling the viewers stats in the database.

package stats

import (
	"database/sql"
	"log"
)

const statsFilePath = "data/stats/"

type dbViewer struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
}

type ViewerManager struct {
	channel string
	viewers map[string]*Viewer
	db      *sql.DB
}

func Init(channel string, db *sql.DB) *ViewerManager {
	manager := ViewerManager{
		channel: channel,
		viewers: map[string]*Viewer{},
		db:      db,
	}

	log.Printf("Loading viewers from database...")
	rows, err := db.Query(`SELECT v.id, v.username, c.count FROM (SELECT id, username FROM viewers WHERE channel=$1) as v `+
		// `JOIN brawlwins as b on b.viewer_id=v.id `+
		`JOIN counts as c on c.viewer_id=v.id and type='money'`, channel)
	if err != nil {
		log.Printf("couldn't read viewers")
	}
	for rows.Next() {
		var id, money int
		var username string
		rows.Scan(&id, &username, &money)
		manager.viewers[username] = &Viewer{
			id:         id,
			updated:    false,
			manager:    &manager,
			Username:   username,
			linesTyped: -1,
			money:      money,
			brawlsWon:  nil,
		}
	}
	rows.Close()
	log.Printf("Done loading viewers")

	go func() {
		log.Printf("Retrieving brawl stats...")
		for _, v := range manager.viewers {
			v.GetBrawlsWon()
		}
		log.Printf("Done retrieving brawl stats")
	}()

	go func() {
		log.Printf("Retrieving money...")
		for _, v := range manager.viewers {
			v.GetMoney()
		}
		log.Printf("Done retrieving money")
	}()

	return &manager
}

func (V *ViewerManager) AllViewers() []*Viewer {
	viewers := []*Viewer{}
	for _, v := range V.viewers {
		viewers = append(viewers, v)
	}
	return viewers
}

func (V *ViewerManager) FindViewer(username string) *Viewer {
	viewer, ok := V.viewers[username]

	if !ok {
		V.viewers[username] = &Viewer{
			id:         -1,
			Username:   username,
			updated:    true,
			manager:    V,
			linesTyped: -1,
			money:      -1,
			brawlsWon:  nil,
		}
		return V.viewers[username]
	}

	return viewer
}

func (V *ViewerManager) FindViewers(usernames []string) []*Viewer {
	v := []*Viewer{}

	for _, u := range usernames {
		v = append(v, V.FindViewer(u))
	}

	return v
}

func (V *ViewerManager) Flush() {
	for _, v := range V.viewers {
		if v.updated {
			v.save()
		}
	}
}
