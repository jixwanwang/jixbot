package command

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/jixwanwang/jixbot/channel"
	"github.com/jixwanwang/jixbot/nlp"
)

type newQuestion struct {
	timestamp time.Time
	username  string
	question  string
}

type questions struct {
	cp *CommandPool

	questionRgx *regexp.Regexp
	qa          *nlp.QuestionAnswerer

	newQuestions map[string]newQuestion
	storeAnswer  *subCommand
}

func (T *questions) Init() {
	qas, err := T.cp.db.RetrieveQuestionAnswers(T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("Failed to lookup question answers!")
	}

	T.questionRgx, _ = regexp.Compile(`[^A-Za-z0-9\?\ ]+`)
	T.qa = nlp.NewQuestionAnswerer()
	for _, qa := range qas {
		T.qa.AddQuestionAndAnswer(qa.Q, qa.A)
	}

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

	// clean up old stuff in questions list
	for k, v := range T.newQuestions {
		if time.Since(v.timestamp) > 3*time.Minute {
			delete(T.newQuestions, k)
		}
	}

	args, err := T.storeAnswer.parse(message, clearance)
	if err == nil {
		user := strings.TrimPrefix(args[0], "@")
		q, ok := T.newQuestions[user]
		if ok {
			qa, err := T.cp.db.AddQuestionAnswer(T.cp.channel.GetChannelName(), q.question, args[1])
			if err == nil {
				T.qa.AddQuestionAndAnswer(qa.Q, qa.A)
				T.cp.Say(fmt.Sprintf("@%s question has been answered!", username))
			}
			return
		}
	}

	// Short circuit with a less strict test to prevent doing all this string manipulation
	if strings.Index(message, "who") < 0 &&
		strings.Index(message, "what") < 0 &&
		strings.Index(message, "when") < 0 &&
		strings.Index(message, "where") < 0 &&
		strings.Index(message, "how") < 0 &&
		!strings.HasSuffix(message, "?") {
		return
	}

	// Remove all mentions
	// TODO: make this a regex
	for i := strings.Index(message, "@"); i >= 0; i = strings.Index(message, "@") {
		part := message[i:]
		space := strings.Index(part, " ")
		if space >= 0 {
			message = message[:i] + part[space+1:]
		}
	}

	// // remove commonly used fillers for questions
	// message = strings.Replace(message, "hi ", "", -1)
	// message = strings.Replace(message, "hello ", "", -1)
	// message = strings.Replace(message, "hey ", "", -1)

	// message = strings.Replace(message, " hi", "", -1)
	// message = strings.Replace(message, " hello", "", -1)
	// message = strings.Replace(message, " hey", "", -1)

	// message = strings.TrimSpace(message)

	// Check for questionness
	if !(strings.HasPrefix(message, "who") ||
		strings.HasPrefix(message, "what") ||
		strings.HasPrefix(message, "when") ||
		strings.HasPrefix(message, "where") ||
		strings.HasPrefix(message, "how") ||
		strings.HasSuffix(message, "?")) {
		return
	}

	// Remove all punctuation
	message = T.questionRgx.ReplaceAllString(message, "")

	if strings.HasSuffix(message, "?") {
		message = strings.TrimSuffix(message, "?")
	}

	T.newQuestions[username] = newQuestion{
		timestamp: time.Now(),
		username:  username,
		question:  message,
	}

	log.Printf("question: %s", message)

	// Only attempt to answer questions 50% of the time
	if rand.Intn(2) == 0 {
		return
	}

	// Check if the question has an answer
	answer, score := T.qa.AnswerQuestion(message)

	if len(answer) > 0 && score > 0.8 {
		T.cp.Say(fmt.Sprintf("@%s %s", username, answer))
		return
	}
}
