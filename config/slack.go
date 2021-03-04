package config

import (
	"os"

	"github.com/slack-go/slack"
)

var SlackBot *slack.Client

func init() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	SlackBot = slack.New(botToken, slack.OptionAppLevelToken(appToken))
}
