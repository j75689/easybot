package server

import (
	"github.com/j75689/easybot/config"
	"io/ioutil"
	"plugin"
	"strings"

	"go.uber.org/zap"
)

// LoadPlugins 加載所有模組
func LoadPlugins(path string) {

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

			if f, ok := function.(func(interface{}, map[string]interface{}, *zap.SugaredLogger) (map[string]interface{}, error)); ok {
				ff := config.PluginFunc(f)
				pluginfuncs.Store(runFuncName, &ff)
			} else {
				logger.Errorf("load plugin [%s] failed.\n", runFuncName)
			}

		}
	}
}
