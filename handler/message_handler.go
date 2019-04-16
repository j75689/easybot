package handler

import (
	"sync"

	"github.com/j75689/easybot/model"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Matcher interface {
	Add(cfg *model.MessageHandlerConfig)
	Remove(cfg *model.MessageHandlerConfig)
	Find(message interface{}) (cft *model.MessageHandlerConfig)
}

type TextMatcher struct {
	store *sync.Map
}

func (m *TextMatcher) Add(cfg *model.MessageHandlerConfig) {
	if v := cfg.Match; v != nil {
		match := v.([]interface{})
		for _, target := range match {
			(*m).store.Store(target, cfg)
		}
	}

}
func (m *TextMatcher) Remove(cfg *model.MessageHandlerConfig) {
	if v := cfg.Match; v != nil {
		match := v.([]interface{})
		for _, target := range match {
			(*m).store.Delete(target)
		}
	}

}
func (m *TextMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	if v, ok := m.store.Load(message.(string)); ok {
		cfg = v.(*model.MessageHandlerConfig)
	}
	return
}

type ImageMatcher struct {
	store *sync.Map
}

func (m *ImageMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *ImageMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *ImageMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

type VideoMatcher struct {
	store *sync.Map
}

func (m *VideoMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *VideoMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *VideoMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

type AudioMatcher struct {
	store *sync.Map
}

func (m *AudioMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *AudioMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *AudioMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

type FileMatcher struct {
	store *sync.Map
}

func (m *FileMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *FileMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *FileMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

type LocationMatcher struct {
	store *sync.Map
}

func (m *LocationMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *LocationMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *LocationMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

type StickerMatcher struct {
	store *sync.Map
}

func (m *StickerMatcher) Add(cfg *model.MessageHandlerConfig) {
}
func (m *StickerMatcher) Remove(cfg *model.MessageHandlerConfig) {
}
func (m *StickerMatcher) Find(message interface{}) (cfg *model.MessageHandlerConfig) {
	return
}

// MessageHandler struct
type MessageHandler struct {
	BaseHandler
	Event             linebot.EventType
	DefaultConfig     *model.MessageHandlerConfig
	MessageTypeMapper map[linebot.MessageType]Matcher
}

func (h *MessageHandler) RegisterConfig(cfg *model.MessageHandlerConfig) (err error) {
	h.BaseHandler.RegisterConfig(cfg)
	if matcher := h.MessageTypeMapper[linebot.MessageType(cfg.MessageType)]; matcher != nil {
		matcher.Add(cfg)
	}
	if cfg.Match == nil && cfg.MessageType == "" { // default Config
		h.DefaultConfig = cfg
	}
	return
}

func (h *MessageHandler) DeregisterConfig(id string) (err error) {
	if cfg := h.GetConfig(id); cfg != nil {
		if matcher := h.MessageTypeMapper[linebot.MessageType(cfg.MessageType)]; matcher != nil {
			matcher.Remove(cfg)
		}
		if h.DefaultConfig != nil {
			if cfg.ConfigID == h.DefaultConfig.ConfigID { // default Config
				h.DefaultConfig = nil
			}
		}
	}
	h.BaseHandler.DeregisterConfig(id)
	return
}

func (h *MessageHandler) Run(event *linebot.Event, variables map[string]interface{}) (reply *CustomMessage, err error) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		reply, err = h.handleTextMessage(message.Text, &variables)
	// not implement
	case *linebot.VideoMessage:
	case *linebot.ImageMessage:
	case *linebot.AudioMessage:
	case *linebot.FileMessage:
	case *linebot.LocationMessage:
	case *linebot.StickerMessage:
	}

	if reply == nil { // find wating queue
		reply, err = h.BaseHandler.Run(event, variables)
	}

	if reply == nil && h.DefaultConfig != nil { // run default config
		// add defaultValue
		for k, v := range h.DefaultConfig.DefaultValues {
			variables[k] = v
		}
		var replyStr string
		replyStr, err = h.runStage(h.DefaultConfig.ConfigID, 0, h.DefaultConfig.Stage, variables)
		reply = &CustomMessage{
			Msg: replyStr,
		}
	}

	return
}

func (h *MessageHandler) handleTextMessage(message string, variables *map[string]interface{}) (reply *CustomMessage, err error) {
	if cfg := h.MessageTypeMapper[linebot.MessageTypeText].Find(message); cfg != nil {
		// add defaultValue
		for k, v := range cfg.DefaultValues {
			(*variables)[k] = v
		}
		var replyStr string
		replyStr, err = h.runStage(cfg.ConfigID, 0, cfg.Stage, (*variables))
		reply = &CustomMessage{
			Msg: replyStr,
		}
	}
	return
}

func newMessageHandler() *MessageHandler {

	return &MessageHandler{
		BaseHandler: BaseHandler{
			Configs: &sync.Map{},
			Wating:  &sync.Map{},
		},
		Event: linebot.EventTypeMessage,
		MessageTypeMapper: map[linebot.MessageType]Matcher{
			linebot.MessageTypeText: &TextMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeImage: &ImageMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeVideo: &VideoMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeAudio: &AudioMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeFile: &FileMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeLocation: &LocationMatcher{
				store: &sync.Map{},
			},
			linebot.MessageTypeSticker: &StickerMatcher{
				store: &sync.Map{},
			},
		},
	}
}
