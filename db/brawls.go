package db

import "log"

func (B *dbImpl) GetBrawlSeason(channel string) (season int, err error) {
	row := B.db.QueryRow("SELECT * FROM (SELECT DISTINCT(season) FROM brawlwins WHERE channel=$1 ORDER BY season DESC) AS seasons LIMIT 1", channel)
	err = row.Scan(&season)
	return
}

func (B *dbImpl) GetBrawlWins(viewerID int) (map[int]int, error) {
	rows, err := B.db.Query("SELECT season, wins FROM brawlwins WHERE viewer_id=$1", viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wins := map[int]int{}
	for rows.Next() {
		var season, numWins int
		if err := rows.Scan(&season, &numWins); err != nil {
			log.Printf("couldn't scan brawl win: %s", err.Error())
		}
		wins[season] = numWins
	}

	return wins, nil
}

func (B *dbImpl) SetBrawlWins(viewerID int, channel string, wins map[int]int) error {
	for season, wins := range wins {
		insert := "INSERT INTO brawlwins (season, viewer_id, wins, channel) SELECT $1, $2, $3, $4"
		upsert := "UPDATE brawlwins SET wins=$3 WHERE season=$1 AND viewer_id=$2 AND channel=$4"
		_, err := B.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", season, viewerID, wins, channel)
		if err != nil {
			return err
		}
	}
	return nil
}
