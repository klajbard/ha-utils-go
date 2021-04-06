package slack

import (
	"errors"
	"log"

	"github.com/klajbard/ha-utils-go/config"
	"github.com/slack-go/slack"
)

// Slack notification uses hooks to send messages
// *channel* - Channel id
// *text* - The sent text formatted with markup
// *emoji* - Emoji to send the message as
func NotifySlack(channel, text, emoji string) {
	if config.Conf.Silence {
		return
	}
	if channel == "" || config.Channels[channel] == "" {
		log.Println(errors.New("Unable to get ENV variable"))
	}
	_, _, err := config.SlackBot.PostMessage(channel, slack.MsgOptionText(text, false), slack.MsgOptionIconEmoji(emoji))
	if err != nil {
		log.Printf("Posting message failed: %v", err)
	}
}
