package slack

import (
	"errors"
	"log"

	"../config"
	"github.com/slack-go/slack"
)

// Slack notification uses hooks to send messages
// *group* - The hash for the channel
// Can be extracted from: 'https://hooks.slack.com/services/*group*'
// *text* - The sent text formatted with markup
func NotifySlack(channel, text, emoji string) {
	if channel == "" || config.Channels[channel] == "" {
		log.Println(errors.New("Unable to get ENV variable"))
	}
	_, _, err := config.SlackBot.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
