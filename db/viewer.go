package db

func (B *dbImpl) NewViewer(username, channel string) (id int, err error) {
	// TODO: This double query is bad, should be a QueryRow with a RETURNING
	B.db.Exec("INSERT INTO viewers (username, channel) VALUES ($1, $2)", username, channel)
	row := B.db.QueryRow("SELECT id FROM viewers WHERE username=$1 AND channel=$2", username, channel)
	err = row.Scan(&id)
	return
}

func (B *dbImpl) FindViewer(username, channel string) (id int, err error) {
	row := B.db.QueryRow(`SELECT id FROM viewers WHERE channel=$1 AND username=$2`, channel, username)
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
