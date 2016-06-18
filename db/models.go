package db

import "time"

type TextCommand struct {
	Clearance int
	Command   string
	Response  string
	Cooldown  time.Duration
}

type Count struct {
	Username string
	Count    int
}

type QuestionAnswer struct {
	ID      int
	Channel string
	Q       string
	A       string
}
