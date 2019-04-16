package model

// MessageHandlerConfig 處理傳入訊息的設定
type MessageHandlerConfig struct {
	ID            string                 `json:"-" bson:"_id"`
	ConfigID      string                 `json:"id" bson:"id"`
	EventType     string                 `json:"eventType" bson:"eventType"`
	MessageType   string                 `json:"messagetype,omitempty" bson:"messagetype,omitempty"`
	DefaultValues map[string]interface{} `json:"defaultValues" bson:"defaultValues"`
	Match         interface{}            `json:"match,omitempty" bson:"match,omitempty"`
	Stage         []*StageConfig         `json:"stage" bson:"stage"`
}

// StageConfig 處理的執行步驟
type StageConfig struct {
	Type        string            `json:"type" bson:"type"` // reply,action,wait
	Timeout     int               `json:"timeout,omitempty" bson:"timeout,omitempty"`
	Target      *Target           `json:"target,omitempty" bson:"target,omitempty"`
	Plugin      string            `json:"plugin,omitempty" bson:"plugin,omitempty"`
	Parameter   interface{}       `json:"parameter,omitempty" bson:"parameter,omitempty"`
	Value       interface{}       `json:"value,omitempty" bson:"value,omitempty"`
	Extract     map[string]string `json:"extract,omitempty" bson:"extract,omitempty"`
	PreFunction []*StageConfig    `json:"prefunction,omitempty" bson:"prefunction,omitempty"`
	Failed      []*StageConfig    `json:"failed,omitempty" bson:"failed,omitempty"`
}

// Target event
type Target struct {
	EventType   string `json:"eventType" bson:"eventType"`
	MessageType string `json:"messagetype,omitempty" bson:"messagetype,omitempty"`
}
