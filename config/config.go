package config

import (
	"log"
	"os"
)

var (
	Nickname        string
	OauthToken      string
	FancyNickname   string
	FancyOauthToken string
	GroupChat       string

	PastebinKey string

	TwilioAccount string
	TwilioSecret  string
	TwilioNumber  string
	JixNumber     string

	ClientID string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
)

func init() {
	Nickname = os.Getenv("NICKNAME")
	OauthToken = os.Getenv("OATH_TOKEN")
	FancyNickname = os.Getenv("FANCY_NICKNAME")
	FancyOauthToken = os.Getenv("FANCY_OAUTH_TOKEN")
	GroupChat = os.Getenv("GROUPCHAT")

	PastebinKey = os.Getenv("PASTEBIN_API_KEY")

	ClientID = os.Getenv("CLIENT_ID")

	TwilioAccount = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioSecret = os.Getenv("TWILIO_SECRET")
	TwilioNumber = os.Getenv("TWILIO_NUMBER")

	JixNumber = os.Getenv("JIX_NUMBER")

	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASS")

	// TODO: verify things have values
	log.Printf("%v %v", FancyNickname, FancyOauthToken)
}
