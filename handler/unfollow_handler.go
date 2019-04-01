package handler

import (
	"sync"

	"github.com/j75689/easybot/config"
	"github.com/line/line-bot-sdk-go/linebot"
)

type UnfollowHandler struct {
	BaseHandler
	Event  linebot.EventType
	Config *config.MessageHandlerConfig
}

func (h *UnfollowHandler) RegisterConfig(cfg *config.MessageHandlerConfig) (err error) {
	h.DeregisterConfig(cfg.ID)

	h.BaseHandler.RegisterConfig(cfg)
	h.Config = cfg
	return
}

func (h *UnfollowHandler) DeregisterConfig(id string) (err error) {
	if h.Config != nil {
		if id == h.Config.ID {
			h.Config = nil
		}
		err = h.BaseHandler.DeregisterConfig(id)
	}

	return
}

func (h *UnfollowHandler) Run(event *linebot.Event, variables map[string]interface{}) (reply *config.CustomMessage, err error) {
	if h.Config != nil {
		// add defaultValue
		for k, v := range h.Config.DefaultValues {
			variables[k] = v
		}
		var replyStr string
		replyStr, err = h.runStage(h.Config.ID, 0, h.Config.Stage, variables)
		reply = &config.CustomMessage{
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
