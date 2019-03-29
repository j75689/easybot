package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/j75689/easybot/config"
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

// RegisterConfig to Handler
func RegisterConfig(cfg *config.MessageHandlerConfig) error {
	if handler := handlers[linebot.EventType(cfg.EventType)]; handler != nil {
		return handler.RegisterConfig(cfg)
	}
	return fmt.Errorf("eventType:[%s] handler not found.", cfg.EventType)
}

// DeregisterConfig from Handler
func DeregisterConfig(cfg *config.MessageHandlerConfig) error {

	if handler := handlers[linebot.EventType(cfg.EventType)]; handler != nil {
		return handler.DeregisterConfig(cfg.ID)
	}
	return fmt.Errorf("eventType:[%s] handler not found.", cfg.EventType)
}

// Excute Function
func Execute(event *linebot.Event) (reply *config.CustomMessage, err error) {
	var (
		variables = make(map[string]interface{})
	)
	variables["event"] = structs.Map(event)

	if handler := handlers[event.Type]; handler != nil {
		reply, err = handler.Run(event, variables)
	}
	return
}

// Handler interface
type Handler interface {
	GetConfig(string) *config.MessageHandlerConfig
	RegisterConfig(*config.MessageHandlerConfig) error
	DeregisterConfig(string) error
	Run(*linebot.Event, map[string]interface{}) (*config.CustomMessage, error)
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
	for _, stage := range stageConfig {
		switch stage.Type {
		case "action":
			var (
				b         []byte
				Parameter interface{}
				next      bool
			)
			// 取代參數中的變數值
			b, _ = json.Marshal(stage.Parameter)
			ParamData := util.ReplaceVariables(string(b), variables)
			json.Unmarshal([]byte(ParamData), &Parameter)
			// 執行
			variables, next, err = plugin.Excute(stage.Plugin, Parameter, variables)
			if !next {
				if stage.Failed != nil {
					reply, err = h.runStage(id, stage.Failed, variables)
					return
				}
			}

		case "reply":
			var b []byte
			b, err = json.Marshal(stage.Value)
			reply = util.ReplaceVariables(string(b), variables)

		}

	}
	return
}
