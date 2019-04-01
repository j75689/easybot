package config

// MessageHandlerConfig 處理傳入訊息的設定
type MessageHandlerConfig struct {
	ID            string            `json:"id"`
	EventType     string            `json:"eventType"`
	MessageType   string            `json:"messagetype,omitempty"`
	DefaultValues map[string]string `json:"defaultValues"`
	Match         interface{}       `json:"match,omitempty"`
	Stage         []*StageConfig    `json:"stage"`
}

// StageConfig 處理的執行步驟
type StageConfig struct {
	Type        string            `json:"type"` // reply,action,wait
	Timeout     int               `json:"timeout,omitempty"`
	Target      *Target           `json:"target,omitempty"`
	Plugin      string            `json:"plugin,omitempty"`
	Parameter   interface{}       `json:"parameter,omitempty"`
	Value       interface{}       `json:"value,omitempty"`
	Extract     map[string]string `json:"extract,omitempty"`
	PreFunction []*StageConfig    `json:"prefunction,omitempty"`
	Failed      []*StageConfig    `json:"failed,omitempty"`
}

// Target event
type Target struct {
	EventType   string `json:"eventType"`
	MessageType string `json:"messagetype,omitempty"`
}
