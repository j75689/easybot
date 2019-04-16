package handler

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

// CustomMessage Reply Message
type CustomMessage struct {
	Msg string
}

func (m *CustomMessage) Message() {
}

func (m *CustomMessage) MarshalJSON() ([]byte, error) {
	return []byte(m.Msg), nil
}

// WithQuickReplies method of CustomMessage
func (m *CustomMessage) WithQuickReplies(items *linebot.QuickReplyItems) linebot.SendingMessage {
	return m
}
