package stats

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

const statsFilePath = "data/stats/"

// Represents a viewer in a single stream.
type Viewer struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
	Money      int    `json:"money"`
	BrawlsWon  int    `json:"brawls_won"`
}

type ViewerManager struct {
	channel string
	viewers map[string]*Viewer
}

func Init(channel string) *ViewerManager {
	// TODO: preload all the viewers whose data is already stored
	// os.MkdirAll(statsFilePath+channel, 0755)

	manager := ViewerManager{
		channel: channel,
		viewers: map[string]*Viewer{},
	}

	statsRaw, _ := ioutil.ReadFile(statsFilePath + channel + "_stats")
	statLines := strings.Split(string(statsRaw), "\n")

	for _, line := range statLines {
		var viewer Viewer

		err := json.Unmarshal([]byte(line), &viewer)

		if err != nil {
			continue
		}

		manager.viewers[viewer.Username] = &viewer
	}

	return &manager
}

func (V *ViewerManager) AllViewers() []*Viewer {
	viewers := []*Viewer{}
	for _, v := range V.viewers {
		viewers = append(viewers, v)
	}
	return viewers
}

func (V *ViewerManager) FindViewer(username string) *Viewer {
	viewer, ok := V.viewers[username]

	if !ok {
		V.viewers[username] = &Viewer{Username: username}
		return V.viewers[username]
	}

	return viewer
}

func (V *ViewerManager) FindViewers(usernames []string) []*Viewer {
	v := []*Viewer{}

	for _, u := range usernames {
		v = append(v, V.FindViewer(u))
	}

	return v
}

func (V *ViewerManager) Flush() {
	output := ""
	for _, v := range V.viewers {
		data, _ := json.Marshal(v)
		log.Printf("Saved %v for %s", string(data), V.channel)
		output = output + string(data) + "\n"
	}

	ioutil.WriteFile(statsFilePath+V.channel+"_stats", []byte(output), 0666)
}
