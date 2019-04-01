package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

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
	handlers[linebot.EventTypeFollow] = newFollowHandler()
	handlers[linebot.EventTypeUnfollow] = newUnfollowHandler()
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

// WatingGroup thread safe process list
type WatingGroup struct {
	sync.RWMutex
	process []*WatingProcess
}

// Get process from list
func (ins *WatingGroup) Get(event *linebot.Event) (item *WatingProcess) {
	ins.RLock()
	defer ins.RUnlock()
	for _, v := range ins.process {
		if len(v.Config.Stage) > v.StageIndex {
			stage := v.Config.Stage[v.StageIndex]
			if linebot.EventType(stage.Target.EventType) == event.Type {
				if event.Type == linebot.EventTypeMessage {
					var messageType linebot.MessageType
					switch event.Message.(type) {
					case *linebot.TextMessage:
						messageType = linebot.MessageTypeText
					case *linebot.VideoMessage:
						messageType = linebot.MessageTypeVideo
					case *linebot.ImageMessage:
						messageType = linebot.MessageTypeImage
					case *linebot.AudioMessage:
						messageType = linebot.MessageTypeAudio
					case *linebot.FileMessage:
						messageType = linebot.MessageTypeFile
					case *linebot.LocationMessage:
						messageType = linebot.MessageTypeLocation
					case *linebot.StickerMessage:
						messageType = linebot.MessageTypeSticker
					}
					if messageType == linebot.MessageType(stage.Target.MessageType) {
						return v
					}
				} else {
					return v
				}
			}
		}
	}
	return
}

// Add process to list
func (ins *WatingGroup) Add(item *WatingProcess, ctx context.Context) {
	ins.Lock()
	defer ins.Unlock()
	// timeout
	go func(ins *WatingGroup, item *WatingProcess) {
		select {
		case <-ctx.Done():
			ins.Delete(item)
			return
		}
	}(ins, item)
	ins.process = append(ins.process, item)
}

// Delete process form list
func (ins *WatingGroup) Delete(item *WatingProcess) {
	ins.Lock()
	defer ins.Unlock()

	for idx, v := range ins.process {
		if reflect.DeepEqual(v, item) {
			if idx+1 == len(ins.process) {
				ins.process = ins.process[0:idx]
			} else {
				ins.process = append(ins.process[0:idx], ins.process[idx+1:len(ins.process)]...)
			}
			break
		}
	}
}

// WatingProcess status
type WatingProcess struct {
	cancel     *context.CancelFunc
	Config     *config.MessageHandlerConfig
	StageIndex int
	Variables  map[string]interface{}
}

func (ins *WatingProcess) Cancel() {
	if ins.cancel != nil {
		(*ins.cancel)()
	}
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
	Wating  *sync.Map // userID::status
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

func (h *BaseHandler) Run(event *linebot.Event, variables map[string]interface{}) (reply *config.CustomMessage, err error) {
	// process wating group
	if v, ok := h.Wating.Load(event.Source.UserID); ok {
		group := v.(*WatingGroup)
		if process := group.Get(event); process != nil {
			stage := process.Config.Stage[process.StageIndex]
			if stage.Extract != nil {
				refs := map[string]interface{}{
					"event": structs.Map(event),
				}
				for k, v := range stage.Extract {
					process.Variables[k] = util.GetJSONValue(v, refs)
				}
			}
			if len(process.Config.Stage) > process.StageIndex+1 {
				var replyStr string
				replyStr, err = h.runStage(process.Config.ID, process.StageIndex+1, process.Config.Stage, process.Variables)
				reply = &config.CustomMessage{
					Msg: replyStr,
				}
			}

			// process done
			process.Cancel()
		}
	}
	return
}

func (h *BaseHandler) runStage(id string, startIndex int, stageConfig []*config.StageConfig, variables map[string]interface{}) (reply string, err error) {
	for idx := startIndex; idx < len(stageConfig); idx++ {
		stage := stageConfig[idx]
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
					reply, err = h.runStage(id, 0, stage.Failed, variables)
					return
				}
			}
		case "wait":
			var (
				group  *WatingGroup
				userID = util.GetJSONValue("event.Source.UserID", variables)
			)
			if v, ok := h.Wating.Load(userID); ok {
				group = v.(*WatingGroup)
			} else {
				group = new(WatingGroup)
				h.Wating.Store(userID, group)
			}
			timeout := stage.Timeout
			if timeout <= 0 {
				timeout = 120 // default 120 second
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
			process := &WatingProcess{
				cancel:     &cancel,
				StageIndex: idx,
				Config:     h.GetConfig(id),
				Variables:  variables,
			}
			group.Add(process, ctx)

			if stage.PreFunction != nil {
				reply, err = h.runStage(id, 0, stage.PreFunction, variables)
			}
			return // abort for wating
		case "reply":
			var b []byte
			b, err = json.Marshal(stage.Value)
			reply = util.ReplaceVariables(string(b), variables)

		}

	}
	return
}
