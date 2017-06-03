package db

func (B *dbImpl) GetAllChannels() ([]string, error) {
	rows, err := B.db.Query("SELECT DISTINCT(channel) FROM commands")
	if err != nil {
		return nil, err
	}

	channels := []string{}
	for rows.Next() {
		var channel string
		err := rows.Scan(&channel)
		if err == nil {
			channels = append(channels, channel)
		}
	}
	rows.Close()

	return channels, nil
}

func (B *dbImpl) GetChannelProperties(channel string) (map[string]string, error) {
	rows, err := B.db.Query("SELECT k, v FROM channel_properties WHERE channel=$1", channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	props := map[string]string{}
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		if err == nil {
			props[k] = v
		}
	}

	return props, nil
}

func (B *dbImpl) SetChannelProperty(channel, k, v string) error {
	insert := "INSERT INTO channel_properties (channel, k, v) SELECT $1, $2, $3"
	upsert := "UPDATE channel_properties SET v=$3 WHERE k=$2 AND channel=$1"
	_, err := B.db.Exec(`WITH upsert AS (`+upsert+` RETURNING *) `+insert+
		` WHERE NOT EXISTS (SELECT * FROM upsert);`,
		channel, k, v)
	return err
}
