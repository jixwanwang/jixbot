package messaging

import (
	"net/http"
	"net/url"
	"strings"
)

const urlBase = "https://api.twilio.com/2010-04-01"

// const urlBase = "http://requestb.in/1ijp55g1"

type Texter struct {
	account  string
	password string
	number   string
	toNumber string
}

type twilioRequest struct {
	Body string `json:Body`
	To   string `json:To`
	From string `json:From`
}

func NewTexter(account, password, number, toNumber string) Texter {
	return Texter{
		account:  account,
		password: password,
		number:   number,
	}
}

func (T Texter) SendText(body string) {
	opts := url.Values{}
	opts.Set("To", T.toNumber)
	opts.Set("Body", body)
	opts.Set("From", T.number)

	req, _ := http.NewRequest("POST", urlBase+"/Accounts/"+T.account+"/Messages.json", strings.NewReader(opts.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(T.account, T.password)

	http.DefaultClient.Do(req)
}
