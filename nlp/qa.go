package nlp

type qa struct {
	question map[string]int
	answer   string
}

type QuestionAnswerer struct {
	qas []*qa
}

func NewQuestionAnswerer() *QuestionAnswerer {
	return &QuestionAnswerer{
		qas: []*qa{},
	}
}

func (Q *QuestionAnswerer) AddQuestionAndAnswer(q, a string) {
	words := cleanWords(q)

	m := countWords(words)
	Q.qas = append(Q.qas, &qa{
		question: m,
		answer:   a,
	})
}

func (Q *QuestionAnswerer) AnswerQuestion(question string) (string, float64) {
	words := cleanWords(question)
	input := countWords(words)

	score := 0.0
	var bestAnswer string
	for _, q := range Q.qas {
		if s := cosineSimilarity(q.question, input); s > score {
			score = s
			bestAnswer = q.answer
		}
	}

	if score > 0.8 {
		return bestAnswer, score
	}
	return "", 0
}
