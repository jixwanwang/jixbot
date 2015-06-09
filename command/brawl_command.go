package command

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type brawlCommand struct {
	cp               *CommandPool
	brawlComm        *subCommand
	newSeasonComm    *subCommand
	pileComm         *subCommand
	statsComm        *subCommand
	alltimeStatsComm *subCommand
	statComm         *subCommand

	season   int
	brawlers map[string]int
	active   bool
}

func (T *brawlCommand) Init() {
	row := T.cp.db.QueryRow("select * from (select distinct(season) from brawlwins order by season desc) as seasons limit 1")
	if err := row.Scan(&T.season); err != nil {
		log.Printf("couldn't determine the brawl season, assuming to be 1")
		T.season = 1
	}
	log.Printf("The current brawl season is %d", T.season)

	T.brawlers = map[string]int{}

	T.brawlComm = &subCommand{
		command:  "!brawl",
		numArgs:  0,
		cooldown: 15 * time.Minute,
	}

	T.newSeasonComm = &subCommand{
		command:  "!newbrawlseason",
		numArgs:  0,
		cooldown: 12 * time.Second,
	}

	T.pileComm = &subCommand{
		command:  "!pileon",
		numArgs:  0,
		cooldown: 0,
	}

	T.statsComm = &subCommand{
		command:  "!brawlstats",
		numArgs:  1,
		cooldown: 15 * time.Second,
	}

	T.statComm = &subCommand{
		command:  "!brawlwins",
		numArgs:  0,
		cooldown: 5 * time.Second,
	}
}

func (T *brawlCommand) ID() string {
	return "brawl"
}

func (T *brawlCommand) endBrawl() {
	T.active = false
	log.Printf("Brawl ended")

	users := []string{}

	for k, v := range T.brawlers {
		for i := 0; i < v; i++ {
			users = append(users, k)
		}
	}

	T.brawlers = map[string]int{}

	if len(users) <= 0 {
		return
	}

	if len(users) == 1 {
		T.cp.irc.Say("#"+T.cp.channel.GetChannelName(), fmt.Sprintf("The brawl is over, but %s was the only one fighting. That was a boring brawl.", users[0]))
		return
	}

	winner := users[rand.Intn(len(users))]
	// 10% chance of broadcaster winning brawl if they join
	if _, ok := T.brawlers[T.cp.channel.GetChannelName()]; ok && rand.Intn(10) == 1 {
		winner = T.cp.channel.GetChannelName()
	}

	message := fmt.Sprintf("The brawl is over, the tavern is a mess, but @%s is the last one standing! They loot 100 %ss from the losers.", winner, T.cp.channel.Currency)

	T.cp.irc.Say("#"+T.cp.channel.GetChannelName(), message)

	winningUser, in := T.cp.channel.InChannel(winner)
	if in {
		winningUser.WinBrawl(T.season)
	}
}

func (T *brawlCommand) startBrawl() {
	T.active = true

	duration := rand.Intn(120) + 60
	timer := time.NewTimer(time.Duration(duration) * time.Second)

	go func() {
		<-timer.C

		T.endBrawl()
	}()

	T.cp.irc.Say("#"+T.cp.channel.GetChannelName(), fmt.Sprintf("A brawl has started in Twitch Chat! Type !pileon to join the fight! Everyone, get in here!"))
}

func (T *brawlCommand) Response(username, message string) string {
	message = strings.TrimSpace(strings.ToLower(message))
	clearance := T.cp.channel.GetLevel(username)

	_, err := T.brawlComm.parse(message)
	if err == nil && clearance >= channel.MOD && T.active == false {
		T.startBrawl()
		return ""
	}

	_, err = T.newSeasonComm.parse(message)
	if err == nil && clearance >= channel.BROADCASTER && T.active == false {
		T.season = T.season + 1
		// TODO: add top winners of the season
		return fmt.Sprintf("The brawl season has ended! We are now in season %d.", T.season)
	}

	_, err = T.pileComm.parse(message)
	if err == nil && T.active == true {
		_, in := T.cp.channel.InChannel(username)
		if !in {
			return ""
		}

		T.brawlers[username] = 1
		return ""
	}

	_, err = T.statComm.parse(message)
	if err == nil {
		user, in := T.cp.channel.InChannel(username)
		if in {
			brawlsWon := user.GetBrawlsWon()
			wins, ok := brawlsWon[T.season]
			if !ok {
				wins = 0
			}
			return fmt.Sprintf("@%s you have won %d brawls this season", username, wins)
		}
	}

	args, err := T.statsComm.parse(message)
	if err == nil {
		season, _ := strconv.Atoi(args[0])
		return T.calculateBrawlStats(season)
	}

	return ""
}

type brawlWin struct {
	username string
	wins     int
}

type viewerBrawlerInterface struct {
	viewers []brawlWin
}

func (V *viewerBrawlerInterface) Len() int {
	return len(V.viewers)
}

func (V *viewerBrawlerInterface) Less(i, j int) bool {
	return V.viewers[i].wins > V.viewers[j].wins
}

func (V *viewerBrawlerInterface) Swap(i, j int) {
	oldi := V.viewers[i]
	V.viewers[i] = V.viewers[j]
	V.viewers[j] = oldi
}

func (T *brawlCommand) calculateBrawlStats(season int) string {
	brawlWins := []brawlWin{}

	foundWinner := false
	for _, u := range T.cp.channel.AllViewers() {
		wins := 0
		if season > 0 {
			w, ok := u.GetBrawlsWon()[season]
			if ok {
				wins = w
			}
		} else {
			wins = u.GetTotalBrawlsWon()
		}
		if wins > 0 {
			foundWinner = true
		}
		brawlWins = append(brawlWins, brawlWin{
			username: u.Username,
			wins:     wins,
		})
	}

	if !foundWinner {
		return fmt.Sprintf("No one has won for season %d yet!", season)
	}

	sorter := &viewerBrawlerInterface{
		viewers: brawlWins,
	}
	sort.Sort(sorter)

	winners := []brawlWin{}
	tiedWinners := []brawlWin{}
	numWins := 10000
	count := -1
	output := ""
	if season > 0 {
		output = output + fmt.Sprintf("Best brawlers for season %d: ", season)
	} else {
		output = output + "All-time best brawlers: "
	}
	for _, w := range sorter.viewers {
		if w.wins < numWins {
			if len(tiedWinners) != 0 {
				for _, winner := range tiedWinners {
					output = fmt.Sprintf("%s%s & ", output, winner.username)
				}
				output = fmt.Sprintf("%s (%d wins), ", output[:len(output)-2], numWins)
			}

			winners = append(winners, tiedWinners...)
			tiedWinners = []brawlWin{}
			numWins = w.wins
			count = count + 1
		}
		if w.wins == numWins {
			tiedWinners = append(tiedWinners, w)
		}
		if numWins == 0 || count > 3 {
			break
		}
	}

	return output[:len(output)-2]
}

func (T *brawlCommand) String() string {
	return ""
}
