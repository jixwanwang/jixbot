package channel

import "github.com/jixwanwang/jixbot/stats"

type Level int

const (
	VIEWER      Level = 0
	MOD         Level = 1
	STAFF       Level = 2
	BROADCASTER Level = 3
	GOD         Level = 4
)

func init() {
	// TODO: Load known staff
}

type ViewerList struct {
	channel             string
	db                  *stats.ViewerManager
	viewers             map[string]*stats.Viewer
	staff               map[string]int
	mods                map[string]int
	lotteryContributers map[string]int
}

func NewViewerList(channel string) *ViewerList {
	stats.Init(channel)

	return &ViewerList{
		channel:             channel,
		db:                  stats.Init(channel),
		viewers:             map[string]*stats.Viewer{},
		staff:               map[string]int{},
		mods:                map[string]int{},
		lotteryContributers: map[string]int{},
	}
}

func (V *ViewerList) GetChannelName() string {
	return V.channel
}

func (V *ViewerList) AddViewer(username string) {
	if _, ok := V.viewers[username]; !ok {
		V.viewers[username] = V.db.FindViewer(username)
	}
}

func (V *ViewerList) AddViewers(usernames []string) {
	for _, u := range usernames {
		V.AddViewer(u)
	}
}

func (V *ViewerList) RemoveViewer(username string) {
	delete(V.viewers, username)
	delete(V.mods, username)
}

func (V *ViewerList) AddMod(username string) {
	V.AddViewers([]string{username})
	if _, ok := V.viewers[username]; ok {
		V.mods[username] = 1
	}
}

func (V *ViewerList) RemoveMod(username string) {
	delete(V.mods, username)
}

func (V *ViewerList) InChannel(username string) (*stats.Viewer, bool) {
	v, ok := V.viewers[username]
	return v, ok
}

func (V *ViewerList) AllViewers() []*stats.Viewer {
	return V.db.AllViewers()
}

func (V *ViewerList) GetLevel(username string) Level {
	if username == "jixwanwang" {
		return GOD
	} else if username == V.channel {
		return BROADCASTER
	} else if _, ok := V.staff[username]; ok {
		return STAFF
	} else if _, ok := V.mods[username]; ok {
		return MOD
	}
	return VIEWER
}

func (V *ViewerList) RecordMessage(username, msg string) {
	v, ok := V.viewers[username]
	if !ok {
		V.AddViewer(username)
		v = V.viewers[username]
	}

	v.LinesTyped = v.LinesTyped + 1
	v.Money = v.Money + 1
}

func (V *ViewerList) Tick() {
	for _, v := range V.viewers {
		v.Money = v.Money + 1
	}
	V.db.Flush()
}

func (V *ViewerList) Close() {
	V.db.Flush()
}

// func (V *ViewerList) AddToLottery(username string, amount int) int {
// 	value, ok := V.lotteryContributers[username]
// 	if ok {
// 		V.lotteryContributers[username] = value + amount
// 		return value + amount
// 	}

// 	V.lotteryContributers[username] = amount
// 	return amount
// }

// func (V *ViewerList) LotteryReady() string {
// 	if len(V.lotteryContributers) < 10 {
// 		return "Not enough people bought tickets to run a lottery."
// 	}
// 	return ""
// }

// func (V *ViewerList) RunLottery() (string, int) {
// 	tickets := []string{}
// 	for username, v := range V.lotteryContributers {
// 		for i := 0; i < v; i++ {
// 			tickets = append(tickets, username)
// 		}
// 	}

// 	log.Printf("%v", tickets)

// 	winner := tickets[rand.Intn(len(tickets))]
// 	winningViewer := V.db.FindViewer(winner)
// 	winningAmount := V.viewers["jixbot"].Money / 2
// 	winningViewer.Money = winningViewer.Money + winningAmount
// 	V.viewers["jixbot"].Money = V.viewers["jixbot"].Money - winningAmount

// 	V.lotteryContributers = map[string]int{}

// 	return winner, winningAmount
// }
