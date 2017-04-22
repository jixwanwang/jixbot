package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

const defaultCooldown = 100 * time.Millisecond

type textCommand struct {
	cp        *CommandPool
	clearance channel.Level

	comm *subCommand

	command  string
	response string
	cooldown time.Duration

	numArgs     int
	numReplaces int
	// Map from output argument index (for the response) to input argument index (from the command invocation)
	argMappings  map[int]int
	userMappings map[int]bool

	urlRegex *regexp.Regexp
}

func (T *textCommand) Init() {
	T.ValidateArguments()

	urlRegex, err := regexp.Compile(`\$url:(\ )?((http|https):\/{2})?([0-9a-zA-Z_-]+\.)+[0-9a-zA-Z_-]+` +
		`(\/([~0-9a-zA-Z\#\+\%\.\/\?=_-]+)?)?(\ )?\$`)
	if err != nil {
		log.Printf("url regex parse: %v", err)
	}

	T.urlRegex = urlRegex

	if T.numArgs > 0 {
		T.comm = &subCommand{
			command:   T.command,
			numArgs:   T.numArgs - 1,
			cooldown:  T.cooldown,
			clearance: T.clearance,
		}
	} else {
		T.comm = &subCommand{
			command:   T.command,
			numArgs:   T.numArgs,
			cooldown:  T.cooldown,
			clearance: T.clearance,
		}
	}
}

func (T *textCommand) ValidateArguments() bool {
	regex, err := regexp.Compile(`\$[0-9u]\$`)
	if err != nil {
		log.Fatalf("regex isn't supposed to error: %v", err)
	}
	matches := regex.FindAllString(T.response, -1)

	T.argMappings = map[int]int{}
	T.userMappings = map[int]bool{}
	// Temporary array of arguments to check for consecutive.
	indices := []bool{false, false, false, false, false, false, false, false, false, false}
	T.numReplaces = 0
	for i, arg := range matches {
		matchType := arg[1 : len(arg)-1]
		if matchType == "u" {
			T.userMappings[i] = true
		} else if val, err := strconv.Atoi(matchType); err == nil {
			indices[val] = true
			T.argMappings[i] = val
		}
		T.numReplaces++
	}

	// Check if the argument references are consecutive and start from zero
	ended := false
	T.numArgs = 0
	for i := range indices {
		if !indices[i] {
			ended = true
		} else if indices[i] && ended {
			return false
		} else if indices[i] {
			T.numArgs++
		}
	}

	return true
}

func (T *textCommand) ID() string {
	return "text"
}

func (T *textCommand) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	clearance := T.cp.channel.GetLevel(username)
	args, err := T.comm.parse(message, clearance)
	if err == nil && len(args) >= T.numArgs-1 {
		responseArgs := []interface{}{}
		for i := 0; i < T.numReplaces; i++ {
			if arg, ok := T.argMappings[i]; ok {
				if arg == len(args) {
					responseArgs = append(responseArgs, "")
				} else {
					responseArgs = append(responseArgs, args[arg])
				}
			} else if _, ok := T.userMappings[i]; ok {
				responseArgs = append(responseArgs, username)
			}
		}

		regex, err := regexp.Compile(`\$[0-9u]\$`)
		if err != nil {
			log.Fatalf("regex isn't supposed to error: %v", err)
		}
		resp := regex.ReplaceAllString(T.response, "%s")

		response := fmt.Sprintf(resp, responseArgs...)

		urlMatches := T.urlRegex.FindAllString(response, -1)
		for _, match := range urlMatches {
			url := strings.TrimSpace(strings.TrimPrefix(match[1:len(match)-1], "url:"))

			apiResponse, err := makeAPICall(url)
			if err != nil {
				T.cp.Say(fmt.Sprintf("There was an error with %s: %v", url, err))
				return
			}

			response = strings.Replace(response, match, apiResponse, -1)
		}

		T.cp.Say(response)
	}
}

func makeAPICall(url string) (string, error) {
	if strings.Index(url, "http") < 0 {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.Replace(string(b), "\n", " ", -1), nil
}

func (T *textCommand) String() string {
	level := "viewer"
	switch T.clearance {
	case channel.VIEWER:
		level = "viewer"
	default:
		level = "mod"
	}
	return fmt.Sprintf("%s,%s,%s", T.command, level, T.response)
}
