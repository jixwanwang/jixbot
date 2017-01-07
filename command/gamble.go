package command

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type gamble struct {
	cp *CommandPool

	gamble *subCommand
	bet    *subCommand

	betters map[string]int
	active  bool
}

func (T *gamble) Init() {
	// username to bet mapping
	T.betters = map[string]int{}

	T.gamble = &subCommand{
		command:   "!gamble",
		numArgs:   0,
		cooldown:  20 * time.Minute,
		clearance: channel.MOD,
	}

	T.bet = &subCommand{
		command:   "!bet",
		numArgs:   1,
		cooldown:  0,
		clearance: channel.VIEWER,
	}
}

func (T *gamble) ID() string {
	return "gamble"
}

func (T *gamble) endBetting() {
	T.active = false

	users := []string{}

	for k := range T.betters {
		users = append(users, k)
	}

	if len(users) <= 0 {
		return
	}

	if len(users) == 1 {
		T.cp.Say(fmt.Sprintf("The betting is over, but %s was the only one betting. LAAAAAME.", users[0]))
		return
	}

	winnerIndex := rand.Intn(len(users))
	winner := users[winnerIndex]
	bet := T.betters[winner]
	winnings := int(.75 * float64(bet))

	message := fmt.Sprintf("Betting is over, %v wins the jackpot! They bet %v and won %v %vs! ", winner, bet, bet+winnings, T.cp.channel.Currency)

	user, in := T.cp.channel.InChannel(winner)
	if in {
		user.AddMoney(winnings + bet)
	} else {
		message = message + "Sadly they left the chat so they don't get any cash FailFish"
	}
	T.cp.Say(message)

	T.betters = map[string]int{}
}

func (T *gamble) startBetting() {
	T.active = true

	duration := rand.Intn(60) + 120
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	go func() {
		<-timer.C

		T.endBetting()
	}()

	T.cp.Say(fmt.Sprintf("A betting round has started! Type !bet to add to the pool for a chance to win the jackpot!"))
}

func (T *gamble) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	message = strings.TrimSpace(message)
	clearance := T.cp.channel.GetLevel(username)

	_, err := T.gamble.parse(message, clearance)
	if err == nil && T.active == false {
		T.startBetting()
		return
	}

	args, err := T.bet.parse(message, clearance)
	if err == nil && T.active == true {
		user, in := T.cp.channel.InChannel(username)
		if !in {
			return
		}

		bet, _ := strconv.Atoi(args[0])
		user.AddMoney(-bet)

		if len(args) > 0 {
			T.betters[username] = bet
		}

		return
	}

	return
}
