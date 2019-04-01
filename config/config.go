package config

// MessageHandlerConfig 處理傳入訊息的設定
type MessageHandlerConfig struct {
	ID            string            `json:"id"`
	EventType     string            `json:"eventType"`
	MessageType   string            `json:"messagetype,omitempty"`
	DefaultValues map[string]string `json:"defaultValues"`
	Match         interface{}       `json:"match,omitempty"`
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
