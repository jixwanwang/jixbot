package command

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/db"
)

type newQuestion struct {
	timestamp time.Time
	username  string
	question  string
}

type questions struct {
	cp *CommandPool

	qas           []db.QuestionAnswer
	newQuestions  map[string]newQuestion
	storeQuestion *subCommand
	storeAnswer   *subCommand
}

func (T *questions) Init() {
	qas, err := T.cp.db.RetrieveQuestionAnswers(T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("Failed to lookup question answers!")
	}

	T.qas = qas
	T.newQuestions = map[string]newQuestion{}

	T.storeAnswer = &subCommand{
		command:   "!q",
		numArgs:   1,
		cooldown:  5 * time.Second,
		clearance: channel.MOD,
	}
}

func (T *questions) ID() string {
	return "questions"
}

func (T *questions) Response(username, message string, whisper bool) {
	if whisper {
		return
	}

	message = strings.TrimSpace(strings.ToLower(message))

	clearance := T.cp.channel.GetLevel(username)

	args, err := T.storeAnswer.parse(message, clearance)
	if err == nil {
		user := strings.TrimPrefix(args[0], "@")
		q, ok := T.newQuestions[user]
		if ok {
			qa, err := T.cp.db.AddQuestionAnswer(T.cp.channel.GetChannelName(), q.question, args[1])
			if err == nil {
				T.qas = append(T.qas, qa)
				T.cp.Say(fmt.Sprintf("@%s question has been answered!", username))
			}
			return
		}
	}

	// Short circuit with a less strict test to prevent doing all this string manipulation
	if !(strings.Index(message, "who") >= 0 ||
		strings.Index(message, "what") >= 0 ||
		strings.Index(message, "when") >= 0 ||
		strings.Index(message, "where") >= 0 ||
		strings.Index(message, "how") >= 0 ||
		strings.HasSuffix(message, "?")) {
		return
	}

	// Remove all mentions
	for i := strings.Index(message, "@"); i >= 0; i = strings.Index(message, "@") {
		part := message[i:]
		space := strings.Index(part, " ")
		if space >= 0 {
			message = message[:i] + part[space+1:]
		}
	}

	// Remove all punctuation
	reg, err := regexp.Compile(`[^A-Za-z0-9\?\ ]+`)
	message = reg.ReplaceAllString(message, "")

	// remove commonly used fillers for questions
	message = strings.Replace(message, "hi ", "", -1)
	message = strings.Replace(message, "hello ", "", -1)
	message = strings.Replace(message, "hey ", "", -1)

	message = strings.Replace(message, " hi", "", -1)
	message = strings.Replace(message, " hello", "", -1)
	message = strings.Replace(message, " hey", "", -1)

	message = strings.TrimSpace(message)

	// Check for questionness
	if !(strings.HasPrefix(message, "who") ||
		strings.HasPrefix(message, "what") ||
		strings.HasPrefix(message, "when") ||
		strings.HasPrefix(message, "where") ||
		strings.HasPrefix(message, "how") ||
		strings.HasSuffix(message, "?")) {
		return
	}

	// clean up old stuff in questions list
	for k, v := range T.newQuestions {
		if time.Since(v.timestamp) > 5*time.Minute {
			delete(T.newQuestions, k)
		}
	}

	if strings.HasSuffix(message, "?") {
		message = strings.TrimSuffix(message, "?")
	}

	log.Printf("question: %s", message)

	// Check if the question has an answer
	score := 0.0
	var bestAnswer string
	for _, q := range T.qas {
		if s := T.similarQuestionScore(q.Q, message); s > score {
			score = s
			bestAnswer = q.A
		}
	}

	if score > 0.8 {
		T.cp.Say(fmt.Sprintf("@%s %s", username, bestAnswer))
		return
	}

	T.newQuestions[username] = newQuestion{
		timestamp: time.Now(),
		username:  username,
		question:  message,
	}
}

func (T *questions) similarQuestionScore(q1, q2 string) float64 {
	// Calculate cosine score using the formula:
	// cos(v1, v2) = ( v1 . v2 )/ ( ||v1|| * ||v2|| )

	// Get word vectors
	count1 := countWords(q1)
	count2 := countWords(q2)

	// Calculate product of lengths
	sum1 := sumSquares(count1)
	sum2 := sumSquares(count2)
	denominator := math.Sqrt(float64(sum1 * sum2))

	// Dot product of vectors
	sharedWords := sharedKeys(count1, count2)
	numerator := 0
	for _, word := range sharedWords {
		numerator += count1[word] * count2[word]
	}

	if denominator == 0 {
		return 0
	}

	log.Printf("score between %s, %s: %v", q1, q2, float64(numerator)/denominator)
	return float64(numerator) / denominator
}

func countWords(s string) map[string]int {
	counts := map[string]int{}
	for _, word := range strings.Split(s, " ") {
		if len(word) == 0 {
			continue
		}
		// TODO: stem the word before incrementing its count
		if _, ok := counts[word]; !ok {
			counts[word] = 0
		}
		counts[word]++
	}

	return counts
}

func sumSquares(m map[string]int) int {
	sum := 0
	for _, v := range m {
		sum += v * v
	}
	return sum
}

func sharedKeys(m1, m2 map[string]int) []string {
	shared := []string{}
	for k := range m1 {
		if _, ok := m2[k]; ok {
			shared = append(shared, k)
		}
	}
	return shared
}
