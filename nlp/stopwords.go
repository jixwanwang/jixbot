package nlp

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var stopwordRegex *regexp.Regexp

func init() {
	stopwords := []string{"the", "is", "at",
		"which", "on", "and", "a", "an",
		"am", "hello", "hey",
		"be", "as", "by",
		"for", "from",
		"he", "her", "hers",
		"him", "his", "us", "we", "were", "to", "too"}

	b, err := ioutil.ReadFile("./nlp/stop-word-list.csv")
	if err != nil {
		log.Printf("failed to read stopword file, defaulting to basic stopwords")
	} else {
		stopwords = strings.Split(string(b), ",")
	}

	reStr := ""

	for i, word := range stopwords {
		if i != 0 {
			reStr += `|`
		}
		reStr += `\A` + word + `\z`
	}
	stopwordRegex = regexp.MustCompile(reStr)
}

func FilterStopwords(message string) string {
	return string(stopwordRegex.ReplaceAll([]byte(message), []byte("")))
}
