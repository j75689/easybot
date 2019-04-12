package plugin

import (
	"encoding/json"

	"go.uber.org/zap"
)

type EqualPluginConfig struct {
	Target interface{} `json:"target"`
	Value  interface{} `json:"value"`
}

func Equal(input interface{}, variables map[string]interface{}, logger *zap.SugaredLogger) (map[string]interface{}, bool, error) {
	var (
		err  error
		next = false
	)
	logger.Info("[plugin] ", "Excute Equal Plugin")
	var config EqualPluginConfig
	param, err := json.Marshal(input)

	if err != nil {
		logger.Error("[plugin] ", err)
	}

	logger.Debug("[plugin] ", string(param))

	err = json.Unmarshal(param, &config)
	if err != nil {
		logger.Error("[plugin] ", err)
		return nil, next, err
	}

	logger.Debug("[plugin] ", config)

	if config.Target == config.Value {
		next = true
	}

	return variables, next, err
}
