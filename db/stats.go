package db

import (
	"database/sql"
	"fmt"
)

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
	query := fmt.Sprintf(`SELECT c.%s as count, v.username FROM better_counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE v.channel=$1 `+
		`ORDER BY count DESC LIMIT 10`, kind)

	rows, err := B.db.Query(query, channel)
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

func (B *dbImpl) HighestRatio(channel, numerator, denominator string) ([]Ratio, error) {
	query := fmt.Sprintf(`SELECT cast(c.%s AS FLOAT)/c.%s AS ratio, v.username FROM better_counts AS c `+
		`JOIN viewers AS v ON v.id = c.viewer_id `+
		`WHERE v.channel=$1 AND c.%s > 0`+
		`ORDER BY ratio DESC LIMIT 10`, numerator, denominator, denominator)

	rows, err := B.db.Query(query, channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewer string
	var ratio float64
	ratios := []Ratio{}
	for rows.Next() {
		rows.Scan(&ratio, &viewer)

		ratios = append(ratios, Ratio{
			Username: viewer,
			Ratio:    ratio,
		})
	}
	return ratios, nil
}
