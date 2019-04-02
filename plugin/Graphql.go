package plugin

import (
	"context"
	"encoding/json"

	"github.com/j75689/easybot/pkg/util"

	"github.com/machinebox/graphql"
	"go.uber.org/zap"
)

type GraphqlPluginConfig struct {
	APIURL    string                 `json:"apiURL"`
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
	Output    map[string]string      `json:"output"`
}

func Graphql(input interface{}, variables map[string]interface{}, logger *zap.SugaredLogger) (map[string]interface{}, bool, error) {
	logger.Info("Excute Graphql Plugin")
	var (
		config GraphqlPluginConfig
		next   = true
		err    error
	)

	err = json.Unmarshal(util.GetJSONBytes(input), &config)
	if err != nil {
		logger.Error(err)
		return nil, next, err
	}

	logger.Debug(config)

	client := graphql.NewClient(config.APIURL)

	req := graphql.NewRequest(config.Query)

	for key, variable := range config.Variables {
		req.Var(key, variable)
	}

	req.Header.Set("Cache-Control", "no-cache")

	ctx := context.Background()
	var resp map[string]interface{}
	err = client.Run(ctx, req, &resp)

	if err != nil {
		logger.Error(err)
	}

	for k, v := range config.Output {
		variables[k] = util.GetJSONValue(v, resp)
	}

	logger.Debug(resp)
	return variables, next, err
}
