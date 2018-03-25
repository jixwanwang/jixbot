package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
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

	isFancy bool

	urlRegex  *regexp.Regexp
	randRegex *regexp.Regexp
}

func (T *textCommand) Init() {
	T.ValidateArguments()

	urlRegex, err := regexp.Compile(`\$url:(\ )?((http|https):\/{2})?([0-9a-zA-Z_-]+\.)+.*(\ )?\$`)
	if err != nil {
		log.Printf("url regex parse: %v", err)
	}

	T.urlRegex = urlRegex

	randRegex, err := regexp.Compile(`\$rand:[0-9]+-[0-9]+\$`)
	if err != nil {
		log.Printf("rand regex parse: %v", err)
	}
	T.randRegex = randRegex

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

	if len(T.cp.channel.Emotes) > 0 {
		for _, emote := range T.cp.channel.Emotes {
			if strings.Index(T.response, emote) >= 0 {
				T.isFancy = !T.cp.channel.BotIsSubbed
				return
			}
		}
	}
}

func (T *textCommand) ValidateArguments() bool {
	// Don't allow commands that let the user choose a prefix
	// IE "!badcomm $0$ this command is bad" could become "!badcomm /ban someinnocentguy this command is bad"
	prefixinput, err := regexp.Compile(`^\$[0-9]\$`)
	if err != nil {
		log.Fatalf("regex isn't supposed to error: %v", err)
	}
	if prefixinput.Find([]byte(T.response)) != nil {
		return false
	}

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
			uri := strings.TrimSpace(strings.TrimPrefix(match[1:len(match)-1], "url:"))
			raw, err := url.Parse(uri)
			if err != nil {
				T.cp.Say(fmt.Sprintf("Malformed url: %s", uri))
				return
			}

			if i := strings.Index(uri, "?"); i >= 0 {
				uri = uri[:i+1] + raw.Query().Encode()
			}

			apiResponse, err := makeAPICall(uri)
			if err != nil {
				T.cp.Say(fmt.Sprintf("There was an error with %s: %v", uri, err))
				return
			}

			response = strings.Replace(response, match, apiResponse, -1)
		}

		randMatches := T.randRegex.FindAllString(response, -1)
		for _, match := range randMatches {
			randRange := strings.TrimSpace(strings.TrimPrefix(match[1:len(match)-1], "rand:"))
			lower, _ := strconv.Atoi(randRange[:strings.Index(randRange, "-")])
			upper, _ := strconv.Atoi(randRange[strings.Index(randRange, "-")+1:])
			if upper-lower < 0 {
				response = strings.Replace(response, match, "", -1)
				continue
			}

			// make upper inclusive
			random := rand.Intn(upper-lower+1) + lower
			response = strings.Replace(response, match, strconv.Itoa(random), -1)
		}

		if T.isFancy {
			T.cp.FancySay(response)
		} else {
			T.cp.Say(response)
		}
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
