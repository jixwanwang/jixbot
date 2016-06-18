package db

func (B *dbImpl) RetrieveQuestionAnswers(channel string) ([]QuestionAnswer, error) {
	rows, err := B.db.Query("SELECT id, channel, question, answer FROM questions WHERE channel=$1", channel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ch, question, answer string
	var id int

	qas := []QuestionAnswer{}
	for rows.Next() {
		err = rows.Scan(&id, &ch, &question, &answer)
		if err == nil {
			qas = append(qas, QuestionAnswer{
				ID:      id,
				Channel: ch,
				Q:       question,
				A:       answer,
			})
		}
	}
	return qas, nil
}

func (B *dbImpl) AddQuestionAnswer(channel, question, answer string) (QuestionAnswer, error) {
	row := B.db.QueryRow("INSERT INTO questions (channel, question, answer) VALUES ($1, $2, $3) RETURNING id", channel, question, answer)
	var id int
	err := row.Scan(&id)
	return QuestionAnswer{
		ID:      id,
		Channel: channel,
		Q:       question,
		A:       answer,
	}, err
}

func (B *dbImpl) UpdateQuestionAnswer(qa QuestionAnswer) error {
	_, err := B.db.Exec("UPDATE questions SET channel=$1, question=$2, answer=$3 WHERE id=$4", qa.Channel, qa.Q, qa.A, qa.ID)
	return err
}
