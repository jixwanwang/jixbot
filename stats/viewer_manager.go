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
	// TODO: preload all the viewers whose data is already stored
	// os.MkdirAll(statsFilePath+channel, 0755)

	manager := ViewerManager{
		channel: channel,
		viewers: map[string]*Viewer{},
		db:      db,
	}

	log.Printf("Loading viewers from database...")
	rows, err := db.Query("SELECT id, username FROM viewers WHERE channel=$1", channel)
	if err != nil {
		log.Printf("couldn't read viewers")
	}
	for rows.Next() {
		var id int
		var username string
		rows.Scan(&id, &username)
		manager.viewers[username] = &Viewer{
			id:         id,
			updated:    false,
			manager:    &manager,
			Username:   username,
			linesTyped: -1,
			money:      -1,
			brawlsWon:  nil,
		}
	}
	rows.Close()
	log.Printf("Done loading viewers")

	log.Printf("Retrieving brawl stats...")
	for _, v := range manager.viewers {
		v.GetBrawlsWon()
	}
	log.Printf("Done retrieving brawl stats")

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

	// output := ""
	// for _, v := range V.viewers {
	// 	if v.updated {
	// 		v.save()
	// 	}
	// 	data, _ := json.Marshal(v)
	// 	log.Printf("Saved %v for %s", string(data), V.channel)
	// 	output = output + string(data) + "\n"
	// }

	// ioutil.WriteFile(statsFilePath+V.channel+"_stats", []byte(output), 0666)
}
