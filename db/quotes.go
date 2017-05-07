package db

import (
	"database/sql"
	"fmt"
)

type Quote struct {
	Quote string
	Kind  string
	Rank  int
}

func (B *dbImpl) GetQuote(channel, kind string, rank int) (quote string, quoteRank int, err error) {
	var row *sql.Row
	if rank == 0 {
		row = B.db.QueryRow(`SELECT quote, rank FROM quotes WHERE channel=$1 AND quote_type=$2 ORDER BY random() LIMIT 1`, channel, kind)
	} else {
		row = B.db.QueryRow(`SELECT quote, rank FROM quotes WHERE channel=$1 AND quote_type=$2 AND rank=$3`, channel, kind, rank)
	}
	err = row.Scan(&quote, &quoteRank)
	return
}

func (B *dbImpl) SearchQuote(channel, kind, term string) (quote string, quoteRank int, err error) {
	row := B.db.QueryRow(`SELECT quote, rank FROM quotes WHERE channel=$1 AND quote_type=$2 `+
		`AND quote ILIKE $3 ORDER BY random() LIMIT 1`, channel, kind, fmt.Sprintf(`%s%v%s`, "%", term, "%"))
	err = row.Scan(&quote, &quoteRank)
	return
}

func (B *dbImpl) AllQuotes(channel, kind string) ([]Quote, error) {
	rows, err := B.db.Query(`SELECT quote, rank FROM quotes WHERE channel=$1 AND quote_type=$2 ORDER BY rank ASC`, channel, kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quote := ""
	rank := 0
	quotes := []Quote{}
	for rows.Next() {
		rows.Scan(&quote, &rank)
		quotes = append(quotes, Quote{
			Quote: quote,
			Kind:  kind,
			Rank:  rank,
		})
	}
	return quotes, nil
}

func (B *dbImpl) AddQuote(channel, kind string, quote string) (rank int, err error) {
	row := B.db.QueryRow(`INSERT INTO quotes(channel, quote_type, rank, quote) `+
		`SELECT channel, $2, MAX(rank)+1 AS rank, $3 AS quote `+
		`FROM quotes WHERE channel=$1 AND quote_type=$2 GROUP BY channel `+
		`UNION ALL SELECT $1, $2, 1, $3 WHERE NOT EXISTS `+
		`(SELECT 1 FROM quotes WHERE channel=$1) `+
		`RETURNING rank`, channel, kind, quote)
	err = row.Scan(&rank)
	return
}

func (B *dbImpl) DeleteQuote(channel, kind string, rank int) (quoteRank int, err error) {
	row := B.db.QueryRow(`DELETE FROM quotes WHERE channel=$1 AND quote_type=$2 AND rank=$3 RETURNING rank`, channel, kind, rank)
	err = row.Scan(&quoteRank)
	return
}
