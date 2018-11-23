package command

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

const entryAmount = 1

type lottery struct {
	cp        *CommandPool
	startComm *subCommand
	endComm   *subCommand
	enterComm *subCommand

	lotteryComm   *subCommand
	seasonComm    *subCommand
	newSeasonComm *subCommand
	pileComm      *subCommand
	statsComm     *subCommand
	statComm      *subCommand

	entries map[string]int
	active  bool
}

func (T *lottery) Init() {
	T.entries = map[string]int{}

	T.startComm = &subCommand{
		command:   "!lottery",
		numArgs:   0,
		cooldown:  5 * time.Second,
		clearance: channel.BROADCASTER,
	}

	T.endComm = &subCommand{
		command:   "!endlottery",
		numArgs:   0,
		cooldown:  1 * time.Second,
		clearance: channel.BROADCASTER,
	}

	T.enterComm = &subCommand{
		command:   "!enter",
		numArgs:   0,
		cooldown:  0,
		clearance: channel.VIEWER,
	}
}

func (T *lottery) ID() string {
	return "lottery"
}

func (T *lottery) endlottery() {
	T.active = false

	users := []string{}

	// add user to user array for every entry they have
	for k, v := range T.entries {
		for i := 0; i < v; i++ {
			users = append(users, k)
		}
	}

	if len(T.entries) <= 1 {
		T.cp.Say("The lottery is over, but not enough people entered to make it interesting")

		// refund users
		for k, v := range T.entries {
			user, in := T.cp.channel.InChannel(k)
			if in {
				user.AddMoney(entryAmount * v)
			}
		}
		return
	}

	winnerIndex := rand.Intn(len(users))
	winner := users[winnerIndex]

	T.cp.Say(fmt.Sprintf("The winner of the lottery is %s! PogChamp They purchased %d tickets. ", winner, entries[winner]))

	T.entries = map[string]int{}
}

func (T *lottery) startlottery() {
	T.active = true

	T.cp.Say(fmt.Sprintf("Hey everyone! @%s has started a lottery. Use !enter <# of tickets> to enter the lottery for your chance to win! Each ticket costs 100 %ss", T.cp.channel.GetChannelName(), T.cp.channel.Currency))
}

func (T *lottery) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	message = strings.TrimSpace(message)
	clearance := T.cp.channel.GetLevel(username)

	_, err := T.startComm.parse(message, clearance)
	if err == nil && T.active == false {
		T.startlottery()
		return
	}

	_, err = T.endComm.parse(message, clearance)
	if err == nil {
		T.endlottery()
		return
	}

	args, err := T.enterComm.parse(message, clearance)
	if err == nil && T.active {
		user, in := T.cp.channel.InChannel(username)
		if !in {
			return
		}

		tickets := 1
		if len(args) > 0 {
			tickets, _ = strconv.Atoi(strings.TrimSpace(args[0]))
			if tickets <= 0 {
				tickets = 1
			}
		}

		if user.GetMoney() < tickets*entryAmount {
			T.cp.Whisper(username, fmt.Sprintf("You don't have enough money to purchase %d tickets. You can purchase up to %d tickets with your money", tickets, user.GetMoney()/entryAmount))
			return
		}

		user.AddMoney(-tickets * entryAmount)

		if val, ok := T.entries[username]; ok {
			T.entries[username] = val + tickets
		} else {
			T.entries[username] = tickets
		}

		T.cp.Whisper(username, fmt.Sprintf("You have purchased %d tickets costing %d %ss", tickets, tickets*entryAmount, T.cp.channel.Currency))

		return
	}

	return
}
