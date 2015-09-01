package command

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

const (
	commandFilePath = "data/textcommands/"
	configFilePath  = "data/config/"
	globalChannel   = "_global"
)

type Command interface {
	ID() string
	WhisperOnly() bool
	Init()
	Response(username, message string) string
	String() string
}

var BadCommandError = fmt.Errorf("Bad command")
var NotPermittedError = fmt.Errorf("Not permitted to use this command")

type subCommand struct {
	command    string
	numArgs    int
	cooldown   time.Duration
	lastCalled time.Time
	clearance  channel.Level
}

func (C *subCommand) parse(message string, clearance channel.Level) ([]string, error) {
	args := []string{}

	if strings.Index(strings.ToLower(message), C.command) >= 0 {
		log.Printf("%s called with clearance %s", C.command, clearance)
	}

	if clearance < C.clearance {
		return args, NotPermittedError
	}

	// Rate limit
	if C.cooldown.Nanoseconds() > 0 &&
		time.Since(C.lastCalled).Nanoseconds() < C.cooldown.Nanoseconds() {
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
			args = append(args, strings.TrimSpace(strings.TrimPrefix(message, prefix)))
		}
		C.lastCalled = time.Now()
		return args, nil
	}

	return args, BadCommandError
}
