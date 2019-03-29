package handler

import (
	"sync"

	"github.com/j75689/easybot/config"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Matcher interface {
	Add(cfg *config.MessageHandlerConfig)
	Remove(cfg *config.MessageHandlerConfig)
	Find(message interface{}) (cft *config.MessageHandlerConfig)
}

type TextMatcher struct {
	store *sync.Map
}

func (m *TextMatcher) Add(cfg *config.MessageHandlerConfig) {
	match := cfg.Match.([]interface{})
	for _, target := range match {
		(*m).store.Store(target, cfg)
	}
}
func (m *TextMatcher) Remove(cfg *config.MessageHandlerConfig) {
	match := cfg.Match.([]string)
	for _, target := range match {
		(*m).store.Delete(target)
	}

}
func (m *TextMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	if v, ok := m.store.Load(message.(string)); ok {
		cfg = v.(*config.MessageHandlerConfig)
	}
	return
}

type ImageMatcher struct {
	store *sync.Map
}

func (m *ImageMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *ImageMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *ImageMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

type VideoMatcher struct {
	store *sync.Map
}

func (m *VideoMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *VideoMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *VideoMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

type AudioMatcher struct {
	store *sync.Map
}

func (m *AudioMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *AudioMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *AudioMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

type FileMatcher struct {
	store *sync.Map
}

func (m *FileMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *FileMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *FileMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

type LocationMatcher struct {
	store *sync.Map
}

func (m *LocationMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *LocationMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *LocationMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

type StickerMatcher struct {
	store *sync.Map
}

func (m *StickerMatcher) Add(cfg *config.MessageHandlerConfig) {
}
func (m *StickerMatcher) Remove(cfg *config.MessageHandlerConfig) {
}
func (m *StickerMatcher) Find(message interface{}) (cfg *config.MessageHandlerConfig) {
	return
}

// MessageHandler struct
type MessageHandler struct {
	BaseHandler
	Event             linebot.EventType
	MessageTypeMapper map[linebot.MessageType]Matcher
}

func (h *MessageHandler) RegisterConfig(cfg *config.MessageHandlerConfig) (err error) {
	h.BaseHandler.RegisterConfig(cfg)
	if matcher := h.MessageTypeMapper[linebot.MessageType(cfg.MessageType)]; matcher != nil {

		matcher.Add(cfg)
	}
	return
}

func (h *MessageHandler) DeregisterConfig(id string) (err error) {
	h.BaseHandler.DeregisterConfig(id)
	if cfg := h.GetConfig(id); cfg != nil {
		if matcher := h.MessageTypeMapper[linebot.MessageType(cfg.MessageType)]; matcher != nil {
			matcher.Remove(cfg)
		}
	}
	return
}

func (h *MessageHandler) Run(event *linebot.Event, variables map[string]interface{}) (reply *config.CustomMessage, err error) {
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
	return
}

func (h *MessageHandler) handleTextMessage(message string, variables *map[string]interface{}) (reply *config.CustomMessage, err error) {
	if cfg := h.MessageTypeMapper[linebot.MessageTypeText].Find(message); cfg != nil {
		// 塞入預設變數值
		for k, v := range cfg.DefaultValues {
			(*variables)[k] = v
		}
		var replyStr string
		replyStr, err = h.runStage(cfg.ID, cfg.Stage, (*variables))
		reply = &config.CustomMessage{
			Msg: replyStr,
		}
	}
	return
}

func newMessageHandler() *MessageHandler {

	return &MessageHandler{
		BaseHandler: BaseHandler{
			Configs: &sync.Map{},
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
