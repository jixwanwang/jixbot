package stats

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const statsFilePath = "data/stats/"

// Represents a viewer in a single stream.
type Viewer struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
	KappaCount int    `json:"kappa_count"`
	Money      int    `json:"money"`
}

func Init(channel string) {
	// TODO: preload all the viewers whose data is already stored
	os.MkdirAll(statsFilePath+channel, 0755)
}

func NewViewer(username, channel string) *Viewer {
	statsRaw, _ := ioutil.ReadFile(statsFilePath + channel + "/" + username)

	var v map[string]interface{}

	viewer := &Viewer{Username: username}

	err := json.Unmarshal(statsRaw, &v)
	if err != nil {
		// No stats available
		return &Viewer{Username: username}
	}

	if l, ok := v["lines_typed"]; ok {
		value, ok := l.(float64)
		if ok {
			viewer.LinesTyped = int(value)
		}
	}
	if l, ok := v["kappa_count"]; ok {
		value, ok := l.(float64)
		if ok {
			viewer.KappaCount = int(value)
		}
	}
	if l, ok := v["money"]; ok {
		value, ok := l.(float64)
		if ok {
			viewer.Money = int(value)
		}
	}

	log.Printf("read viewer %v", v)

	return viewer
}

func NewViewers(usernames []string, channel string) []*Viewer {
	v := []*Viewer{}

	for _, u := range usernames {
		v = append(v, NewViewer(u, channel))
	}

	return v
}

func SaveViewer(v *Viewer, channel string) {
	data, _ := json.Marshal(v)
	log.Printf("Saved %v for %s", string(data), channel)
	ioutil.WriteFile(statsFilePath+channel+"/"+v.Username, data, 0666)
}
