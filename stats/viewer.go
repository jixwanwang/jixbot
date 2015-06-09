package stats

import "log"

// Represents a viewer in a single stream.
type Viewer struct {
	id         int
	updated    bool
	manager    *ViewerManager
	Username   string      `json:"username"`
	linesTyped int         `json:"lines_typed"`
	money      int         `json:"money"`
	brawlsWon  map[int]int `json:"brawls_won"`
}

func (V *Viewer) GetBrawlsWon() map[int]int {
	if V.brawlsWon == nil {
		brawlsWon := V.lookupBrawlWins()
		V.brawlsWon = brawlsWon
	}
	return V.brawlsWon
}

func (V *Viewer) GetTotalBrawlsWon() int {
	total := 0
	for _, wins := range V.GetBrawlsWon() {
		total = total + wins
	}
	return total
}

func (V *Viewer) lookupBrawlWins() map[int]int {
	wins := map[int]int{}
	if V.id < 0 {
		return wins
	}

	rows, err := V.manager.db.Query("SELECT season, wins FROM brawlwins WHERE viewer_id=$1", V.id)
	if err != nil {
		log.Printf("couldn't find brawlwins for viewer with id %d", V.id)
		return wins
	}
	for rows.Next() {
		var season, numWins int
		if err := rows.Scan(&season, &numWins); err != nil {
			log.Printf("coudln't scan: %s", err.Error())
		}
		wins[season] = numWins
	}
	rows.Close()
	return wins
}

func (V *Viewer) WinBrawl(season int) {
	V.GetBrawlsWon()

	if _, ok := V.brawlsWon[season]; ok {
		V.brawlsWon[season] = V.brawlsWon[season] + 1
	} else {
		V.brawlsWon[season] = 1
	}
	V.AddMoney(100)

	V.updated = true
}

func (V *Viewer) GetLinesTyped() int {
	if V.linesTyped < 0 {
		if V.id > 0 {
			row := V.manager.db.QueryRow("SELECT count FROM counts WHERE type='lines_typed' AND viewer_id=$1", V.id)

			if err := row.Scan(&V.linesTyped); err != nil {
				log.Printf("couldn't find lines typed for viewer with id %d", V.id)
				V.linesTyped = 0
			}
		} else {
			V.linesTyped = 0
		}
	}
	return V.linesTyped
}

func (V *Viewer) AddLineTyped() {
	V.linesTyped = V.GetLinesTyped() + 1
	V.updated = true
}

func (V *Viewer) GetMoney() int {
	if V.money < 0 {
		if V.id > 0 {
			row := V.manager.db.QueryRow("SELECT count FROM counts WHERE type='money' AND viewer_id=$1", V.id)

			if err := row.Scan(&V.money); err != nil {
				log.Printf("couldn't find money for viewer with id %d", V.id)
				V.money = 0
			}
		} else {
			V.money = 0
		}
	}
	return V.money
}

func (V *Viewer) AddMoney(amount int) {
	V.money = V.GetMoney() + amount
	if V.money < 0 {
		V.money = 0
	}
	V.updated = true
}

func (V *Viewer) save() {
	if V.id == -1 {
		// Write new user to db
		V.manager.db.Exec("INSERT INTO viewers (username, channel) VALUES ($1, $2)", V.Username, V.manager.channel)
		row := V.manager.db.QueryRow("SELECT id FROM viewers WHERE username=$1 AND channel=$2", V.Username, V.manager.channel)
		var id int
		if err := row.Scan(&id); err != nil {
			log.Printf("failed to get id of new user: %s", err.Error())
		}
		V.id = id
		log.Printf("created new viewer with id %d, username %s", V.id, V.Username)
	}
	// TODO: write brawls, money, lines typed to db
	if V.brawlsWon != nil {
		for season, wins := range V.brawlsWon {
			insert := "INSERT INTO brawlwins (season, viewer_id, wins) SELECT $1, $2, $3"
			upsert := "UPDATE brawlwins SET wins=$3 WHERE season=$1 AND viewer_id=$2"
			V.manager.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", season, V.id, wins)
		}
	}
	if V.money > 0 {
		insert := "INSERT INTO counts (type, viewer_id, count) SELECT 'money', $1, $2"
		upsert := "UPDATE counts SET count=$2 WHERE type='money' AND viewer_id=$1"
		V.manager.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", V.id, V.money)
	}
	if V.linesTyped > 0 {
		insert := "INSERT INTO counts (type, viewer_id, count) SELECT 'lines_typed', $1, $2"
		upsert := "UPDATE counts SET count=$2 WHERE type='lines_typed' AND viewer_id=$1"
		query := "WITH upsert AS (" + upsert + " RETURNING *) " + insert + " WHERE NOT EXISTS (SELECT * FROM upsert);"
		V.manager.db.Exec(query, V.id, V.linesTyped)
	}
	V.updated = false
	log.Printf("saved user %s", V.Username)
}
