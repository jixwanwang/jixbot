package db

import "database/sql"

type Count struct {
	Username string
	Count    int
}

func (B *dbImpl) BrawlStats(channel string, season int) ([]Count, error) {
	var rows *sql.Rows
	var err error
	if season > 0 {
		rows, err = B.db.Query(`SELECT sum(wins) totalwins, username FROM brawlwins AS b `+
			`JOIN viewers AS v ON v.id=b.viewer_id `+
			`WHERE b.channel=$1 AND b.season=$2 `+
			`GROUP BY username ORDER BY totalwins DESC LIMIT 50`, channel, season)
	} else {
		rows, err = B.db.Query(`SELECT sum(wins) totalwins, username FROM brawlwins AS b `+
			`JOIN viewers AS v ON v.id=b.viewer_id `+
			`WHERE b.channel=$1 `+
			`GROUP BY username ORDER BY totalwins DESC LIMIT 50`, channel)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	brawlWins := []Count{}
	var username string
	var wins int
	for rows.Next() {
		err := rows.Scan(&wins, &username)
		if err != nil {
			continue
		}

		brawlWins = append(brawlWins, Count{
			Username: username,
			Count:    wins,
		})
	}

	return brawlWins, nil
}

func (B *dbImpl) HighestCount(channel, kind string) ([]Count, error) {
	rows, err := B.db.Query(`SELECT sum(c.count) as lines, v.username FROM counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE c.type=$2 AND v.channel=$1 `+
		`GROUP BY v.username ORDER BY lines DESC LIMIT 10`, channel, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewer string
	var count int
	counts := []Count{}
	for rows.Next() {
		rows.Scan(&count, &viewer)

		counts = append(counts, Count{
			Username: viewer,
			Count:    count,
		})
	}
	return counts, nil
}
