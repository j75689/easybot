package config

import "go.uber.org/zap"

// PluginFunc 統一插件進入點
type PluginFunc func(interface{}, map[string]interface{}, *zap.SugaredLogger) (map[string]interface{}, error)

// MessageHandlerConfig 處理傳入訊息的設定
type MessageHandlerConfig struct {
	ID            string            `json:"id"`
	EventType     string            `json:"eventType"`
	MessageType   string            `json:"messagetype"`
	DefaultValues map[string]string `json:"defaultValues"`
	Match         interface{}       `json:"match"`
	TimeOut       int               `json:"timeout"`
	Stage         []*StageConfig    `json:"stage"`
}

// StageConfig 處理的執行步驟
type StageConfig struct {
	Type      string         `json:"type"`
	Plugin    string         `json:"plugin"`
	Parameter interface{}    `json:"parameter"`
	Value     interface{}    `json:"value"`
	Failed    []*StageConfig `json:"failed"`
}
