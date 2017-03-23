package command

import (
	"fmt"
	"log"
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
	cleanRgx    *regexp.Regexp
	mentionRgx  *regexp.Regexp
	qa          *nlp.QuestionAnswerer

	newQuestions         map[string]newQuestion
	storeAnswer          *subCommand
	lastQuestionAnswered time.Time
}

func (T *questions) Init() {
	qas, err := T.cp.db.RetrieveQuestionAnswers(T.cp.channel.GetChannelName())
	if err != nil {
		log.Printf("Failed to lookup question answers!")
	}

	T.questionRgx, _ = regexp.Compile(`^(who|what|when|where|why)`)
	T.cleanRgx, _ = regexp.Compile(`[^A-Za-z0-9\?\ ]+`)
	T.mentionRgx, _ = regexp.Compile(`@[a-zA-Z0-9_]+`)

	T.qa = nlp.NewQuestionAnswerer()
	for _, qa := range qas {
		T.qa.AddQuestionAndAnswer(qa.Q, qa.A)
	}

	T.lastQuestionAnswered = time.Now()
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

	message = strings.TrimSpace(message)

	clearance := T.cp.channel.GetLevel(username)

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

	message = strings.ToLower(message)

	// Check for questionness
	if !(T.questionRgx.MatchString(message) || message[len(message)-1:] == "?") {
		return
	}

	// clean up old stuff in questions list
	for k, v := range T.newQuestions {
		if time.Since(v.timestamp) > 1*time.Minute {
			delete(T.newQuestions, k)
		}
	}

	// Remove all mentions and punctuation
	message = T.mentionRgx.ReplaceAllString(message, "")
	message = T.cleanRgx.ReplaceAllString(message, "")

	// Trim ?, don't want that in the question body
	message = strings.TrimSuffix(message, "?")

	T.newQuestions[username] = newQuestion{
		timestamp: time.Now(),
		username:  username,
		question:  message,
	}

	if time.Since(T.lastQuestionAnswered) > 15*time.Second {
		// Check if the question has an answer
		answer, score := T.qa.AnswerQuestion(message)

		if len(answer) > 0 && score > 0.8 {
			T.lastQuestionAnswered = time.Now()
			T.cp.Say(fmt.Sprintf("@%s %s", username, answer))
			return
		}
	}
}
