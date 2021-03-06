package context

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/j75689/easybot/handler"
	messagehandler "github.com/j75689/easybot/handler"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

// HandleLineHook process line event.
func HandleLineHook(handler *httphandler.WebhookHandler, bot *linebot.Client) func(*gin.Context) {

	if handler != nil && bot != nil {

		// Setup HTTP Server for receiving requests from LINE platform
		handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {

			for _, event := range events {
				logger.Debug("[bot] ", structs.Map(event))
				if msg, err := messagehandler.Execute(event); msg != nil {
					if err != nil {
						logger.Warn("[bot] ", err)
					}
					msgData, _ := msg.MarshalJSON()
					logger.Debug("[bot] ", string(msgData))
					if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						logger.Error("[bot] ", err)
					}
				}
			}
		})

		return gin.WrapH(handler)
	}
	return nil
}

// HandlePushMessage process push api
func HandlePushMessage(bot *linebot.Client) func(*gin.Context) {
	return func(c *gin.Context) {

		if bot == nil {
			logger.Error("[api] ", "missing bot instance")
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Server Exception"})
			return
		}

		var (
			postdata []byte
			err      error
		)

		if postdata, err = ioutil.ReadAll(c.Request.Body); err == nil {
			logger.Info("[api] ", "push ", c.Param("userID"))
			if _, err = bot.PushMessage(c.Param("userID"), &handler.CustomMessage{Msg: string(postdata)}).Do(); err == nil {
				c.JSON(http.StatusOK, gin.H{"success": true})
				return
			}
		}
		logger.Error("[api] ", "Push Cause Error: ", err)
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "Post Line Push Falid."})
	}
}

// HandleMulticastMessage process multicast api
func HandleMulticastMessage(bot *linebot.Client) func(*gin.Context) {
	return func(c *gin.Context) {

		if bot == nil {
			logger.Error("[api] ", "missing bot instance")
			c.JSON(http.StatusOK, gin.H{"success": false, "error": "Server Exception"})
			return
		}

		type MulticastBody struct {
			UserIDs []string      `json:"to"`
			Message []interface{} `json:"messages"`
		}
		var (
			postdata      []byte
			multicastBody MulticastBody
			Messages      []linebot.SendingMessage
			err           error
		)
		if postdata, err = ioutil.ReadAll(c.Request.Body); err == nil {
			if err = json.Unmarshal(postdata, &multicastBody); err == nil {
				for _, data := range multicastBody.Message {
					if msg, err := json.Marshal(data); err == nil {
						Messages = append(Messages, &handler.CustomMessage{Msg: string(msg)})
					} else {
						logger.Error("[api] ", "Muticast Cause Error: ", err)
					}
				}

				if _, err = bot.Multicast(multicastBody.UserIDs, Messages...).Do(); err == nil {
					c.JSON(http.StatusOK, gin.H{"success": true})
					return
				} else {
					logger.Error("[api] ", "Muticast Cause Error: ", err)
					c.JSON(http.StatusOK, gin.H{"success": false, "error": "Post Line Multicast Falid."})
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error()})

	}
}
