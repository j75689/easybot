package plugin

import (
	"io/ioutil"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func Request(input interface{}, variables map[string]interface{}, logger *zap.SugaredLogger) (map[string]interface{}, error) {
	req, err := http.NewRequest(variables["method"].(string), variables["url"].(string), strings.NewReader(variables["data"].(string)))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	ioutil.ReadAll(res.Body)
	return variables, nil
}
