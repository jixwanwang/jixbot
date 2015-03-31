package channel

import "github.com/jixwanwang/jixbot/stats"

type Level int

const (
	VIEWER Level = 0
	MOD    Level = 1
	STAFF  Level = 2
	GOD    Level = 3
)

func init() {
	// TODO: Load known staff
}

type ViewerList struct {
	channel string
	viewers map[string]*stats.Viewer
	staff   map[string]int
	mods    map[string]int
}

func NewViewerList(channel string) *ViewerList {
	stats.Init(channel)

	return &ViewerList{
		channel: channel,
		viewers: map[string]*stats.Viewer{},
		staff:   map[string]int{},
		mods:    map[string]int{},
	}
}

func (V *ViewerList) GetChannelName() string {
	return V.channel
}

func (V *ViewerList) AddViewer(username string) {
	if _, ok := V.viewers[username]; !ok {
		V.viewers[username] = stats.NewViewer(username, V.channel)
	}
}

func (V *ViewerList) AddViewers(usernames []string) {
	for _, u := range usernames {
		V.AddViewer(u)
	}
}

func (V *ViewerList) RemoveViewer(username string) {
	if v, ok := V.viewers[username]; ok {
		stats.SaveViewer(v, V.GetChannelName())
	}
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

func (V *ViewerList) GetLevel(username string) Level {
	if username == "jixwanwang" {
		return GOD
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

	stats.SaveViewer(v, V.channel)
}

func (V *ViewerList) Tick() {
	for _, v := range V.viewers {
		v.Money = v.Money + 1
		stats.SaveViewer(v, V.channel)
	}
}
