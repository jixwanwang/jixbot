package command

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var compliments = []string{
	`%s you look nice today! HeyGuys`,
	`You grace us with your presence, %s`,
	`/me hopes %s is having a great day!`,
	`%s, you are awesome!`,
	`Thanks for watching the stream, %s!`,
}

var responses = []string{
	`Nice to meet you %s, I'm a nice bot MrDestructoid`,
	`What's up %s`,
	`/me waits for %s to say more`,
	`Cool story %s, tell it again Kappa`,
}

type conversation struct {
	cp *CommandPool

	lastResponse time.Time
}

func (T *conversation) Init() {

}

func (T *conversation) ID() string {
	return "conversation"
}

func (T *conversation) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	if time.Since(T.lastResponse).Seconds() < 10 || rand.Intn(2) == 0 {
		return
	}

	if strings.Index(strings.ToLower(message), "jixbot") >= 0 && strings.Index(message, "!") != 0 {
		T.lastResponse = time.Now()
		T.cp.Say(fmt.Sprintf(responses[rand.Intn(len(responses))], username))
	} else if strings.Index(message, "HeyGuys") >= 0 {
		T.lastResponse = time.Now()
		T.cp.Say(fmt.Sprintf(compliments[rand.Intn(len(compliments))], username))
	}
}
