package channel

import "log"

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

	linesTyped int
	timeSpent  int
	money      int
	brawlsWon  map[int]int
}

func (V *Viewer) Reset() {
	V.money = 0
	V.linesTyped = 0
	V.timeSpent = 0
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
	V.AddMoney(500)

	V.updated = true
}

func (V *Viewer) GetLinesTyped() int {
	if V.linesTyped < 0 {
		if V.id > 0 {
			lines, err := V.manager.db.GetCount(V.id, "lines_typed")
			if err == nil {
				V.linesTyped = lines
			}
		} else {
			V.linesTyped = 0
		}
	}
	return V.linesTyped
}

func (V *Viewer) AddLineTyped() {
	V.linesTyped = V.GetLinesTyped() + 1
	V.updated = true
}

func (V *Viewer) GetTimeSpent() int {
	if V.timeSpent < 0 {
		if V.id > 0 {
			time, err := V.manager.db.GetCount(V.id, "time")
			if err == nil {
				V.timeSpent = time
			}
		} else {
			V.timeSpent = 0
		}
	}
	return V.timeSpent
}

func (V *Viewer) AddTimeSpent(minutes int) {
	V.timeSpent = V.GetTimeSpent() + minutes
	V.updated = true
}

func (V *Viewer) GetMoney() int {
	if V.money < 0 {
		if V.id > 0 {
			money, err := V.manager.db.GetCount(V.id, "money")
			if err == nil {
				V.money = money
			}
		} else {
			V.money = 0
		}
	}
	return V.money
}

func (V *Viewer) AddMoney(amount int) {
	V.money = V.GetMoney() + amount
	if V.money < 0 {
		V.money = 0
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
		log.Printf("%s is updated", V.Username)
		if V.id == -1 {
			id, err := V.manager.db.NewViewer(V.Username, V.manager.channel)
			if err == nil {
				V.id = id
			}
		}
		if V.brawlsWon != nil {
			V.manager.db.SetBrawlWins(V.id, V.manager.channel, V.brawlsWon)
		}
		if V.money > 0 {
			V.manager.db.SetCount(V.id, "money", V.money)
		}
		if V.linesTyped > 0 {
			V.manager.db.SetCount(V.id, "lines_typed", V.linesTyped)
		}
		if V.timeSpent > 0 {
			V.manager.db.SetCount(V.id, "time", V.timeSpent)
		}
		V.updated = false
	}
}
