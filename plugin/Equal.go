package plugin

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type EqualPluginConfig struct {
	Target string `json:"target"`
	Value  string `json:"value"`
}

func Equal(input interface{}, variables map[string]interface{}, logger *zap.SugaredLogger) (map[string]interface{}, error) {
	var err error
	logger.Info("Excute Equal Plugin")
	var config EqualPluginConfig
	param, err := json.Marshal(input)

	if err != nil {
		logger.Error(err)
	}

	logger.Debug(string(param))

	err = json.Unmarshal(param, &config)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Debug(config)

	if config.Target != config.Value {
		err = fmt.Errorf("No Match target:[%v],value:[%v]", config.Target, config.Value)
	}

	return variables, err
}
