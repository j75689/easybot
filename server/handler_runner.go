package server

import (
	"github.com/j75689/easybot/config"
	"github.com/j75689/easybot/pkg/util"
	"encoding/json"
	"fmt"
)

func registerMessageHandlerConfig(configID string, data []byte) (err error) {
	var messageConfig config.MessageHandlerConfig
	err = json.Unmarshal(data, &messageConfig)
	if err != nil {
		return
	}
	if configID != messageConfig.ID {
		return fmt.Errorf("Config ID not match input:[%s],file:[%s]", configID, messageConfig.ID)
	}
	switch messageConfig.MessageType {
	case "TextMessage":
		TextMessageHandleConfig.Store(messageConfig.Match, messageConfig) // MessageText 為key
	}
	return
}

func unregisterMessageHandlerConfig(messageConfig *config.MessageHandlerConfig) {
	TextMessageHandleConfig.Delete(messageConfig.Match)
}

func runhandler(messageConfig *config.MessageHandlerConfig, variables map[string]interface{}) (reply string) {
	var (
		err error
	)
	// 塞入預設變數值
	for k, v := range messageConfig.DefaultValues {
		variables[k] = v
	}

	reply, err = runStage(messageConfig.ID, messageConfig.Stage, variables)
	if err != nil {
		logger.Errorw(err.Error(), "id", messageConfig.ID)
	}

	return
}

func runStage(id string, stageConfig []*config.StageConfig, variables map[string]interface{}) (reply string, err error) {
	for idx, stage := range stageConfig {
		switch stage.Type {
		case "action":
			if v, ok := pluginfuncs.Load(stage.Plugin); ok {
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
				plugin := v.(*config.PluginFunc)
				variables, err = (*plugin)(Parameter, variables, logger)
				if err != nil {
					logger.Errorw(err.Error(), "id", id, "stage", idx, "type", stage.Type, "plugin", stage.Plugin)
					// Stage 執行失敗，切換到Failed執行的階段
					if stage.Failed != nil {
						reply, err = runStage(id, stage.Failed, variables)
						return
					}
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
