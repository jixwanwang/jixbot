package db

type Count struct {
	Username string
	Count    int
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
