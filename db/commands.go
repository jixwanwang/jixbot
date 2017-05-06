package db

import (
	"database/sql"
	"time"
)

func (B *dbImpl) GetCommands(channel string) (map[string]bool, error) {
	rows, err := B.db.Query("SELECT command FROM commands WHERE channel=$1", channel)
	if err != nil {
		return nil, err
	}

	allowed := map[string]bool{}
	for rows.Next() {
		var comm string
		if err := rows.Scan(&comm); err == nil {
			allowed[comm] = true
		}
	}

	return allowed, nil
}

func (B *dbImpl) AddCommand(channel, command string) error {
	_, err := B.db.Exec("INSERT INTO commands (channel, command) VALUES ($1, $2)", channel, command)
	return err
}

func (B *dbImpl) DeleteCommand(channel, command string) error {
	_, err := B.db.Exec("DELETE FROM commands WHERE channel=$1 AND command=$2", channel, command)
	return err
}

func (B *dbImpl) GetTextCommands(channel string) ([]TextCommand, error) {
	rows, err := B.db.Query("SELECT command, message, clearance, cooldown FROM textcommands WHERE channel=$1", channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commands := []TextCommand{}
	for rows.Next() {
		var comm, message string
		var clearance, cd int
		rows.Scan(&comm, &message, &clearance, &cd)

		cooldown := time.Duration(cd) * time.Second

		command := TextCommand{
			Clearance: clearance,
			Command:   comm,
			Response:  message,
			Cooldown:  cooldown,
		}
		commands = append(commands, command)
	}
	return commands, nil
}

func (B *dbImpl) AddTextCommand(channel string, comm TextCommand) error {
	_, err := B.db.Exec("INSERT INTO textcommands (channel, command, message, clearance, cooldown) VALUES ($1,$2,$3,$4,$5)",
		channel,
		comm.Command,
		comm.Response,
		comm.Clearance,
		int(comm.Cooldown.Seconds()))
	return err
}

func (B *dbImpl) UpdateTextCommand(channel string, comm TextCommand) error {
	_, err := B.db.Exec("UPDATE textcommands SET message=$1, clearance=$2, cooldown=$3 WHERE channel=$4 AND command=$5",
		comm.Response,
		comm.Clearance,
		int(comm.Cooldown.Seconds()),
		channel,
		comm.Command)
	return err
}

func (B *dbImpl) DeleteTextCommand(channel, comm string) error {
	_, err := B.db.Exec("DELETE FROM textcommands WHERE channel=$1 AND command=$2", channel, comm)
	return err
}

func (B *dbImpl) GetQuote(channel string, rank int) (quote string, quoteRank int, err error) {
	var row *sql.Row
	if rank == 0 {
		row = B.db.QueryRow(`SELECT quote, rank FROM quotes WHERE channel=$1 ORDER BY random() LIMIT 1`, channel)
	} else {
		row = B.db.QueryRow(`SELECT quote, rank FROM quotes WHERE channel=$1 AND rank=$2`, channel, rank)
	}
	err = row.Scan(&quote, &quoteRank)
	return
}

func (B *dbImpl) AllQuotes(channel string) ([]string, error) {
	rows, err := B.db.Query(`SELECT quote, rank FROM quotes WHERE channel=$1 ORDER BY rank ASC`, channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quote := ""
	quotes := []string{}
	for rows.Next() {
		rows.Scan(&quote)
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func (B *dbImpl) AddQuote(channel string, quote string) (rank int, err error) {
	row := B.db.QueryRow(`INSERT INTO quotes(rank, channel, quote) `+
		`SELECT MAX(rank)+1 AS rank, channel, $2 AS quote `+
		`FROM quotes WHERE channel=$1 GROUP BY channel `+
		`UNION ALL SELECT 1, $1, $2 WHERE NOT EXISTS `+
		`(SELECT 1 FROM quotes WHERE channel=$1) `+
		`RETURNING rank`, channel, quote)
	err = row.Scan(&rank)
	return
}
