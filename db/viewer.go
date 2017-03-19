package db

func (B *dbImpl) NewViewer(username, channel string) (id int, err error) {
	row := B.db.QueryRow("INSERT INTO viewers (username, channel) VALUES ($1, $2) RETURNING id", username, channel)
	err = row.Scan(&id)
	return
}

func (B *dbImpl) FindViewer(username, channel string) (id int, err error) {
	row := B.db.QueryRow(`SELECT id FROM viewers WHERE channel=$1 AND username=$2 ORDER BY id ASC LIMIT 1`, channel, username)
	err = row.Scan(&id)
	return
}

func (B *dbImpl) GetCount(viewerID int, kind string) (count int, err error) {
	row := B.db.QueryRow("SELECT count FROM counts WHERE type=$2 AND viewer_id=$1", viewerID, kind)
	err = row.Scan(&count)
	return
}

func (B *dbImpl) SetCount(viewerID int, kind string, count int) error {
	insert := "INSERT INTO counts (type, viewer_id, count) SELECT $2, $1, $3"
	upsert := "UPDATE counts SET count=$3 WHERE type=$2 AND viewer_id=$1"
	_, err := B.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", viewerID, kind, count)
	return err
}

func (B *dbImpl) GetCounts(viewerID int) (*Counts, error) {
	row := B.db.QueryRow("SELECT money, lines_typed, time_spent FROM better_counts WHERE viewer_id=$1", viewerID)
	counts := &Counts{
		ViewerID: viewerID,
	}
	err := row.Scan(&counts.Money, &counts.LinesTyped, &counts.TimeSpent)
	return counts, err
}

func (B *dbImpl) SetCounts(counts *Counts) error {
	insert := "INSERT INTO better_counts (viewer_id, money, lines_typed, time_spent) SELECT $1, $2, $3, $4"
	upsert := "UPDATE better_counts SET money=$2, lines_typed=$3, time_spent=$4 WHERE viewer_id=$1"
	_, err := B.db.Exec("WITH upsert AS ("+upsert+" RETURNING *) "+insert+" WHERE NOT EXISTS (SELECT * FROM upsert);", counts.ViewerID, counts.Money, counts.LinesTyped, counts.TimeSpent)
	return err
}
