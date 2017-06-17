package command

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
)

type brawl struct {
	cp            *CommandPool
	brawlComm     *subCommand
	seasonComm    *subCommand
	newSeasonComm *subCommand
	pileComm      *subCommand
	statsComm     *subCommand
	statComm      *subCommand

	season   int
	brawlers map[string]string
	betters  map[string]int
	active   bool
}

func (T *brawl) Init() {
	season, err := T.cp.db.GetBrawlSeason(T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("couldn't determine the brawl season, assuming to be 1")
		T.season = 1
	} else {
		T.season = season
	}

	// username to weapon mapping
	T.brawlers = map[string]string{}
	T.betters = map[string]int{}

	T.brawlComm = &subCommand{
		command:   "!brawl",
		numArgs:   0,
		cooldown:  15 * time.Minute,
		clearance: channel.MOD,
	}

	T.seasonComm = &subCommand{
		command:   "!brawlseason",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.VIEWER,
	}

	T.newSeasonComm = &subCommand{
		command:   "!newbrawlseason",
		numArgs:   0,
		cooldown:  1 * time.Hour,
		clearance: channel.BROADCASTER,
	}

	T.pileComm = &subCommand{
		command:   "!pileon",
		numArgs:   0,
		cooldown:  0,
		clearance: channel.VIEWER,
	}

	T.statsComm = &subCommand{
		command:   "!brawlstats",
		numArgs:   0,
		cooldown:  30 * time.Second,
		clearance: channel.VIEWER,
	}

	T.statComm = &subCommand{
		command:   "!brawlwins",
		numArgs:   0,
		cooldown:  200 * time.Millisecond,
		clearance: channel.VIEWER,
	}
}

func (T *brawl) ID() string {
	return "brawl"
}

func (T *brawl) endBrawl() {
	T.active = false

	users := []string{}

	for k := range T.brawlers {
		users = append(users, k)
	}

	if len(users) <= 0 {
		return
	}

	if len(users) == 1 {
		T.cp.Say(fmt.Sprintf("The brawl is over, but %s was the only one fighting. That was boring.", users[0]))
		// refund bet
		user, in := T.cp.channel.InChannel(users[0])
		if in {
			user.AddMoney(T.betters[users[0]])
		}
		return
	} else if len(users) < T.cp.channel.MinBrawlers {
		T.cp.Say(fmt.Sprintf("Only a few people joined the brawl, while others just sat around and watched. That was really boring."))
		// refund bets
		for _, u := range users {
			user, in := T.cp.channel.InChannel(u)
			if in {
				user.AddMoney(T.betters[u])
			}
		}
		return
	}

	winnerIndex := rand.Intn(len(users))
	winner := users[winnerIndex]

	// Default winnings for no bet
	winnings := 500

	// If everyone bets, the expected winnings is same as the bet, so no one makes any money.
	// However if not everyone bets, the expected winnings is less than the bet. Gambling always causes a loss ;P
	if bet, ok := T.betters[winner]; ok {
		winnings = int(float64(bet*len(T.brawlers)) * 0.9)
	}

	weapon := T.brawlers[winner]

	if len(weapon) > 0 {
		var message string
		message = fmt.Sprintf(T.cp.channel.BrawlEndMessageWithWeapon, winner, weapon, winnings, T.cp.channel.Currency)
		T.cp.Say(message)
	} else {
		message := fmt.Sprintf(T.cp.channel.BrawlEndMessageNoWeapon, winner, winnings, T.cp.channel.Currency)
		T.cp.Say(message)
	}

	T.brawlers = map[string]string{}
	T.betters = map[string]int{}

	winningUser, in := T.cp.channel.InChannel(winner)
	if in {
		winningUser.AddMoney(winnings)
		winningUser.WinBrawl(T.season)
	}

	T.cp.channel.Flush()
}

func (T *brawl) startBrawl() {
	T.active = true

	duration := rand.Intn(120) + 60
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	go func() {
		<-timer.C

		T.endBrawl()
	}()

	T.cp.Say(T.cp.channel.BrawlStartMessage)
}

func (T *brawl) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	message = strings.TrimSpace(message)
	clearance := T.cp.channel.GetLevel(username)

	_, err := T.brawlComm.parse(message, clearance)
	if err == nil && T.active == false {
		T.startBrawl()
		return
	}

	_, err = T.seasonComm.parse(message, clearance)
	if err == nil {
		T.cp.Say(fmt.Sprintf("The current brawl season is %d", T.season))
		return
	}

	_, err = T.newSeasonComm.parse(message, clearance)
	if err == nil && T.active == false {
		T.cp.Say(T.calculateBrawlStats(T.season))
		T.season = T.season + 1
		T.cp.Say(fmt.Sprintf("The brawl season has ended! We are now in season %d.", T.season))

		// allow brawling right away to persist the new season
		T.brawlComm.lastCalled = time.Now().Add(-1 * time.Hour)
		return
	}

	args, err := T.pileComm.parse(message, clearance)
	if err == nil && T.active == true {
		user, in := T.cp.channel.InChannel(username)
		if !in {
			return
		}

		if len(args) > 0 {
			weapon := args[0]
			firstArg := strings.Split(weapon, " ")[0]
			if strings.Index(firstArg, "bet=") == 0 {
				bet, _ := strconv.Atoi(strings.TrimPrefix(firstArg, "bet="))
				weapon = strings.TrimSpace(strings.TrimPrefix(weapon, firstArg))
				if user.GetMoney() >= bet && bet > 0 {
					user.AddMoney(-bet)
					T.betters[username] = T.betters[username] + bet
				}
			}

			T.brawlers[username] = weapon
		} else {
			T.brawlers[username] = ""
		}

		return
	}

	args, err = T.statComm.parse(message, clearance)
	if err == nil {
		user, in := T.cp.channel.InChannel(username)
		if in {
			season := T.season
			if len(args) > 0 {
				s, err := strconv.Atoi(args[0])
				if err == nil {
					season = s
				}
			}

			brawlsWon := user.GetBrawlsWon()
			wins := brawlsWon[season]
			if season == T.season {
				T.cp.Say(fmt.Sprintf("@%s you have won %d brawls this season", username, wins))
			} else {
				T.cp.Say(fmt.Sprintf("@%s you won %d brawls in season %v", username, wins, season))
			}
			return
		}
	}

	args, err = T.statsComm.parse(message, clearance)
	if err == nil {
		season := T.season
		if len(args) > 0 {
			if args[0] == "all" {
				season = 0
			}

			s, err := strconv.Atoi(args[0])
			if err == nil {
				season = s
			}
		}

		T.cp.Say(T.calculateBrawlStats(season))
		return
	}

	return
}

type brawlWin struct {
	username string
	wins     int
}

func (T *brawl) calculateBrawlStats(season int) string {
	output := ""
	brawlWins, err := T.cp.db.BrawlStats(T.cp.channel.GetChannelName(), season)
	if err != nil {
		return ""
	}

	if len(brawlWins) == 0 {
		return fmt.Sprintf("No one has won for season %d yet!", season)
	}

	tiedWinners := []db.Count{}
	numWins := 10000
	count := -1
	for _, w := range brawlWins {
		if w.Count < numWins {
			if len(tiedWinners) != 0 {
				for _, winner := range tiedWinners {
					output = fmt.Sprintf("%s%s & ", output, winner.Username)
				}
				output = fmt.Sprintf("%s (%d wins), ", output[:len(output)-2], numWins)
			}

			tiedWinners = []db.Count{}
			numWins = w.Count
			count = count + 1
		}

		// Don't display 0 wins, only keep top 5 win groups, and ignore groups with more than 7 people in it.
		if numWins == 0 || count > 4 || len(tiedWinners) > 7 {
			break
		}

		if w.Count == numWins {
			tiedWinners = append(tiedWinners, w)
		}
	}

	if len(tiedWinners) != 0 && len(tiedWinners) <= 7 {
		for _, winner := range tiedWinners {
			output = fmt.Sprintf("%s%s & ", output, winner.Username)
		}
		output = fmt.Sprintf("%s (%d wins), ", output[:len(output)-2], numWins)
	}

	return output[:len(output)-2]
}
