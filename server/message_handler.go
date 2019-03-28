package server

import (
	"github.com/j75689/easybot/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

func handleMessage(source *linebot.EventSource, input linebot.Message) *config.CustomMessage {
	var variables = map[string]interface{}{
		"source.RoomID":  source.RoomID,
		"source.GroupID": source.GroupID,
		"source.UserID":  source.UserID,
	}

	switch message := input.(type) {
	case *linebot.TextMessage:
		return handleTextMessage(message.Text, &variables)
	}
	return nil
}

func handleTextMessage(message string, variables *map[string]interface{}) *config.CustomMessage {
	logger.Debug(message, *variables)
	if value, ok := TextMessageHandleConfig.Load(message); ok {
		(*variables)["source.Message"] = message
		messageConfig := value.(config.MessageHandlerConfig)
		return &config.CustomMessage{
			Msg: runhandler(&messageConfig, *variables),
		}
	}
	return nil
}
