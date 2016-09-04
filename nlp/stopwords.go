package nlp

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var stopwordRegex *regexp.Regexp
var stopwordMap map[string]bool

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

	stopwordMap = map[string]bool{}
	for _, w := range stopwords {
		stopwordMap[w] = true
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

func IsStopword(word string) bool {
	_, ok := stopwordMap[word]
	return ok
}

func FilterStopwords(message string) string {
	return string(stopwordRegex.ReplaceAll([]byte(message), []byte("")))
}
