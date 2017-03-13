package channel

import (
	"log"

	"github.com/jixwanwang/jixbot/db"
)

const JIXBOT_CREATOR = "jix_bot"

var blacklistedUsers = map[string]int{
	"nightbot":     0,
	"moobot":       0,
	"jixbot":       0,
	"twitchnotify": 0,
	"revlobot":     0,
	"xanbot":       0,
}

// Represents a viewer in a single stream.
type Viewer struct {
	Username string

	id      int
	updated bool
	manager *ViewerList

	counts    *db.Counts
	brawlsWon map[int]int
}

func (V *Viewer) Reset() {
	V.counts = &db.Counts{}
	V.brawlsWon = map[int]int{}
	V.updated = true
}

func (V *Viewer) GetBrawlsWon() map[int]int {
	if V.brawlsWon == nil {
		brawlsWon := V.lookupBrawlWins()
		V.brawlsWon = brawlsWon
	}
	return V.brawlsWon
}

func (V *Viewer) GetTotalBrawlsWon() int {
	total := 0
	for _, wins := range V.GetBrawlsWon() {
		total = total + wins
	}
	return total
}

func (V *Viewer) lookupBrawlWins() map[int]int {
	if V.id < 0 {
		return map[int]int{}
	}

	wins, err := V.manager.db.GetBrawlWins(V.id)
	if err != nil {
		return map[int]int{}
	}

	return wins
}

func (V *Viewer) WinBrawl(season int) {
	V.GetBrawlsWon()

	if _, ok := V.brawlsWon[season]; ok {
		V.brawlsWon[season] = V.brawlsWon[season] + 1
	} else {
		V.brawlsWon[season] = 1
	}

	V.updated = true
}

func (V *Viewer) getCounts() {
	V.counts = &db.Counts{}

	if V.id > 0 {
		counts, err := V.manager.db.GetCounts(V.id)
		if err == nil {
			V.counts = counts
		}
	}
}

func (V *Viewer) GetLinesTyped() int {
	if V.counts == nil {
		V.getCounts()
	}
	return V.counts.LinesTyped
}

func (V *Viewer) AddLineTyped() {
	V.counts.LinesTyped = V.GetLinesTyped() + 1
	V.updated = true
}

func (V *Viewer) GetTimeSpent() int {
	if V.counts == nil {
		V.getCounts()
	}
	return V.counts.TimeSpent
}

func (V *Viewer) AddTimeSpent(minutes int) {
	V.counts.TimeSpent = V.GetTimeSpent() + minutes
	V.updated = true
}

func (V *Viewer) GetMoney() int {
	if V.counts == nil {
		V.getCounts()
	}
	return V.counts.Money
}

func (V *Viewer) AddMoney(amount int) {
	V.counts.Money = V.GetMoney() + amount
	if V.counts.Money < 0 {
		V.counts.Money = 0
	}
	V.updated = true
}

// TO DELETE DUPED BRAWL WINS:
/*
delete from brawlwins where id in
	(select id from
		(select *, row_number() OVER (ORDER BY viewer_id ASC) as row from brawlwins otr where
			(select count(*) from brawlwins inr where otr.viewer_id=inr.viewer_id and otr.season=inr.season and otr.channel=inr.channel) > 1
		ORDER BY viewer_id ASC) as dupes
	where dupes.row %2=0);
*/

func (V *Viewer) save() {
	if _, ok := blacklistedUsers[V.Username]; ok {
		V.Reset()
	}
	if V.updated {
		if V.id == -1 {
			id, err := V.manager.db.NewViewer(V.Username, V.manager.channel)
			if err == nil {
				V.id = id
			}
		}
		if V.brawlsWon != nil {
			V.manager.db.SetBrawlWins(V.id, V.manager.channel, V.brawlsWon)
		}

		if V.counts == nil {
			V.counts = &db.Counts{}
		}
		V.counts.ViewerID = V.id
		err := V.manager.db.SetCounts(V.counts)
		if err != nil {
			log.Printf("%v", err)
		}

		V.updated = false
	}
}
