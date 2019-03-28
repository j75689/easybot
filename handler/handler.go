package handler

import (
	"encoding/json"
	"sync"

	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/pkg/logger"
	"github.com/j75689/easybot/pkg/util"
	"github.com/j75689/easybot/plugin"

	"github.com/fatih/structs"

	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	handlers = make(map[linebot.EventType]Handler)
)

func init() {
	handlers[linebot.EventTypeMessage] = newMessageHandler()
	handlers[linebot.EventTypeFollow] = nil
	handlers[linebot.EventTypeUnfollow] = nil
	handlers[linebot.EventTypeJoin] = nil
	handlers[linebot.EventTypeLeave] = nil
	handlers[linebot.EventTypeMemberJoined] = nil
	handlers[linebot.EventTypeMemberLeft] = nil
	handlers[linebot.EventTypePostback] = nil
	handlers[linebot.EventTypeBeacon] = nil
	handlers[linebot.EventTypeAccountLink] = nil
	handlers[linebot.EventTypeThings] = nil
}

// Excute 執行
func Excute(event linebot.Event) (reply *config.CustomMessage) {
	var (
		variables = make(map[string]interface{})
	)
	variables["event"] = structs.Map(event)

	if handler := handlers[event.Type]; handler != nil {

		reply = handler.Run(event, variables)
	}
	return
}

// Handler interface
type Handler interface {
	GetConfig(string) *config.MessageHandlerConfig
	RegisterConfig(*config.MessageHandlerConfig) error
	DeregisterConfig(string) error
	Run(linebot.Event, map[string]interface{}) *config.CustomMessage
}

// BaseHandler basic implement
type BaseHandler struct {
	Configs *sync.Map
}

func (h *BaseHandler) GetConfig(id string) (cfg *config.MessageHandlerConfig) {
	if v, ok := h.Configs.Load(id); ok {
		cfg = v.(*config.MessageHandlerConfig)
	}
	return
}

func (h *BaseHandler) RegisterConfig(cfg *config.MessageHandlerConfig) (err error) {
	h.Configs.Store(cfg.ID, cfg)
	return
}

func (h *BaseHandler) DeregisterConfig(id string) (err error) {
	h.Configs.Delete(id)
	return
}

func (h *BaseHandler) Run(event linebot.Event, variables map[string]interface{}) (reply *config.CustomMessage) {
	return
}

func (h *BaseHandler) runStage(id string, stageConfig []*config.StageConfig, variables map[string]interface{}) (reply string, err error) {
	for idx, stage := range stageConfig {
		switch stage.Type {
		case "action":

			// 取代參數中的變數值
			var b []byte
			if b, err = json.Marshal(stage.Parameter); err != nil {
				logger.Errorw(err.Error(), "id", id, "stage", idx, "type", stage.Type, "plugin", stage.Plugin)
			}
			ParamData := util.ReplaceVariables(string(b), variables)
			var Parameter interface{}
			if err = json.Unmarshal([]byte(ParamData), &Parameter); err != nil {
				logger.Errorw(err.Error(), "id", id, "stage", idx, "type", stage.Type, "plugin", stage.Plugin)
			}
			// 執行
			variables, err = plugin.Excute(stage.Plugin, Parameter, variables)
			if err != nil {
				logger.Errorw(err.Error(), "id", id, "stage", idx, "type", stage.Type, "plugin", stage.Plugin)
				// Stage 執行失敗，切換到Failed執行的階段
				if stage.Failed != nil {
					reply, err = h.runStage(id, stage.Failed, variables)
					return
				}
			}

		case "reply":
			b, err := json.Marshal(stage.Value)
			if err != nil {
				logger.Errorw(err.Error(), "id", id, "stage", idx, "type", stage.Type)
			}
			reply = util.ReplaceVariables(string(b), variables)

		}

	}
	return
}
