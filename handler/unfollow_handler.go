package handler

import (
	"sync"

	"github.com/j75689/easybot/model"
	"github.com/line/line-bot-sdk-go/linebot"
)

type UnfollowHandler struct {
	BaseHandler
	Event  linebot.EventType
	Config *model.MessageHandlerConfig
}

func (h *UnfollowHandler) RegisterConfig(cfg *model.MessageHandlerConfig) (err error) {
	h.DeregisterConfig(cfg.ConfigID)

	h.BaseHandler.RegisterConfig(cfg)
	h.Config = cfg
	return
}

func (h *UnfollowHandler) DeregisterConfig(id string) (err error) {
	if h.Config != nil {
		if id == h.Config.ConfigID {
			h.Config = nil
		}
		err = h.BaseHandler.DeregisterConfig(id)
	}

	return
}

func (h *UnfollowHandler) Run(event *linebot.Event, variables map[string]interface{}) (reply *CustomMessage, err error) {
	reply, err = h.BaseHandler.Run(event, variables)

	if h.Config != nil && reply == nil {
		// add defaultValue
		for k, v := range h.Config.DefaultValues {
			variables[k] = v
		}
		var replyStr string
		replyStr, err = h.runStage(h.Config.ConfigID, 0, h.Config.Stage, variables)
		reply = &CustomMessage{
			Msg: replyStr,
		}
	}

	return
}

func newUnfollowHandler() *UnfollowHandler {
	return &UnfollowHandler{
		BaseHandler: BaseHandler{
			Configs: &sync.Map{},
			Wating:  &sync.Map{},
		},
		Event:  linebot.EventTypeUnfollow,
		Config: nil,
	}
}
