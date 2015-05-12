package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const statsFilePath = "data/stats/"

type oldStruct struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
	KappaCount int    `json:"kappa_count"`
	Money      int    `json:"money"`
}

type newStruct struct {
	Username   string `json:"username"`
	LinesTyped int    `json:"lines_typed"`
	Money      int    `json:"money"`
	BrawlsWon  int    `json:"brawls_won"`
}

func main() {
	channels := strings.Split(os.Getenv("CHANNELS"), ",")

	for _, channel := range channels {
		statsRaw, _ := ioutil.ReadFile(statsFilePath + channel + "_stats")
		statLines := strings.Split(string(statsRaw), "\n")

		newStats := []newStruct{}
		for _, line := range statLines {
			var old oldStruct
			newData := newStruct{}

			err := json.Unmarshal([]byte(line), &old)

			if err != nil {
				continue
			}

			newData.Username = old.Username
			newData.LinesTyped = old.LinesTyped
			newData.Money = old.Money
			newData.BrawlsWon = 0

			newStats = append(newStats, newData)
		}

		output := ""
		for _, v := range newStats {
			data, _ := json.Marshal(v)
			log.Printf("Saved %v for %s", string(data), channel)
			output = output + string(data) + "\n"
		}

		ioutil.WriteFile(statsFilePath+channel+"_stats", []byte(output), 0666)
	}
}
