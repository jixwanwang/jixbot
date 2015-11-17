package command

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
)

type brawl struct {
	cp               *CommandPool
	brawlComm        *subCommand
	seasonComm       *subCommand
	newSeasonComm    *subCommand
	pileComm         *subCommand
	statsComm        *subCommand
	alltimeStatsComm *subCommand
	statComm         *subCommand

	season   int
	brawlers map[string]string
	active   bool
}

func (T *brawl) Init() {
	row := T.cp.db.QueryRow("select * from (select distinct(season) from brawlwins where channel=$1 order by season desc) as seasons limit 1", T.cp.channel.GetChannelName())
	if err := row.Scan(&T.season); err != nil {
		log.Printf("couldn't determine the brawl season, assuming to be 1")
		T.season = 1
	}
	log.Printf("The current brawl season is %d", T.season)

	// username to weapon mapping
	T.brawlers = map[string]string{}

	T.brawlComm = &subCommand{
		command:   "!brawl",
		numArgs:   0,
		cooldown:  15 * time.Minute,
		clearance: channel.MOD,
	}

	T.seasonComm = &subCommand{
		command:   "!brawlseason",
		numArgs:   0,
		cooldown:  15 * time.Second,
		clearance: channel.VIEWER,
	}

	T.newSeasonComm = &subCommand{
		command:   "!newbrawlseason",
		numArgs:   0,
		cooldown:  12 * time.Second,
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
		numArgs:   1,
		cooldown:  15 * time.Second,
		clearance: channel.VIEWER,
	}

	T.statComm = &subCommand{
		command:   "!brawlwins",
		numArgs:   0,
		cooldown:  5 * time.Second,
		clearance: channel.VIEWER,
	}
}

func (T *brawl) ID() string {
	return "brawl"
}

func (T *brawl) endBrawl() {
	T.active = false
	log.Printf("Brawl ended")

	users := []string{}

	for k := range T.brawlers {
		users = append(users, k)
	}

	if len(users) <= 0 {
		return
	}

	if len(users) == 1 {
		T.cp.Say(fmt.Sprintf("The brawl is over, but %s was the only one fighting. That was a boring brawl.", users[0]))
		return
	} else if len(users) < 5 {
		T.cp.Say(fmt.Sprintf("Only a few people joined the brawl, while others just sat around and watched. That was really boring."))
		return
	}

	winnerIndex := rand.Intn(len(users))
	winner := users[winnerIndex]
	// Broadcaster has higher chance of winning if they piled on
	if _, ok := T.brawlers[T.cp.channel.GetChannelName()]; ok && winnerIndex < 2 {
		winner = T.cp.channel.GetChannelName()
	}
	weapon := T.brawlers[winner]

	T.brawlers = map[string]string{}

	if len(weapon) > 0 {
		var message string
		message = fmt.Sprintf("The brawl is over, the tavern is a mess! @%s has defeated everyone with their %s ! They loot 500 %ss from the losers.", winner, weapon, T.cp.channel.Currency)
		T.cp.Say(message)
	} else {
		message := fmt.Sprintf("The brawl is over, the tavern is a mess, but @%s is the last one standing! They loot 500 %ss from the losers.", winner, T.cp.channel.Currency)
		T.cp.Say(message)
	}

	winningUser, in := T.cp.channel.InChannel(winner)
	if in {
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

	T.cp.Say(fmt.Sprintf("PogChamp A brawl has started in Twitch Chat! Type !pileon to join the fight! You can also use a weapon using !pileon <weapon>! Everyone, get in here! PogChamp"))
}

func (T *brawl) Response(username, message string, whisper bool) {
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
		return
	}

	args, err := T.pileComm.parse(message, clearance)
	if err == nil && T.active == true {
		_, in := T.cp.channel.InChannel(username)
		if !in {
			return
		}

		if len(args) > 0 {
			T.brawlers[username] = args[0]
		} else {
			T.brawlers[username] = ""
		}

		return
	}

	_, err = T.statComm.parse(message, clearance)
	if err == nil {
		user, in := T.cp.channel.InChannel(username)
		if in {
			brawlsWon := user.GetBrawlsWon()
			wins, ok := brawlsWon[T.season]
			if !ok {
				wins = 0
			}
			T.cp.Say(fmt.Sprintf("@%s you have won %d brawls this season", username, wins))
			return
		}
	}

	args, err = T.statsComm.parse(message, clearance)
	if err == nil {
		season, _ := strconv.Atoi(args[0])
		T.cp.Say(T.calculateBrawlStats(season))
		return
	}

	return
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

func (T *brawl) calculateBrawlStats(season int) string {
	output := ""
	var rows *sql.Rows
	var err error
	if season > 0 {
		output = output + fmt.Sprintf("Best brawlers for season %d: ", season)
		rows, err = T.cp.db.Query(`SELECT sum(wins) totalwins, username FROM brawlwins AS b `+
			`JOIN viewers AS v ON v.id=b.viewer_id `+
			`WHERE b.channel=$1 AND b.season=$2 `+
			`GROUP BY username ORDER BY totalwins DESC`, T.cp.channel.GetChannelName(), season)
	} else {
		output = output + "All-time best brawlers: "
		rows, err = T.cp.db.Query(`SELECT sum(wins) totalwins, username FROM brawlwins AS b `+
			`JOIN viewers AS v ON v.id=b.viewer_id `+
			`WHERE b.channel=$1 `+
			`GROUP BY username ORDER BY totalwins DESC`, T.cp.channel.GetChannelName())
	}
	if err != nil {
		return ""
	}

	brawlWins := []brawlWin{}
	foundWinner := false
	var username string
	var wins int
	for rows.Next() {
		err := rows.Scan(&wins, &username)
		log.Printf("%d, %s", wins, username)
		if err != nil {
			continue
		}

		if wins > 0 {
			foundWinner = true
		}

		brawlWins = append(brawlWins, brawlWin{
			username: username,
			wins:     wins,
		})
	}
	rows.Close()

	if !foundWinner {
		return fmt.Sprintf("No one has won for season %d yet!", season)
	}

	tiedWinners := []brawlWin{}
	numWins := 10000
	count := -1
	for _, w := range brawlWins {
		if w.wins < numWins {
			if len(tiedWinners) != 0 {
				for _, winner := range tiedWinners {
					output = fmt.Sprintf("%s%s & ", output, winner.username)
				}
				output = fmt.Sprintf("%s (%d wins), ", output[:len(output)-2], numWins)
			}

			tiedWinners = []brawlWin{}
			numWins = w.wins
			count = count + 1
		}

		// Don't display 0 wins, only keep top 5 win groups, and ignore groups with more than 7 people in it.
		if numWins == 0 || count > 4 || len(tiedWinners) > 7 {
			break
		}

		if w.wins == numWins {
			tiedWinners = append(tiedWinners, w)
		}
	}

	if len(tiedWinners) != 0 && len(tiedWinners) <= 7 {
		for _, winner := range tiedWinners {
			output = fmt.Sprintf("%s%s & ", output, winner.username)
		}
		output = fmt.Sprintf("%s (%d wins), ", output[:len(output)-2], numWins)
	}

	return output[:len(output)-2]
}
