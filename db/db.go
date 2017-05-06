package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DB interface {
	GetAllChannels() ([]string, error)

	GetChannelProperties(channel string) (map[string]string, error)
	SetChannelProperty(channel, k, v string) error

	GetChannelEmotes(channel string) ([]string, error)
	AddChannelEmote(channel, emote string) error
	DeleteChannelEmote(channel, emote string) error

	GetCommands(channel string) (map[string]bool, error)
	AddCommand(channel, command string) error
	DeleteCommand(channel, command string) error

	GetTextCommands(channel string) ([]TextCommand, error)
	AddTextCommand(channel string, comm TextCommand) error
	UpdateTextCommand(channel string, comm TextCommand) error
	DeleteTextCommand(channel, comm string) error

	NewViewer(username, channel string) (id int, err error)
	FindViewer(username, channel string) (id int, err error)

	GetCounts(viewerID int) (*Counts, error)
	SetCounts(counts *Counts) error

	GetCount(viewerID int, kind string) (count int, err error)
	SetCount(viewerID int, kind string, count int) error
	HighestCount(channel, kind string) ([]Count, error)

	GetBrawlWins(viewerID int) (map[int]int, error)
	SetBrawlWins(viewerID int, channel string, wins map[int]int) error
	GetBrawlSeason(channel string) (season int, err error)
	BrawlStats(channel string, season int) ([]Count, error)

	RetrieveQuestionAnswers(channel string) ([]QuestionAnswer, error)
	AddQuestionAnswer(channel, question, answer string) (QuestionAnswer, error)
	UpdateQuestionAnswer(qa QuestionAnswer) error

	GetQuote(channel string, rank int) (string, int, error)
	AllQuotes(channel string) ([]string, error)
	AddQuote(channel, quote string) (int, error)
}

type dbImpl struct {
	db *sql.DB
}

var _ DB = new(dbImpl)

func New(host, port, name, user, password string) (DB, error) {
	pgConnect := fmt.Sprintf("dbname=%s user=%s host=%s port=%s",
		name, user, host, port)
	if password != "" {
		pgConnect = fmt.Sprintf("%s password=%s", pgConnect, password)
	} else {
		pgConnect = pgConnect + " sslmode=disable"
	}

	db, err := sql.Open("postgres", pgConnect)

	if err != nil {
		log.Printf("couldn't connect to db: %s", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("couldn't ping db: %s", err.Error())
		return nil, err
	}

	return &dbImpl{db: db}, nil
}
