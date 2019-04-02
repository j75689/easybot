package plugin

import (
	"fmt"
	"io/ioutil"
	"plugin"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// PluginFunc 統一插件進入點
type PluginFunc func(interface{}, map[string]interface{}, *zap.SugaredLogger) (map[string]interface{}, bool, error)

var (
	pluginfuncs = &sync.Map{}
	logger      *zap.SugaredLogger
)

func Load(path string, log *zap.SugaredLogger) {
	logger = log
	logger.Info("load plugin")
	// add default plugin
	{
		graphql := PluginFunc(Graphql)
		pluginfuncs.Store("Graphql", &graphql)
		equal := PluginFunc(Equal)
		pluginfuncs.Store("Equal", &equal)
		curl := PluginFunc(Curl)
		pluginfuncs.Store("Curl", &curl)
	}
	// load addition plugin
	load(path)
}

// load all .so plugin file
func load(path string) {

	// Fix Path
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// 讀取目錄
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err)
		return
	}

	// 載入
	for _, f := range files {
		if !f.IsDir() {
			var runFuncName = f.Name()

			if !strings.HasSuffix(f.Name(), ".so") {
				continue
			}

			if strings.LastIndexAny(runFuncName, ".") > -1 {
				runFuncName = runFuncName[0:strings.LastIndexAny(runFuncName, ".")]
			}

			p, err := plugin.Open(path + f.Name())
			if err != nil {
				logger.Error(err)
				continue
			}

			function, err := p.Lookup(runFuncName)
			if err != nil {
				logger.Error(err)
				continue
			}

			if f, ok := function.(func(interface{}, map[string]interface{}, *zap.SugaredLogger) (map[string]interface{}, bool, error)); ok {
				ff := PluginFunc(f)
				pluginfuncs.Store(runFuncName, &ff)
			} else {
				logger.Errorf("load plugin [%s] failed.\n", runFuncName)
			}

		}
	}
}

// Excute plugin
func Excute(pluginName string, input interface{}, variables map[string]interface{}) (map[string]interface{}, bool, error) {
	if v, ok := pluginfuncs.Load(pluginName); ok {
		// 執行
		plugin := v.(*PluginFunc)
		return (*plugin)(input, variables, logger)
	}
	return variables, true, fmt.Errorf("plugin [%s] not found.\n", pluginName)
}
