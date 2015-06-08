package command

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/irc"
	"github.com/jixwanwang/jixbot/messaging"
)

const (
	commandFilePath = "data/textcommands/"
	configFilePath  = "data/config/"
	globalChannel   = "_global"
)

type Command interface {
	ID() string
	Init()
	Response(username, message string) string
	String() string
}

func NewCommandPool(channel *channel.ViewerList, broadcaster *channel.Broadcaster, irc *irc.Client, texter messaging.Texter, db *sql.DB) *CommandPool {
	cp := &CommandPool{
		channel:     channel,
		broadcaster: broadcaster,
		irc:         irc,
		db:          db,
		texter:      texter,
	}

	globals := loadTextCommands(db, globalChannel)
	channels := loadTextCommands(db, channel.GetChannelName())

	specials := cp.specialCommands()
	filtered := filterCommands(db, channel.GetChannelName(), specials)

	cp.specials = filtered
	for _, c := range cp.specials {
		c.Init()
	}
	cp.commands = channels
	cp.globalcommands = globals

	return cp
}

func loadTextCommands(db *sql.DB, channelName string) []*textCommand {
	commands := []*textCommand{}

	rows, err := db.Query("SELECT command, message, clearance FROM textcommands WHERE channel=$1", channelName)
	if err != nil {
		log.Printf("Couldn't read text commands")
	}
	for rows.Next() {
		var comm, message string
		var clearance int
		rows.Scan(&comm, &message, &clearance)

		commands = append(commands, &textCommand{
			clearance: channel.Level(clearance),
			command:   comm,
			response:  message,
		})
	}
	rows.Close()

	return commands
}

func filterCommands(db *sql.DB, channelName string, commands []Command) []Command {
	newcommands := []Command{}

	rows, err := db.Query("SELECT command FROM commands WHERE channel=$1", channelName)
	if err != nil {
		log.Printf("Couldn't read commands")
	}

	ids := map[string]int{}
	for rows.Next() {
		var comm string
		if err := rows.Scan(&comm); err == nil {
			ids[comm] = 1
		}
	}

	for _, c := range commands {
		if _, ok := ids[c.ID()]; ok {
			newcommands = append(newcommands, c)
		}
	}

	return newcommands
}

var BadCommandError = fmt.Errorf("Bad command")

type subCommand struct {
	command    string
	numArgs    int
	cooldown   time.Duration
	lastCalled time.Time
}

func (C *subCommand) parse(message string) ([]string, error) {
	args := []string{}

	// Rate limit
	if C.cooldown.Nanoseconds() > 0 && time.Since(C.lastCalled).Nanoseconds() < C.cooldown.Nanoseconds() {
		return args, BadCommandError
	}

	parts := strings.Split(message, " ")
	if parts[0] != C.command {
		return args, BadCommandError
	}

	if len(parts)-1 < C.numArgs {
		return args, BadCommandError
	}

	prefix := strings.Join(parts[:C.numArgs+1], " ")
	// There will be a trailing space on the prefix if the message has more than enough parts
	if C.numArgs < len(parts)-1 {
		prefix = prefix + " "
	}

	if strings.HasPrefix(message, prefix) {
		args = parts[1 : C.numArgs+1]
		remaining := strings.TrimPrefix(message, prefix)
		if len(remaining) > 0 {
			args = append(args, strings.TrimPrefix(message, prefix))
		}
		C.lastCalled = time.Now()
		return args, nil
	}

	return args, BadCommandError
}
